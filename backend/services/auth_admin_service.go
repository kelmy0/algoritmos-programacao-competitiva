package services

import (
	"context"
	"errors"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type AuthAdminRepository interface {
	GetAdminByEmail(ctx context.Context, email string) (*models.Administrator, error)
}

type AuthAdminService struct {
	repo AuthAdminRepository
}

func NewAuthAdminService(repo AuthAdminRepository) *AuthAdminService {
	return &AuthAdminService{repo: repo}
}

func (s *AuthAdminService) AuthAdmin(ctx context.Context, data dto.AuthAdminRequest) (*models.Administrator, error) {
	admin, err := s.repo.GetAdminByEmail(ctx, data.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	isValid, err := utils.VerifyPassword(data.Password, admin.PasswordHash)
	if err != nil || !isValid {
		return nil, errors.New("invalid email or password")
	}

	return admin, nil
}
