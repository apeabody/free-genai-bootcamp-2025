package api

import (
	"net/http"
	"time"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/db"
	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
	"github.com/gin-gonic/gin"
)

// GetLastStudySession returns details about the user's most recent study session
func GetLastStudySession(c *gin.Context) {
	db := db.GetDB()
	var session models.StudySessionWithStats

	err := db.QueryRow(`
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
		ORDER BY ss.created_at DESC
		LIMIT 1
	`).Scan(
		&session.ID,
		&session.ActivityName,
		&session.GroupName,
		&session.StartTime,
		&session.ReviewItemsCount,
		&session.CorrectCount,
		&session.IncorrectCount,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch last study session"})
		return
	}

	// Calculate score as percentage
	if session.ReviewItemsCount > 0 {
		session.Score = (session.CorrectCount * 100) / session.ReviewItemsCount
	}

	// Set end time to 30 minutes after start time for this example
	session.EndTime = session.StartTime.Add(30 * time.Minute)

	c.JSON(http.StatusOK, session)
}

// GetStudyProgress returns the user's overall study progress
func GetStudyProgress(c *gin.Context) {
	db := db.GetDB()

	var progress struct {
		TotalWords     int     `json:"total_words"`
		WordsStudied   int     `json:"words_studied"`
		AverageMastery float64 `json:"average_mastery"`
	}

	// Get word counts and mastery level
	err := db.QueryRow(`
		SELECT 
			(SELECT COUNT(*) FROM words) as total_words,
			(SELECT COUNT(DISTINCT word_id) FROM word_review_items) as words_studied,
			(SELECT COALESCE(AVG(CASE WHEN correct THEN 1.0 ELSE 0.0 END), 0) FROM word_review_items) as average_mastery
	`).Scan(
		&progress.TotalWords,
		&progress.WordsStudied,
		&progress.AverageMastery,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch word statistics"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetQuickStats returns a quick overview of the user's study statistics
func GetQuickStats(c *gin.Context) {
	db := db.GetDB()

	var stats struct {
		TotalSessions  int     `json:"total_sessions"`
		TotalReviews   int     `json:"total_reviews"`
		TotalWords     int     `json:"total_words"`
		WordsStudied   int     `json:"words_studied"`
		AverageMastery float64 `json:"average_mastery"`
	}

	// Get session and review counts
	err := db.QueryRow(`
		SELECT 
			(SELECT COUNT(*) FROM study_sessions) as total_sessions,
			(SELECT COUNT(*) FROM word_review_items) as total_reviews,
			(SELECT COUNT(*) FROM words) as total_words,
			(SELECT COUNT(DISTINCT word_id) FROM word_review_items) as words_studied,
			(SELECT COALESCE(AVG(CASE WHEN correct THEN 1.0 ELSE 0.0 END), 0) FROM word_review_items) as average_mastery
	`).Scan(
		&stats.TotalSessions,
		&stats.TotalReviews,
		&stats.TotalWords,
		&stats.WordsStudied,
		&stats.AverageMastery,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quick statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
