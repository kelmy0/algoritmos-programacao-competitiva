package models

import "time"

type Difficulty string

const (
	DifficultyBeginner     Difficulty = "beginner"
	DifficultyIntermediate Difficulty = "intermediate"
	DifficultyAdvanced     Difficulty = "advanced"
	DifficultyExpert       Difficulty = "expert"
)

type Status string

const (
	StatusApproved Status = "approved"
	StatusPeding   Status = "pending"
	StatusRejected Status = "rejected"
)

type Algorithm struct {
	Id         string     `db:"id"`
	PublicId   string     `db:"public_id"`
	Slug       string     `db:"slug"`
	Name       string     `db:"name"`
	Category   string     `db:"category"`
	Difficulty Difficulty `db:"difficulty"`
	Content    string     `db:"content"`
	Status     string     `db:"status,oneof=approved pending rejected"`
	AuthorId   string     `db:"author_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

type NewAlgorithm struct {
	PublicId   string     `db:"public_id"`
	Slug       string     `db:"slug"`
	Name       string     `db:"name"`
	Category   string     `db:"category"`
	Difficulty Difficulty `db:"difficulty,oneof=beginner intermediate advanced expert"`
	Content    string     `db:"content"`
	AuthorId   string     `db:"author_id"`
}

type PutAlgorithm struct {
	PublicId   string     `db:"public_id"`
	Slug       string     `db:"slug"`
	Name       string     `db:"name"`
	Category   string     `db:"category"`
	Difficulty Difficulty `db:"difficulty,oneof=beginner intermediate advanced expert"`
	Content    string     `db:"content"`
}
