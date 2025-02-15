package testutil

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/db"
	"github.com/gin-gonic/gin"
)

// SetupTestDB creates a test database and returns a cleanup function
func SetupTestDB(t *testing.T) func() {
	t.Helper()

	// Create a temporary database file
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Initialize the test database
	if err := db.InitDB(dbPath); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	// Run migrations
	if err := runMigrations(db.GetDB()); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Run test data seeding
	if err := seedTestData(db.GetDB()); err != nil {
		t.Fatalf("Failed to seed test data: %v", err)
	}

	// Return cleanup function
	return func() {
		if err := db.CloseDB(); err != nil {
			fmt.Printf("Error closing DB: %v\n", err)
		}
		os.Remove(dbPath)
	}
}

// SetupTestRouter creates a test Gin router
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// MakeRequest performs a test HTTP request and returns the response
func MakeRequest(r *gin.Engine, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ParseResponse parses the JSON response into the given struct
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder, v interface{}) {
	t.Helper()
	if err := json.NewDecoder(w.Body).Decode(v); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
}

// AssertStatus checks if the response status code matches the expected code
func AssertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if w.Code != expected {
		t.Errorf("Expected status %d, got %d", expected, w.Code)
	}
}

// AssertJSON checks if the response body matches the expected JSON
func AssertJSON(t *testing.T, w *httptest.ResponseRecorder, expected interface{}) {
	t.Helper()

	var got interface{}
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Failed to marshal expected JSON: %v", err)
	}

	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("Failed to marshal response JSON: %v", err)
	}

	if string(expectedJSON) != string(gotJSON) {
		t.Errorf("Expected JSON %s, got %s", expectedJSON, gotJSON)
	}
}

// runMigrations runs the database migrations
func runMigrations(db *sql.DB) error {
	// Read and execute the migration SQL
	migrationSQL, err := os.ReadFile("../db/migrations/001_create_tables.up.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(migrationSQL))
	return err
}

// seedTestData seeds the test database with sample data
func seedTestData(db *sql.DB) error {
	// Insert test data
	_, err := db.Exec(`
		-- Insert test groups
		INSERT INTO groups (id, name) VALUES
		(1, 'Test Group 1'),
		(2, 'Test Group 2');

		-- Insert test words
		INSERT INTO words (id, english, spanish, level) VALUES
		(1, 'hello', 'hola', 'beginner'),
		(2, 'goodbye', 'adios', 'beginner'),
		(3, 'thank you', 'gracias', 'beginner');

		-- Insert test word groups
		INSERT INTO word_groups (word_id, group_id) VALUES
		(1, 1),
		(2, 1);

		-- Insert test study activities
		INSERT INTO study_activities (id, name, description) VALUES
		(1, 'Flashcards', 'Practice with flashcards'),
		(2, 'Quiz', 'Test your knowledge');

		-- Insert test study activity groups
		INSERT INTO study_activity_groups (study_activity_id, group_id) VALUES
		(1, 1),
		(1, 2),
		(2, 1),
		(2, 2);

		-- Insert test study sessions
		INSERT INTO study_sessions (id, study_activity_id, group_id, created_at) VALUES
		(1, 1, 1, '2025-02-14T21:49:02-08:00'),
		(2, 2, 2, '2025-02-14T21:49:02-08:00');

		-- Insert test word review items
		INSERT INTO word_review_items (word_id, study_activity_id, correct, response_time) VALUES
		(1, 1, 1, 1.5),
		(2, 1, 0, 2.0),
		(3, 1, 1, 1.0);
	`)
	return err
}
