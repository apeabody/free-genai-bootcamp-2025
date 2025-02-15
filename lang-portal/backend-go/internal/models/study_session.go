package models

import "time"

type StudySession struct {
	ID              int       `json:"id" db:"id"`
	GroupID         int       `json:"group_id" db:"group_id"`
	StudyActivityID int       `json:"study_activity_id" db:"study_activity_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type StudySessionWithStats struct {
	StudySession
	ActivityName     string    `json:"activity_name"`
	GroupName        string    `json:"group_name"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	ReviewItemsCount int       `json:"review_items_count"`
	Score            int       `json:"score"`
	CorrectCount     int       `json:"correct_count"`
	IncorrectCount   int       `json:"incorrect_count"`
}
