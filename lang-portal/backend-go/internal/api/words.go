package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/db"
	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
	"github.com/gin-gonic/gin"
)

// GetWords returns a paginated list of words with their statistics
func GetWords(c *gin.Context) {
	db := db.GetDB()

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 100
	}
	offset := (page - 1) * perPage

	// Get total count
	var totalItems int
	err := db.QueryRow("SELECT COUNT(*) FROM words").Scan(&totalItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count words"})
		return
	}

	// Get words with stats
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
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		GROUP BY w.id
		ORDER BY w.id
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch words"})
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

	// Calculate pagination metadata
	totalPages := (totalItems + perPage - 1) / perPage

	c.JSON(http.StatusOK, gin.H{
		"words": words,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  totalItems,
			"per_page":     perPage,
		},
	})
}

// GetWord returns details for a specific word
func GetWord(c *gin.Context) {
	db := db.GetDB()

	// Parse word ID
	wordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	// Get word with stats
	var word models.WordWithStats
	err = db.QueryRow(`
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
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE w.id = ?
		GROUP BY w.id
	`, wordID).Scan(
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

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch word"})
		return
	}

	c.JSON(http.StatusOK, word)
}
