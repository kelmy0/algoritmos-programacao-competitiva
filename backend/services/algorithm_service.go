package services

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type AlgorithmRepository interface {
	List(ctx context.Context, limit, offset int) (list []dto.ListAlgorithmDTO, err error)
	ListAdmin(ctx context.Context, limit, offset int, userId, status string) (list []dto.ListAlgorithmDTO, err error)
	ListModeration(ctx context.Context, limit, offset int, status string) (list []dto.ListAlgorithmDTO, err error)
	GetByPublicID(ctx context.Context, publicId string) (algo *models.Algorithm, err error)
	GetAdminAlgorithmById(ctx context.Context, algoId, userId string) (algo *models.Algorithm, err error)
	PostAlgorithm(ctx context.Context, data models.PostAlgorithm) (algo dto.PostAlgorithmResponse, err error)
	DeleteAlgorithm(ctx context.Context, publicId, userId string) error
	RestoreAlgorithm(ctx context.Context, publicId, userId string) error
	PutAlgorithm(ctx context.Context, data models.PutAlgorithm, userId string) (algo dto.PutAlgorithmResponse, err error)
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

func (s *AlgorithmService) List(ctx context.Context, page, limit int) (data []dto.ListAlgorithmDTO, currentPage int, hasMore bool, err error) {
	currentPage, limit, offset := normalizePagination(page, limit)
	algorithms, err := s.AlgoRepo.List(ctx, limit+1, offset)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmsNotFound) {
			return nil, currentPage, false, models.ErrAlgorithmsNotFound
		}

		slog.Error("failed to query algorithms", "page", currentPage, "limit", limit, "offset", offset, "error", err)
		return nil, currentPage, false, models.ErrFailQueryingAlgorithm
	}

	if len(algorithms) > limit {
		hasMore = true
		data = algorithms[:limit]
	} else {
		data = algorithms
	}

	return data, currentPage, hasMore, nil
}

func (s *AlgorithmService) ListAdmin(ctx context.Context, page, limit int, idUser, status string) (data []dto.ListAlgorithmDTO, currentPage int, hasMore bool, err error) {
	currentPage, limit, offset := normalizePagination(page, limit)

	if !slices.Contains(models.AllStatuses, models.Status(status)) && status != "" {
		slog.Warn("invalid algorithm status provided, defaulting to approved", "providedStatus", status, "userId", idUser)
		status = "approved"
	}

	algorithms, err := s.AlgoRepo.ListAdmin(ctx, limit+1, offset, idUser, status)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmsNotFound) {
			return nil, currentPage, false, models.ErrAlgorithmsNotFound
		}

		slog.Error("failed to query admin algorithms", "id", idUser, "page", currentPage, "limit", limit, "offset", offset, "error", err)
		return nil, currentPage, false, models.ErrFailQueryingAlgorithm
	}

	if len(algorithms) > limit {
		hasMore = true
		data = algorithms[:limit]
	} else {
		data = algorithms
	}

	return data, currentPage, hasMore, nil
}

func (s *AlgorithmService) ListModeration(ctx context.Context, page, limit int, userId, status string) (data []dto.ListAlgorithmDTO, currentPage int, hasMore bool, err error) {
	currentPage, limit, offset := normalizePagination(page, limit)

	if (!slices.Contains(models.AllStatuses, models.Status(status)) && status != "") || status == "deleted" {
		slog.Warn("invalid algorithm status provided, defaulting to approved", "providedStatus", status, "userId", userId)
		status = "approved"
	}

	user, err := s.UserRepo.GetUserByIdForAuth(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, currentPage, false, models.ErrUserNotFound
		}
		slog.Error("database error querying user ID during list moderation", "userId", userId, "error", err)
		return nil, currentPage, false, models.ErrFailQueryUser
	}

	if !user.Enable {
		return nil, currentPage, false, models.ErrUserNotEnabled
	}

	if !slices.Contains(user.Permissions, "moderate:algorithms") {
		return nil, currentPage, false, models.ErrAlgorithmNoModeratePermission
	}

	algorithms, err := s.AlgoRepo.ListModeration(ctx, limit+1, offset, status)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmsNotFound) {
			return nil, currentPage, false, models.ErrAlgorithmsNotFound
		}

		slog.Error("failed to query moderation algorithms", "id", userId, "page", currentPage, "limit", limit, "offset", offset, "error", err)
		return nil, currentPage, false, models.ErrFailQueryingAlgorithm
	}

	if len(algorithms) > limit {
		hasMore = true
		data = algorithms[:limit]
	} else {
		data = algorithms
	}

	return data, currentPage, hasMore, nil
}

func (s *AlgorithmService) GetAlgorithmByPublicID(ctx context.Context, publicId string) (algo *dto.AlgorithmDTO, err error) {
	rawAlgo, err := s.AlgoRepo.GetByPublicID(ctx, publicId)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		slog.Error("database error querying algorithm by public id", "publicId", publicId, "error", err)
		return nil, models.ErrFailQueryingAlgorithm
	}

	textDescompressed, err := utils.DecompressText(rawAlgo.Content)
	if err != nil {
		slog.Error("failed to decompress algorithm content", "publicId", publicId, "error", err)
		return nil, models.ErrFailQueryingAlgorithm
	}

	algo = &dto.AlgorithmDTO{
		PublicId:   rawAlgo.PublicId,
		Slug:       rawAlgo.Slug,
		Name:       rawAlgo.Name,
		Category:   rawAlgo.Category,
		Difficulty: rawAlgo.Difficulty,
		Content:    textDescompressed,
		CreatedAt:  rawAlgo.CreatedAt,
		UpdatedAt:  rawAlgo.UpdatedAt,
	}

	return algo, nil
}

func (s *AlgorithmService) GetAdminAlgorithm(ctx context.Context, publicId, userId string) (algo *dto.AlgorithmDTO, err error) {
	rawAlgo, err := s.AlgoRepo.GetAdminAlgorithmById(ctx, publicId, userId)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		slog.Error("error querying admin algorithm", "id", publicId, "userId", userId, "error", err)
		return nil, models.ErrFailQueryingAlgorithm
	}

	textDescompressed, err := utils.DecompressText(rawAlgo.Content)
	if err != nil {
		slog.Error("failed to decompress algorithm content", "publicId", publicId, "userId", userId, "error", err)
		return nil, models.ErrFailQueryingAlgorithm
	}

	algo = &dto.AlgorithmDTO{
		PublicId:   rawAlgo.PublicId,
		Slug:       rawAlgo.Slug,
		Name:       rawAlgo.Name,
		Category:   rawAlgo.Category,
		Difficulty: rawAlgo.Difficulty,
		Content:    textDescompressed,
		AuthorId:   rawAlgo.AuthorId,
		Status:     rawAlgo.Status,
		CreatedAt:  rawAlgo.CreatedAt,
		UpdatedAt:  rawAlgo.UpdatedAt,
	}

	return algo, nil
}

func (s *AlgorithmService) PostAlgorithm(ctx context.Context, data dto.PostAlgorithmRequest, userId string) (algo dto.PostAlgorithmResponse, err error) {
	user, err := s.getAndValidateAuthor(ctx, userId, "post algorithm")
	if err != nil {
		return algo, err
	}

	name, category, compressedContent, err := validateAndSanitizeAlgorithmFields(data.Name, data.Category, data.Content)
	if err != nil {
		return algo, err
	}

	publicId, err := utils.GeneratePublicID()
	if err != nil {
		slog.Error("failed to generate secure public ID", "userId", userId, "error", err)
		return algo, models.ErrFailGeneratePublicId
	}

	algorithm := models.PostAlgorithm{
		PublicId:   publicId,
		Name:       name,
		Slug:       utils.Slug(name),
		Category:   category,
		Difficulty: data.Difficulty,
		Content:    compressedContent,
		AuthorId:   user.Id,
	}

	algo, err = s.AlgoRepo.PostAlgorithm(ctx, algorithm)
	if err != nil {
		slog.Error("repository failed to save algorithm", "slug", algorithm.Slug, "userId", userId, "error", err)
		return algo, models.ErrFailPostingAlgorithm
	}

	return algo, nil
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

func (s *AlgorithmService) PutAlgorithm(ctx context.Context, data dto.PutAlgorithmRequest, publicId, userId string) (algo dto.PutAlgorithmResponse, err error) {
	user, err := s.getAndValidateAuthor(ctx, userId, "update algorithm")
	if err != nil {
		return algo, err
	}

	name, category, compressedContent, err := validateAndSanitizeAlgorithmFields(data.Name, data.Category, data.Content)
	if err != nil {
		return algo, err
	}

	algoId := utils.SanitizeTitle(publicId)

	if _, err := s.validateOwnership(ctx, algoId, user.Id, "update"); err != nil {
		return algo, err
	}

	algorithm := models.PutAlgorithm{
		PublicId:   algoId,
		Name:       name,
		Slug:       utils.Slug(name),
		Category:   category,
		Difficulty: data.Difficulty,
		Content:    compressedContent,
	}

	algo, err = s.AlgoRepo.PutAlgorithm(ctx, algorithm, user.Id)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return algo, models.ErrAlgorithmNotFound
		}
		slog.Error("database error during update execution of algorithm", "publicId", algoId, "userId", userId, "error", err)
		return algo, models.ErrFailQueryingAlgorithm
	}

	return algo, nil
}

func (s *AlgorithmService) SitemapAlgorithms(ctx context.Context) (data []dto.SitemapItem, err error) {
	data, err = s.AlgoRepo.SitemapAlgorithms(ctx)
	if err != nil {
		slog.Error("failed to generate algorithms sitemap", "error", err)
		return data, models.ErrSitemapAlgorithms
	}

	return data, nil
}

func validateAndSanitizeAlgorithmFields(name, category, content string) (string, string, []byte, error) {
	nameSanitized := utils.SanitizeTitle(name)
	categorySanitized := utils.SanitizeTitle(category)

	if nameSanitized == "" || utf8.RuneCountInString(nameSanitized) < 3 {
		return "", "", nil, models.ErrInvalidAlgorithmName
	}

	if categorySanitized == "" || utf8.RuneCountInString(categorySanitized) < 3 {
		return "", "", nil, models.ErrInvalidAlgorithmCategory
	}

	trimmedContent := strings.TrimSpace(content)
	if trimmedContent == "" {
		return "", "", nil, models.ErrInvalidAlgorithmContent
	}

	compressed, err := utils.CompressText(trimmedContent)
	if err != nil {
		slog.Error("failed to compress algorithm content", "error", err)
		return "", "", nil, models.ErrInvalidAlgorithmContent
	}

	return nameSanitized, categorySanitized, compressed, nil
}

func normalizePagination(page, limit int) (cPage int, cLimit int, offset int) {
	cPage = max(1, page)
	cLimit = max(1, min(100, limit))

	offset = (page - 1) * limit
	return cPage, cLimit, offset
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
	rawAlgo, err := s.AlgoRepo.GetAdminAlgorithmById(ctx, publicId, userId)
	if err != nil {
		if errors.Is(err, models.ErrAlgorithmNotFound) {
			return nil, models.ErrAlgorithmNotFound
		}
		slog.Error("database error during "+action+" of algorithm", "publicId", publicId, "userId", userId, "error", err)
		return nil, models.ErrFailQueryingAlgorithm
	}

	if rawAlgo.AuthorId != userId {
		return nil, models.ErrAlgorithmAuthorMismatch
	}

	textDescompressed, err := utils.DecompressText(rawAlgo.Content)
	if err != nil {
		return nil, models.ErrFailQueryingAlgorithm
	}

	algo := &dto.AlgorithmDTO{
		PublicId:   publicId,
		Slug:       rawAlgo.Slug,
		Name:       rawAlgo.Name,
		Category:   rawAlgo.Category,
		Difficulty: rawAlgo.Difficulty,
		Content:    textDescompressed,
		AuthorId:   rawAlgo.AuthorId,
		Status:     rawAlgo.Status,
		CreatedAt:  rawAlgo.CreatedAt,
		UpdatedAt:  rawAlgo.UpdatedAt,
	}

	return algo, nil
}
