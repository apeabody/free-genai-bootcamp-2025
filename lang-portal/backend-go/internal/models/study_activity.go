package models

import "time"

type StudyActivity struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type StudyActivityWithStats struct {
	StudyActivity
	TotalSessions int `json:"total_sessions"`
}
