package auth

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"github.com/redis/go-redis/v9"
)

type sessionAuthRepository interface {
	RotateRefreshToken(ctx context.Context, oldTokenId, newTokenId, userId, familyId string, newExpiresAt time.Time) error
	RevokeFamily(ctx context.Context, familyId string) error
	GetRefreshTokenById(ctx context.Context, id string) (*models.RefreshToken, error)
	DeleteFamily(ctx context.Context, familyId string) error
	DeleteOtherFamilies(ctx context.Context, userId, currentFamilyId string) error
	DeleteAllUserRefreshTokens(ctx context.Context, userId string) error
}

type sessionUserRepository interface {
	GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error)
}

type sessionRedis interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
}

type SessionService struct {
	authRepo             sessionAuthRepository
	userRepo             sessionUserRepository
	redisClient          sessionRedis
	jwtRefreshPublicKey  ed25519.PublicKey
	jwtAccessPrivateKey  ed25519.PrivateKey
	jwtRefreshPrivateKey ed25519.PrivateKey
	jwtAccessExpiration  int
	jwtRefreshExpiration int
	appDomain            string
}

func NewSessionService(authRepo sessionAuthRepository, userRepo sessionUserRepository,
	redisClient sessionRedis, jwtRefreshPublicKey ed25519.PublicKey,
	jwtAccessPrivateKey, jwtRefreshPrivateKey ed25519.PrivateKey,
	jwtAccessExpiration, jwtRefreshExpiration int, appDomain string) *SessionService {
	return &SessionService{
		authRepo:             authRepo,
		userRepo:             userRepo,
		redisClient:          redisClient,
		jwtRefreshPublicKey:  jwtRefreshPublicKey,
		jwtAccessPrivateKey:  jwtAccessPrivateKey,
		jwtRefreshPrivateKey: jwtRefreshPrivateKey,
		jwtAccessExpiration:  jwtAccessExpiration,
		jwtRefreshExpiration: jwtRefreshExpiration,
		appDomain:            appDomain,
	}
}

func (s *SessionService) RefreshToken(ctx context.Context, refreshTokenString, deviceHash string) (RefreshTokenResult, error) {
	claims, err := utils.ValidateRefreshToken(refreshTokenString, s.jwtRefreshPublicKey, s.appDomain)
	if err != nil {
		slog.WarnContext(ctx, "refresh token JWT validation failed", slog.Any("error", err))
		return RefreshTokenResult{}, models.ErrInvalidOrExpiredRefresh
	}

	if claims.DeviceHash != deviceHash {
		slog.WarnContext(ctx, "[SECURITY ALERT] Device hash mismatch! Revoking entire family.",
			slog.String("user_id", claims.Subject),
			slog.String("family_id", claims.FamilyId),
			slog.String("token_id", claims.ID),
			slog.String("token_dvh", claims.DeviceHash),
			slog.String("dvh", deviceHash),
		)
		_ = s.authRepo.RevokeFamily(ctx, claims.FamilyId)
		return RefreshTokenResult{}, models.ErrInvalidOrExpiredRefresh
	}

	dbToken, err := s.authRepo.GetRefreshTokenById(ctx, claims.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error querying session database for refresh token",
			slog.String("token_id", claims.ID),
			slog.Any("error", err),
		)
		return RefreshTokenResult{}, models.ErrInvalidOrExpiredRefresh
	}

	if dbToken == nil {
		return RefreshTokenResult{}, models.ErrInvalidOrExpiredRefresh
	}

	if dbToken.IsRevoked {
		slog.WarnContext(ctx, "[SECURITY ALERT] Revoked refresh token reused! Revoking entire family.",
			slog.String("user_id", dbToken.UserId),
			slog.String("family_id", dbToken.FamilyId),
			slog.String("token_id", dbToken.Id),
		)
		_ = s.authRepo.RevokeFamily(ctx, dbToken.FamilyId)
		return RefreshTokenResult{}, models.ErrInvalidOrExpiredRefresh
	}

	if dbToken.UserId != claims.Subject {
		slog.WarnContext(ctx, "token integrity mismatch: DB UserId does not match token Subject",
			slog.String("db_user_id", dbToken.UserId),
			slog.String("token_subject", claims.Subject),
		)
		return RefreshTokenResult{}, models.ErrTokenMetadataMisMatch
	}

	user, err := s.userRepo.GetUserByIdForAuth(ctx, claims.Subject)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return RefreshTokenResult{}, models.ErrUserNotFound
		}
		slog.ErrorContext(ctx, "error retrieving user during session refresh",
			slog.String("user_id", claims.Subject),
			slog.Any("error", err),
		)
		return RefreshTokenResult{}, models.ErrUserNotFound
	}

	if !user.Enable {
		slog.WarnContext(ctx, "disabled user attempted token refresh", slog.String("user_id", user.Id))
		return RefreshTokenResult{}, models.ErrUserNotEnabled
	}

	hasPassword := user.PasswordHash != nil

	_, newAccessToken, err := utils.GenerateAccessToken(
		user.Id, user.Name, user.Username, user.Email, s.appDomain, user.Permissions,
		s.jwtAccessPrivateKey, user.Role.IsEmployee, user.TwoFactorAuthentication,
		hasPassword, time.Now().Add(time.Duration(s.jwtAccessExpiration)*time.Minute),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to sign new access token during refresh",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return RefreshTokenResult{}, models.ErrGeneratingToken
	}

	newRefreshExpiresAt := time.Now().AddDate(0, 0, s.jwtRefreshExpiration)
	newTokenId, _, newRefreshToken, err := utils.GenerateRefreshToken(
		user.Id, s.appDomain, dbToken.FamilyId, deviceHash, s.jwtRefreshPrivateKey, newRefreshExpiresAt,
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to sign new refresh token during rotation",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return RefreshTokenResult{}, models.ErrGeneratingToken
	}

	err = s.authRepo.RotateRefreshToken(ctx, dbToken.Id, newTokenId, user.Id, dbToken.FamilyId, newRefreshExpiresAt)
	if err != nil {
		slog.ErrorContext(ctx, "failed to rotate refresh token in database",
			slog.String("user_id", user.Id),
			slog.String("old_token_id", dbToken.Id),
			slog.Any("error", err),
		)
		return RefreshTokenResult{}, models.ErrGeneratingToken
	}

	return RefreshTokenResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *SessionService) Logout(ctx context.Context, userId, refreshTokenString, accessJti string, accessExpiresAt time.Time) error {
	claims, err := utils.ValidateRefreshToken(refreshTokenString, s.jwtRefreshPublicKey, s.appDomain)
	if err != nil {
		return models.ErrInvalidOrExpiredRefresh
	}

	if claims.Subject != userId {
		slog.WarnContext(ctx, "security mismatch during logout",
			slog.String("token_subject", claims.Subject),
			slog.String("user_id_param", userId),
		)
		return models.ErrTokenMetadataMisMatch
	}

	err = s.authRepo.DeleteFamily(ctx, claims.FamilyId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete session family during logout",
			slog.String("family_id", claims.FamilyId),
			slog.String("user_id", userId),
			slog.Any("error", err),
		)
		return models.ErrUnexpectedLogout
	}

	ttl := time.Until(accessExpiresAt)
	if ttl > 0 {
		err = s.redisClient.Set(ctx, "blacklist:jti:"+accessJti, "revoked", ttl).Err()
		if err != nil {
			slog.WarnContext(ctx, "failed to blacklist access token in redis",
				slog.String("access_jti", accessJti),
				slog.String("user_id", userId),
				slog.Any("error", err),
			)
		}
	}

	return nil
}

func (s *SessionService) LogoutOtherDevices(ctx context.Context, userId, refreshTokenString, accessJti, deviceHash string) error {
	claims, err := utils.ValidateRefreshToken(refreshTokenString, s.jwtRefreshPublicKey, s.appDomain)
	if err != nil {
		return models.ErrInvalidOrExpiredRefresh
	}

	if claims.Subject != userId {
		slog.WarnContext(ctx, "security mismatch during logout other devices",
			slog.String("token_subject", claims.Subject),
			slog.String("user_id_param", userId),
		)
		return models.ErrTokenMetadataMisMatch
	}

	if claims.DeviceHash != deviceHash {
		slog.WarnContext(ctx, "[SECURITY ALERT] Device hash mismatch! Revoking entire family.",
			slog.String("user_id", claims.Subject),
			slog.String("family_id", claims.FamilyId),
			slog.String("token_id", claims.ID),
			slog.String("token_dvh", claims.DeviceHash),
			slog.String("dvh", deviceHash),
		)
		_ = s.authRepo.RevokeFamily(ctx, claims.FamilyId)
		return models.ErrInvalidOrExpiredRefresh
	}

	err = s.authRepo.DeleteOtherFamilies(ctx, userId, claims.FamilyId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete other families during logout",
			slog.String("user_id", userId),
			slog.String("current_family_id", claims.FamilyId),
			slog.Any("error", err),
		)
		return models.ErrUnexpectedLogout
	}

	nowTimestamp := time.Now().Unix()
	redisTTL := time.Duration(s.jwtAccessExpiration) * time.Minute

	redisValue := fmt.Sprintf("%d:%s", nowTimestamp, accessJti)
	err = s.redisClient.Set(ctx, "logout_other:"+userId, redisValue, redisTTL).Err()

	if err != nil {
		slog.ErrorContext(ctx, "failed to set logout_other timestamp in redis",
			slog.String("user_id", userId),
			slog.Any("error", err),
		)
		return models.ErrUnexpectedLogout
	}

	return nil
}

func (s *SessionService) LogoutAllDevices(ctx context.Context, userId, refreshTokenString, deviceHash string) error {
	claims, err := utils.ValidateRefreshToken(refreshTokenString, s.jwtRefreshPublicKey, s.appDomain)
	if err != nil {
		return models.ErrInvalidOrExpiredRefresh
	}

	if claims.Subject != userId {
		slog.WarnContext(ctx, "security mismatch during logout all",
			slog.String("token_subject", claims.Subject),
			slog.String("user_id_param", userId),
		)
		return models.ErrTokenMetadataMisMatch
	}

	if claims.DeviceHash != deviceHash {
		slog.WarnContext(ctx, "[SECURITY ALERT] Device hash mismatch! Revoking entire family.",
			slog.String("user_id", claims.Subject),
			slog.String("family_id", claims.FamilyId),
			slog.String("token_id", claims.ID),
			slog.String("token_dvh", claims.DeviceHash),
			slog.String("dvh", deviceHash),
		)
		_ = s.authRepo.RevokeFamily(ctx, claims.FamilyId)
		return models.ErrInvalidOrExpiredRefresh
	}

	err = s.authRepo.DeleteAllUserRefreshTokens(ctx, userId)
	if err != nil {
		slog.ErrorContext(ctx, "database error revoking all tokens for user",
			slog.String("user_id", userId),
			slog.Any("error", err),
		)
		return models.ErrUnexpectedLogout
	}

	nowTimestamp := time.Now().Unix()
	redisTTL := time.Duration(s.jwtAccessExpiration) * time.Minute

	redisValue := fmt.Sprintf("%d", nowTimestamp)
	err = s.redisClient.Set(ctx, "logout_all:"+userId, redisValue, redisTTL).Err()
	if err != nil {
		slog.ErrorContext(ctx, "failed to set logout_all timestamp in redis",
			slog.String("user_id", userId),
			slog.Any("error", err),
		)
		return models.ErrUnexpectedLogout
	}

	return nil
}
