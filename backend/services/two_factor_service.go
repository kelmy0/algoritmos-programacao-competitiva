package services

import (
	"context"
	"errors"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/repositories"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"github.com/pquerna/otp/totp"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
	Save2FASecret(ctx context.Context, userId, secret string) error
	Enable2FA(ctx context.Context, userId string) error
	Disable2FA(ctx context.Context, userId string) error
	GetAuthData(ctx context.Context, userId string) (*repositories.UserAuthData, error)
}

type TwoFactorService struct {
	Repo          UserRepository
	EncryptSecret string
	AppName       string
}

func NewTwoFactorService(repo UserRepository, encryptSecret, appName string) *TwoFactorService {
	return &TwoFactorService{Repo: repo, EncryptSecret: encryptSecret, AppName: appName}
}

func (s *TwoFactorService) Generate2FA(ctx context.Context, userId, email string) (*dto.TwoFactorGenerateResponse, error) {
	twoFactorData, err := s.Repo.GetAuthData(ctx, userId)
	if err != nil {
		return nil, errors.New("Error retrieving 2FA data")
	}

	if twoFactorData.IsEnabled {
		return nil, errors.New("2FA already enabled")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.AppName,
		AccountName: email,
		SecretSize:  32,
	})
	if err != nil {
		return nil, errors.New("Error generating TOTP")
	}

	encryptedSecret, err := utils.Encrypt(key.Secret(), s.EncryptSecret)
	if err != nil {
		return nil, errors.New("Error encrypting TOTP")
	}

	err = s.Repo.Save2FASecret(ctx, userId, encryptedSecret)
	if err != nil {
		return nil, errors.New("Error saving TOTP")
	}

	return &dto.TwoFactorGenerateResponse{
		Secret: key.Secret(),
		QRCode: key.URL(),
	}, nil
}

func (s *TwoFactorService) Enable2FA(ctx context.Context, userId, code string) error {
	twoFactorData, err := s.Repo.GetAuthData(ctx, userId)
	if err != nil {
		return errors.New("Error retrieving 2FA data")
	}

	if twoFactorData.IsEnabled {
		return errors.New("2FA already enabled")
	}

	if twoFactorData.Secret == "" {
		return errors.New("2FA setup has not been initiated for this user")
	}

	decryptedSecret, err := utils.Decrypt(twoFactorData.Secret, s.EncryptSecret)
	if err != nil {
		return errors.New("Error processing authentication security")
	}

	isValid := totp.Validate(code, decryptedSecret)
	if !isValid {
		return errors.New("2FA code is invalid or expired")
	}

	err = s.Repo.Enable2FA(ctx, userId)
	if err != nil {
		return errors.New("Error activating 2FA")
	}

	return nil
}

func (s *TwoFactorService) Disable2FA(ctx context.Context, userId, password string) error {
	twoFactorData, err := s.Repo.GetAuthData(ctx, userId)

	if err != nil {
		return errors.New("Error retrieving 2FA data")
	}

	if !twoFactorData.IsEnabled {
		return errors.New("2FA already disabled")
	}

	isValid, err := utils.VerifyPassword(password, twoFactorData.PasswordHash)
	if err != nil || !isValid {
		return errors.New("Invalid password")
	}

	err = s.Repo.Disable2FA(ctx, userId)
	if err != nil {
		return errors.New("Error disabling 2FA")
	}

	return nil
}
