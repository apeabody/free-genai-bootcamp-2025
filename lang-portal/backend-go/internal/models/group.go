package models

import "time"

type Group struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type GroupWithStats struct {
	Group
	Statistics struct {
		TotalWordCount  int     `json:"total_word_count"`
		MasteredWords   int     `json:"mastered_words"`
		InProgressWords int     `json:"in_progress_words"`
		AverageMastery  float64 `json:"average_mastery"`
	} `json:"statistics"`
}
