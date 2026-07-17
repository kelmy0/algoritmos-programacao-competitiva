package services

import (
	"context"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
)

type UserRepo interface {
	GetPublicCredentialsUser(ctx context.Context, id string) (*models.User, error)
}

type UsersService struct {
	UserRepo UserRepo
}

func NewUsersService(userRepo UserRepo) *UsersService {
	return &UsersService{
		UserRepo: userRepo,
	}
}
