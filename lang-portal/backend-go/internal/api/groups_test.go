package api

import (
	"testing"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/db"
	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/testutil"
)

func TestGetGroups(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/groups", GetGroups)

	// Test
	w := testutil.MakeRequest(r, "GET", "/api/groups", nil)

	// Assert
	testutil.AssertStatus(t, w, 200)

	var response []map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	if len(response) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(response))
	}

	// Check first group
	firstGroup := response[0]
	if firstGroup["name"] != "Test Group 1" {
		t.Errorf("Expected group name 'Test Group 1', got %v", firstGroup["name"])
	}
}

func TestGetGroup(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/groups/:id", GetGroup)

	// Test existing group
	w := testutil.MakeRequest(r, "GET", "/api/groups/1", nil)
	testutil.AssertStatus(t, w, 200)

	var response map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	if response["name"] != "Test Group 1" {
		t.Errorf("Expected group name 'Test Group 1', got %v", response["name"])
	}

	// Test group statistics
	stats := response["statistics"].(map[string]interface{})
	if stats["total_word_count"] != float64(2) {
		t.Errorf("Expected total_word_count to be 2, got %v", stats["total_word_count"])
	}

	// Test non-existent group
	w = testutil.MakeRequest(r, "GET", "/api/groups/999", nil)
	testutil.AssertStatus(t, w, 404)
}

func TestGetGroupWords(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/groups/:id/words", GetGroupWords)

	// Test successful case
	w := testutil.MakeRequest(r, "GET", "/api/groups/1/words", nil)
	testutil.AssertStatus(t, w, 200)

	var response []map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	if len(response) != 2 {
		t.Errorf("Expected 2 words in group, got %d", len(response))
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
			if word["spanish"] != "hola" {
				t.Errorf("Expected Spanish translation 'hola', got %v", word["spanish"])
			}
			if word["correct_count"] != float64(1) {
				t.Errorf("Expected 'hello' to have 1 correct review, got %v", word["correct_count"])
			}
			if word["mastery_level"] == float64(0) {
				t.Error("Expected non-zero mastery level for 'hello'")
			}
		case "goodbye":
			foundGoodbye = true
			if word["spanish"] != "adios" {
				t.Errorf("Expected Spanish translation 'adios', got %v", word["spanish"])
			}
			if word["incorrect_count"] != float64(1) {
				t.Errorf("Expected 'goodbye' to have 1 incorrect review, got %v", word["incorrect_count"])
			}
			if word["mastery_level"] != float64(0) {
				t.Error("Expected zero mastery level for 'goodbye'")
			}
		}
	}

	if !foundHello || !foundGoodbye {
		t.Error("Expected to find both 'hello' and 'goodbye' in group words")
	}

	// Test non-existent group
	w = testutil.MakeRequest(r, "GET", "/api/groups/999/words", nil)
	testutil.AssertStatus(t, w, 404)

	// Test invalid group ID format
	w = testutil.MakeRequest(r, "GET", "/api/groups/invalid/words", nil)
	testutil.AssertStatus(t, w, 400)

	// Test database error
	db.GetDB().Close()
	w = testutil.MakeRequest(r, "GET", "/api/groups/1/words", nil)
	testutil.AssertStatus(t, w, 500)
}

func TestGetGroupStudySessions(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := testutil.SetupTestRouter()
	r.GET("/api/groups/:id/study_sessions", GetGroupStudySessions)

	// Test
	w := testutil.MakeRequest(r, "GET", "/api/groups/1/study_sessions", nil)
	testutil.AssertStatus(t, w, 200)

	var response []map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	if len(response) != 1 {
		t.Errorf("Expected 1 study session for group, got %d", len(response))
	}

	// Check session details
	session := response[0]
	if session["activity_name"] != "Flashcards" {
		t.Errorf("Expected activity name 'Flashcards', got %v", session["activity_name"])
	}
}
