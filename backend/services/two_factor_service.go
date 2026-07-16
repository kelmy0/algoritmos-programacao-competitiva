package services

import (
	"context"
	"errors"
	"log"

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

type TwoFactorService struct {
	Repo          TwoFactorUserRepository
	EncryptSecret string
	AppName       string
}

func NewTwoFactorService(repo TwoFactorUserRepository, encryptSecret, appName string) *TwoFactorService {
	return &TwoFactorService{Repo: repo, EncryptSecret: encryptSecret, AppName: appName}
}

func (s *TwoFactorService) Generate2FA(ctx context.Context, userId, email string) (*dto.TwoFactorGenerateResponse, error) {
	twoFactorData, err := s.Repo.GetAuthData(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}

		log.Printf("[Generate2FA] database query error for user %s: %v", userId, err)
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
		log.Printf("[Generate2FA] failed to generate TOTP key for user %s: %v", userId, err)
		return nil, models.ErrGeneratingToken
	}

	encryptedSecret, err := utils.Encrypt(key.Secret(), s.EncryptSecret)
	if err != nil {
		log.Printf("[Generate2FA] AES encryption failed for user %s secret: %v", userId, err)
		return nil, models.ErrCryptTokenFailed
	}

	err = s.Repo.Save2FASecret(ctx, userId, encryptedSecret)
	if err != nil {
		log.Printf("[Generate2FA] failed to save encrypted 2FA secret to DB for user %s: %v", userId, err)
		return nil, models.Err2FASaveFailed
	}

	return &dto.TwoFactorGenerateResponse{
		Secret: key.Secret(),
		QRCode: key.URL(),
	}, nil
}

func (s *TwoFactorService) Enable2FA(ctx context.Context, userId, code string) error {
	twoFactorData, err := s.Repo.GetAuthData(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return models.ErrUserNotFound
		}

		log.Printf("[Enable2FA] database query error for user %s: %v", userId, err)
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
		log.Printf("[Enable2FA] AES decryption of 2FA secret failed for user %s: %v", userId, err)
		return models.ErrDecryptTokenFailed
	}

	isValid := totp.Validate(code, decryptedSecret)
	if !isValid {
		return models.Err2FAInvalid
	}

	err = s.Repo.Enable2FA(ctx, userId)
	if err != nil {
		log.Printf("[Enable2FA] failed to update 2FA status to enabled in DB for user %s: %v", userId, err)
		return models.Err2FAUpdateFailed
	}

	return nil
}

func (s *TwoFactorService) Disable2FA(ctx context.Context, userId, password string) error {
	twoFactorData, err := s.Repo.GetAuthData(ctx, userId)
	if err != nil {
		log.Printf("[Disable2FA] database query error for user %s: %v", userId, err)
		return models.Err2FAGetDataFailed
	}

	if !twoFactorData.IsEnabled {
		return models.Err2FAAlreadyDisabled
	}

	isValid, err := utils.VerifyPassword(password, twoFactorData.PasswordHash)
	if err != nil {
		log.Printf("[Disable2FA] Argon2 verification system error for user %s: %v", userId, err)
		return models.ErrPasswordVerificationFailed
	}

	if !isValid {
		return models.ErrIncorrectPassword
	}

	err = s.Repo.Disable2FA(ctx, userId)
	if err != nil {
		log.Printf("[Disable2FA] failed to update 2FA status to disabled (clean secret/enable flag) in DB for user %s: %v", userId, err)
		return models.Err2FAUpdateFailed
	}

	return nil
}
