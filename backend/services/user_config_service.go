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
	DefinePassword(ctx context.Context, id, newPassword string) error
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

	user, token, err := s.validateUserSession(ctx, userIdContext, refreshTokenString)
	if err != nil {
		return err
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
		return models.ErrPasswordChangeFailed
	}

	err = s.UserRepo.ChangePassword(ctx, user.Id, newPasswordHash)
	if err != nil {
		return models.ErrPasswordChangeFailed
	}

	err = s.AuthRepo.DeleteAllRefreshToken(ctx, user.Id, token.Id)
	if err != nil {
		return models.ErrPasswordChangeButNotLogout
	}

	return nil
}

func (s *UserConfigService) DefinePassword(ctx context.Context, userIdContext, refreshTokenString string, data dto.DefinePasswordRequest) error {
	if data.NewPassword != data.ConfirmNewPassword {
		return models.ErrPasswordsDontMatch
	}

	user, token, err := s.validateUserSession(ctx, userIdContext, refreshTokenString)
	if err != nil {
		return err
	}

	if user.PasswordHash != nil {
		return models.ErrPasswordSet
	}

	newPasswordHash, err := utils.HashPassword(data.NewPassword, s.ArgonParams)
	if err != nil {
		return models.ErrPasswordSetFailed
	}

	err = s.UserRepo.DefinePassword(ctx, user.Id, newPasswordHash)
	if err != nil {
		return models.ErrPasswordSetFailed
	}

	err = s.AuthRepo.DeleteAllRefreshToken(ctx, user.Id, token.Id)
	if err != nil {
		return models.ErrPasswordSetButNotLogout
	}

	return nil
}

func (s *UserConfigService) validateUserSession(ctx context.Context, userIdContext, refreshTokenString string) (*models.User, *models.RefreshToken, error) {
	claims, err := utils.ValidateToken(refreshTokenString, s.JwtRefreshSecret, s.AppName)
	if err != nil {
		return nil, nil, models.ErrInvalidOrExpiredRefresh
	}

	tokenExists, err := s.AuthRepo.GetRefreshTokenById(ctx, claims.ID)
	if err != nil || tokenExists == nil {
		return nil, nil, models.ErrInvalidOrExpiredRefresh
	}

	if tokenExists.UserId != claims.Subject || userIdContext != claims.Subject {
		return nil, nil, models.ErrTokenMetadataMisMatch
	}

	user, err := s.UserRepo.GetUserByIdForAuth(ctx, claims.Subject)
	if err != nil {
		return nil, nil, models.ErrUserNotFound
	}

	return user, tokenExists, nil
}
