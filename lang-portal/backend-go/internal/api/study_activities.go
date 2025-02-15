package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/db"
	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
	"github.com/gin-gonic/gin"
)

// GetStudyActivity returns details for a specific study activity
func GetStudyActivity(c *gin.Context) {
	db := db.GetDB()

	activityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	// Get activity details
	var activity models.StudyActivityWithStats
	err = db.QueryRow(`
		SELECT 
			sa.id,
			sa.name,
			sa.description,
			sa.created_at,
			sa.updated_at,
			COUNT(DISTINCT ss.id) as total_sessions
		FROM study_activities sa
		LEFT JOIN study_sessions ss ON sa.id = ss.study_activity_id
		WHERE sa.id = ?
		GROUP BY sa.id
	`, activityID).Scan(
		&activity.ID,
		&activity.Name,
		&activity.Description,
		&activity.CreatedAt,
		&activity.UpdatedAt,
		&activity.TotalSessions,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study activity not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch study activity"})
		return
	}

	// Get associated group IDs
	rows, err := db.Query(`
		SELECT group_id
		FROM study_activity_groups
		WHERE study_activity_id = ?
	`, activityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity groups"})
		return
	}
	defer rows.Close()

	var groupIDs []int
	for rows.Next() {
		var groupID int
		if err := rows.Scan(&groupID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan group ID"})
			return
		}
		groupIDs = append(groupIDs, groupID)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             activity.ID,
		"name":           activity.Name,
		"description":    activity.Description,
		"created_at":     activity.CreatedAt,
		"updated_at":     activity.UpdatedAt,
		"total_sessions": activity.TotalSessions,
		"group_ids":      groupIDs,
	})
}

// GetStudyActivitySessions returns all study sessions for a specific activity
func GetStudyActivitySessions(c *gin.Context) {
	db := db.GetDB()

	activityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	rows, err := db.Query(`
		SELECT 
			ss.id,
			sa.name as activity_name,
			g.name as group_name,
			ss.created_at as start_time,
			(SELECT COUNT(*) FROM word_review_items wri WHERE wri.study_activity_id = ss.study_activity_id) as review_items_count,
			(SELECT COUNT(*) FROM word_review_items wri WHERE wri.study_activity_id = ss.study_activity_id AND wri.correct = 1) as correct_count,
			(SELECT COUNT(*) FROM word_review_items wri WHERE wri.study_activity_id = ss.study_activity_id AND wri.correct = 0) as incorrect_count
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		WHERE ss.study_activity_id = ?
		ORDER BY ss.created_at DESC
	`, activityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity study sessions"})
		return
	}
	defer rows.Close()

	var sessions []models.StudySessionWithStats
	for rows.Next() {
		var session models.StudySessionWithStats
		err := rows.Scan(
			&session.ID,
			&session.ActivityName,
			&session.GroupName,
			&session.StartTime,
			&session.ReviewItemsCount,
			&session.CorrectCount,
			&session.IncorrectCount,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan study session"})
			return
		}

		// Calculate score as percentage
		if session.ReviewItemsCount > 0 {
			session.Score = (session.CorrectCount * 100) / session.ReviewItemsCount
		}

		sessions = append(sessions, session)
	}

	c.JSON(http.StatusOK, sessions)
}

// CreateStudyActivity creates a new study activity
func CreateStudyActivity(c *gin.Context) {
	db := db.GetDB()

	// Parse request
	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		GroupIDs    []int  `json:"group_ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Create study activity
	result, err := tx.Exec(`
		INSERT INTO study_activities (name, description)
		VALUES (?, ?)
	`, request.Name, request.Description)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			fmt.Printf("Error rolling back transaction: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create study activity"})
		return
	}

	activityID, err := result.LastInsertId()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			fmt.Printf("Error rolling back transaction: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get created activity ID"})
		return
	}

	// Link activity to groups
	for _, groupID := range request.GroupIDs {
		_, err = tx.Exec(`
			INSERT INTO study_activity_groups (study_activity_id, group_id)
			VALUES (?, ?)
		`, activityID, groupID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				fmt.Printf("Error rolling back transaction: %v\n", err)
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link activity to group"})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			fmt.Printf("Error rolling back transaction: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          activityID,
		"name":        request.Name,
		"description": request.Description,
		"group_ids":   request.GroupIDs,
	})
}
