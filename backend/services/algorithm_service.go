package services

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"unicode/utf8"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type AlgorithmRepository interface {
	List(ctx context.Context, limit, offset int) ([]dto.AlgorithmDTO, error)
	ListAdmin(ctx context.Context, limit, offset int, userId, status string) ([]dto.AlgorithmDTO, error)
	ListModeration(ctx context.Context, limit, offset int, status string) ([]dto.AlgorithmDTO, error)
	GetByPublicID(ctx context.Context, publicId string) (*dto.AlgorithmDTO, error)
	GetAdminAlgorithmById(ctx context.Context, algoId, userId string) (*dto.AlgorithmDTO, error)
	PostAlgorithm(ctx context.Context, data models.NewAlgorithm) (*dto.AlgorithmDTO, error)
	DeleteAlgorithm(ctx context.Context, publicId, userId string) error
	RestoreAlgorithm(ctx context.Context, publicId, userId string) error
	PutAlgorithm(ctx context.Context, data models.PutAlgorithm, userId string) (*dto.AlgorithmDTO, error)
	SitemapAlgorithms(ctx context.Context) ([]dto.SitemapItem, error)
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

func (s *AlgorithmService) List(ctx context.Context, page, limit int) ([]dto.AlgorithmDTO, int, error) {
	page, limit, offset := normalizePagination(page, limit, 10)

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

func (s *AlgorithmService) ListAdmin(ctx context.Context, page, limit int, idUser, status string) ([]dto.AlgorithmDTO, int, error) {
	page, limit, offset := normalizePagination(page, limit, 1)

	if !slices.Contains(models.AllStatuses, models.Status(status)) && status != "" {
		slog.Warn("invalid algorithm status provided, defaulting to approved", "providedStatus", status, "userId", idUser)
		status = "approved"
	}

	algorithms, err := s.AlgoRepo.ListAdmin(ctx, limit, offset, idUser, status)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmsNotFound) {
			return nil, page, models.ErrAlgorithmsNotFound
		}

		slog.Error("failed to query admin algorithms", "id", idUser, "page", page, "limit", limit, "offset", offset, "error", err)
		return nil, page, models.ErrFailQueryingAlgorithm
	}

	return algorithms, page, nil
}

func (s *AlgorithmService) GetAdminAlgorithm(ctx context.Context, algoId, userId string) (*dto.AlgorithmDTO, error) {
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

func (s *AlgorithmService) GetAlgorithmByPublicID(ctx context.Context, publicId string) (*dto.AlgorithmDTO, error) {
	algo, err := s.AlgoRepo.GetByPublicID(ctx, publicId)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		slog.Error("database error querying algorithm by public id", "publicId", publicId, "error", err)
		return nil, models.ErrFailQueryingAlgorithm
	}
	return algo, nil
}

func (s *AlgorithmService) PostAlgorithm(ctx context.Context, data dto.PostAlgorithmRequest, userId string) (*dto.AlgorithmDTO, error) {
	user, err := s.getAndValidateAuthor(ctx, userId, "post algorithm")
	if err != nil {
		return nil, err
	}

	name, category, content, err := validateAndSanitizeAlgorithmFields(data.Name, data.Category, data.Content)
	if err != nil {
		return nil, err
	}

	publicId, err := utils.GeneratePublicID()
	if err != nil {
		slog.Error("failed to generate secure public ID", "userId", userId, "error", err)
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
		slog.Error("repository failed to save algorithm", "slug", algorithm.Slug, "userId", userId, "error", err)
		return nil, models.ErrFailPostingAlgorithm
	}
	return res, nil
}

func (s *AlgorithmService) DeleteAlgorithm(ctx context.Context, algoId, userId string) error {
	user, err := s.getAndValidateAuthor(ctx, userId, "delete algorithm")
	if err != nil {
		return err
	}

	publicId := utils.SanitizeTitle(algoId)

	algo, err := s.validateOwnership(ctx, publicId, user.Id, "delete")
	if err != nil {
		return err
	}

	if algo.Status == "deleted" {
		return nil
	}

	err = s.AlgoRepo.DeleteAlgorithm(ctx, publicId, user.Id)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return models.ErrAlgorithmNotFound
		}
		slog.Error("database error during delete execution of algorithm", "publicId", publicId, "userId", userId, "error", err)
		return models.ErrFailQueryingAlgorithm
	}

	return nil
}

func (s *AlgorithmService) RestoreAlgorithm(ctx context.Context, algoId, userId string) error {
	user, err := s.getAndValidateAuthor(ctx, userId, "restore algorithm")
	if err != nil {
		return err
	}

	publicId := utils.SanitizeTitle(algoId)

	algo, err := s.validateOwnership(ctx, publicId, user.Id, "delete")
	if err != nil {
		return err
	}

	if algo.Status != "deleted" {
		return nil
	}

	err = s.AlgoRepo.RestoreAlgorithm(ctx, publicId, user.Id)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return models.ErrAlgorithmNotFound
		}
		slog.Error("database error during restore execution of algorithm", "publicId", publicId, "userId", userId, "error", err)
		return models.ErrFailQueryingAlgorithm
	}

	return nil
}

func (s *AlgorithmService) PutAlgorithm(ctx context.Context, data dto.PutAlgorithmRequest, publicId, userId string) (*dto.AlgorithmDTO, error) {
	user, err := s.getAndValidateAuthor(ctx, userId, "update algorithm")
	if err != nil {
		return nil, err
	}

	name, category, content, err := validateAndSanitizeAlgorithmFields(data.Name, data.Category, data.Content)
	if err != nil {
		return nil, err
	}

	algoId := utils.SanitizeTitle(publicId)

	if _, err := s.validateOwnership(ctx, algoId, user.Id, "update"); err != nil {
		return nil, err
	}

	algorithm := models.PutAlgorithm{
		PublicId:   algoId,
		Name:       name,
		Slug:       utils.Slug(name),
		Category:   category,
		Difficulty: data.Difficulty,
		Content:    content,
	}

	res, err := s.AlgoRepo.PutAlgorithm(ctx, algorithm, user.Id)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		slog.Error("database error during update execution of algorithm", "publicId", algoId, "userId", userId, "error", err)
		return nil, models.ErrFailQueryingAlgorithm
	}

	return res, nil
}

func (s *AlgorithmService) SitemapAlgorithms(ctx context.Context) ([]dto.SitemapItem, error) {
	algorithms, err := s.AlgoRepo.SitemapAlgorithms(ctx)
	if err != nil {
		return nil, models.ErrSitemapAlgorithms
	}

	return algorithms, nil
}

func (s *AlgorithmService) ListModeration(ctx context.Context, page, limit int, userId, status string) ([]dto.AlgorithmDTO, int, error) {
	if !slices.Contains(models.AllStatuses, models.Status(status)) && status != "" || status == "deleted" {
		slog.Warn("invalid algorithm status provided, defaulting to approved", "providedStatus", status, "userId", userId)
		status = "approved"
	}

	user, err := s.UserRepo.GetUserByIdForAuth(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, page, models.ErrUserNotFound
		}
		slog.Error("database error querying user ID during list moderation", "userId", userId, "error", err)
		return nil, page, models.ErrFailQueryUser
	}

	if !user.Enable {
		return nil, page, models.ErrUserNotEnabled
	}

	if !slices.Contains(user.Permissions, "moderate:algorithms") {
		return nil, page, models.ErrAlgorithmNoCreatePermission
	}

	page, limit, offset := normalizePagination(page, limit, 1)

	if !slices.Contains(models.AllStatuses, models.Status(status)) && status != "" {
		slog.Warn("invalid algorithm status provided, defaulting to approved", "providedStatus", status, "userId", userId)
		status = "approved"
	}

	algorithms, err := s.AlgoRepo.ListModeration(ctx, limit, offset, status)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmsNotFound) {
			return nil, page, models.ErrAlgorithmsNotFound
		}

		slog.Error("failed to query moderation algorithms", "id", userId, "page", page, "limit", limit, "offset", offset, "error", err)
		return nil, page, models.ErrFailQueryingAlgorithm
	}

	return algorithms, page, nil
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

func normalizePagination(page, limit, defaultLimit int) (int, int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = defaultLimit
	}
	offset := (page - 1) * limit
	return page, limit, offset
}

func (s *AlgorithmService) getAndValidateAuthor(ctx context.Context, userId, action string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByIdForAuth(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}
		slog.Error("database error querying user ID during "+action, "userId", userId, "error", err)
		return nil, models.ErrFailQueryUser
	}

	if !user.Enable {
		return nil, models.ErrUserNotEnabled
	}

	if !slices.Contains(user.Permissions, "create:algorithms") {
		return nil, models.ErrAlgorithmNoCreatePermission
	}

	return user, nil
}

func (s *AlgorithmService) validateOwnership(ctx context.Context, publicId, userId, action string) (*dto.AlgorithmDTO, error) {
	algo, err := s.AlgoRepo.GetAdminAlgorithmById(ctx, publicId, userId)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		slog.Error("database error during "+action+" of algorithm", "publicId", publicId, "userId", userId, "error", err)
		return nil, models.ErrFailQueryingAlgorithm
	}

	if algo.AuthorId != userId {
		return nil, models.ErrAlgorithmAuthorMismatch
	}

	return algo, nil
}
