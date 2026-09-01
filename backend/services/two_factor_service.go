package services

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/repositories"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
)

type TwoFactorUserRepository interface {
	Save2FASecret(ctx context.Context, userId, secret string) error
	Enable2FA(ctx context.Context, userId string) error
	Disable2FA(ctx context.Context, userId string) error
	GetAuthData(ctx context.Context, userId string) (repositories.UserAuthData, error)
	GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error)
}

type TwoFactorAuthRepository interface {
	GetRefreshTokenById(ctx context.Context, id string) (*models.RefreshToken, error)
	DeleteAllUserRefreshTokens(ctx context.Context, userId string) error
	RevokeFamily(ctx context.Context, familyId string) error
	SaveRefreshToken(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error
}

type TwoFactorService struct {
	UserRepo             TwoFactorUserRepository
	AuthRepo             TwoFactorAuthRepository
	RedisClient          *redis.Client
	EncryptSecret        string
	AppName              string
	AppDomain            string
	JwtAccessPrivateKey  ed25519.PrivateKey
	JwtRefreshPrivateKey ed25519.PrivateKey
	JwtAccessPublicKey   ed25519.PublicKey
	JwtRefreshPublicKey  ed25519.PublicKey
	JwtAccessExpiration  int
	JwtRefreshExpiration int
}

type Enable2FAResult struct {
	AccessToken  string
	RefreshToken string
}

func NewTwoFactorService(userRepo TwoFactorUserRepository, authRepo TwoFactorAuthRepository, redisClient *redis.Client,
	encryptSecret, appName, appDomain string, accessPrivate, refreshPrivate ed25519.PrivateKey,
	accessPublic, refreshPublic ed25519.PublicKey, accessDuration, refreshDuration int) *TwoFactorService {
	return &TwoFactorService{
		UserRepo:             userRepo,
		AuthRepo:             authRepo,
		RedisClient:          redisClient,
		EncryptSecret:        encryptSecret,
		AppName:              appName,
		AppDomain:            appDomain,
		JwtAccessPrivateKey:  accessPrivate,
		JwtRefreshPrivateKey: refreshPrivate,
		JwtAccessPublicKey:   accessPublic,
		JwtRefreshPublicKey:  refreshPublic,
		JwtRefreshExpiration: refreshDuration,
		JwtAccessExpiration:  accessDuration,
	}
}

func (s *TwoFactorService) Generate2FA(ctx context.Context, data dto.TwoFactorGenerateRequest) (dto.TwoFactorGenerateResponse, error) {
	twoFactorData, err := s.UserRepo.GetAuthData(ctx, data.UserId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return dto.TwoFactorGenerateResponse{}, models.ErrUserNotFound
		}

		slog.ErrorContext(ctx, "database query error",
			"op", "Generate2FA",
			"user_id", data.UserId,
			"error", err,
		)
		return dto.TwoFactorGenerateResponse{}, models.ErrFailQueryUser
	}

	if twoFactorData.IsEnabled {
		return dto.TwoFactorGenerateResponse{}, models.Err2FAAlreadyEnabled
	}

	if twoFactorData.PasswordHash != "" {
		isValid, err := utils.VerifyPassword(data.Password, twoFactorData.PasswordHash)
		if err != nil {
			slog.ErrorContext(ctx, "Argon2 verification system error",
				"op", "Disable2FA",
				"user_id", data.UserId,
				"error", err,
			)
			return dto.TwoFactorGenerateResponse{}, models.ErrPasswordVerificationFailed
		}

		if !isValid {
			return dto.TwoFactorGenerateResponse{}, models.ErrIncorrectPassword
		}
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.AppName,
		AccountName: data.Email,
		SecretSize:  32,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate TOTP key",
			"op", "Generate2FA",
			"user_id", data.UserId,
			"error", err,
		)
		return dto.TwoFactorGenerateResponse{}, models.ErrGeneratingToken
	}

	encryptedSecret, err := utils.Encrypt(key.Secret(), s.EncryptSecret)
	if err != nil {
		slog.ErrorContext(ctx, "AES encryption failed for secret",
			"op", "Generate2FA",
			"user_id", data.UserId,
			"error", err,
		)
		return dto.TwoFactorGenerateResponse{}, models.ErrCryptTokenFailed
	}

	err = s.UserRepo.Save2FASecret(ctx, data.UserId, encryptedSecret)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save encrypted 2FA secret to DB",
			"op", "Generate2FA",
			"user_id", data.UserId,
			"error", err,
		)
		return dto.TwoFactorGenerateResponse{}, models.Err2FASaveFailed
	}

	return dto.TwoFactorGenerateResponse{
		Secret: key.Secret(),
		QRCode: key.URL(),
	}, nil
}

func (s *TwoFactorService) Enable2FA(ctx context.Context, data dto.TwoFactorEnableRequest) (Enable2FAResult, error) {
	claims, err := utils.ValidateRefreshToken(data.RefreshToken, s.JwtRefreshPublicKey, s.AppDomain)

	if err != nil {
		slog.WarnContext(ctx, "refresh token JWT validation failed", slog.Any("error", err))
		return Enable2FAResult{}, models.ErrInvalidOrExpiredRefresh
	}

	if claims.Subject != data.UserId {
		slog.WarnContext(ctx, "security mismatch during 2fa enable",
			slog.String("token_subject", claims.Subject),
			slog.String("user_id_param", data.UserId),
		)
		return Enable2FAResult{}, models.ErrTokenMetadataMisMatch
	}

	if claims.DeviceHash != data.DeviceHash {
		slog.WarnContext(ctx, "[SECURITY ALERT] Device hash mismatch! Revoking entire family.",
			slog.String("user_id", claims.Subject),
			slog.String("family_id", claims.FamilyId),
			slog.String("token_id", claims.ID),
			slog.String("token_dvh", claims.DeviceHash),
			slog.String("dvh", data.DeviceHash),
		)
		_ = s.AuthRepo.RevokeFamily(ctx, claims.FamilyId)
		return Enable2FAResult{}, models.ErrInvalidOrExpiredRefresh
	}

	dbToken, err := s.AuthRepo.GetRefreshTokenById(ctx, claims.ID)
	if err != nil || dbToken == nil || dbToken.IsRevoked {
		return Enable2FAResult{}, models.ErrInvalidOrExpiredRefresh
	}

	user, err := s.UserRepo.GetUserByIdForAuth(ctx, data.UserId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return Enable2FAResult{}, models.ErrUserNotFound
		}

		slog.ErrorContext(ctx, "database query error",
			slog.String("op", "Enable2FA"),
			slog.String("user_id", data.UserId),
			slog.Any("error", err),
		)
		return Enable2FAResult{}, models.Err2FAGetDataFailed
	}

	if user.TwoFactorAuthentication {
		return Enable2FAResult{}, models.Err2FAAlreadyEnabled
	}

	if user.TwoFactorSecret == nil {
		return Enable2FAResult{}, models.Err2FANotInitiated
	}

	decryptedSecret, err := utils.Decrypt(*user.TwoFactorSecret, s.EncryptSecret)
	if err != nil {
		slog.ErrorContext(ctx, "AES decryption of 2FA secret failed",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return Enable2FAResult{}, models.ErrDecryptTokenFailed
	}

	isValid := totp.Validate(data.Code, decryptedSecret)
	if !isValid {
		return Enable2FAResult{}, models.Err2FAInvalid
	}

	if err = s.UserRepo.Enable2FA(ctx, user.Id); err != nil {
		slog.ErrorContext(ctx, "failed to update 2FA status to enabled in DB",
			slog.String("op", "Enable2FA"),
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return Enable2FAResult{}, models.Err2FAUpdateFailed
	}
	user.TwoFactorAuthentication = true

	err = s.AuthRepo.DeleteAllUserRefreshTokens(ctx, user.Id)
	if err != nil {
		slog.ErrorContext(ctx, "database error revoking all tokens for user",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return Enable2FAResult{}, models.ErrUnexpectedLogout
	}

	nowTimestamp := time.Now().Add(-1 * time.Second).Unix()
	redisTTL := time.Duration(s.JwtAccessExpiration) * time.Minute

	redisValue := fmt.Sprintf("%d", nowTimestamp)
	err = s.RedisClient.Set(ctx, "logout_all:"+user.Id, redisValue, redisTTL).Err()
	if err != nil {
		slog.ErrorContext(ctx, "failed to set logout_all timestamp in redis",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
	}

	hasPassword := user.PasswordHash != nil

	_, accessToken, err := utils.GenerateAccessToken(
		user.Id, user.Name, user.Username, user.Email, s.AppDomain, user.Permissions,
		s.JwtAccessPrivateKey, user.Role.IsEmployee, user.TwoFactorAuthentication,
		hasPassword, time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate access token",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return Enable2FAResult{}, models.ErrGeneratingToken
	}

	refreshExpiresAt := time.Now().AddDate(0, 0, s.JwtRefreshExpiration)
	idToken, familyId, refreshToken, err := utils.GenerateRefreshToken(
		user.Id, s.AppDomain, "", data.DeviceHash, s.JwtRefreshPrivateKey, refreshExpiresAt,
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate refresh token",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return Enable2FAResult{}, models.ErrGeneratingToken
	}

	err = s.AuthRepo.SaveRefreshToken(ctx, idToken, user.Id, familyId, refreshExpiresAt)
	if err != nil {
		slog.ErrorContext(ctx, "failed to persist refresh token into database",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return Enable2FAResult{}, models.ErrGeneratingToken
	}

	return Enable2FAResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *TwoFactorService) Disable2FA(ctx context.Context, data dto.TwoFactorDisableRequest) error {
	claims, err := utils.ValidateRefreshToken(data.RefreshToken, s.JwtRefreshPublicKey, s.AppDomain)
	if err != nil {
		slog.WarnContext(ctx, "refresh token validation failed during 2FA disable",
			slog.String("op", "Disable2FA"),
			slog.String("user_id", data.UserId),
			slog.Any("error", err),
		)
		return models.ErrInvalidOrExpiredRefresh
	}

	if claims.Subject != data.UserId {
		slog.WarnContext(ctx, "security mismatch during 2fa disable",
			slog.String("token_subject", claims.Subject),
			slog.String("user_id_param", data.UserId),
		)
		return models.ErrTokenMetadataMisMatch
	}

	if claims.DeviceHash != data.DeviceHash {
		slog.WarnContext(ctx, "[SECURITY ALERT] Device hash mismatch during 2FA disable! Revoking family",
			slog.String("op", "Disable2FA"),
			slog.String("user_id", data.UserId),
			slog.String("family_id", claims.FamilyId),
		)
		_ = s.AuthRepo.RevokeFamily(ctx, claims.FamilyId)
		return models.ErrInvalidOrExpiredRefresh
	}

	dbToken, err := s.AuthRepo.GetRefreshTokenById(ctx, claims.ID)
	if err != nil || dbToken == nil || dbToken.IsRevoked {
		return models.ErrInvalidOrExpiredRefresh
	}

	twoFactorData, err := s.UserRepo.GetAuthData(ctx, data.UserId)
	if err != nil {
		slog.ErrorContext(ctx, "database query error",
			slog.String("op", "Disable2FA"),
			slog.String("user_id", data.UserId),
			slog.Any("error", err),
		)
		return models.Err2FAGetDataFailed
	}

	if !twoFactorData.IsEnabled {
		return models.Err2FAAlreadyDisabled
	}

	if twoFactorData.PasswordHash != "" {
		isValid, err := utils.VerifyPassword(data.Password, twoFactorData.PasswordHash)
		if err != nil {
			slog.ErrorContext(ctx, "Argon2 verification system error",
				slog.String("op", "Disable2FA"),
				slog.String("user_id", data.UserId),
				slog.Any("error", err),
			)
			return models.ErrPasswordVerificationFailed
		}

		if !isValid {
			return models.ErrIncorrectPassword
		}
	}

	err = s.UserRepo.Disable2FA(ctx, data.UserId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update 2FA status to disabled in DB",
			slog.String("op", "Disable2FA"),
			slog.String("user_id", data.UserId),
			slog.Any("error", err),
		)
		return models.Err2FAUpdateFailed
	}

	if err = s.AuthRepo.DeleteAllUserRefreshTokens(ctx, data.UserId); err != nil {
		slog.WarnContext(ctx, "failed to revoke refresh tokens",
			slog.String("op", "Disable2FA"),
			slog.String("user_id", data.UserId),
			slog.Any("error", err),
		)
	}

	nowTimestamp := time.Now().Add(-1 * time.Second).Unix()
	redisTTL := time.Duration(s.JwtAccessExpiration) * time.Minute

	if err = s.RedisClient.Set(ctx, "logout_all:"+data.UserId, fmt.Sprintf("%d", nowTimestamp), redisTTL).Err(); err != nil {
		slog.ErrorContext(ctx, "failed to set logout_all in redis",
			slog.String("op", "Disable2FA"),
			slog.String("user_id", data.UserId),
			slog.Any("error", err),
		)
	}

	return nil
}
