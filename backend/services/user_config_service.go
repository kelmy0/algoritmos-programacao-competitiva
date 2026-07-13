package services

import (
	"context"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type UserConfigRepo interface {
	GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error)
	ChangePassword(ctx context.Context, id, newPassword string) error
}

type AuthConfigRepo interface {
	DeleteAllRefreshToken(ctx context.Context, userId, tokenId string) error
	GetRefreshTokenById(ctx context.Context, id string) (*models.RefreshToken, error)
}

type UserConfigService struct {
	UserRepo         UserConfigRepo
	AuthRepo         AuthConfigRepo
	ArgonParams      utils.ArgonParams
	JwtRefreshSecret string
	AppName          string
}

func NewUserConfigService(userRepo UserConfigRepo, authRepo AuthConfigRepo, argonParams utils.ArgonParams, jwtRSecret, appName string) *UserConfigService {
	return &UserConfigService{
		UserRepo:         userRepo,
		ArgonParams:      argonParams,
		AuthRepo:         authRepo,
		JwtRefreshSecret: jwtRSecret,
		AppName:          appName,
	}
}

func (s *UserConfigService) ChangePassword(ctx context.Context, userIdContext, refreshTokenString string, data dto.ChangePasswordRequest) error {
	if data.NewPassword != data.ConfirmNewPassword {
		return models.ErrPasswordsDontMatch
	}

	claims, err := utils.ValidateToken(refreshTokenString, s.JwtRefreshSecret, s.AppName)
	if err != nil {
		return models.ErrInvalidOrExpiredRefresh
	}

	tokenExists, err := s.AuthRepo.GetRefreshTokenById(ctx, claims.ID)
	if err != nil || tokenExists == nil {
		return models.ErrInvalidOrExpiredRefresh
	}

	if tokenExists.UserId != claims.Subject || userIdContext != claims.Subject {
		return models.ErrTokenMetadataMisMatch
	}

	userId := claims.Subject

	user, err := s.UserRepo.GetUserByIdForAuth(ctx, userId)
	if err != nil {
		return models.ErrUserNotFound
	}

	if user.PasswordHash == nil {
		return models.ErrPasswordNotSet
	}

	ok, err := utils.VerifyPassword(data.OldPassword, *user.PasswordHash)
	if !ok || err != nil {
		return models.ErrIncorrectPassword
	}

	newPasswordHash, err := utils.HashPassword(data.NewPassword, s.ArgonParams)
	if err != nil {
		return models.ErrFailChangingPassword
	}

	err = s.UserRepo.ChangePassword(ctx, userId, newPasswordHash)
	if err != nil {
		return models.ErrFailChangingPassword
	}

	err = s.AuthRepo.DeleteAllRefreshToken(ctx, userId, tokenExists.Id)
	if err != nil {
		return models.ErrPasswordChangeButNotLogout
	}

	return nil
}
