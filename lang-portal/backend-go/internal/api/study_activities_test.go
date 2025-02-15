package api

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/testutil"
)

func TestGetStudyActivity(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/study_activities/:id", GetStudyActivity)

	// Test existing activity
	w := testutil.MakeRequest(r, "GET", "/api/study_activities/1", nil)
	testutil.AssertStatus(t, w, 200)

	var response map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	if response["name"] != "Flashcards" {
		t.Errorf("Expected activity name 'Flashcards', got %v", response["name"])
	}
	if response["description"] != "Practice with flashcards" {
		t.Errorf("Expected description 'Practice with flashcards', got %v", response["description"])
	}

	// Test non-existent activity
	w = testutil.MakeRequest(r, "GET", "/api/study_activities/999", nil)
	testutil.AssertStatus(t, w, 404)
}

func TestGetStudyActivitySessions(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/study_activities/:id/study_sessions", GetStudyActivitySessions)

	// Test
	w := testutil.MakeRequest(r, "GET", "/api/study_activities/1/study_sessions", nil)
	testutil.AssertStatus(t, w, 200)

	var response []map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	if len(response) != 1 {
		t.Errorf("Expected 1 study session, got %d", len(response))
	}

	// Check session details
	session := response[0]
	if session["activity_name"] != "Flashcards" {
		t.Errorf("Expected activity name 'Flashcards', got %v", session["activity_name"])
	}
}

func TestCreateStudyActivity(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.POST("/api/study_activities", CreateStudyActivity)

	// Create test activity data
	activity := models.StudyActivity{
		Name:        "New Activity",
		Description: "Test description",
	}
	body, err := json.Marshal(activity)
	if err != nil {
		t.Fatalf("Failed to marshal activity: %v", err)
	}

	// Test
	w := testutil.MakeRequest(r, "POST", "/api/study_activities", bytes.NewBuffer(body))
	testutil.AssertStatus(t, w, 201)

	var response models.StudyActivity
	testutil.ParseResponse(t, w, &response)

	if response.Name != activity.Name {
		t.Errorf("Expected activity name '%s', got '%s'", activity.Name, response.Name)
	}
	if response.Description != activity.Description {
		t.Errorf("Expected description '%s', got '%s'", activity.Description, response.Description)
	}
	if response.ID == 0 {
		t.Error("Expected non-zero ID for created activity")
	}

	// Test invalid request
	w = testutil.MakeRequest(r, "POST", "/api/study_activities", bytes.NewBufferString("invalid json"))
	testutil.AssertStatus(t, w, 400)
}
