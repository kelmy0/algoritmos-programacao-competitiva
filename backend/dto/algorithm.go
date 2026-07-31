package dto

import (
	"time"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
)

type AlgorithmDTO struct {
	PublicId   string            `json:"publicId"`
	Slug       string            `json:"slug"`
	Name       string            `json:"name"`
	Category   string            `json:"category"`
	Difficulty models.Difficulty `json:"difficulty"`
	Content    string            `json:"content"`
	Status     models.Status     `json:"status"`
	AuthorId   string            `json:"authorId"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

type ListAlgorithmsResponse struct {
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Data  []AlgorithmDTO `json:"data"`
}

type ListAdminAlgorithmsResponse struct {
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Algorithms []AlgorithmDTO `json:"algorithms"`
}

type AlgorithmResponse struct {
	Data *AlgorithmDTO `json:"data"`
}

type PostAlgorithmRequest struct {
	Name       string            `json:"name" binding:"required,min=3"`
	Category   string            `json:"category" binding:"required,min=3"`
	Difficulty models.Difficulty `json:"difficulty" binding:"required,oneof=beginner intermediate advanced expert"`
	Content    string            `json:"content" binding:"required,min=10"`
}

type PutAlgorithmRequest struct {
	Name       string            `json:"name" binding:"required,min=3"`
	Category   string            `json:"category" binding:"required,min=3"`
	Difficulty models.Difficulty `json:"difficulty" binding:"required,oneof=beginner intermediate advanced expert"`
	Content    string            `json:"content" binding:"required,min=10"`
}

type SitemapItem struct {
	Slug      string    `json:"slug"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SitemapResponse struct {
	Data []SitemapItem `json:"data"`
}
