package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api"
	"github.com/apeabody/free-genai-bootcamp-2025/lang-portal/backend-go/internal/testutil"
	"github.com/gin-gonic/gin"
)

// encodeJSON is a helper function to encode data to JSON and handle errors
func encodeJSON(t *testing.T, data interface{}) *bytes.Buffer {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		t.Fatalf("Failed to encode JSON data: %v", err)
	}
	return buf
}

func setupRouter() *gin.Engine {
	r := testutil.SetupTestRouter()

	// Groups
	r.GET("/api/groups", api.GetGroups)
	r.GET("/api/groups/:id", api.GetGroup)
	r.POST("/api/groups", api.CreateGroup)

	// Study Activities
	r.GET("/api/study_activities/:id", api.GetStudyActivity)
	r.POST("/api/study_activities", api.CreateStudyActivity)
	r.GET("/api/study_activities/:id/sessions", api.GetStudyActivitySessions)

	// Study Sessions
	r.GET("/api/study_sessions", api.GetStudySessions)
	r.POST("/api/study_sessions", api.CreateStudySession)
	r.GET("/api/study_sessions/:id", api.GetStudySession)
	r.GET("/api/study_sessions/:id/words", api.GetStudySessionWords)
	r.POST("/api/study_sessions/:id/words/:word_id/review", api.CreateWordReview)

	// Dashboard
	r.GET("/api/dashboard/last_session", api.GetLastStudySession)
	r.GET("/api/dashboard/progress", api.GetStudyProgress)
	r.GET("/api/dashboard/quick_stats", api.GetQuickStats)

	// System
	r.POST("/api/reset_history", api.ResetHistory)

	return r
}

func TestCompleteUserFlow(t *testing.T) {
	// Setup
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	r := setupRouter()

	// Step 1: Create a study activity
	activityPayload := map[string]interface{}{
		"name":        "Spanish Basics",
		"description": "Learn basic Spanish words",
		"group_ids":   []int{1}, // Assuming group 1 exists in test data
	}
	activityBody, _ := json.Marshal(activityPayload)
	w := testutil.MakeRequest(r, "POST", "/api/study_activities", bytes.NewBuffer(activityBody))
	testutil.AssertStatus(t, w, http.StatusCreated)
	if w.Code != http.StatusCreated {
		t.Logf("Activity creation response: %s", w.Body.String())
	}

	var activity struct {
		ID int `json:"id"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &activity)
	if err != nil {
		t.Fatalf("Failed to parse activity response: %v, body: %s", err, w.Body.String())
	}
	t.Logf("Created activity with ID: %d", activity.ID)

	// Step 2: Get study activity details
	w = testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/study_activities/%d", activity.ID), nil)
	testutil.AssertStatus(t, w, http.StatusOK)
	if w.Code != http.StatusOK {
		t.Logf("Response body: %s", w.Body.String())
	}

	// Step 3: Create a study session
	sessionPayload := map[string]interface{}{
		"study_activity_id": activity.ID,
		"start_time":        "2025-02-14T21:49:02-08:00",
	}
	sessionBody, _ := json.Marshal(sessionPayload)
	w = testutil.MakeRequest(r, "POST", "/api/study_sessions", bytes.NewBuffer(sessionBody))
	testutil.AssertStatus(t, w, http.StatusCreated)
	if w.Code != http.StatusCreated {
		t.Logf("Response body: %s", w.Body.String())
	}

	var session struct {
		ID int `json:"id"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &session)
	if err != nil {
		t.Fatalf("Failed to parse session response: %v", err)
	}

	// Step 4: Get study sessions for the activity
	w = testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/study_activities/%d/sessions", activity.ID), nil)
	testutil.AssertStatus(t, w, http.StatusOK)

	var sessions []struct {
		ID int `json:"id"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &sessions)
	if err != nil {
		t.Fatalf("Failed to parse sessions response: %v", err)
	}
	if len(sessions) == 0 {
		t.Fatal("Expected at least one study session")
	}

	sessionID := session.ID // Use the session we just created

	// Step 4: Get words for the study session
	w = testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/study_sessions/%d/words", sessionID), nil)
	testutil.AssertStatus(t, w, http.StatusOK)

	var words []struct {
		ID int `json:"id"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &words)
	if err != nil {
		t.Fatalf("Failed to parse words response: %v", err)
	}
	if len(words) == 0 {
		t.Fatal("Expected at least one word")
	}

	// Step 5: Submit word reviews
	for _, word := range words {
		reviewPayload := map[string]interface{}{
			"correct":       true,
			"response_time": 1.5,
		}
		reviewBody, _ := json.Marshal(reviewPayload)
		w = testutil.MakeRequest(r, "POST",
			fmt.Sprintf("/api/study_sessions/%d/words/%d/review", sessionID, word.ID),
			bytes.NewBuffer(reviewBody))
		testutil.AssertStatus(t, w, http.StatusOK)
	}

	// Step 6: Check study progress
	w = testutil.MakeRequest(r, "GET", "/api/dashboard/progress", nil)
	testutil.AssertStatus(t, w, http.StatusOK)

	var progress struct {
		TotalWords     int     `json:"total_words"`
		WordsStudied   int     `json:"words_studied"`
		AverageMastery float64 `json:"average_mastery"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &progress)
	if err != nil {
		t.Fatalf("Failed to parse progress response: %v", err)
	}
	if progress.WordsStudied == 0 {
		t.Error("Expected some words to be studied")
	}
	if progress.AverageMastery == 0 {
		t.Error("Expected non-zero mastery level")
	}

	// Step 7: Check quick stats
	w = testutil.MakeRequest(r, "GET", "/api/dashboard/quick_stats", nil)
	testutil.AssertStatus(t, w, http.StatusOK)

	var stats struct {
		TotalSessions int `json:"total_sessions"`
		TotalReviews  int `json:"total_reviews"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &stats)
	if err != nil {
		t.Fatalf("Failed to parse stats response: %v", err)
	}
	if stats.TotalSessions == 0 {
		t.Error("Expected at least one session")
	}
	if stats.TotalReviews == 0 {
		t.Error("Expected some reviews")
	}

	// Step 8: Reset history
	w = testutil.MakeRequest(r, "POST", "/api/reset_history", nil)
	testutil.AssertStatus(t, w, http.StatusOK)

	// Step 9: Verify reset
	w = testutil.MakeRequest(r, "GET", "/api/dashboard/quick_stats", nil)
	testutil.AssertStatus(t, w, http.StatusOK)

	err = json.Unmarshal(w.Body.Bytes(), &stats)
	if err != nil {
		t.Fatalf("Failed to parse stats response after reset: %v", err)
	}
	if stats.TotalSessions != 0 {
		t.Error("Expected no sessions after reset")
	}
	if stats.TotalReviews != 0 {
		t.Error("Expected no reviews after reset")
	}
}

func TestE2E_MultipleSessions(t *testing.T) {
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	// Reset history to clear any existing sessions
	router := setupRouter()
	resetResponse := testutil.MakeRequest(router, "POST", "/api/reset_history", bytes.NewReader(nil))
	testutil.AssertStatus(t, resetResponse, http.StatusOK)

	// Get existing group from database
	groupResponse := testutil.MakeRequest(router, "GET", "/api/groups", bytes.NewReader(nil))
	testutil.AssertStatus(t, groupResponse, http.StatusOK)

	var groups []struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(groupResponse.Body.Bytes(), &groups); err != nil {
		t.Fatalf("Failed to parse groups response: %v", err)
	}

	if len(groups) == 0 {
		t.Fatal("No groups found in database")
	}
	group := groups[0] // Use the first group

	// Then create a study activity
	activityData := map[string]interface{}{
		"name":        "Test Activity",
		"description": "Test Description",
		"group_ids":   []int{group.ID},
	}
	activityBody := encodeJSON(t, activityData)
	activityResponse := testutil.MakeRequest(router, "POST", "/api/study_activities", activityBody)
	testutil.AssertStatus(t, activityResponse, http.StatusCreated)

	var activity struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(activityResponse.Body.Bytes(), &activity); err != nil {
		t.Logf("Response body: %s", activityResponse.Body.String())
		t.Fatalf("Failed to parse activity response: %v", err)
	}

	if activity.ID == 0 {
		t.Fatal("Activity ID is 0")
	}

	// Create multiple study sessions
	for i := 0; i < 3; i++ {
		sessionData := map[string]interface{}{
			"study_activity_id": activity.ID,
			"start_time":        time.Now().Format(time.RFC3339),
		}
		body := encodeJSON(t, sessionData)
		response := testutil.MakeRequest(router, "POST", "/api/study_sessions", body)
		testutil.AssertStatus(t, response, http.StatusCreated)
		if response.Code != http.StatusCreated {
			t.Logf("Response body: %s", response.Body.String())
			t.Fatalf("Failed to create study session %d", i+1)
		}
	}

	// Retrieve study sessions and verify the count
	response := testutil.MakeRequest(router, "GET", "/api/study_sessions", nil)
	testutil.AssertStatus(t, response, http.StatusOK)
	if response.Code != http.StatusOK {
		t.Logf("Response body: %s", response.Body.String())
		t.Fatal("Failed to get study sessions")
	}

	var sessions []map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &sessions); err != nil {
		t.Logf("Response body: %s", response.Body.String())
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(sessions) != 3 {
		t.Errorf("Expected 3 study sessions, got %d", len(sessions))
	}
}

func TestE2E_InvalidActivityData(t *testing.T) {
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	router := setupRouter()

	// Create study activity with invalid data (missing required fields)
	activityData := map[string]interface{}{
		"group_id": 1,
	}
	body := encodeJSON(t, activityData)
	response := testutil.MakeRequest(router, "POST", "/api/study_activities", body)
	fmt.Printf("Response: %v\n", response.Body.String()) // Log the response
}

func TestE2E_PaginationLimits(t *testing.T) {
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	router := setupRouter()

	// Get existing group from database
	groupResponse := testutil.MakeRequest(router, "GET", "/api/groups", nil)
	testutil.AssertStatus(t, groupResponse, http.StatusOK)

	var groups []struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(groupResponse.Body.Bytes(), &groups); err != nil {
		t.Fatalf("Failed to parse groups response: %v", err)
	}

	if len(groups) == 0 {
		t.Fatal("No groups found in database")
	}
	group := groups[0] // Use the first group

	// Create a study activity
	activityData := map[string]interface{}{
		"name":        "Test Activity",
		"description": "Test Description",
		"group_ids":   []int{group.ID},
	}
	activityBody := encodeJSON(t, activityData)
	activityResponse := testutil.MakeRequest(router, "POST", "/api/study_activities", activityBody)
	testutil.AssertStatus(t, activityResponse, http.StatusCreated)

	var activity struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(activityResponse.Body.Bytes(), &activity); err != nil {
		t.Logf("Response body: %s", activityResponse.Body.String())
		t.Fatalf("Failed to parse activity response: %v", err)
	}

	if activity.ID == 0 {
		t.Fatal("Activity ID is 0")
	}

	// Create more than the default pagination limit study sessions
	for i := 0; i < 30; i++ {
		sessionData := map[string]interface{}{
			"study_activity_id": activity.ID,
			"start_time":        time.Now().Format(time.RFC3339),
		}
		body := encodeJSON(t, sessionData)
		response := testutil.MakeRequest(router, "POST", "/api/study_sessions", body)
		testutil.AssertStatus(t, response, http.StatusCreated)
		if response.Code != http.StatusCreated {
			t.Logf("Response body: %s", response.Body.String())
			t.Fatalf("Failed to create study session %d", i+1)
		}
	}

	// Retrieve study sessions with default pagination and verify the count
	response := testutil.MakeRequest(router, "GET", "/api/study_sessions", nil)
	testutil.AssertStatus(t, response, http.StatusOK)
	if response.Code != http.StatusOK {
		t.Logf("Response body: %s", response.Body.String())
		t.Fatal("Failed to get study sessions with default pagination")
	}

	var sessions []map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &sessions); err != nil {
		t.Logf("Response body: %s", response.Body.String())
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Assuming default pagination limit is 20
	if len(sessions) != 20 {
		t.Errorf("Expected 20 study sessions, got %d", len(sessions))
	}

	// Retrieve study sessions with a higher limit
	response = testutil.MakeRequest(router, "GET", "/api/study_sessions?limit=25", nil)
	testutil.AssertStatus(t, response, http.StatusOK)
	if response.Code != http.StatusOK {
		t.Logf("Response body: %s", response.Body.String())
		t.Fatal("Failed to get study sessions with higher limit")
	}

	if err := json.Unmarshal(response.Body.Bytes(), &sessions); err != nil {
		t.Logf("Response body: %s", response.Body.String())
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(sessions) != 25 {
		t.Errorf("Expected 25 study sessions, got %d", len(sessions))
	}
}

func TestE2E_ResetHistoryNoData(t *testing.T) {
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	router := setupRouter()

	// Reset history when there is no data
	response := testutil.MakeRequest(router, "POST", "/api/reset_history", bytes.NewReader(nil))
	testutil.AssertStatus(t, response, http.StatusOK)
	if response.Code != http.StatusOK {
		t.Logf("Response body: %s", response.Body.String())
		t.Fatal("Failed to reset history")
	}

	// Verify that there are no study sessions or word review items
	response = testutil.MakeRequest(router, "GET", "/api/study_sessions", nil)
	testutil.AssertStatus(t, response, http.StatusOK)
	if response.Code != http.StatusOK {
		t.Logf("Response body: %s", response.Body.String())
		t.Fatal("Failed to get study sessions after reset")
	}

	var sessions []map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &sessions); err != nil {
		t.Logf("Response body: %s", response.Body.String())
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(sessions) != 0 {
		t.Errorf("Expected 0 study sessions, got %d", len(sessions))
	}
}

func TestE2E_GroupOperations(t *testing.T) {
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	router := setupRouter()

	// Create a new group
	groupData := map[string]interface{}{
		"name": "Test Group",
	}
	body := encodeJSON(t, groupData)
	response := testutil.MakeRequest(router, "POST", "/api/groups", bytes.NewReader(body.Bytes()))
	testutil.AssertStatus(t, response, http.StatusCreated)

	var group struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &group); err != nil {
		t.Fatalf("Failed to parse group response: %v", err)
	}

	// Get the group by ID
	response = testutil.MakeRequest(router, "GET", fmt.Sprintf("/api/groups/%d", group.ID), bytes.NewReader(nil))
	testutil.AssertStatus(t, response, http.StatusOK)

	var groupDetails map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &groupDetails); err != nil {
		t.Fatalf("Failed to parse group details: %v", err)
	}

	if groupDetails["name"] != "Test Group" {
		t.Errorf("Expected group name 'Test Group', got %v", groupDetails["name"])
	}
}

func TestE2E_StudyActivityOperations(t *testing.T) {
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	router := setupRouter()

	// Create a new group first
	groupData := map[string]interface{}{
		"name": "Test Group",
	}
	body := encodeJSON(t, groupData)
	response := testutil.MakeRequest(router, "POST", "/api/groups", body)
	testutil.AssertStatus(t, response, http.StatusCreated)

	var group struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &group); err != nil {
		t.Fatalf("Failed to parse group response: %v", err)
	}

	// Create a study activity
	activityData := map[string]interface{}{
		"name":        "Test Activity",
		"description": "Test Description",
		"group_ids":   []int{group.ID},
	}
	body = encodeJSON(t, activityData)
	response = testutil.MakeRequest(router, "POST", "/api/study_activities", bytes.NewReader(body.Bytes()))
	testutil.AssertStatus(t, response, http.StatusCreated)

	var activity struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &activity); err != nil {
		t.Fatalf("Failed to parse activity response: %v", err)
	}

	// Create some study sessions for the activity
	for i := 0; i < 3; i++ {
		sessionData := map[string]interface{}{
			"study_activity_id": activity.ID,
			"start_time":        time.Now().Format(time.RFC3339),
		}
		body = encodeJSON(t, sessionData)
		response = testutil.MakeRequest(router, "POST", "/api/study_sessions", bytes.NewReader(body.Bytes()))
		testutil.AssertStatus(t, response, http.StatusCreated)
	}

	// Get activity sessions
	response = testutil.MakeRequest(router, "GET", fmt.Sprintf("/api/study_activities/%d/sessions", activity.ID), bytes.NewReader(nil))
	testutil.AssertStatus(t, response, http.StatusOK)

	var sessions []map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &sessions); err != nil {
		t.Fatalf("Failed to parse sessions response: %v", err)
	}

	if len(sessions) != 3 {
		t.Errorf("Expected 3 sessions, got %d", len(sessions))
	}
}

func TestE2E_StudySessionOperations(t *testing.T) {
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	router := setupRouter()

	// Create a group and activity first
	groupData := map[string]interface{}{
		"name": "Test Group",
	}
	body := encodeJSON(t, groupData)
	response := testutil.MakeRequest(router, "POST", "/api/groups", bytes.NewReader(body.Bytes()))
	testutil.AssertStatus(t, response, http.StatusCreated)

	var group struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &group); err != nil {
		t.Fatalf("Failed to parse group response: %v", err)
	}

	activityData := map[string]interface{}{
		"name":        "Test Activity",
		"description": "Test Description",
		"group_ids":   []int{group.ID},
	}
	body = new(bytes.Buffer)
	body = encodeJSON(t, activityData)
	response = testutil.MakeRequest(router, "POST", "/api/study_activities", bytes.NewReader(body.Bytes()))
	testutil.AssertStatus(t, response, http.StatusCreated)

	var activity struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &activity); err != nil {
		t.Fatalf("Failed to parse activity response: %v", err)
	}

	// Create a study session
	sessionData := map[string]interface{}{
		"study_activity_id": activity.ID,
		"start_time":        time.Now().Format(time.RFC3339),
	}
	body = new(bytes.Buffer)
	body = encodeJSON(t, sessionData)
	response = testutil.MakeRequest(router, "POST", "/api/study_sessions", bytes.NewReader(body.Bytes()))
	testutil.AssertStatus(t, response, http.StatusCreated)

	var session struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &session); err != nil {
		t.Fatalf("Failed to parse session response: %v", err)
	}

	// Get session by ID
	response = testutil.MakeRequest(router, "GET", fmt.Sprintf("/api/study_sessions/%d", session.ID), bytes.NewReader(nil))
	testutil.AssertStatus(t, response, http.StatusOK)

	var sessionDetails map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &sessionDetails); err != nil {
		t.Fatalf("Failed to parse session details: %v", err)
	}

	// Get session words
	response = testutil.MakeRequest(router, "GET", fmt.Sprintf("/api/study_sessions/%d/words", session.ID), bytes.NewReader(nil))
	testutil.AssertStatus(t, response, http.StatusOK)

	var words []map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &words); err != nil {
		t.Fatalf("Failed to parse words response: %v", err)
	}

	// Create a word review
	if len(words) > 0 {
		wordID := int(words[0]["id"].(float64))
		reviewData := map[string]interface{}{
			"correct":       true,
			"time_taken_ms": 1000,
		}
		body = new(bytes.Buffer)
		body = encodeJSON(t, reviewData)
		response = testutil.MakeRequest(router, "POST", fmt.Sprintf("/api/study_sessions/%d/words/%d/review", session.ID, wordID), bytes.NewReader(body.Bytes()))
		testutil.AssertStatus(t, response, http.StatusCreated)
	}
}

func TestE2E_DashboardEndpoints(t *testing.T) {
	cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	router := setupRouter()

	// Create test data first
	groupData := map[string]interface{}{
		"name": "Test Group",
	}
	body := encodeJSON(t, groupData)
	response := testutil.MakeRequest(router, "POST", "/api/groups", bytes.NewReader(body.Bytes()))
	testutil.AssertStatus(t, response, http.StatusCreated)

	var group struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &group); err != nil {
		t.Fatalf("Failed to parse group response: %v", err)
	}

	activityData := map[string]interface{}{
		"name":        "Test Activity",
		"description": "Test Description",
		"group_ids":   []int{group.ID},
	}
	body = new(bytes.Buffer)
	body = encodeJSON(t, activityData)
	response = testutil.MakeRequest(router, "POST", "/api/study_activities", bytes.NewReader(body.Bytes()))
	testutil.AssertStatus(t, response, http.StatusCreated)

	var activity struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &activity); err != nil {
		t.Fatalf("Failed to parse activity response: %v", err)
	}

	// Create a study session
	sessionData := map[string]interface{}{
		"study_activity_id": activity.ID,
		"start_time":        time.Now().Format(time.RFC3339),
	}
	body = new(bytes.Buffer)
	body = encodeJSON(t, sessionData)
	response = testutil.MakeRequest(router, "POST", "/api/study_sessions", bytes.NewReader(body.Bytes()))
	testutil.AssertStatus(t, response, http.StatusCreated)

	// Test last session endpoint
	response = testutil.MakeRequest(router, "GET", "/api/dashboard/last_session", bytes.NewReader(nil))
	testutil.AssertStatus(t, response, http.StatusOK)

	var lastSession map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &lastSession); err != nil {
		t.Fatalf("Failed to parse last session response: %v", err)
	}

	// Test progress endpoint
	response = testutil.MakeRequest(router, "GET", "/api/dashboard/progress", bytes.NewReader(nil))
	testutil.AssertStatus(t, response, http.StatusOK)

	var progress map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &progress); err != nil {
		t.Fatalf("Failed to parse progress response: %v", err)
	}

	// Test quick stats endpoint
	response = testutil.MakeRequest(router, "GET", "/api/dashboard/quick_stats", bytes.NewReader(nil))
	testutil.AssertStatus(t, response, http.StatusOK)

	var stats map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &stats); err != nil {
		t.Fatalf("Failed to parse stats response: %v", err)
	}
}
