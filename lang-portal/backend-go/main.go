package main

import (
	"fmt"
	"os"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api"
	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	if err := db.InitDB("words.db"); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.CloseDB(); err != nil {
			fmt.Printf("Error closing DB: %v\n", err)
		}
	}()

	// Create Gin router
	r := gin.Default()

	// Setup CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup routes
	setupRoutes(r)

	// Start server
	fmt.Println("Server starting on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

func setupRoutes(r *gin.Engine) {
	// Dashboard routes
	r.GET("/api/dashboard/last_study_session", api.GetLastStudySession)
	r.GET("/api/dashboard/study_progress", api.GetStudyProgress)
	r.GET("/api/dashboard/quick_stats", api.GetQuickStats)

	// Study activities routes
	r.GET("/api/study_activities/:id", api.GetStudyActivity)
	r.GET("/api/study_activities/:id/study_sessions", api.GetStudyActivitySessions)
	r.POST("/api/study_activities", api.CreateStudyActivity)

	// Words routes
	r.GET("/api/words", api.GetWords)
	r.GET("/api/words/:id", api.GetWord)

	// Groups routes
	r.GET("/api/groups", api.GetGroups)
	r.GET("/api/groups/:id", api.GetGroup)
	r.GET("/api/groups/:id/words", api.GetGroupWords)
	r.GET("/api/groups/:id/study_sessions", api.GetGroupStudySessions)

	// Study sessions routes
	r.GET("/api/study_sessions", api.GetStudySessions)
	r.GET("/api/study_sessions/:id", api.GetStudySession)
	r.GET("/api/study_sessions/:id/words", api.GetStudySessionWords)
	r.POST("/api/study_sessions/:id/words/:word_id/review", api.CreateWordReview)

	// System management routes
	r.POST("/api/reset_history", api.ResetHistory)
	r.POST("/api/full_reset", api.FullReset)
}
