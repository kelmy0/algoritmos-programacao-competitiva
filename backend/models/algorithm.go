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
	PublicId   string     `json:"public_id" db:"public_id"`
	Slug       string     `json:"slug" db:"slug"`
	Name       string     `json:"name" db:"name"`
	Category   string     `json:"category" db:"category"`
	Difficulty Difficulty `json:"difficulty" db:"difficulty"`
	Content    string     `json:"content" db:"content"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

type NewAlgorithm struct {
	PublicId   string     `db:"public_id"`
	Slug       string     `db:"slug"`
	Name       string     `db:"name"`
	Category   string     `db:"category"`
	Difficulty Difficulty `db:"difficulty,oneof=beginner intermediate advanced expert"`
	Content    string     `db:"content"`
}
