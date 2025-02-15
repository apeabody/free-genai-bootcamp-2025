package api

import (
	"testing"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/testutil"
)

func TestGetWords(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/words", GetWords)

	// Test
	w := testutil.MakeRequest(r, "GET", "/api/words", nil)

	// Assert
	testutil.AssertStatus(t, w, 200)

	var response struct {
		Words []map[string]interface{} `json:"words"`
	}
	testutil.ParseResponse(t, w, &response)

	if len(response.Words) != 3 {
		t.Errorf("Expected 3 words, got %d", len(response.Words))
	}

	// Check first word
	firstWord := response.Words[0]
	if firstWord["english"] != "hello" || firstWord["spanish"] != "hola" {
		t.Errorf("Unexpected word data: %v", firstWord)
	}
}

func TestGetWord(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/words/:id", GetWord)

	// Test existing word
	w := testutil.MakeRequest(r, "GET", "/api/words/1", nil)
	testutil.AssertStatus(t, w, 200)

	var response map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	if response["english"] != "hello" || response["spanish"] != "hola" {
		t.Errorf("Unexpected word data: %v", response)
	}

	// Test non-existent word
	w = testutil.MakeRequest(r, "GET", "/api/words/999", nil)
	testutil.AssertStatus(t, w, 404)
}

func TestGetWordStats(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/words/:id", GetWord)

	// Test word with review history
	w := testutil.MakeRequest(r, "GET", "/api/words/1", nil)
	testutil.AssertStatus(t, w, 200)

	var response map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	// Word 1 has one correct review
	if response["correct_count"] != float64(1) {
		t.Errorf("Expected correct_count to be 1, got %v", response["correct_count"])
	}
	if response["incorrect_count"] != float64(0) {
		t.Errorf("Expected incorrect_count to be 0, got %v", response["incorrect_count"])
	}
	if response["mastery_level"] != float64(1) {
		t.Errorf("Expected mastery_level to be 1, got %v", response["mastery_level"])
	}
}
