package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/db"
	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
	"github.com/gin-gonic/gin"
)

// GetGroups returns a list of all word groups
func GetGroups(c *gin.Context) {
	db := db.GetDB()

	rows, err := db.Query(`
		SELECT 
			g.id,
			g.name,
			g.created_at,
			g.updated_at,
			COUNT(DISTINCT wg.word_id) as total_word_count,
			COUNT(DISTINCT CASE WHEN wri.correct = 1 THEN w.id END) as mastered_words,
			COUNT(DISTINCT CASE WHEN wri.correct = 0 THEN w.id END) as in_progress_words,
			COALESCE(AVG(CASE WHEN wri.correct THEN 1.0 ELSE 0.0 END), 0) as average_mastery
		FROM groups g
		LEFT JOIN word_groups wg ON g.id = wg.group_id
		LEFT JOIN words w ON wg.word_id = w.id
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		GROUP BY g.id
		ORDER BY g.name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}
	defer rows.Close()

	var groups []models.GroupWithStats
	for rows.Next() {
		var group models.GroupWithStats
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.Statistics.TotalWordCount,
			&group.Statistics.MasteredWords,
			&group.Statistics.InProgressWords,
			&group.Statistics.AverageMastery,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan group"})
			return
		}
		groups = append(groups, group)
	}

	c.JSON(http.StatusOK, groups)
}

// GetGroup returns details for a specific group
func GetGroup(c *gin.Context) {
	db := db.GetDB()

	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var group models.GroupWithStats
	err = db.QueryRow(`
		SELECT 
			g.id,
			g.name,
			g.created_at,
			g.updated_at,
			COUNT(DISTINCT wg.word_id) as total_word_count,
			COUNT(DISTINCT CASE WHEN wri.correct = 1 THEN w.id END) as mastered_words,
			COUNT(DISTINCT CASE WHEN wri.correct = 0 THEN w.id END) as in_progress_words,
			COALESCE(AVG(CASE WHEN wri.correct THEN 1.0 ELSE 0.0 END), 0) as average_mastery
		FROM groups g
		LEFT JOIN word_groups wg ON g.id = wg.group_id
		LEFT JOIN words w ON wg.word_id = w.id
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE g.id = ?
		GROUP BY g.id
	`, groupID).Scan(
		&group.ID,
		&group.Name,
		&group.CreatedAt,
		&group.UpdatedAt,
		&group.Statistics.TotalWordCount,
		&group.Statistics.MasteredWords,
		&group.Statistics.InProgressWords,
		&group.Statistics.AverageMastery,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch group"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetGroupWords returns all words in a specific group
func GetGroupWords(c *gin.Context) {
	db := db.GetDB()

	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Check if group exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM groups WHERE id = ?)", groupID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check group existence"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	rows, err := db.Query(`
		SELECT 
			w.id,
			w.english,
			w.spanish,
			w.level,
			w.created_at,
			w.updated_at,
			COUNT(CASE WHEN wri.correct = 1 THEN 1 END) as correct_count,
			COUNT(CASE WHEN wri.correct = 0 THEN 1 END) as incorrect_count,
			COALESCE(AVG(CASE WHEN wri.correct THEN 1.0 ELSE 0.0 END), 0) as mastery_level
		FROM words w
		JOIN word_groups wg ON w.id = wg.word_id
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wg.group_id = ?
		GROUP BY w.id
		ORDER BY w.english
	`, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch group words"})
		return
	}
	defer rows.Close()

	var words []models.WordWithStats
	for rows.Next() {
		var word models.WordWithStats
		err := rows.Scan(
			&word.ID,
			&word.English,
			&word.Spanish,
			&word.Level,
			&word.CreatedAt,
			&word.UpdatedAt,
			&word.CorrectCount,
			&word.IncorrectCount,
			&word.MasteryLevel,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan word"})
			return
		}
		words = append(words, word)
	}

	c.JSON(http.StatusOK, words)
}

// GetGroupStudySessions returns all study sessions for a specific group
func GetGroupStudySessions(c *gin.Context) {
	db := db.GetDB()

	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
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
		WHERE ss.group_id = ?
		ORDER BY ss.created_at DESC
	`, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch group study sessions"})
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

// CreateGroup creates a new word group
func CreateGroup(c *gin.Context) {
	db := db.GetDB()

	// Parse request body
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Insert the group
	result, err := db.Exec(`
		INSERT INTO groups (name, created_at, updated_at)
		VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, group.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	// Get the inserted ID
	groupID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get group ID"})
		return
	}

	// Return the created group
	group.ID = int(groupID)
	c.JSON(http.StatusCreated, group)
}
