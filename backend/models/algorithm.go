package models

import "time"

type Difficulty string

const (
	DifficultyBeginner     Difficulty = "beginner"
	DifficultyIntermediate Difficulty = "intermediate"
	DifficultyAdvanced     Difficulty = "advanced"
	DifficultyExpert       Difficulty = "expert"
)

type Algorithm struct {
	Id         string     `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	Category   string     `json:"category" db:"category"`
	Difficulty Difficulty `json:"difficulty" db:"difficulty"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}
