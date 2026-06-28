package services

import (
	"context"
	"errors"
	"time"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type AuthService struct {
	Repo                 AuthRepository
	JwtAccessSecret      string
	JwtRefreshSecret     string
	AppName              string
	JwtAccessExpiration  int
	JwtRefreshExpiration int
}

func NewAuthService(repo AuthRepository, jwtAccessSecret, jwtRefreshSecret, appName string, jwtAccessExpiration int, jwtRefreshExpiration int) *AuthService {
	return &AuthService{
		Repo:                 repo,
		JwtAccessSecret:      jwtAccessSecret,
		JwtRefreshSecret:     jwtRefreshSecret,
		AppName:              appName,
		JwtAccessExpiration:  jwtAccessExpiration,
		JwtRefreshExpiration: jwtRefreshExpiration,
	}
}

// Returns access token, Refresh token, errors
func (s *AuthService) Auth(ctx context.Context, data dto.AuthRequest, restrict bool) (*dto.LoginResponse, string, int, error) {
	user, err := s.Repo.GetUserByEmail(ctx, data.Email)
	if err != nil {
		return nil, "", 0, errors.New("invalid email or password")
	}

	isValid, err := utils.VerifyPassword(data.Password, user.PasswordHash)
	if err != nil || !isValid {
		return nil, "", 0, errors.New("invalid email or password")
	}

	if restrict && !user.Role.IsEmployee {
		return nil, "", 0, errors.New("invalid email or password")
	}

	// Minutes
	accessToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.RoleId, s.JwtAccessSecret, s.AppName, time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute))
	if err != nil {
		return nil, "", 0, errors.New("Error generating Token.")
	}

	// Days
	refreshToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.RoleId, s.JwtRefreshSecret, s.AppName, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		return nil, "", 0, errors.New("Error generating Token.")
	}

	response := &dto.LoginResponse{
		AcessToken: accessToken,
	}

	return response, refreshToken, s.JwtRefreshExpiration, nil
}
