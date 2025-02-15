package api

import (
	"fmt"
	"net/http"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/db"
	"github.com/gin-gonic/gin"
)

// ResetHistory deletes all study sessions and word review items
func ResetHistory(c *gin.Context) {
	db := db.GetDB()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Delete all word review items
	_, err = tx.Exec("DELETE FROM word_review_items")
	if err != nil {
		if err := tx.Rollback(); err != nil {
			fmt.Printf("Error rolling back transaction: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete word review items"})
		return
	}

	// Delete all study sessions
	_, err = tx.Exec("DELETE FROM study_sessions")
	if err != nil {
		if err := tx.Rollback(); err != nil {
			fmt.Printf("Error rolling back transaction: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete study sessions"})
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			fmt.Printf("Error rolling back transaction: %v\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully reset study history",
		"details": "Deleted all study sessions and word review items",
	})
}

// FullReset resets the entire database to its initial state
func FullReset(c *gin.Context) {
	db := db.GetDB()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Delete all data from tables in the correct order
	tables := []string{
		"word_review_items",
		"study_sessions",
		"word_groups",
		"words",
		"groups",
		"study_activities",
	}

	for _, table := range tables {
		_, err = tx.Exec("DELETE FROM " + table)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				fmt.Printf("Error rolling back transaction: %v\n", err)
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data from " + table})
			return
		}
	}

	// Reset auto-increment counters
	for _, table := range tables {
		_, err = tx.Exec("DELETE FROM sqlite_sequence WHERE name = ?", table)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				fmt.Printf("Error rolling back transaction: %v\n", err)
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset auto-increment for " + table})
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

	// Note: Database will need to be manually reseeded using the setup script

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully reset database",
		"details": "Database has been reset to initial state with seed data",
	})
}
