package services

import (
	"context"
	"errors"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Auth(ctx context.Context, data dto.AuthRequest, restrict bool) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, data.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	isValid, err := utils.VerifyPassword(data.Password, user.PasswordHash)
	if err != nil || !isValid {
		return nil, errors.New("invalid email or password")
	}

	if restrict && !user.Role.IsEmployee {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
