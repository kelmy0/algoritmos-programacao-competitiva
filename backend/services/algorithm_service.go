package services

import (
	"context"
	"errors"
	"log"
	"unicode/utf8"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type AlgorithmRepository interface {
	List(ctx context.Context, limit, offset int) ([]models.Algorithm, error)
	GetByPublicID(ctx context.Context, publicId string) (*models.Algorithm, error)
	PostAlgorithm(ctx context.Context, data models.NewAlgorithm) (*models.Algorithm, error)
	DeleteAlgorithm(ctx context.Context, publicId string) (*models.Algorithm, error)
	PutAlgorithm(ctx context.Context, data models.PutAlgorithm) (*models.Algorithm, error)
}

type AlgorithmService struct {
	repo AlgorithmRepository
}

func NewAlgorithmService(repo AlgorithmRepository) *AlgorithmService {
	return &AlgorithmService{repo: repo}
}

func (s *AlgorithmService) List(ctx context.Context, page, limit int) ([]models.Algorithm, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	data, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		log.Printf("[AlgorithmService.List] failed to retrieve algorithms (page: %d, limit: %d): %v", page, limit, err)
		return nil, page, models.ErrFailQueryingAlgorithm
	}

	return data, page, err
}

func (s *AlgorithmService) GetAlgorithmByPublicID(ctx context.Context, publicId string) (*models.Algorithm, error) {
	algo, err := s.repo.GetByPublicID(ctx, publicId)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		log.Printf("[AlgorithmService.GetAlgorithmByPublicID] database error for public_id %s: %v", publicId, err)
		return nil, models.ErrFailQueryingAlgorithm
	}
	return algo, nil
}

func (s *AlgorithmService) PostAlgorithm(ctx context.Context, data dto.PostAlgorithmRequest) (*models.Algorithm, error) {
	name, category, content, err := validateAndSanitizeAlgorithmFields(data.Name, data.Category, data.Content)
	if err != nil {
		return nil, err
	}

	publicId, err := utils.GeneratePublicID()
	if err != nil {
		log.Printf("[AlgorithmService.PostAlgorithm] failed to generate secure public ID: %v", err)
		return nil, models.ErrFailGeneratePublicId
	}

	algorithm := &models.NewAlgorithm{
		PublicId:   publicId,
		Name:       name,
		Slug:       utils.Slug(name),
		Category:   category,
		Difficulty: data.Difficulty,
		Content:    content,
	}

	res, err := s.repo.PostAlgorithm(ctx, *algorithm)
	if err != nil {
		log.Printf("[AlgorithmService.PostAlgorithm] repository failed to save algorithm (slug: %s): %v", algorithm.Slug, err)
		return nil, models.ErrFailQueryingAlgorithm
	}
	return res, nil
}

func (s *AlgorithmService) DeleteAlgorithm(ctx context.Context, publicId string) (*models.Algorithm, error) {
	algo, err := s.repo.DeleteAlgorithm(ctx, publicId)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		log.Printf("[AlgorithmService.DeleteAlgorithm] database error during deletion of %s: %v", publicId, err)
		return nil, models.ErrFailQueryingAlgorithm
	}
	return algo, nil
}

func (s *AlgorithmService) PutAlgorithm(ctx context.Context, data dto.PutAlgorithmRequest) (*models.Algorithm, error) {
	name, category, content, err := validateAndSanitizeAlgorithmFields(data.Name, data.Category, data.Content)
	if err != nil {
		return nil, err
	}

	publicId := utils.SanitizeTitle(data.PublicId)

	algorithm := &models.PutAlgorithm{
		PublicId:   publicId,
		Name:       name,
		Slug:       utils.Slug(name),
		Category:   category,
		Difficulty: data.Difficulty,
		Content:    content,
	}

	res, err := s.repo.PutAlgorithm(ctx, *algorithm)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		log.Printf("[AlgorithmService.PutAlgorithm] database error during update of %s: %v", publicId, err)
		return nil, models.ErrFailQueryingAlgorithm
	}

	return res, nil
}

func validateAndSanitizeAlgorithmFields(name, category, content string) (string, string, string, error) {
	nameSanitized := utils.SanitizeTitle(name)
	categorySanitized := utils.SanitizeTitle(category)
	contentSanitized := utils.SanitizeMarkDown(content)

	if nameSanitized == "" || contentSanitized == "" || categorySanitized == "" ||
		utf8.RuneCountInString(nameSanitized) < 3 ||
		utf8.RuneCountInString(categorySanitized) < 3 ||
		utf8.RuneCountInString(contentSanitized) < 10 {
		return "", "", "", models.ErrInvalidNameCategoryContent
	}

	return nameSanitized, categorySanitized, contentSanitized, nil
}
