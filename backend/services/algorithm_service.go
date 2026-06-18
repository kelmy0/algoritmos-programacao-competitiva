package services

import (
	"context"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
)

type AlgorithmRepository interface {
	List(ctx context.Context, limit, offset int) ([]models.Algorithm, error)
	GetById(ctx context.Context, id string) (*models.Algorithm, error)
}

type AlgorithmService struct {
	repo AlgorithmRepository
}

func NewAlgorithmService(repo AlgorithmRepository) *AlgorithmService {
	return &AlgorithmService{repo: repo}
}

// List some Algorithms
func (s *AlgorithmService) List(ctx context.Context, page, limit int) ([]models.Algorithm, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	data, err := s.repo.List(ctx, limit, offset)
	return data, page, err
}

// Get a specific algorithm by id
func (s *AlgorithmService) GetAlgorithmById(ctx context.Context, id string) (*models.Algorithm, error) {
	return s.repo.GetById(ctx, id)
}
