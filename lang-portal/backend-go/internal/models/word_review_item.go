package models

import "time"

type WordReviewItem struct {
	ID              int       `json:"id" db:"id"`
	WordID          int       `json:"word_id" db:"word_id"`
	StudyActivityID int       `json:"study_activity_id" db:"study_activity_id"`
	Correct         bool      `json:"correct" db:"correct"`
	ResponseTime    float64   `json:"response_time" db:"response_time"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type WordReviewResult struct {
	WordID          int     `json:"word_id"`
	SessionID       int     `json:"session_id"`
	Correct         bool    `json:"correct"`
	ResponseTime    float64 `json:"response_time"`
	NewMasteryLevel float64 `json:"new_mastery_level"`
}
