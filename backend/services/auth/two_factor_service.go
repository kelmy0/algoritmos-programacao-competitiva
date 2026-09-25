package auth

import (
	"context"
	"crypto/ed25519"
	"errors"
	"log/slog"
	"time"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
)

type twoFactorUserRepository interface {
	GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error)
}

type twoFactorSessionIssuer interface {
	IssueSession(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (AuthResult, error)
}

type twoFactorRedis interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
}

type TwoFactorService struct {
	userRepo            twoFactorUserRepository
	sessionIssuer       twoFactorSessionIssuer
	redisClient         twoFactorRedis
	jwtAccessPublicKey  ed25519.PublicKey
	jwtAccessPrivateKey ed25519.PrivateKey
	appDomain           string
	encryptSecret       string
}

func NewTwoFactorService(userRepo twoFactorUserRepository, sessionIssuer twoFactorSessionIssuer,
	redisClient twoFactorRedis, jwtAccessPublicKey ed25519.PublicKey,
	jwtAccessPrivateKey ed25519.PrivateKey, appDomain, encryptSecret string) *TwoFactorService {
	return &TwoFactorService{
		userRepo:            userRepo,
		sessionIssuer:       sessionIssuer,
		redisClient:         redisClient,
		jwtAccessPublicKey:  jwtAccessPublicKey,
		jwtAccessPrivateKey: jwtAccessPrivateKey,
		appDomain:           appDomain,
		encryptSecret:       encryptSecret,
	}
}

func (s *TwoFactorService) VerifyLogin2FA(ctx context.Context, data dto.Verify2FARequest) (AuthResult, error) {
	claims, err := utils.ValidateAccessToken(data.PreAuthToken, s.jwtAccessPublicKey, s.appDomain)
	if err != nil {
		slog.WarnContext(ctx, "pre-auth token validation failed during 2FA", slog.Any("error", err))
		return AuthResult{}, models.ErrSessionExpired
	}

	userId := claims.Subject
	if userId == "" {
		slog.WarnContext(ctx, "pre-auth token claims missing Subject field")
		return AuthResult{}, models.ErrSessionData
	}

	if claims.DeviceHash != data.DeviceHash {
		slog.WarnContext(ctx, "[SECURITY ALERT] Device hash mismatch!",
			slog.String("user_id", claims.Subject),
			slog.String("token_id", claims.ID),
			slog.String("token_dvh", claims.DeviceHash),
			slog.String("dvh", data.DeviceHash),
		)

		return AuthResult{}, models.ErrSessionExpired
	}

	if claims.ID != "" {
		blacklisted, _ := s.redisClient.Exists(ctx, "blacklist:jti:"+claims.ID).Result()
		if blacklisted > 0 {
			slog.WarnContext(ctx, "attempted re-use of pre-auth token", slog.String("jti", claims.ID))
			return AuthResult{}, models.ErrSessionExpired
		}
	}

	user, err := s.userRepo.GetUserByIdForAuth(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return AuthResult{}, models.ErrUserNotFound
		}
		slog.ErrorContext(ctx, "database error querying user during 2FA",
			slog.String("user_id", userId),
			slog.Any("error", err),
		)
		return AuthResult{}, models.ErrFailQueryUser
	}

	if !user.Enable {
		slog.WarnContext(ctx, "disabled user attempted 2FA login", slog.String("user_id", user.Id))
		return AuthResult{}, models.ErrUserNotEnabled
	}

	if user.TwoFactorSecret == nil || *user.TwoFactorSecret == "" {
		slog.WarnContext(ctx, "user attempted 2FA verification without secret configured", slog.String("user_id", user.Id))
		return AuthResult{}, models.Err2FANotInitiated
	}

	decryptedSecret, err := utils.Decrypt(*user.TwoFactorSecret, s.encryptSecret)
	if err != nil {
		slog.ErrorContext(ctx, "AES decryption of 2FA secret failed",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return AuthResult{}, models.ErrUnexpectedLogin
	}

	isValid := totp.Validate(data.Code, decryptedSecret)
	if !isValid {
		return AuthResult{}, models.Err2FAInvalid
	}

	if claims.ID != "" && claims.ExpiresAt != nil {
		ttl := time.Until(claims.ExpiresAt.Time)
		if ttl > 0 {
			_ = s.redisClient.Set(ctx, "blacklist:jti:"+claims.ID, "used", ttl).Err()
		}
	}

	hasPassword := user.PasswordHash != nil

	return s.sessionIssuer.IssueSession(ctx, user, data.DeviceHash, hasPassword)
}
