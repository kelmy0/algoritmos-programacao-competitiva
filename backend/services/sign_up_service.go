package services

import (
	"context"
	"crypto/ed25519"
	"errors"
	"log/slog"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type SignUpUserRepository interface {
	CheckAvailability(ctx context.Context, email, username string) (emailTaken bool, usernameTaken bool, err error)
	CreateUser(ctx context.Context, data models.NewUser) (string, error)
}

type SignUpAuthRepository interface {
	SaveRefreshToken(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error
}

type SignUpService struct {
	UserRepo             SignUpUserRepository
	AuthRepo             SignUpAuthRepository
	ArgonParams          utils.ArgonParams
	JwtAccessPrivateKey  ed25519.PrivateKey
	JwtRefreshPrivateKey ed25519.PrivateKey
	JwtAccessExpiration  int
	JwtRefreshExpiration int
	AppDomain            string
}

type SignUpResult struct {
	SignUpResponse *dto.SignUpResponse
	RefreshToken   string
}

func NewSignUpService(userRepo SignUpUserRepository, authRepo SignUpAuthRepository, argonParams utils.ArgonParams, jwtaccessPrivateKey, jwtRefreshPrivateKey ed25519.PrivateKey, appDomain string, jwtAccessExpiration, jwtRefreshExpiration int) *SignUpService {
	return &SignUpService{
		UserRepo:             userRepo,
		AuthRepo:             authRepo,
		ArgonParams:          argonParams,
		JwtAccessPrivateKey:  jwtaccessPrivateKey,
		JwtRefreshPrivateKey: jwtRefreshPrivateKey,
		JwtAccessExpiration:  jwtAccessExpiration,
		JwtRefreshExpiration: jwtRefreshExpiration,
		AppDomain:            appDomain,
	}
}

func (s *SignUpService) SignUp(ctx context.Context, data dto.SignUpRequest) (*SignUpResult, error) {
	if data.Password != data.ConfirmPassword {
		return nil, models.ErrPasswordsDontMatch
	}

	if !utils.IsPasswordValid(data.Password) {
		return nil, models.ErrPasswordIsNotValid
	}

	sanitizedName := utils.SanitizeHumanName(data.Name)
	sanitizedUsername := utils.SanitizeUsername(data.Username)
	sanitizedEmail := strings.ToLower(strings.TrimSpace(data.Email))

	if sanitizedName == "" || utf8.RuneCountInString(sanitizedName) < 6 {
		return nil, models.ErrInvalidRegistrationName
	}

	if sanitizedUsername == "" || utf8.RuneCountInString(sanitizedUsername) < 6 {
		return nil, models.ErrInvalidRegistrationUsername
	}

	_, err := mail.ParseAddress(sanitizedEmail)
	if err != nil || !strings.Contains(sanitizedEmail, "@") || strings.LastIndex(sanitizedEmail, ".") < strings.LastIndex(sanitizedEmail, "@") {
		return nil, models.ErrInvalidEmailFormat
	}

	emailTaken, usernameTaken, err := s.UserRepo.CheckAvailability(ctx, sanitizedEmail, sanitizedUsername)
	if err != nil {
		slog.ErrorContext(ctx, "failed to verify user availability during sign up", slog.Any("error", err))
		return nil, models.ErrFailQueryUser
	}

	if emailTaken {
		return nil, models.ErrEmailAlreadyUsed
	}
	if usernameTaken {
		return nil, models.ErrUsernameAlreadyUsed
	}

	passwordHash, err := utils.HashPassword(data.Password, s.ArgonParams)
	if err != nil {
		slog.ErrorContext(ctx, "Argon2 hashing failed for new user registration", slog.Any("error", err))
		return nil, models.ErrCryptTokenFailed
	}

	dataUser := models.NewUser{
		Name:         sanitizedName,
		Username:     sanitizedUsername,
		Email:        sanitizedEmail,
		PasswordHash: passwordHash,
	}

	userId, err := s.UserRepo.CreateUser(ctx, dataUser)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEmailAlreadyUsed),
			errors.Is(err, models.ErrUsernameAlreadyUsed),
			errors.Is(err, models.ErrUserAlreadyExists):
			return nil, err

		default:
			slog.ErrorContext(ctx, "failed to register user",
				slog.String("email", utils.MaskEmail(dataUser.Email)),
				slog.Any("error", err),
			)
			return nil, models.ErrUserRegistrationFailed
		}
	}

	_, accessToken, errAccess := utils.GenerateAccessToken(
		userId, sanitizedName, sanitizedUsername, sanitizedEmail, s.AppDomain, []string{},
		s.JwtAccessPrivateKey, false, false, true,
		time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute),
	)

	refreshExpiresAt := time.Now().AddDate(0, 0, s.JwtRefreshExpiration)
	idToken, familyId, refreshToken, errRefresh := utils.GenerateRefreshToken(
		userId, s.AppDomain, "", data.DeviceHash, s.JwtRefreshPrivateKey, refreshExpiresAt,
	)

	var errSave error
	if errRefresh == nil {
		errSave = s.AuthRepo.SaveRefreshToken(ctx, idToken, userId, familyId, refreshExpiresAt)
	}

	if errAccess != nil || errRefresh != nil || errSave != nil {
		slog.WarnContext(ctx, "user registered successfully, but auto-login session generation failed",
			slog.String("user_id", userId),
			slog.Any("access_err", errAccess),
			slog.Any("refresh_err", errRefresh),
			slog.Any("save_err", errSave),
		)

		response := &dto.SignUpResponse{
			AccessToken: "",
			Success:     true,
			AutoLogin:   false,
		}

		return &SignUpResult{SignUpResponse: response, RefreshToken: ""}, nil
	}

	response := &dto.SignUpResponse{
		AccessToken: accessToken,
		Success:     true,
		AutoLogin:   true,
	}

	return &SignUpResult{SignUpResponse: response, RefreshToken: refreshToken}, nil
}
