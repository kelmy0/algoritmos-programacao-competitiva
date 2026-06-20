package dto

import "github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"

type PostAlgorithmRequest struct {
	Name       string            `json:"name" binding:"required,min=3"`
	Category   string            `json:"category" binding:"required"`
	Difficulty models.Difficulty `json:"difficulty" binding:"required,oneof=beginner intermediate advanced expert"`
	Content    string            `json:"content" binding:"required,min=10"`
}

type PutAlgorithmRequest struct {
	PublicId   string            `json:"public_id" binding:"required,min=8"`
	Name       string            `json:"name" binding:"required,min=3"`
	Category   string            `json:"category" binding:"required"`
	Difficulty models.Difficulty `json:"difficulty" binding:"required,oneof=beginner intermediate advanced expert"`
	Content    string            `json:"content" binding:"required,min=10"`
}
