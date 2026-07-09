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

type AuthService struct {
	AuthRepo             AuthRepository
	UserRepo             UserRepository
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

func NewAuthService(authRepo AuthRepository, userRepo UserRepository, jwtAccessSecret, jwtRefreshSecret, appName, encryptSecret string, jwtAccessExpiration int, jwtRefreshExpiration int) *AuthService {
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
	user, err := s.UserRepo.GetUserByEmail(ctx, data.Email)
	if err != nil || !user.Enable || user.PasswordHash == nil {
		return nil, errors.New("invalid email or password")
	}

	isValid, err := utils.VerifyPassword(data.Password, *user.PasswordHash)
	if err != nil || !isValid {
		return nil, errors.New("invalid email or password")
	}

	if user.TwoFactorAuthentication {
		_, preAuthToken, err := utils.GenerateToken(user.Id, "", "", nil, s.JwtAccessSecret, s.AppName, false, time.Now().Add(5*time.Minute))
		if err != nil {
			return nil, errors.New("error processing login")
		}

		response := &dto.LoginResponse{
			Requires2FA:  true,
			PreAuthToken: preAuthToken,
		}

		return &AuthResult{response, ""}, nil
	}

	// Minutes
	_, accessToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtAccessSecret, s.AppName, user.Role.IsEmployee, time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute))
	if err != nil {
		return nil, errors.New("Error generating Token.")
	}

	// Days
	idToken, refreshToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtRefreshSecret, s.AppName, user.Role.IsEmployee, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		return nil, errors.New("Error generating Token.")
	}

	err = s.AuthRepo.SaveRefreshToken(ctx, idToken, user.Id, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		return nil, errors.New("Error generating Token.")
	}

	response := &dto.LoginResponse{
		AccessToken: accessToken,
		Requires2FA: false,
	}

	return &AuthResult{response, refreshToken}, nil
}

func (s *AuthService) VerifyLogin2FA(ctx context.Context, data dto.Verify2FARequest) (*AuthResult, error) {
	claims, err := utils.ValidadeToken(data.PreAuthToken, s.JwtAccessSecret, s.AppName)
	if err != nil {
		return nil, errors.New("session expired. please log in again")
	}

	userId := claims.Subject
	if userId == "" {
		return nil, errors.New("invalid session data")
	}

	user, err := s.UserRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.TwoFactorSecret == nil || *user.TwoFactorSecret == "" {
		return nil, errors.New("2FA setup has not been initiated for this user")
	}

	decryptedSecret, err := utils.Decrypt(*user.TwoFactorSecret, s.EncryptSecret)
	if err != nil {
		println(err.Error())
		return nil, errors.New("error processing security")
	}

	isValid := totp.Validate(data.Code, decryptedSecret)
	if !isValid {
		return nil, errors.New("2FA code is invalid or expired")
	}
	_, accessToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtAccessSecret, s.AppName, user.Role.IsEmployee, time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute))
	if err != nil {
		return nil, errors.New("Error generating Token.")
	}

	idToken, refreshToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtRefreshSecret, s.AppName, user.Role.IsEmployee, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		return nil, errors.New("Error generating Token.")
	}

	err = s.AuthRepo.SaveRefreshToken(ctx, idToken, user.Id, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		return nil, errors.New("Error generating Token.")
	}

	response := &dto.LoginResponse{
		AccessToken: accessToken,
		Requires2FA: false,
	}

	return &AuthResult{response, refreshToken}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (string, error) {
	claims, err := utils.ValidadeToken(refreshTokenString, s.JwtRefreshSecret, s.AppName)
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	tokenExists, err := s.AuthRepo.GetRefreshTokenById(ctx, claims.ID)
	if err != nil || tokenExists == nil {
		return "", errors.New("invalid or expired refresh token")
	}

	if tokenExists.UserId != claims.Subject {
		return "", errors.New("token metadata mismatch: security violation")
	}

	user, err := s.UserRepo.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !user.Enable {
		return "", errors.New("user account is disabled")
	}

	_, accessToken, err := utils.GenerateToken(user.Id, user.Username, user.Email, user.Permissions, s.JwtAccessSecret, s.AppName, user.Role.IsEmployee, time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute))

	if err != nil {
		return "", errors.New("error generating new access token")
	}

	return accessToken, nil
}

func (s *AuthService) Logout(ctx context.Context, userId, refreshTokenString string) error {
	claims, err := utils.ValidadeToken(refreshTokenString, s.JwtRefreshSecret, s.AppName)
	if err != nil {
		return errors.New("invalid or expired refresh token")
	}

	return s.AuthRepo.DeleteRefreshTokenById(ctx, userId, claims.ID)
}

func (s *AuthService) LogoutAll(ctx context.Context, userId, refreshTokenString string) error {
	claims, err := utils.ValidadeToken(refreshTokenString, s.JwtRefreshSecret, s.AppName)
	if err != nil {
		return errors.New("invalid or expired refresh token")
	}

	if claims.Subject != userId {
		return errors.New("token mismatch: security violation")
	}

	return s.AuthRepo.DeleteAllRefreshToken(ctx, userId, claims.ID)
}
