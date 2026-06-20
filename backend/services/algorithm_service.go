package services

import (
	"context"
	"errors"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type AlgorithmRepository interface {
	List(ctx context.Context, limit, offset int) ([]models.Algorithm, error)
	GetByPublicID(ctx context.Context, publicId string) (*models.Algorithm, error)
	PostAlgorithm(ctx context.Context, data models.NewAlgorithm) (*models.Algorithm, error)
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
func (s *AlgorithmService) GetAlgorithmByPublicID(ctx context.Context, id string) (*models.Algorithm, error) {
	return s.repo.GetByPublicID(ctx, id)
}

func (s *AlgorithmService) PostAlgorithm(ctx context.Context, data dto.PostAlgorithmRequest) (*models.Algorithm, error) {
	nameSanitized := utils.SanitizeName(data.Name)
	content := utils.SanitizeMarkDown(data.Content)
	categorySanitized := utils.SanitizeName(data.Category)

	if nameSanitized == "" || content == "" {
		return nil, errors.New("Invalid name or content!")
	}

	publicId, err := utils.GeneratePublicID()
	if err != nil {
		return nil, errors.New("Error to generate public id")
	}

	slug := utils.Slug(nameSanitized)

	println(publicId)
	println(nameSanitized)
	println(slug)
	println(data.Difficulty)
	println(content)

	algorithm := &models.NewAlgorithm{
		PublicId:   publicId,
		Name:       nameSanitized,
		Slug:       slug,
		Category:   categorySanitized,
		Difficulty: data.Difficulty,
		Content:    content,
	}

	return s.repo.PostAlgorithm(ctx, *algorithm)
}
