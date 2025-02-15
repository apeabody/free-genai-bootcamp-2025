package api

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/db"
	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/testutil"
)

func TestGetStudySessions(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/study_sessions", GetStudySessions)

	// Test
	w := testutil.MakeRequest(r, "GET", "/api/study_sessions", nil)
	testutil.AssertStatus(t, w, 200)

	var response []map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	if len(response) != 2 {
		t.Errorf("Expected 2 study sessions, got %d", len(response))
	}

	// Check first session
	firstSession := response[0]
	if firstSession["activity_name"] != "Flashcards" { // Most recent session
		t.Errorf("Expected activity name 'Flashcards', got %v", firstSession["activity_name"])
	}
}

func TestGetStudySession(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/study_sessions/:id", GetStudySession)

	// Test existing session
	w := testutil.MakeRequest(r, "GET", "/api/study_sessions/1", nil)
	testutil.AssertStatus(t, w, 200)

	var response map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	if response["activity_name"] != "Flashcards" {
		t.Errorf("Expected activity name 'Flashcards', got %v", response["activity_name"])
	}

	// Test non-existent session
	w = testutil.MakeRequest(r, "GET", "/api/study_sessions/999", nil)
	testutil.AssertStatus(t, w, 404)
}

func TestGetStudySessionWords(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/study_sessions/:id/words", GetStudySessionWords)

	// Test successful case
	w := testutil.MakeRequest(r, "GET", "/api/study_sessions/1/words", nil)
	testutil.AssertStatus(t, w, 200)

	var response []map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	if len(response) != 2 {
		t.Errorf("Expected 2 words in session, got %d", len(response))
	}

	// Check required fields for each word
	requiredFields := []string{
		"id",
		"english",
		"spanish",
		"correct_count",
		"incorrect_count",
		"mastery_level",
	}

	for _, word := range response {
		for _, field := range requiredFields {
			if word[field] == nil {
				t.Errorf("Expected %s in word response", field)
			}
		}
	}

	// Check specific word data
	foundHello := false
	foundGoodbye := false
	for _, word := range response {
		switch word["english"] {
		case "hello":
			foundHello = true
			if word["correct_count"] != float64(1) {
				t.Errorf("Expected 'hello' to have 1 correct review, got %v", word["correct_count"])
			}
			if word["spanish"] != "hola" {
				t.Errorf("Expected Spanish translation 'hola', got %v", word["spanish"])
			}
		case "goodbye":
			foundGoodbye = true
			if word["incorrect_count"] != float64(1) {
				t.Errorf("Expected 'goodbye' to have 1 incorrect review, got %v", word["incorrect_count"])
			}
			if word["spanish"] != "adios" {
				t.Errorf("Expected Spanish translation 'adios', got %v", word["spanish"])
			}
		}
	}

	if !foundHello || !foundGoodbye {
		t.Error("Expected to find both 'hello' and 'goodbye' in session words")
	}

	// Test non-existent session
	w = testutil.MakeRequest(r, "GET", "/api/study_sessions/999/words", nil)
	testutil.AssertStatus(t, w, 404)

	// Test invalid session ID format
	w = testutil.MakeRequest(r, "GET", "/api/study_sessions/invalid/words", nil)
	testutil.AssertStatus(t, w, 400)

	// Test database error
	db.GetDB().Close()
	w = testutil.MakeRequest(r, "GET", "/api/study_sessions/1/words", nil)
	testutil.AssertStatus(t, w, 500)
}

func TestCreateWordReview(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.POST("/api/study_sessions/:id/words/:word_id/review", CreateWordReview)

	// Test successful case - correct answer
	review := struct {
		Correct      bool    `json:"correct"`
		ResponseTime float64 `json:"response_time"`
	}{
		Correct:      true,
		ResponseTime: 1.2,
	}
	body, err := json.Marshal(review)
	if err != nil {
		t.Fatalf("Failed to marshal review: %v", err)
	}

	w := testutil.MakeRequest(r, "POST", "/api/study_sessions/1/words/1/review", bytes.NewBuffer(body))
	testutil.AssertStatus(t, w, 200)

	var response map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	// Check all required fields
	requiredFields := []string{
		"session_id",
		"word_id",
		"correct",
		"response_time",
	}

	for _, field := range requiredFields {
		if response[field] == nil {
			t.Errorf("Expected %s in response", field)
		}
	}

	// Check specific values
	if !response["correct"].(bool) {
		t.Error("Expected review to be marked as correct")
	}
	if response["response_time"] != float64(1.2) {
		t.Errorf("Expected response time 1.2, got %v", response["response_time"])
	}
	if response["session_id"] != float64(1) {
		t.Errorf("Expected session_id 1, got %v", response["session_id"])
	}
	if response["word_id"] != float64(1) {
		t.Errorf("Expected word_id 1, got %v", response["word_id"])
	}

	// Test incorrect answer
	review.Correct = false
	review.ResponseTime = 2.5
	body, _ = json.Marshal(review)
	w = testutil.MakeRequest(r, "POST", "/api/study_sessions/1/words/2/review", bytes.NewBuffer(body))
	testutil.AssertStatus(t, w, 200)

	testutil.ParseResponse(t, w, &response)
	if response["correct"].(bool) {
		t.Error("Expected review to be marked as incorrect")
	}
	if response["response_time"] != float64(2.5) {
		t.Errorf("Expected response time 2.5, got %v", response["response_time"])
	}

	// Test invalid session ID
	w = testutil.MakeRequest(r, "POST", "/api/study_sessions/999/words/1/review", bytes.NewBuffer(body))
	testutil.AssertStatus(t, w, 404) // Session not found

	// Test invalid word ID
	w = testutil.MakeRequest(r, "POST", "/api/study_sessions/1/words/999/review", bytes.NewBuffer(body))
	testutil.AssertStatus(t, w, 404) // Word not found

	// Test invalid session ID format
	w = testutil.MakeRequest(r, "POST", "/api/study_sessions/invalid/words/1/review", bytes.NewBuffer(body))
	testutil.AssertStatus(t, w, 400)

	// Test invalid word ID format
	w = testutil.MakeRequest(r, "POST", "/api/study_sessions/1/words/invalid/review", bytes.NewBuffer(body))
	testutil.AssertStatus(t, w, 400)

	// Test invalid request body
	w = testutil.MakeRequest(r, "POST", "/api/study_sessions/1/words/1/review", bytes.NewBufferString("invalid json"))
	testutil.AssertStatus(t, w, 400)

	// Test missing required fields
	invalidReview := struct {
		Correct bool `json:"correct"`
	}{
		Correct: true,
	}
	body, _ = json.Marshal(invalidReview)
	w = testutil.MakeRequest(r, "POST", "/api/study_sessions/1/words/1/review", bytes.NewBuffer(body))
	testutil.AssertStatus(t, w, 400)

	// Test database error
	db.GetDB().Close()
	w = testutil.MakeRequest(r, "POST", "/api/study_sessions/1/words/1/review", bytes.NewBuffer(body))
	testutil.AssertStatus(t, w, 500)
}
