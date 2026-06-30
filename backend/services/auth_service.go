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
	SaveRefreshToken(ctx context.Context, tokenId, userId string, expiresAt time.Time) error
	GetRefreshTokenById(ctx context.Context, id string) (*models.RefreshToken, error)
	DeleteRefreshTokenById(ctx context.Context, userId, tokenId string) error
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
func (s *AuthService) Auth(ctx context.Context, data dto.AuthRequest) (*dto.LoginResponse, string, int, error) {
	user, err := s.Repo.GetUserByEmail(ctx, data.Email)
	if err != nil || !user.Enable {
		return nil, "", 0, errors.New("invalid email or password")
	}

	isValid, err := utils.VerifyPassword(data.Password, user.PasswordHash)
	if err != nil || !isValid {
		return nil, "", 0, errors.New("invalid email or password")
	}

	// Minutes
	_, accessToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtAccessSecret, s.AppName, user.Role.IsEmployee, time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute))
	if err != nil {
		return nil, "", 0, errors.New("Error generating Token.")
	}

	// Days
	idToken, refreshToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtRefreshSecret, s.AppName, user.Role.IsEmployee, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		return nil, "", 0, errors.New("Error generating Token.")
	}

	err = s.Repo.SaveRefreshToken(ctx, idToken, user.Id, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		return nil, "", 0, errors.New("Error generating Token.")
	}

	response := &dto.LoginResponse{
		AcessToken: accessToken,
	}

	return response, refreshToken, s.JwtRefreshExpiration, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (string, error) {
	claims, err := utils.ValidadeToken(refreshTokenString, s.JwtRefreshSecret, s.AppName)
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	tokenExists, err := s.Repo.GetRefreshTokenById(ctx, claims.ID)
	if err != nil || tokenExists == nil {
		return "", errors.New("invalid or expired refresh token")
	}

	user, err := s.Repo.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !user.Enable {
		return "", errors.New("user account is disabled")
	}

	_, accessToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtAccessSecret, s.AppName, user.Role.IsEmployee, time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute))

	if err != nil {
		return "", errors.New("error generating new access token")
	}

	return accessToken, nil
}

func (s *AuthService) Logout(ctx context.Context, userId, refreshTokenString string) error {
	claims, err := utils.ValidadeToken(refreshTokenString, s.JwtRefreshSecret, s.AppName)
	if err != nil {
		return errors.New("invalid or expired refresh token")
	}

	return s.Repo.DeleteRefreshTokenById(ctx, userId, claims.ID)
}
