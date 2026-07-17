package services

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"time"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"github.com/pquerna/otp/totp"
)

type AuthRepository interface {
	SaveRefreshToken(ctx context.Context, tokenId, userId string, expiresAt time.Time) error
	GetRefreshTokenById(ctx context.Context, id string) (*models.RefreshToken, error)
	DeleteRefreshTokenById(ctx context.Context, userId, tokenId string) error
	DeleteAllRefreshToken(ctx context.Context, userId, tokenId string) error
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
	JwtAccessSecret      string
	JwtRefreshSecret     string
	JwtAccessExpiration  int
	JwtRefreshExpiration int
	AppName              string
	EncryptSecret        string
}

type AuthResult struct {
	LoginResponse *dto.LoginResponse
	RefreshToken  string
}

func NewAuthService(authRepo AuthRepository, userRepo AuthUserRepository, jwtAccessSecret, jwtRefreshSecret, appName, encryptSecret string, jwtAccessExpiration int, jwtRefreshExpiration int) *AuthService {
	return &AuthService{
		AuthRepo:             authRepo,
		UserRepo:             userRepo,
		JwtAccessSecret:      jwtAccessSecret,
		JwtRefreshSecret:     jwtRefreshSecret,
		AppName:              appName,
		JwtAccessExpiration:  jwtAccessExpiration,
		JwtRefreshExpiration: jwtRefreshExpiration,
		EncryptSecret:        encryptSecret,
	}
}

func (s *AuthService) Auth(ctx context.Context, data dto.AuthRequest) (*AuthResult, error) {
	user, err := s.UserRepo.GetUserByEmailForAuth(ctx, data.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrInvalidEmailOrPassword
		}
		log.Printf("[Auth] database query error for %s: %v", utils.MaskEmail(data.Email), err)
		return nil, models.ErrInvalidEmailOrPassword
	}

	if !user.Enable {
		log.Printf("[Auth] login blocked: account for user %s (%s) is disabled", user.Id, utils.MaskEmail(user.Email))
		return nil, models.ErrInvalidEmailOrPassword
	}

	if user.PasswordHash == nil {
		log.Printf("[Auth] login blocked: user %s (%s) does not have a local password set", user.Id, utils.MaskEmail(user.Email))
		return nil, models.ErrInvalidEmailOrPassword
	}

	isValid, err := utils.VerifyPassword(data.Password, *user.PasswordHash)
	if err != nil {
		log.Printf("[Auth] Argon2 verification failed for user %s: %v", user.Id, err)
		return nil, models.ErrPasswordVerificationFailed
	}
	if !isValid {
		return nil, models.ErrInvalidEmailOrPassword
	}

	if user.TwoFactorAuthentication {
		_, preAuthToken, err := utils.GenerateToken(user.Id, "", "", nil, s.JwtAccessSecret, s.AppName, false, time.Now().Add(5*time.Minute))
		if err != nil {
			log.Printf("[Auth] failed to generate 2FA pre-auth token for user %s: %v", user.Id, err)
			return nil, models.ErrUnexpectedLogin
		}

		response := &dto.LoginResponse{
			Requires2FA:  true,
			PreAuthToken: preAuthToken,
		}

		return &AuthResult{response, ""}, nil
	}

	return s.issueSession(ctx, user)
}

func (s *AuthService) VerifyLogin2FA(ctx context.Context, data dto.Verify2FARequest) (*AuthResult, error) {
	claims, err := utils.ValidateToken(data.PreAuthToken, s.JwtAccessSecret, s.AppName)
	if err != nil {
		log.Printf("[VerifyLogin2FA] pre-auth token validation failed: %v", err)
		return nil, models.ErrSessionExpired
	}

	userId := claims.Subject
	if userId == "" {
		log.Printf("[VerifyLogin2FA] pre-auth token claims missing Subject field")
		return nil, models.ErrSessionData
	}

	user, err := s.UserRepo.GetUserByIdForAuth(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}
		log.Printf("[VerifyLogin2FA] database query error for user %s: %v", userId, err)
		return nil, models.ErrUserNotFound
	}

	if user.TwoFactorSecret == nil || *user.TwoFactorSecret == "" {
		log.Printf("[VerifyLogin2FA] user %s attempted 2FA verification but has no secret configured", user.Id)
		return nil, models.Err2FANotInitiated
	}

	decryptedSecret, err := utils.Decrypt(*user.TwoFactorSecret, s.EncryptSecret)
	if err != nil {
		log.Printf("[VerifyLogin2FA] AES decryption of 2FA secret failed for user %s: %v", user.Id, err)
		return nil, models.ErrUnexpectedLogin
	}

	isValid := totp.Validate(data.Code, decryptedSecret)
	if !isValid {
		return nil, models.Err2FAInvalid
	}

	return s.issueSession(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (string, error) {
	claims, err := utils.ValidateToken(refreshTokenString, s.JwtRefreshSecret, s.AppName)
	if err != nil {
		log.Printf("[RefreshToken] refresh token validation failed: %v", err)
		return "", models.ErrInvalidOrExpiredRefresh
	}

	tokenExists, err := s.AuthRepo.GetRefreshTokenById(ctx, claims.ID)
	if err != nil {
		log.Printf("[RefreshToken] error querying session database for token %s: %v", claims.ID, err)
		return "", models.ErrInvalidOrExpiredRefresh
	}
	if tokenExists == nil {
		return "", models.ErrInvalidOrExpiredRefresh
	}

	if tokenExists.UserId != claims.Subject {
		log.Printf("[RefreshToken] Token integrity warning! DB UserId (%s) does not match token Subject (%s)", tokenExists.UserId, claims.Subject)
		return "", models.ErrTokenMetadataMisMatch
	}

	user, err := s.UserRepo.GetUserByEmailForAuth(ctx, claims.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return "", models.ErrUserNotFound
		}
		log.Printf("[RefreshToken] error retrieving user %s during refresh session: %v", claims.Subject, err)
		return "", models.ErrUserNotFound
	}

	if !user.Enable {
		log.Printf("[RefreshToken] user %s is disabled. Blocking token generation", user.Id)
		return "", models.ErrUserNotEnabled
	}

	_, accessToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtAccessSecret, s.AppName, user.Role.IsEmployee, time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute))

	if err != nil {
		log.Printf("[RefreshToken] failed to sign new access token for user %s: %v", user.Id, err)
		return "", models.ErrGeneratingToken
	}

	return accessToken, nil
}

func (s *AuthService) Logout(ctx context.Context, userId, refreshTokenString string) error {
	claims, err := utils.ValidateToken(refreshTokenString, s.JwtRefreshSecret, s.AppName)
	if err != nil {
		return models.ErrInvalidOrExpiredRefresh
	}

	err = s.AuthRepo.DeleteRefreshTokenById(ctx, userId, claims.ID)
	if err != nil {
		log.Printf("[Logout] failed to delete session token %s for user %s: %v", claims.ID, userId, err)
		return models.ErrUnexpectedLogout
	}
	return nil
}

func (s *AuthService) LogoutAll(ctx context.Context, userId, refreshTokenString string) error {
	claims, err := utils.ValidateToken(refreshTokenString, s.JwtRefreshSecret, s.AppName)
	if err != nil {
		return models.ErrInvalidOrExpiredRefresh
	}

	if claims.Subject != userId {
		log.Printf("[LogoutAll] Security mismatch. Subject in token (%s) does not match parameter user (%s)", claims.Subject, userId)
		return models.ErrTokenMetadataMisMatch
	}

	err = s.AuthRepo.DeleteAllRefreshToken(ctx, userId, claims.ID)
	if err != nil {
		log.Printf("[LogoutAll] database error revoking all tokens for user %s: %v", userId, err)
		return models.ErrUnexpectedLogout
	}
	return nil
}

func (s *AuthService) AuthWithSocialProvider(ctx context.Context, provider, socialUserId, email, name string) (*AuthResult, error) {
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
						slog.Error("failed to register social user", "email", utils.MaskEmail(email), "provider", provider, "error", err)
						return nil, models.ErrRegisterSocialUser
					}
				} else {
					slog.Error("database error checking email existence during social auth", "email", utils.MaskEmail(email), "error", err)
					return nil, models.ErrFailQueryUser
				}
			} else {
				slog.Warn("social login block: email already exists with a different provider or password",
					"email", utils.MaskEmail(email), "attempted_provider", provider)
				return nil, models.ErrUserAlreadyExists
			}
		} else {
			slog.Error("database query error fetching user by social ID", "provider", provider, "error", err)
			return nil, models.ErrFailQueryUser
		}
	}

	if !user.Enable {
		slog.Warn("disabled user tried to sign in", "userId", user.Id, "provider", provider)
		return nil, models.ErrUserNotEnabled
	}

	if user.TwoFactorAuthentication {
		_, preAuthToken, err := utils.GenerateToken(user.Id, "", "", nil, s.JwtAccessSecret, s.AppName, false, time.Now().Add(5*time.Minute))
		if err != nil {
			slog.Error("failed to generate pre-auth token", "userId", user.Id, "error", err)
			return nil, models.ErrUnexpectedLogin
		}

		response := &dto.LoginResponse{
			Requires2FA:  true,
			PreAuthToken: preAuthToken,
		}

		return &AuthResult{response, ""}, nil
	}

	_, accessToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtAccessSecret, s.AppName, user.Role.IsEmployee, time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute))
	if err != nil {
		slog.Error("failed to generate access token during social auth", "userId", user.Id, "error", err)
		return nil, models.ErrGeneratingToken
	}

	idToken, refreshToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtRefreshSecret, s.AppName, user.Role.IsEmployee, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		slog.Error("failed to generate refresh token during social auth", "userId", user.Id, "error", err)
		return nil, models.ErrGeneratingToken
	}

	err = s.AuthRepo.SaveRefreshToken(ctx, idToken, user.Id, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		slog.Error("failed to persist refresh token to database during social auth", "userId", user.Id, "error", err)
		return nil, models.ErrGeneratingToken
	}

	response := &dto.LoginResponse{
		AccessToken: accessToken,
		Requires2FA: false,
	}

	return &AuthResult{response, refreshToken}, nil
}

func (s *AuthService) LinkSocialAccount(ctx context.Context, currentUserId, provider, socialUserId, email string) error {
	existingUser, err := s.UserRepo.GetUserBySocialID(ctx, provider, socialUserId)
	if err == nil && existingUser != nil {
		slog.Warn("social ID already linked to another account", "socialUserId", socialUserId, "provider", provider)
		return models.ErrSocialAccountAlreadyLinked
	}

	if err != nil && !errors.Is(err, models.ErrUserNotFound) {
		slog.Error("database error checking social ID existence during linking", "provider", provider, "error", err)
		return models.ErrFailQueryUser
	}

	currentUser, err := s.UserRepo.GetUserByIdForAuth(ctx, currentUserId)
	if err != nil {
		slog.Error("failed to find current user for linking", "userId", currentUserId, "error", err)
		return models.ErrUserNotFound
	}

	if currentUser.Email != email {
		slog.Warn("email mismatch during social link attempt",
			"userId", currentUserId,
			"accountEmail", utils.MaskEmail(currentUser.Email),
			"socialEmail", utils.MaskEmail(email),
		)
		return models.ErrEmailMismatchForSocialLink
	}

	err = s.UserRepo.CreateSocialLink(ctx, currentUser.Id, provider, socialUserId)
	if err != nil {
		slog.Error("failed to create social link for logged user", "userId", currentUser.Id, "provider", provider, "error", err)
		return models.ErrLinkSocialAccount
	}

	return nil
}

func (s *AuthService) issueSession(ctx context.Context, user *models.User) (*AuthResult, error) {
	// Access Token
	_, accessToken, err := utils.GenerateToken(
		user.Id, user.Username, user.Email, user.Permissions,
		s.JwtAccessSecret, s.AppName, user.Role.IsEmployee,
		time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute),
	)
	if err != nil {
		log.Printf("[issueSession] failed to generate access token for user %s: %v", user.Id, err)
		return nil, models.ErrGeneratingToken
	}

	// Refresh Token
	idToken, refreshToken, err := utils.GenerateToken(
		user.Id, user.Username, user.Email, user.Permissions,
		s.JwtRefreshSecret, s.AppName, user.Role.IsEmployee,
		time.Now().AddDate(0, 0, s.JwtRefreshExpiration),
	)
	if err != nil {
		log.Printf("[issueSession] failed to generate refresh token for user %s: %v", user.Id, err)
		return nil, models.ErrGeneratingToken
	}

	err = s.AuthRepo.SaveRefreshToken(ctx, idToken, user.Id, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		log.Printf("[issueSession] failed to save refresh token for user %s into database: %v", user.Id, err)
		return nil, models.ErrGeneratingToken
	}

	response := &dto.LoginResponse{
		AccessToken: accessToken,
		Requires2FA: false,
	}

	return &AuthResult{response, refreshToken}, nil
}
