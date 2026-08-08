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
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
)

type AuthRepository interface {
	SaveRefreshToken(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error
	GetRefreshTokenById(ctx context.Context, id string) (*models.RefreshToken, error)
	RotateRefreshToken(ctx context.Context, oldTokenId, newTokenId, userId, familyId string, newExpiresAt time.Time) error
	RevokeFamily(ctx context.Context, familyId string) error
	DeleteFamily(ctx context.Context, familyId string) error
	DeleteOtherFamilies(ctx context.Context, userId, currentFamilyId string) error
	DeleteAllUserRefreshTokens(ctx context.Context, userId string) error
}

type AuthUserRepository interface {
	GetUserByEmailForAuth(ctx context.Context, email string) (*models.User, error)
	GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error)
	GetUserBySocialID(ctx context.Context, provider, socialId string) (*models.User, error)
	CreateSocialUser(ctx context.Context, newUser models.NewUserGoogle, provider, socialId string) (*models.User, error)
	CreateSocialLink(ctx context.Context, id, provider, socialUserId string) error
}

type AuthService struct {
	AuthRepo             AuthRepository
	UserRepo             AuthUserRepository
	RedisClient          *redis.Client
	JwtAccessPrivateKey  ed25519.PrivateKey
	JwtRefreshPrivateKey ed25519.PrivateKey
	JwtAccessPublicKey   ed25519.PublicKey
	JwtRefreshPublicKey  ed25519.PublicKey
	JwtAccessExpiration  int
	JwtRefreshExpiration int
	AppDomain            string
	EncryptSecret        string
}

type AuthResult struct {
	LoginResponse *dto.LoginResponse
	RefreshToken  string
}

type RefreshTokenResult struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthService(authRepo AuthRepository, userRepo AuthUserRepository, redisClient *redis.Client, jwtAccessPrivateKey, jwtRefreshPrivateKey ed25519.PrivateKey, jwtAccessPublicKey, jwtRefreshPublicKey ed25519.PublicKey, appDomain, encryptSecret string, jwtAccessExpiration int, jwtRefreshExpiration int) *AuthService {
	return &AuthService{
		AuthRepo:             authRepo,
		UserRepo:             userRepo,
		RedisClient:          redisClient,
		JwtAccessPrivateKey:  jwtAccessPrivateKey,
		JwtRefreshPrivateKey: jwtRefreshPrivateKey,
		JwtAccessPublicKey:   jwtAccessPublicKey,
		JwtRefreshPublicKey:  jwtRefreshPublicKey,
		AppDomain:            appDomain,
		JwtAccessExpiration:  jwtAccessExpiration,
		JwtRefreshExpiration: jwtRefreshExpiration,
		EncryptSecret:        encryptSecret,
	}
}

func (s *AuthService) Auth(ctx context.Context, data dto.AuthRequest) (*AuthResult, error) {
	maskedEmail := utils.MaskEmail(data.Email)
	user, err := s.UserRepo.GetUserByEmailForAuth(ctx, data.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrInvalidEmailOrPassword
		}
		slog.ErrorContext(ctx, "database query error fetching user by email during auth",
			slog.String("email", maskedEmail),
			slog.Any("error", err),
		)
		return nil, models.ErrInvalidEmailOrPassword
	}

	if !user.Enable {
		slog.WarnContext(ctx, "login blocked: account is disabled",
			slog.String("user_id", user.Id),
			slog.String("email", maskedEmail),
		)
		return nil, models.ErrInvalidEmailOrPassword
	}

	if user.PasswordHash == nil {
		slog.WarnContext(ctx, "login blocked: user has no local password configured",
			slog.String("user_id", user.Id),
			slog.String("email", maskedEmail),
		)
		return nil, models.ErrInvalidEmailOrPassword
	}

	isValid, err := utils.VerifyPassword(data.Password, *user.PasswordHash)
	if err != nil {
		slog.ErrorContext(ctx, "Argon2 password verification failed",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return nil, models.ErrPasswordVerificationFailed
	}
	if !isValid {
		return nil, models.ErrInvalidEmailOrPassword
	}

	if user.TwoFactorAuthentication {
		_, preAuthToken, err := utils.GeneratePreAuthToken(user.Id, s.AppDomain, data.DeviceHash, s.JwtAccessPrivateKey, time.Now().Add(5*time.Minute))
		if err != nil {
			slog.ErrorContext(ctx, "failed to generate 2FA pre-auth token",
				slog.String("user_id", user.Id),
				slog.Any("error", err),
			)
			return nil, models.ErrUnexpectedLogin
		}

		response := &dto.LoginResponse{
			Requires2FA:  true,
			PreAuthToken: preAuthToken,
		}

		return &AuthResult{LoginResponse: response, RefreshToken: ""}, nil
	}

	return s.issueSession(ctx, user, data.DeviceHash, user.TwoFactorAuthentication)
}

func (s *AuthService) VerifyLogin2FA(ctx context.Context, data dto.Verify2FARequest) (*AuthResult, error) {
	claims, err := utils.ValidateAccessToken(data.PreAuthToken, s.JwtAccessPublicKey, s.AppDomain)
	if err != nil {
		slog.WarnContext(ctx, "pre-auth token validation failed during 2FA", slog.Any("error", err))
		return nil, models.ErrSessionExpired
	}

	userId := claims.Subject
	if userId == "" {
		slog.WarnContext(ctx, "pre-auth token claims missing Subject field")
		return nil, models.ErrSessionData
	}

	/*
		if claims.DeviceHash != data.DeviceHash {
			slog.WarnContext(ctx, "[SECURITY ALERT] Device hash mismatch!",
				slog.String("user_id", claims.Subject),
				slog.String("token_id", claims.ID),
				slog.String("token_dvh", claims.DeviceHash),
				slog.String("dvh", data.DeviceHash),
			)

			return nil, models.ErrSessionExpired
		}*/

	if claims.ID != "" {
		blacklisted, _ := s.RedisClient.Exists(ctx, "blacklist:jti:"+claims.ID).Result()
		if blacklisted > 0 {
			slog.WarnContext(ctx, "attempted re-use of pre-auth token", slog.String("jti", claims.ID))
			return nil, models.ErrSessionExpired
		}
	}

	user, err := s.UserRepo.GetUserByIdForAuth(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}
		slog.ErrorContext(ctx, "database error querying user during 2FA",
			slog.String("user_id", userId),
			slog.Any("error", err),
		)
		return nil, models.ErrFailQueryUser
	}

	if !user.Enable {
		slog.WarnContext(ctx, "disabled user attempted 2FA login", slog.String("user_id", user.Id))
		return nil, models.ErrUserNotEnabled
	}

	if user.TwoFactorSecret == nil || *user.TwoFactorSecret == "" {
		slog.WarnContext(ctx, "user attempted 2FA verification without secret configured", slog.String("user_id", user.Id))
		return nil, models.Err2FANotInitiated
	}

	decryptedSecret, err := utils.Decrypt(*user.TwoFactorSecret, s.EncryptSecret)
	if err != nil {
		slog.ErrorContext(ctx, "AES decryption of 2FA secret failed",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return nil, models.ErrUnexpectedLogin
	}

	isValid := totp.Validate(data.Code, decryptedSecret)
	if !isValid {
		return nil, models.Err2FAInvalid
	}

	if claims.ID != "" && claims.ExpiresAt != nil {
		ttl := time.Until(claims.ExpiresAt.Time)
		if ttl > 0 {
			_ = s.RedisClient.Set(ctx, "blacklist:jti:"+claims.ID, "used", ttl).Err()
		}
	}

	return s.issueSession(ctx, user, data.DeviceHash, user.TwoFactorAuthentication)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString, deviceHash string) (*RefreshTokenResult, error) {
	claims, err := utils.ValidateRefreshToken(refreshTokenString, s.JwtRefreshPublicKey, s.AppDomain)
	if err != nil {
		slog.WarnContext(ctx, "refresh token JWT validation failed", slog.Any("error", err))
		return nil, models.ErrInvalidOrExpiredRefresh
	}

	if claims.DeviceHash != deviceHash {
		slog.WarnContext(ctx, "[SECURITY ALERT] Device hash mismatch! Revoking entire family.",
			slog.String("user_id", claims.Subject),
			slog.String("family_id", claims.FamilyId),
			slog.String("token_id", claims.ID),
			slog.String("token_dvh", claims.DeviceHash),
			slog.String("dvh", deviceHash),
		)
		_ = s.AuthRepo.RevokeFamily(ctx, claims.FamilyId)
		return nil, models.ErrInvalidOrExpiredRefresh
	}

	dbToken, err := s.AuthRepo.GetRefreshTokenById(ctx, claims.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error querying session database for refresh token",
			slog.String("token_id", claims.ID),
			slog.Any("error", err),
		)
		return nil, models.ErrInvalidOrExpiredRefresh
	}

	if dbToken == nil {
		return nil, models.ErrInvalidOrExpiredRefresh
	}

	if dbToken.IsRevoked {
		slog.WarnContext(ctx, "[SECURITY ALERT] Revoked refresh token reused! Revoking entire family.",
			slog.String("user_id", dbToken.UserId),
			slog.String("family_id", dbToken.FamilyId),
			slog.String("token_id", dbToken.Id),
		)
		_ = s.AuthRepo.RevokeFamily(ctx, dbToken.FamilyId)
		return nil, models.ErrInvalidOrExpiredRefresh
	}

	if dbToken.UserId != claims.Subject {
		slog.WarnContext(ctx, "token integrity mismatch: DB UserId does not match token Subject",
			slog.String("db_user_id", dbToken.UserId),
			slog.String("token_subject", claims.Subject),
		)
		return nil, models.ErrTokenMetadataMisMatch
	}

	user, err := s.UserRepo.GetUserByIdForAuth(ctx, claims.Subject)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}
		slog.ErrorContext(ctx, "error retrieving user during session refresh",
			slog.String("user_id", claims.Subject),
			slog.Any("error", err),
		)
		return nil, models.ErrUserNotFound
	}

	if !user.Enable {
		slog.WarnContext(ctx, "disabled user attempted token refresh", slog.String("user_id", user.Id))
		return nil, models.ErrUserNotEnabled
	}

	_, newAccessToken, err := utils.GenerateAccessToken(
		user.Id, user.Username, user.Email, s.AppDomain, user.Permissions,
		s.JwtAccessPrivateKey, user.Role.IsEmployee, claims.Is2FAEnabled,
		time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to sign new access token during refresh",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return nil, models.ErrGeneratingToken
	}

	newRefreshExpiresAt := time.Now().AddDate(0, 0, s.JwtRefreshExpiration)
	newTokenId, _, newRefreshToken, err := utils.GenerateRefreshToken(
		user.Id, s.AppDomain, dbToken.FamilyId, deviceHash, s.JwtRefreshPrivateKey, newRefreshExpiresAt,
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to sign new refresh token during rotation",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return nil, models.ErrGeneratingToken
	}

	err = s.AuthRepo.RotateRefreshToken(ctx, dbToken.Id, newTokenId, user.Id, dbToken.FamilyId, newRefreshExpiresAt)
	if err != nil {
		slog.ErrorContext(ctx, "failed to rotate refresh token in database",
			slog.String("user_id", user.Id),
			slog.String("old_token_id", dbToken.Id),
			slog.Any("error", err),
		)
		return nil, models.ErrGeneratingToken
	}

	return &RefreshTokenResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userId, refreshTokenString, accessJti string, accessExpiresAt time.Time) error {
	claims, err := utils.ValidateRefreshToken(refreshTokenString, s.JwtRefreshPublicKey, s.AppDomain)
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

	err = s.AuthRepo.DeleteFamily(ctx, claims.FamilyId)
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
		err = s.RedisClient.Set(ctx, "blacklist:jti:"+accessJti, "revoked", ttl).Err()
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

func (s *AuthService) LogoutOtherDevices(ctx context.Context, userId, refreshTokenString, accessJti, deviceHash string) error {
	claims, err := utils.ValidateRefreshToken(refreshTokenString, s.JwtRefreshPublicKey, s.AppDomain)
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
		_ = s.AuthRepo.RevokeFamily(ctx, claims.FamilyId)
		return models.ErrInvalidOrExpiredRefresh
	}

	err = s.AuthRepo.DeleteOtherFamilies(ctx, userId, claims.FamilyId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete other families during logout",
			slog.String("user_id", userId),
			slog.String("current_family_id", claims.FamilyId),
			slog.Any("error", err),
		)
		return models.ErrUnexpectedLogout
	}

	nowTimestamp := time.Now().Unix()
	redisTTL := time.Duration(s.JwtAccessExpiration) * time.Minute

	redisValue := fmt.Sprintf("%d:%s", nowTimestamp, accessJti)
	err = s.RedisClient.Set(ctx, "logout_other:"+userId, redisValue, redisTTL).Err()

	if err != nil {
		slog.ErrorContext(ctx, "failed to set logout_other timestamp in redis",
			slog.String("user_id", userId),
			slog.Any("error", err),
		)
		return models.ErrUnexpectedLogout
	}

	return nil
}

func (s *AuthService) LogoutAllDevices(ctx context.Context, userId, refreshTokenString, deviceHash string) error {
	claims, err := utils.ValidateRefreshToken(refreshTokenString, s.JwtRefreshPublicKey, s.AppDomain)
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
		_ = s.AuthRepo.RevokeFamily(ctx, claims.FamilyId)
		return models.ErrInvalidOrExpiredRefresh
	}

	err = s.AuthRepo.DeleteAllUserRefreshTokens(ctx, userId)
	if err != nil {
		slog.ErrorContext(ctx, "database error revoking all tokens for user",
			slog.String("user_id", userId),
			slog.Any("error", err),
		)
		return models.ErrUnexpectedLogout
	}

	nowTimestamp := time.Now().Unix()
	redisTTL := time.Duration(s.JwtAccessExpiration) * time.Minute

	redisValue := fmt.Sprintf("%d", nowTimestamp)
	err = s.RedisClient.Set(ctx, "logout_all:"+userId, redisValue, redisTTL).Err()
	if err != nil {
		slog.ErrorContext(ctx, "failed to set logout_all timestamp in redis",
			slog.String("user_id", userId),
			slog.Any("error", err),
		)
		return models.ErrUnexpectedLogout
	}

	return nil
}

func (s *AuthService) AuthWithSocialProvider(ctx context.Context, provider, socialUserId, email, name, deviceHash string) (*AuthResult, error) {
	maskedEmail := utils.MaskEmail(email)
	user, err := s.UserRepo.GetUserBySocialID(ctx, provider, socialUserId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			user, err = s.UserRepo.GetUserByEmailForAuth(ctx, email)

			if err != nil {
				if errors.Is(err, models.ErrUserNotFound) {
					username := utils.NormalizeUsername(name)
					newUser := models.NewUserGoogle{
						Name:         name,
						Username:     username,
						Email:        email,
						Provider:     provider,
						SocialUserId: socialUserId,
					}

					user, err = s.UserRepo.CreateSocialUser(ctx, newUser, provider, socialUserId)
					if err != nil {
						switch {
						case errors.Is(err, models.ErrEmailAlreadyUsed),
							errors.Is(err, models.ErrUsernameAlreadyUsed),
							errors.Is(err, models.ErrUserAlreadyExists):
							return nil, err

						default:
							slog.ErrorContext(ctx, "failed to register social user",
								slog.String("email", maskedEmail),
								slog.String("provider", provider),
								slog.Any("error", err),
							)
							return nil, models.ErrRegisterSocialUser
						}
					}
				} else {
					slog.ErrorContext(ctx, "database error checking email existence during social auth",
						slog.String("email", maskedEmail),
						slog.Any("error", err),
					)
					return nil, models.ErrFailQueryUser
				}
			} else {
				slog.WarnContext(ctx, "social login blocked: email already exists with a different auth strategy",
					slog.String("email", maskedEmail),
					slog.String("attempted_provider", provider),
				)
				return nil, models.ErrUserAlreadyExists
			}
		} else {
			slog.ErrorContext(ctx, "database error fetching user by social ID",
				slog.String("provider", provider),
				slog.Any("error", err),
			)
			return nil, models.ErrFailQueryUser
		}
	}

	if !user.Enable {
		slog.WarnContext(ctx, "disabled user attempted social login",
			slog.String("user_id", user.Id),
			slog.String("provider", provider),
		)
		return nil, models.ErrUserNotEnabled
	}

	if user.TwoFactorAuthentication {
		_, preAuthToken, err := utils.GeneratePreAuthToken(
			user.Id, s.AppDomain, deviceHash, s.JwtAccessPrivateKey, time.Now().Add(5*time.Minute),
		)
		if err != nil {
			slog.ErrorContext(ctx, "failed to generate 2FA pre-auth token during social auth",
				slog.String("user_id", user.Id),
				slog.Any("error", err),
			)
			return nil, models.ErrUnexpectedLogin
		}

		response := &dto.LoginResponse{
			Requires2FA:  true,
			PreAuthToken: preAuthToken,
		}

		return &AuthResult{LoginResponse: response, RefreshToken: ""}, nil
	}

	return s.issueSession(ctx, user, deviceHash, user.TwoFactorAuthentication)
}

func (s *AuthService) LinkSocialAccount(ctx context.Context, currentUserId, provider, socialUserId, email string) error {
	existingUser, err := s.UserRepo.GetUserBySocialID(ctx, provider, socialUserId)
	if err == nil && existingUser != nil {
		slog.WarnContext(ctx, "social ID already linked to another account",
			slog.String("social_user_id", socialUserId),
			slog.String("provider", provider),
		)
		return models.ErrSocialAccountAlreadyLinked
	}

	if err != nil && !errors.Is(err, models.ErrUserNotFound) {
		slog.ErrorContext(ctx, "database error checking social ID existence during linking",
			slog.String("provider", provider),
			slog.Any("error", err),
		)
		return models.ErrFailQueryUser
	}

	currentUser, err := s.UserRepo.GetUserByIdForAuth(ctx, currentUserId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find current user for social linking",
			slog.String("user_id", currentUserId),
			slog.Any("error", err),
		)
		return models.ErrUserNotFound
	}

	if currentUser.Email != email {
		slog.WarnContext(ctx, "email mismatch during social account link attempt",
			slog.String("user_id", currentUserId),
			slog.String("account_email", utils.MaskEmail(currentUser.Email)),
			slog.String("social_email", utils.MaskEmail(email)),
		)
		return models.ErrEmailMismatchForSocialLink
	}

	err = s.UserRepo.CreateSocialLink(ctx, currentUser.Id, provider, socialUserId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create social link in database",
			slog.String("user_id", currentUser.Id),
			slog.String("provider", provider),
			slog.Any("error", err),
		)
		return models.ErrLinkSocialAccount
	}

	return nil
}

func (s *AuthService) issueSession(ctx context.Context, user *models.User, deviceHash string, is2FAEnabled bool) (*AuthResult, error) {
	// Access Token
	_, accessToken, err := utils.GenerateAccessToken(
		user.Id, user.Username, user.Email, s.AppDomain, user.Permissions,
		s.JwtAccessPrivateKey, user.Role.IsEmployee, is2FAEnabled,
		time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate access token",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return nil, models.ErrGeneratingToken
	}

	// Refresh Token
	refreshExpiresAt := time.Now().AddDate(0, 0, s.JwtRefreshExpiration)
	idToken, familyId, refreshToken, err := utils.GenerateRefreshToken(
		user.Id, s.AppDomain, "", deviceHash, s.JwtRefreshPrivateKey, refreshExpiresAt,
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate refresh token",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return nil, models.ErrGeneratingToken
	}

	err = s.AuthRepo.SaveRefreshToken(ctx, idToken, user.Id, familyId, refreshExpiresAt)
	if err != nil {
		slog.ErrorContext(ctx, "failed to persist refresh token into database",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return nil, models.ErrGeneratingToken
	}

	response := &dto.LoginResponse{
		AccessToken: accessToken,
		Requires2FA: false,
	}

	return &AuthResult{LoginResponse: response, RefreshToken: refreshToken}, nil
}
