package auth

import (
	"context"
	"crypto/ed25519"
	"log/slog"
	"time"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type sessionIssuerRepository interface {
	SaveRefreshToken(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error
}

type SessionIssuer struct {
	authRepo             sessionIssuerRepository
	appDomain            string
	jwtAccessPrivateKey  ed25519.PrivateKey
	jwtRefreshPrivateKey ed25519.PrivateKey
	jwtAccessExpiration  int
	jwtRefreshExpiration int
}

func NewSessionIssuer(authRepo sessionIssuerRepository,
	jwtAccessPrivateKey, jwtRefreshPrivateKey ed25519.PrivateKey,
	jwtAccessExpiration, jwtRefreshExpiration int, appDomain string) *SessionIssuer {
	return &SessionIssuer{
		authRepo:             authRepo,
		jwtAccessPrivateKey:  jwtAccessPrivateKey,
		jwtRefreshPrivateKey: jwtRefreshPrivateKey,
		jwtAccessExpiration:  jwtAccessExpiration,
		jwtRefreshExpiration: jwtRefreshExpiration,
		appDomain:            appDomain,
	}
}

func (s *SessionIssuer) IssueSession(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (AuthResult, error) {
	_, accessToken, err := utils.GenerateAccessToken(
		user.Id, user.Name, user.Username, user.Email, s.appDomain, user.Permissions,
		s.jwtAccessPrivateKey, user.Role.IsEmployee, user.TwoFactorAuthentication,
		hasPassword, time.Now().Add(time.Duration(s.jwtAccessExpiration)*time.Minute),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate access token",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return AuthResult{}, models.ErrGeneratingToken
	}

	refreshExpiresAt := time.Now().AddDate(0, 0, s.jwtRefreshExpiration)
	idToken, familyId, refreshToken, err := utils.GenerateRefreshToken(
		user.Id, s.appDomain, "", deviceHash, s.jwtRefreshPrivateKey, refreshExpiresAt,
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate refresh token",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return AuthResult{}, models.ErrGeneratingToken
	}

	err = s.authRepo.SaveRefreshToken(ctx, idToken, user.Id, familyId, refreshExpiresAt)
	if err != nil {
		slog.ErrorContext(ctx, "failed to persist refresh token into database",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return AuthResult{}, models.ErrGeneratingToken
	}

	response := dto.LoginResponse{
		AccessToken: accessToken,
		Requires2FA: false,
	}

	return AuthResult{LoginResponse: response, RefreshToken: refreshToken}, nil
}
