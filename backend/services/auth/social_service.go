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

type socialUserRepository interface {
	GetUserByEmailForAuth(ctx context.Context, email string) (*models.User, error)
	GetUserBySocialID(ctx context.Context, provider, socialId string) (*models.User, error)
	CreateSocialUser(ctx context.Context, newUser models.NewUserGoogle, provider, socialId string) (*models.User, error)
}

type socialSessionIssuer interface {
	IssueSession(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (AuthResult, error)
}

type SocialService struct {
	userRepo            socialUserRepository
	sessionIssuer       socialSessionIssuer
	jwtAccessPrivateKey ed25519.PrivateKey
	appDomain           string
}

func NewSocialService(userRepo socialUserRepository, sessionIssuer socialSessionIssuer,
	jwtAccessPrivateKey ed25519.PrivateKey, appDomain string) *SocialService {
	return &SocialService{
		userRepo:            userRepo,
		sessionIssuer:       sessionIssuer,
		jwtAccessPrivateKey: jwtAccessPrivateKey,
		appDomain:           appDomain,
	}
}

func (s *SocialService) AuthWithSocialProvider(ctx context.Context, provider, socialUserId, email, name, deviceHash string) (AuthResult, error) {
	maskedEmail := utils.MaskEmail(email)
	user, err := s.userRepo.GetUserBySocialID(ctx, provider, socialUserId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			user, err = s.userRepo.GetUserByEmailForAuth(ctx, email)

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

					user, err = s.userRepo.CreateSocialUser(ctx, newUser, provider, socialUserId)
					if err != nil {
						switch {
						case errors.Is(err, models.ErrEmailAlreadyUsed),
							errors.Is(err, models.ErrUsernameAlreadyUsed),
							errors.Is(err, models.ErrUserAlreadyExists):
							return AuthResult{}, err

						default:
							slog.ErrorContext(ctx, "failed to register social user",
								slog.String("email", maskedEmail),
								slog.String("provider", provider),
								slog.Any("error", err),
							)
							return AuthResult{}, models.ErrRegisterSocialUser
						}
					}
				} else {
					slog.ErrorContext(ctx, "database error checking email existence during social auth",
						slog.String("email", maskedEmail),
						slog.Any("error", err),
					)
					return AuthResult{}, models.ErrFailQueryUser
				}
			} else {
				slog.WarnContext(ctx, "social login blocked: email already exists with a different auth strategy",
					slog.String("email", maskedEmail),
					slog.String("attempted_provider", provider),
				)
				return AuthResult{}, models.ErrUserAlreadyExists
			}
		} else {
			slog.ErrorContext(ctx, "database error fetching user by social ID",
				slog.String("provider", provider),
				slog.Any("error", err),
			)
			return AuthResult{}, models.ErrFailQueryUser
		}
	}

	if !user.Enable {
		slog.WarnContext(ctx, "disabled user attempted social login",
			slog.String("user_id", user.Id),
			slog.String("provider", provider),
		)
		return AuthResult{}, models.ErrUserNotEnabled
	}

	if user.TwoFactorAuthentication {
		_, preAuthToken, err := utils.GeneratePreAuthToken(
			user.Id, s.appDomain, deviceHash, s.jwtAccessPrivateKey, time.Now().Add(5*time.Minute),
		)
		if err != nil {
			slog.ErrorContext(ctx, "failed to generate 2FA pre-auth token during social auth",
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

	hasPassword := user.PasswordHash != nil

	return s.sessionIssuer.IssueSession(ctx, user, deviceHash, hasPassword)
}

/*
func (s *SocialService) LinkSocialAccount(ctx context.Context, currentUserId, provider, socialUserId, email string) error {
	existingUser, err := s.userRepo.GetUserBySocialID(ctx, provider, socialUserId)
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

	currentUser, err := s.userRepo.GetUserByIdForAuth(ctx, currentUserId)
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

	err = s.userRepo.CreateSocialLink(ctx, currentUser.Id, provider, socialUserId)
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
*/
