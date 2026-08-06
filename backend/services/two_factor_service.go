package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/repositories"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"github.com/pquerna/otp/totp"
)

type TwoFactorUserRepository interface {
	Save2FASecret(ctx context.Context, userId, secret string) error
	Enable2FA(ctx context.Context, userId string) error
	Disable2FA(ctx context.Context, userId string) error
	GetAuthData(ctx context.Context, userId string) (*repositories.UserAuthData, error)
}

type TwoFactorAuthRepository interface {
	DeleteAllUserRefreshTokens(ctx context.Context, userId string) error
}

type TwoFactorService struct {
	UserRepo      TwoFactorUserRepository
	AuthRepo      TwoFactorAuthRepository
	EncryptSecret string
	AppName       string
}

func NewTwoFactorService(userRepo TwoFactorUserRepository, authRepo TwoFactorAuthRepository, encryptSecret, appName string) *TwoFactorService {
	return &TwoFactorService{UserRepo: userRepo, AuthRepo: authRepo, EncryptSecret: encryptSecret, AppName: appName}
}

func (s *TwoFactorService) Generate2FA(ctx context.Context, userId, email string) (*dto.TwoFactorGenerateResponse, error) {
	twoFactorData, err := s.UserRepo.GetAuthData(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}

		slog.ErrorContext(ctx, "database query error",
			"op", "Generate2FA",
			"user_id", userId,
			"error", err,
		)
		return nil, models.ErrFailQueryUser
	}

	if twoFactorData.IsEnabled {
		return nil, models.Err2FAAlreadyEnabled
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.AppName,
		AccountName: email,
		SecretSize:  32,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate TOTP key",
			"op", "Generate2FA",
			"user_id", userId,
			"error", err,
		)
		return nil, models.ErrGeneratingToken
	}

	encryptedSecret, err := utils.Encrypt(key.Secret(), s.EncryptSecret)
	if err != nil {
		slog.ErrorContext(ctx, "AES encryption failed for secret",
			"op", "Generate2FA",
			"user_id", userId,
			"error", err,
		)
		return nil, models.ErrCryptTokenFailed
	}

	err = s.UserRepo.Save2FASecret(ctx, userId, encryptedSecret)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save encrypted 2FA secret to DB",
			"op", "Generate2FA",
			"user_id", userId,
			"error", err,
		)
		return nil, models.Err2FASaveFailed
	}

	return &dto.TwoFactorGenerateResponse{
		Secret: key.Secret(),
		QRCode: key.URL(),
	}, nil
}

func (s *TwoFactorService) Enable2FA(ctx context.Context, userId, code string) error {
	twoFactorData, err := s.UserRepo.GetAuthData(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return models.ErrUserNotFound
		}

		slog.ErrorContext(ctx, "database query error",
			"op", "Enable2FA",
			"user_id", userId,
			"error", err,
		)
		return models.Err2FAGetDataFailed
	}

	if twoFactorData.IsEnabled {
		return models.Err2FAAlreadyEnabled
	}

	if twoFactorData.Secret == "" {
		return models.Err2FANotInitiated
	}

	decryptedSecret, err := utils.Decrypt(twoFactorData.Secret, s.EncryptSecret)
	if err != nil {
		slog.ErrorContext(ctx, "AES decryption of 2FA secret failed",
			"op", "Enable2FA",
			"user_id", userId,
			"error", err,
		)
		return models.ErrDecryptTokenFailed
	}

	isValid := totp.Validate(code, decryptedSecret)
	if !isValid {
		return models.Err2FAInvalid
	}

	err = s.UserRepo.Enable2FA(ctx, userId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update 2FA status to enabled in DB",
			"op", "Enable2FA",
			"user_id", userId,
			"error", err,
		)
		return models.Err2FAUpdateFailed
	}

	if err = s.AuthRepo.DeleteAllUserRefreshTokens(ctx, userId); err != nil {
		slog.WarnContext(ctx, "failed to revoke refresh tokens",
			"op", "Enable2FA",
			"user_id", userId,
			"error", err,
		)
	}

	return nil
}

func (s *TwoFactorService) Disable2FA(ctx context.Context, userId, password string) error {
	twoFactorData, err := s.UserRepo.GetAuthData(ctx, userId)
	if err != nil {
		slog.ErrorContext(ctx, "database query error",
			"op", "Disable2FA",
			"user_id", userId,
			"error", err,
		)
		return models.Err2FAGetDataFailed
	}

	if !twoFactorData.IsEnabled {
		return models.Err2FAAlreadyDisabled
	}

	isValid, err := utils.VerifyPassword(password, twoFactorData.PasswordHash)
	if err != nil {
		slog.ErrorContext(ctx, "Argon2 verification system error",
			"op", "Disable2FA",
			"user_id", userId,
			"error", err,
		)
		return models.ErrPasswordVerificationFailed
	}

	if !isValid {
		return models.ErrIncorrectPassword
	}

	err = s.UserRepo.Disable2FA(ctx, userId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update 2FA status to disabled in DB",
			"op", "Disable2FA",
			"user_id", userId,
			"error", err,
		)
		return models.Err2FAUpdateFailed
	}

	if err = s.AuthRepo.DeleteAllUserRefreshTokens(ctx, userId); err != nil {
		slog.WarnContext(ctx, "failed to revoke refresh tokens",
			"op", "Disable2FA",
			"user_id", userId,
			"error", err,
		)
	}

	return nil
}
