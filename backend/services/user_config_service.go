package services

import (
	"context"
	"log"
	"time"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type UserConfigRepo interface {
	GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error)
	ChangePassword(ctx context.Context, id, newPassword string) error
	DefinePassword(ctx context.Context, id, newPassword string) error
	GetUserByEmailForAuth(ctx context.Context, email string) (*models.User, error)
	UpdateRecoveryToken(ctx context.Context, userId, tokenHash string, expiresAt time.Time) error
	GetUserByRecoveryToken(ctx context.Context, tokenHash string) (*models.User, error)
}

type AuthConfigRepo interface {
	DeleteAllRefreshToken(ctx context.Context, userId, tokenId string) error
	GetRefreshTokenById(ctx context.Context, id string) (*models.RefreshToken, error)
}

type UserConfigService struct {
	UserRepo         UserConfigRepo
	AuthRepo         AuthConfigRepo
	EmailService     EmailService
	ArgonParams      utils.ArgonParams
	JwtRefreshSecret string
	AppName          string
}

func NewUserConfigService(userRepo UserConfigRepo, authRepo AuthConfigRepo, emailService EmailService, argonParams utils.ArgonParams, jwtRSecret, appName string) *UserConfigService {
	return &UserConfigService{
		UserRepo:         userRepo,
		ArgonParams:      argonParams,
		EmailService:     emailService,
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

func (s *UserConfigService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.UserRepo.GetUserByEmailForAuth(ctx, email)
	if err != nil || user == nil || !user.Enable {
		return nil
	}

	if user.RecoveryTokenExpiresAt != nil {
		now := time.Now()

		if now.Before(*user.RecoveryTokenExpiresAt) {
			minTimeRemainingForNewSend := 5 * time.Minute
			timeLeft := user.RecoveryTokenExpiresAt.Sub(now)

			if timeLeft > minTimeRemainingForNewSend {
				waitTime := timeLeft - minTimeRemainingForNewSend

				log.Printf("[ForgotPassword] Email sending blocked by cooldown for %s. Remaining wait time: %v",
					email, waitTime.Round(time.Second))
				return nil
			}
		}
	}

	token, err := utils.GenerateCustomId(32)
	if err != nil {
		return models.ErrGeneratingToken
	}

	tokenHash := utils.HashSHA512(token)
	expiresAt := time.Now().Add(15 * time.Minute)

	err = s.UserRepo.UpdateRecoveryToken(ctx, user.Id, tokenHash, expiresAt)
	if err != nil {
		return models.ErrGeneratingToken
	}

	utils.GoSafe(func() {
		_ = s.EmailService.SendRecoveryEmail(user.Email, token)
	})

	return nil
}

func (s *UserConfigService) ResetPassword(ctx context.Context, data dto.ResetPasswordRequest) error {
	if data.NewPassword != data.ConfirmNewPassword {
		return models.ErrPasswordsDontMatch
	}

	hashToken := utils.HashSHA512(data.Token)
	user, err := s.UserRepo.GetUserByRecoveryToken(ctx, hashToken)
	if err != nil || user == nil || user.RecoveryTokenExpiresAt == nil || !time.Now().Before(*user.RecoveryTokenExpiresAt) {
		return models.ErrInvalidOrExpiredToken
	}

	if !user.Enable {
		return models.ErrUserNotEnabled
	}

	newPasswordHash, err := utils.HashPassword(data.NewPassword, s.ArgonParams)
	if err != nil {
		return models.ErrPasswordChangeFailed
	}

	err = s.UserRepo.ChangePassword(ctx, user.Id, newPasswordHash)
	if err != nil {
		return models.ErrPasswordChangeFailed
	}

	err = s.AuthRepo.DeleteAllRefreshToken(ctx, user.Id, "")
	if err != nil {
		return models.ErrPasswordChangeButNotLogout
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
