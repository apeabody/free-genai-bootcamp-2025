package api

import (
	"testing"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/db"
	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/testutil"
)

func TestGetLastStudySession(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/dashboard/last_study_session", GetLastStudySession)

	// Test successful case
	w := testutil.MakeRequest(r, "GET", "/api/dashboard/last_study_session", nil)
	testutil.AssertStatus(t, w, 200)

	var response map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	// Check all required fields
	requiredFields := []string{
		"id",
		"activity_name",
		"group_name",
		"start_time",
		"end_time",
		"review_items_count",
		"correct_count",
		"incorrect_count",
		"score",
	}

	for _, field := range requiredFields {
		if response[field] == nil {
			t.Errorf("Expected %s in response", field)
		}
	}

	// Test specific values
	if response["id"] != float64(1) {
		t.Errorf("Expected last session ID to be 1, got %v", response["id"])
	}

	if response["activity_name"] != "Flashcards" {
		t.Errorf("Expected activity name to be 'Flashcards', got %v", response["activity_name"])
	}

	if response["score"] == float64(0) {
		t.Error("Expected non-zero score")
	}

	// Test database error
	db.GetDB().Close()
	w = testutil.MakeRequest(r, "GET", "/api/dashboard/last_study_session", nil)
	testutil.AssertStatus(t, w, 500)
}

func TestGetStudyProgress(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/dashboard/study_progress", GetStudyProgress)

	// Test successful case
	w := testutil.MakeRequest(r, "GET", "/api/dashboard/study_progress", nil)
	testutil.AssertStatus(t, w, 200)

	var response map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	// Check all required fields
	requiredFields := []string{
		"total_words",
		"words_studied",
		"average_mastery",
	}

	for _, field := range requiredFields {
		if response[field] == nil {
			t.Errorf("Expected %s in response", field)
		}
	}

	// Test specific values
	if response["total_words"] != float64(3) {
		t.Errorf("Expected total words to be 6, got %v", response["total_words"])
	}

	if response["words_studied"] == float64(0) {
		t.Error("Expected non-zero words studied")
	}

	if response["average_mastery"] == float64(0) {
		t.Error("Expected non-zero average mastery")
	}

	// Test database error
	db.GetDB().Close()
	w = testutil.MakeRequest(r, "GET", "/api/dashboard/study_progress", nil)
	testutil.AssertStatus(t, w, 500)
}

func TestGetQuickStats(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/dashboard/quick_stats", GetQuickStats)

	// Test successful case
	w := testutil.MakeRequest(r, "GET", "/api/dashboard/quick_stats", nil)
	testutil.AssertStatus(t, w, 200)

	var response map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	// Check all required fields
	requiredFields := []string{
		"total_sessions",
		"total_reviews",
		"total_words",
		"words_studied",
		"average_mastery",
	}

	for _, field := range requiredFields {
		if response[field] == nil {
			t.Errorf("Expected %s in response", field)
		}
	}

	// Test specific values
	if response["total_words"] != float64(3) {
		t.Errorf("Expected total words to be 6, got %v", response["total_words"])
	}

	if response["words_studied"] == float64(0) {
		t.Error("Expected non-zero words studied")
	}

	if response["average_mastery"] == float64(0) {
		t.Error("Expected non-zero average mastery")
	}

	if response["total_sessions"] == float64(0) {
		t.Error("Expected non-zero total sessions")
	}

	if response["total_reviews"] == float64(0) {
		t.Error("Expected non-zero total reviews")
	}

	// Test database error
	db.GetDB().Close()
	w = testutil.MakeRequest(r, "GET", "/api/dashboard/quick_stats", nil)
	testutil.AssertStatus(t, w, 500)
}
