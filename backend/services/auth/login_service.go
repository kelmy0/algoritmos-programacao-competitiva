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
)

type loginUserRepository interface {
	GetUserByEmailForAuth(ctx context.Context, email string) (*models.User, error)
}

type loginSessionIssuer interface {
	IssueSession(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (AuthResult, error)
}

type LoginService struct {
	userRepo            loginUserRepository
	sessionIssuer       loginSessionIssuer
	jwtAccessPrivateKey ed25519.PrivateKey
	appDomain           string
}

func NewLoginService(userRepo loginUserRepository, sessionIssuer loginSessionIssuer,
	jwtAccessPrivateKey ed25519.PrivateKey, appDomain string) *LoginService {
	return &LoginService{
		userRepo:            userRepo,
		sessionIssuer:       sessionIssuer,
		jwtAccessPrivateKey: jwtAccessPrivateKey,
		appDomain:           appDomain,
	}
}

func (s *LoginService) Login(ctx context.Context, data dto.AuthRequest) (AuthResult, error) {
	maskedEmail := utils.MaskEmail(data.Email)
	user, err := s.userRepo.GetUserByEmailForAuth(ctx, data.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return AuthResult{}, models.ErrInvalidEmailOrPassword
		}
		slog.ErrorContext(ctx, "database query error fetching user by email during auth",
			slog.String("email", maskedEmail),
			slog.Any("error", err),
		)
		return AuthResult{}, models.ErrInvalidEmailOrPassword
	}

	if !user.Enable {
		slog.WarnContext(ctx, "login blocked: account is disabled",
			slog.String("user_id", user.Id),
			slog.String("email", maskedEmail),
		)
		return AuthResult{}, models.ErrInvalidEmailOrPassword
	}

	if user.PasswordHash == nil {
		slog.WarnContext(ctx, "login blocked: user has no local password configured",
			slog.String("user_id", user.Id),
			slog.String("email", maskedEmail),
		)
		return AuthResult{}, models.ErrInvalidEmailOrPassword
	}

	isValid, err := utils.VerifyPassword(data.Password, *user.PasswordHash)
	if err != nil {
		slog.ErrorContext(ctx, "Argon2 password verification failed",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return AuthResult{}, models.ErrPasswordVerificationFailed
	}
	if !isValid {
		return AuthResult{}, models.ErrInvalidEmailOrPassword
	}

	if user.TwoFactorAuthentication {
		_, preAuthToken, err := utils.GeneratePreAuthToken(user.Id, s.appDomain, data.DeviceHash, s.jwtAccessPrivateKey, time.Now().Add(5*time.Minute))
		if err != nil {
			slog.ErrorContext(ctx, "failed to generate 2FA pre-auth token",
				slog.String("user_id", user.Id),
				slog.Any("error", err),
			)
			return AuthResult{}, models.ErrUnexpectedLogin
		}

		response := dto.LoginResponse{
			Requires2FA:  true,
			PreAuthToken: preAuthToken,
		}

		return AuthResult{LoginResponse: response, RefreshToken: ""}, nil
	}

	return s.sessionIssuer.IssueSession(ctx, user, data.DeviceHash, true)
}
