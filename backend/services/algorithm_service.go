package services

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"slices"
	"unicode/utf8"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type AlgorithmRepository interface {
	List(ctx context.Context, limit, offset int) ([]models.Algorithm, error)
	ListAdmin(ctx context.Context, limit, offset int, userId string) ([]models.Algorithm, error)
	GetByPublicID(ctx context.Context, publicId string) (*models.Algorithm, error)
	GetAdminAlgorithmById(ctx context.Context, algoId, userId string) (*models.Algorithm, error)
	PostAlgorithm(ctx context.Context, data models.NewAlgorithm) (*models.Algorithm, error)
	DeleteAlgorithm(ctx context.Context, publicId string) (*models.Algorithm, error)
	PutAlgorithm(ctx context.Context, data models.PutAlgorithm, userId string) (*models.Algorithm, error)
}

type AlgorithmUserRepository interface {
	GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error)
}

type AlgorithmService struct {
	AlgoRepo AlgorithmRepository
	UserRepo AlgorithmUserRepository
}

func NewAlgorithmService(algoRepo AlgorithmRepository, userRepo AlgorithmUserRepository) *AlgorithmService {
	return &AlgorithmService{AlgoRepo: algoRepo, UserRepo: userRepo}
}

func (s *AlgorithmService) List(ctx context.Context, page, limit int) ([]models.Algorithm, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	data, err := s.AlgoRepo.List(ctx, limit, offset)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmsNotFound) {
			return nil, page, models.ErrAlgorithmsNotFound
		}

		slog.Error("failed to query algorithms", "page", page, "limit", limit, "offset", offset, "error", err)
		return nil, page, models.ErrFailQueryingAlgorithm
	}

	return data, page, err
}

func (s *AlgorithmService) ListAdmin(ctx context.Context, page, limit int, idUser string) ([]models.Algorithm, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}
	offset := (page - 1) * limit

	algorithms, err := s.AlgoRepo.ListAdmin(ctx, limit, offset, idUser)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmsNotFound) {
			return nil, page, models.ErrAlgorithmsNotFound
		}

		slog.Error("failed to query admin algorithms", "id", idUser, "page", page, "limit", limit, "offset", offset, "error", err)
		return nil, page, models.ErrFailQueryingAlgorithm
	}

	return algorithms, page, nil
}

func (s *AlgorithmService) GetAdminAlgorithm(ctx context.Context, algoId, userId string) (*models.Algorithm, error) {
	algo, err := s.AlgoRepo.GetAdminAlgorithmById(ctx, algoId, userId)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		slog.Error("error querying admin algorithm", "id", algoId, "userId", userId, "error", err)
		return nil, models.ErrFailQueryingAlgorithm
	}
	return algo, nil
}

func (s *AlgorithmService) GetAlgorithmByPublicID(ctx context.Context, publicId string) (*models.Algorithm, error) {
	algo, err := s.AlgoRepo.GetByPublicID(ctx, publicId)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		log.Printf("[AlgorithmService.GetAlgorithmByPublicID] database error for public_id %s: %v", publicId, err)
		return nil, models.ErrFailQueryingAlgorithm
	}
	return algo, nil
}

func (s *AlgorithmService) PostAlgorithm(ctx context.Context, data dto.PostAlgorithmRequest, userId string) (*models.Algorithm, error) {
	user, err := s.UserRepo.GetUserByIdForAuth(ctx, userId)

	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}
		slog.Error("database error querying user ID during post algoritm", "error", err)
		return nil, models.ErrFailQueryUser
	}

	if !user.Enable {
		return nil, models.ErrUserNotEnabled
	}

	if !slices.Contains(user.Permissions, "create:algorithms") {
		return nil, models.ErrAlgorithmNoCreatePermission
	}

	name, category, content, err := validateAndSanitizeAlgorithmFields(data.Name, data.Category, data.Content)
	if err != nil {
		return nil, err
	}

	publicId, err := utils.GeneratePublicID()
	if err != nil {
		log.Printf("[AlgorithmService.PostAlgorithm] failed to generate secure public ID: %v", err)
		return nil, models.ErrFailGeneratePublicId
	}

	algorithm := models.NewAlgorithm{
		PublicId:   publicId,
		Name:       name,
		Slug:       utils.Slug(name),
		Category:   category,
		Difficulty: data.Difficulty,
		Content:    content,
		AuthorId:   user.Id,
	}

	res, err := s.AlgoRepo.PostAlgorithm(ctx, algorithm)
	if err != nil {
		log.Printf("[AlgorithmService.PostAlgorithm] repository failed to save algorithm (slug: %s): %v", algorithm.Slug, err)
		return nil, models.ErrFailPostingAlgorithm
	}
	return res, nil
}

func (s *AlgorithmService) DeleteAlgorithm(ctx context.Context, publicId string) (*models.Algorithm, error) {
	algo, err := s.AlgoRepo.DeleteAlgorithm(ctx, publicId)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		log.Printf("[AlgorithmService.DeleteAlgorithm] database error during deletion of %s: %v", publicId, err)
		return nil, models.ErrFailQueryingAlgorithm
	}
	return algo, nil
}

func (s *AlgorithmService) PutAlgorithm(ctx context.Context, data dto.PutAlgorithmRequest, userId string) (*models.Algorithm, error) {
	user, err := s.UserRepo.GetUserByIdForAuth(ctx, userId)

	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}
		slog.Error("database error querying user ID during update algoritm", "error", err)
		return nil, models.ErrFailQueryUser
	}

	if !user.Enable {
		return nil, models.ErrUserNotEnabled
	}

	if !slices.Contains(user.Permissions, "create:algorithms") {
		return nil, models.ErrAlgorithmNoCreatePermission
	}

	name, category, content, err := validateAndSanitizeAlgorithmFields(data.Name, data.Category, data.Content)
	if err != nil {
		return nil, err
	}

	publicId := utils.SanitizeTitle(data.PublicId)

	algorithm := models.PutAlgorithm{
		PublicId:   publicId,
		Name:       name,
		Slug:       utils.Slug(name),
		Category:   category,
		Difficulty: data.Difficulty,
		Content:    content,
	}

	algo, err := s.AlgoRepo.GetAdminAlgorithmById(ctx, publicId, user.Id)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		slog.Error("database error during update of algorithm", "publicId", publicId, "error", err)
		return nil, models.ErrFailQueryingAlgorithm
	}

	if algo.AuthorId != user.Id {
		return nil, models.ErrAlgorithmAuthorMismatch
	}

	res, err := s.AlgoRepo.PutAlgorithm(ctx, algorithm, user.Id)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		slog.Error("database error during update of algorithm", "publicId", publicId, "error", err)
		return nil, models.ErrFailQueryingAlgorithm
	}

	return res, nil
}

func validateAndSanitizeAlgorithmFields(name, category, content string) (string, string, string, error) {
	nameSanitized := utils.SanitizeTitle(name)
	categorySanitized := utils.SanitizeTitle(category)
	contentSanitized := content

	if nameSanitized == "" || utf8.RuneCountInString(nameSanitized) < 3 {
		return "", "", "", models.ErrInvalidAlgorithmName
	}

	if categorySanitized == "" || utf8.RuneCountInString(categorySanitized) < 3 {
		return "", "", "", models.ErrInvalidAlgorithmCategory
	}

	if contentSanitized == "" || utf8.RuneCountInString(contentSanitized) < 10 {
		return "", "", "", models.ErrInvalidAlgorithmContent
	}

	return nameSanitized, categorySanitized, contentSanitized, nil
}
