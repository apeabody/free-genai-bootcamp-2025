package models

import "time"

type Word struct {
	ID        int       `json:"id" db:"id"`
	English   string    `json:"english" db:"english"`
	Spanish   string    `json:"spanish" db:"spanish"`
	Level     string    `json:"level" db:"level"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type WordWithStats struct {
	Word
	CorrectCount   int     `json:"correct_count"`
	IncorrectCount int     `json:"incorrect_count"`
	MasteryLevel   float64 `json:"mastery_level"`
}
