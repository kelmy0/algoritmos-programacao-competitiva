package dto

import (
	"time"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
)

type ListAlgorithmDTO struct {
	PublicId   string            `json:"publicId"`
	Slug       string            `json:"slug"`
	Name       string            `json:"name"`
	Category   string            `json:"category"`
	Difficulty models.Difficulty `json:"difficulty"`
	Status     models.Status     `json:"status"`
}

type AlgorithmDTO struct {
	PublicId   string            `json:"publicId"`
	Slug       string            `json:"slug"`
	Name       string            `json:"name"`
	Category   string            `json:"category"`
	Difficulty models.Difficulty `json:"difficulty"`
	Content    string            `json:"content"`
	AuthorId   string            `json:"authorId"`
	Status     models.Status     `json:"status"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

type ListAlgorithmsResponse struct {
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	HasMore    bool               `json:"hasMore"`
	Algorithms []ListAlgorithmDTO `json:"algorithms"`
}

type PostAlgorithmRequest struct {
	Name       string            `json:"name" binding:"required,min=3"`
	Category   string            `json:"category" binding:"required,min=3"`
	Difficulty models.Difficulty `json:"difficulty" binding:"required,oneof=beginner intermediate advanced expert"`
	Content    string            `json:"content" binding:"required,min=10"`
}

type PostAlgorithmResponse struct {
	PublicId string `json:"publicId"`
	Slug     string `json:"slug"`
}

type PutAlgorithmRequest struct {
	Name       string            `json:"name" binding:"required,min=3"`
	Category   string            `json:"category" binding:"required,min=3"`
	Difficulty models.Difficulty `json:"difficulty" binding:"required,oneof=beginner intermediate advanced expert"`
	Content    string            `json:"content" binding:"required,min=10"`
}

type PutAlgorithmResponse = PostAlgorithmResponse

type SitemapItem struct {
	Slug      string    `json:"slug"`
	UpdatedAt time.Time `json:"updatedAt"`
}
