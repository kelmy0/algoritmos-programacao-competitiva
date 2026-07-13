package services

import (
	"context"
	"errors"
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
	if err != nil || !user.Enable || user.PasswordHash == nil {
		return nil, models.ErrInvalidEmailOrPassword
	}

	isValid, err := utils.VerifyPassword(data.Password, *user.PasswordHash)
	if err != nil || !isValid {
		return nil, models.ErrInvalidEmailOrPassword
	}

	if user.TwoFactorAuthentication {
		_, preAuthToken, err := utils.GenerateToken(user.Id, "", "", nil, s.JwtAccessSecret, s.AppName, false, time.Now().Add(5*time.Minute))
		if err != nil {
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
		return nil, models.ErrSessionExpired
	}

	userId := claims.Subject
	if userId == "" {
		return nil, models.ErrSessionData
	}

	user, err := s.UserRepo.GetUserByIdForAuth(ctx, userId)
	if err != nil {
		return nil, models.ErrUserNotFound
	}

	if user.TwoFactorSecret == nil || *user.TwoFactorSecret == "" {
		return nil, models.Err2FANotInitiated
	}

	decryptedSecret, err := utils.Decrypt(*user.TwoFactorSecret, s.EncryptSecret)
	if err != nil {
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
		return "", models.ErrInvalidOrExpiredRefresh
	}

	tokenExists, err := s.AuthRepo.GetRefreshTokenById(ctx, claims.ID)
	if err != nil || tokenExists == nil {
		return "", models.ErrInvalidOrExpiredRefresh
	}

	if tokenExists.UserId != claims.Subject {
		return "", models.ErrTokenMetadataMisMatch
	}

	user, err := s.UserRepo.GetUserByEmailForAuth(ctx, claims.Email)
	if err != nil {
		return "", models.ErrUserNotFound
	}

	if !user.Enable {
		return "", models.ErrUserNotEnabled
	}

	_, accessToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtAccessSecret, s.AppName, user.Role.IsEmployee, time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute))

	if err != nil {
		return "", models.ErrGeneratingToken
	}

	return accessToken, nil
}

func (s *AuthService) Logout(ctx context.Context, userId, refreshTokenString string) error {
	claims, err := utils.ValidateToken(refreshTokenString, s.JwtRefreshSecret, s.AppName)
	if err != nil {
		return models.ErrInvalidOrExpiredRefresh
	}

	return s.AuthRepo.DeleteRefreshTokenById(ctx, userId, claims.ID)
}

func (s *AuthService) LogoutAll(ctx context.Context, userId, refreshTokenString string) error {
	claims, err := utils.ValidateToken(refreshTokenString, s.JwtRefreshSecret, s.AppName)
	if err != nil {
		return models.ErrInvalidOrExpiredRefresh
	}

	if claims.Subject != userId {
		return models.ErrTokenMetadataMisMatch
	}

	return s.AuthRepo.DeleteAllRefreshToken(ctx, userId, claims.ID)
}

func (s *AuthService) AuthWithGoogle(ctx context.Context, provider, socialUserId, email, name string) (*AuthResult, error) {
	user, err := s.UserRepo.GetUserBySocialID(ctx, provider, socialUserId)

	if err != nil {

		if errors.Is(err, models.ErrUserNotFound) {
			user, err = s.UserRepo.GetUserByEmailForAuth(ctx, email)

			if err != nil {
				username := utils.NormalizeUsername(name)

				newUser := models.NewUserGoogle{
					Name:         name,
					Username:     username,
					Email:        email,
					Provider:     "google",
					SocialUserId: socialUserId,
				}

				user, err = s.UserRepo.CreateSocialUser(ctx, newUser, provider, socialUserId)
				if err != nil {
					return nil, models.ErrRegisterSocialUser
				}
			} else {
				err = s.UserRepo.CreateSocialLink(ctx, user.Id, provider, socialUserId)
				if err != nil {
					return nil, models.ErrLinkGoogleAccount
				}

				user, err = s.UserRepo.GetUserByIdForAuth(ctx, user.Id)
				if err != nil {
					return nil, models.ErrReloadUser
				}
			}
		} else {
			return nil, models.ErrFailQueryUser
		}
	}
	if !user.Enable {
		return nil, models.ErrUserNotEnabled
	}

	if user.TwoFactorAuthentication {
		_, preAuthToken, err := utils.GenerateToken(user.Id, "", "", nil, s.JwtAccessSecret, s.AppName, false, time.Now().Add(5*time.Minute))
		if err != nil {
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
		return nil, models.ErrGeneratingToken
	}

	idToken, refreshToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtRefreshSecret, s.AppName, user.Role.IsEmployee, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		return nil, models.ErrGeneratingToken
	}

	err = s.AuthRepo.SaveRefreshToken(ctx, idToken, user.Id, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		return nil, models.ErrGeneratingToken
	}

	response := &dto.LoginResponse{
		AccessToken: accessToken,
		Requires2FA: false,
	}

	return &AuthResult{response, refreshToken}, nil
}

func (s *AuthService) issueSession(ctx context.Context, user *models.User) (*AuthResult, error) {
	// Access Token
	_, accessToken, err := utils.GenerateToken(
		user.Id, user.Username, user.Email, user.Permissions,
		s.JwtAccessSecret, s.AppName, user.Role.IsEmployee,
		time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute),
	)
	if err != nil {
		return nil, models.ErrGeneratingToken
	}

	// Refresh Token
	idToken, refreshToken, err := utils.GenerateToken(
		user.Id, user.Username, user.Email, user.Permissions,
		s.JwtRefreshSecret, s.AppName, user.Role.IsEmployee,
		time.Now().AddDate(0, 0, s.JwtRefreshExpiration),
	)
	if err != nil {
		return nil, models.ErrGeneratingToken
	}

	err = s.AuthRepo.SaveRefreshToken(ctx, idToken, user.Id, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		return nil, models.ErrGeneratingToken
	}

	response := &dto.LoginResponse{
		AccessToken: accessToken,
		Requires2FA: false,
	}

	return &AuthResult{response, refreshToken}, nil
}
