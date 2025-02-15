package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/db"
	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateStudySession creates a new study session
func CreateStudySession(c *gin.Context) {
	db := db.GetDB()

	// Parse request body
	var request struct {
		StudyActivityID int       `json:"study_activity_id"`
		StartTime       time.Time `json:"start_time"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if study activity exists and get its groups
	rows, err := db.Query(`
		SELECT group_id
		FROM study_activity_groups
		WHERE study_activity_id = ?
	`, request.StudyActivityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch study activity groups"})
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

	if len(groupIDs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study activity not found or has no groups"})
		return
	}

	// Create study session for the first group
	result, err := db.Exec(`
		INSERT INTO study_sessions (study_activity_id, group_id, created_at)
		VALUES (?, ?, ?)
	`, request.StudyActivityID, groupIDs[0], request.StartTime)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create study session"})
		return
	}

	// Get the ID of the created session
	sessionID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get session ID"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":                sessionID,
		"study_activity_id": request.StudyActivityID,
		"group_id":          groupIDs[0],
		"start_time":        request.StartTime,
	})
}

// GetStudySessions returns a list of all study sessions with pagination
func GetStudySessions(c *gin.Context) {
	db := db.GetDB()

	// Parse pagination parameters
	limit := 20 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
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
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch study sessions"})
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

		// Set end time to 30 minutes after start time
		session.EndTime = session.StartTime.Add(30 * time.Minute)

		sessions = append(sessions, session)
	}

	c.JSON(http.StatusOK, sessions)
}

// GetStudySession returns details for a specific study session
func GetStudySession(c *gin.Context) {
	db := db.GetDB()

	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var session models.StudySessionWithStats
	err = db.QueryRow(`
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
		WHERE ss.id = ?
	`, sessionID).Scan(
		&session.ID,
		&session.ActivityName,
		&session.GroupName,
		&session.StartTime,
		&session.ReviewItemsCount,
		&session.CorrectCount,
		&session.IncorrectCount,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study session not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch study session"})
		return
	}

	// Calculate score as percentage
	if session.ReviewItemsCount > 0 {
		session.Score = (session.CorrectCount * 100) / session.ReviewItemsCount
	}

	// Set end time to 30 minutes after start time
	session.EndTime = session.StartTime.Add(30 * time.Minute)

	c.JSON(http.StatusOK, session)
}

// GetStudySessionWords returns all words reviewed in a specific study session
func GetStudySessionWords(c *gin.Context) {
	db := db.GetDB()

	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	// Get session details
	var studyActivityID, groupID int
	err = db.QueryRow(`
		SELECT study_activity_id, group_id
		FROM study_sessions
		WHERE id = ?
	`, sessionID).Scan(&studyActivityID, &groupID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study session not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch session"})
		return
	}

	// Get words for the group
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
		LEFT JOIN word_review_items wri ON w.id = wri.word_id AND wri.study_activity_id = ?
		WHERE wg.group_id = ?
		GROUP BY w.id
		ORDER BY w.english
	`, studyActivityID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch session words"})
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

// CreateWordReview creates a new word review for a study session
func CreateWordReview(c *gin.Context) {
	db := db.GetDB()

	// Parse parameters
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	wordID, err := strconv.Atoi(c.Param("word_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	// Get session and check if word belongs to the group
	var studyActivityID, groupID int
	err = db.QueryRow(`
		SELECT ss.study_activity_id, ss.group_id
		FROM study_sessions ss
		WHERE ss.id = ?
	`, sessionID).Scan(&studyActivityID, &groupID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study session not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch study session"})
		return
	}

	// Check if word exists and belongs to the group
	var wordExists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM words w
			JOIN word_groups wg ON w.id = wg.word_id
			WHERE w.id = ? AND wg.group_id = ?
		)
	`, wordID, groupID).Scan(&wordExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check word existence"})
		return
	}
	if !wordExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found or does not belong to the group"})
		return
	}

	// Parse request body
	var review struct {
		Correct      bool    `json:"correct"`
		ResponseTime float64 `json:"response_time"`
	}
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate response time
	if review.ResponseTime <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Response time must be positive"})
		return
	}

	// Create word review item
	result, err := db.Exec(`
		INSERT INTO word_review_items (word_id, study_activity_id, correct, response_time)
		VALUES (?, ?, ?, ?)
	`, wordID, studyActivityID, review.Correct, review.ResponseTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create word review"})
		return
	}

	// Check if the insert was successful
	if _, err := result.LastInsertId(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	// Calculate new mastery level
	var masteryLevel float64
	err = db.QueryRow(`
		SELECT COALESCE(AVG(CASE WHEN correct THEN 1.0 ELSE 0.0 END), 0) as mastery_level
		FROM word_review_items
		WHERE word_id = ?
	`, wordID).Scan(&masteryLevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate mastery level"})
		return
	}

	// Return review result
	reviewResult := models.WordReviewResult{
		WordID:          wordID,
		SessionID:       sessionID,
		Correct:         review.Correct,
		ResponseTime:    review.ResponseTime,
		NewMasteryLevel: masteryLevel,
	}

	c.JSON(http.StatusOK, reviewResult)
}

func SetupStudySessionAPI(router *gin.RouterGroup) {
	router.POST("/study-sessions", CreateStudySession)
	router.GET("/study-sessions", GetStudySessions)
	router.GET("/study-sessions/:id", GetStudySession)
	router.GET("/study-sessions/:id/words", GetStudySessionWords)
	router.POST("/study-sessions/:id/words/:word_id/reviews", CreateWordReview)
}
