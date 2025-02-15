package api

import (
	"testing"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/db"
	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/testutil"
)

func TestResetHistory(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.POST("/api/reset_history", ResetHistory)

	// Verify initial state
	var initialSessionCount int
	err := db.GetDB().QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&initialSessionCount)
	if err != nil {
		t.Fatalf("Failed to count initial sessions: %v", err)
	}
	if initialSessionCount == 0 {
		t.Error("Expected some initial study sessions")
	}

	var initialReviewCount int
	err = db.GetDB().QueryRow("SELECT COUNT(*) FROM word_review_items").Scan(&initialReviewCount)
	if err != nil {
		t.Fatalf("Failed to count initial reviews: %v", err)
	}
	if initialReviewCount == 0 {
		t.Error("Expected some initial word reviews")
	}

	// Test reset
	w := testutil.MakeRequest(r, "POST", "/api/reset_history", nil)
	testutil.AssertStatus(t, w, 200)

	// Verify reset state
	var finalSessionCount int
	err = db.GetDB().QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&finalSessionCount)
	if err != nil {
		t.Fatalf("Failed to count final sessions: %v", err)
	}
	if finalSessionCount != 0 {
		t.Errorf("Expected 0 study sessions after reset, got %d", finalSessionCount)
	}

	var finalReviewCount int
	err = db.GetDB().QueryRow("SELECT COUNT(*) FROM word_review_items").Scan(&finalReviewCount)
	if err != nil {
		t.Fatalf("Failed to count final reviews: %v", err)
	}
	if finalReviewCount != 0 {
		t.Errorf("Expected 0 word reviews after reset, got %d", finalReviewCount)
	}

	// Verify other data remains intact
	tables := []string{"words", "groups", "study_activities"}
	for _, table := range tables {
		var count int
		err = db.GetDB().QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to count %s: %v", table, err)
		}
		if count == 0 {
			t.Errorf("Expected %s to remain after history reset", table)
		}
	}

	// Test resetting an already empty history
	w = testutil.MakeRequest(r, "POST", "/api/reset_history", nil)
	testutil.AssertStatus(t, w, 200)

	// Test transaction rollback on word_review_items deletion error
	db := db.GetDB()
	_, err = db.Exec("DROP TABLE word_review_items")
	if err != nil {
		t.Fatalf("Failed to drop word_review_items table: %v", err)
	}
	w = testutil.MakeRequest(r, "POST", "/api/reset_history", nil)
	testutil.AssertStatus(t, w, 500)

	// Recreate word_review_items table
	_, err = db.Exec(`CREATE TABLE word_review_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		word_id INTEGER NOT NULL,
		study_activity_id INTEGER NOT NULL,
		correct BOOLEAN NOT NULL,
		response_time FLOAT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (word_id) REFERENCES words(id),
		FOREIGN KEY (study_activity_id) REFERENCES study_activities(id)
	)`)
	if err != nil {
		t.Fatalf("Failed to recreate word_review_items table: %v", err)
	}

	// Test transaction rollback on study_sessions deletion error
	_, err = db.Exec("DROP TABLE study_sessions")
	if err != nil {
		t.Fatalf("Failed to drop study_sessions table: %v", err)
	}
	w = testutil.MakeRequest(r, "POST", "/api/reset_history", nil)
	testutil.AssertStatus(t, w, 500)

	// Test database error by closing the connection
	db.Close()
	w = testutil.MakeRequest(r, "POST", "/api/reset_history", nil)
	testutil.AssertStatus(t, w, 500)
}

func TestFullReset(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.POST("/api/full_reset", FullReset)

	// Verify initial state
	tables := []string{
		"word_review_items",
		"study_sessions",
		"word_groups",
		"words",
		"groups",
		"study_activities",
	}

	initialCounts := make(map[string]int)
	for _, table := range tables {
		var count int
		err := db.GetDB().QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to count initial %s: %v", table, err)
		}
		initialCounts[table] = count
		if count == 0 {
			t.Errorf("Expected some initial data in %s", table)
		}
	}

	// Test reset
	w := testutil.MakeRequest(r, "POST", "/api/full_reset", nil)
	testutil.AssertStatus(t, w, 200)

	// Verify reset state
	for _, table := range tables {
		var count int
		err := db.GetDB().QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to count final %s: %v", table, err)
		}
		if count != 0 {
			t.Errorf("Expected 0 records in %s after reset, got %d", table, count)
		}
	}

	// Test transaction rollback by breaking the database schema
	_, err := db.GetDB().Exec("ALTER TABLE words RENAME TO words_temp")
	if err != nil {
		t.Fatalf("Failed to rename words table: %v", err)
	}

	w = testutil.MakeRequest(r, "POST", "/api/full_reset", nil)
	testutil.AssertStatus(t, w, 500)

	// Restore the table
	_, err = db.GetDB().Exec("ALTER TABLE words_temp RENAME TO words")
	if err != nil {
		t.Fatalf("Failed to restore words table: %v", err)
	}

	// Test resetting an already empty database
	w = testutil.MakeRequest(r, "POST", "/api/full_reset", nil)
	testutil.AssertStatus(t, w, 200)

	// Test database error by closing the connection
	db.GetDB().Close()
	w = testutil.MakeRequest(r, "POST", "/api/full_reset", nil)
	testutil.AssertStatus(t, w, 500)
}
