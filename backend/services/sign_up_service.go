package services

import (
	"context"
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
	SaveRefreshToken(ctx context.Context, tokenId, userId string, expiresAt time.Time) error
}

type SignUpService struct {
	UserRepo             SignUpUserRepository
	AuthRepo             SignUpAuthRepository
	ArgonParams          utils.ArgonParams
	JwtAccessSecret      string
	JwtRefreshSecret     string
	JwtAccessExpiration  int
	JwtRefreshExpiration int
	AppName              string
}

type SignUpResult struct {
	SignUpResponse *dto.SignUpResponse
	RefreshToken   string
}

func NewSignUpService(userRepo SignUpUserRepository, authRepo SignUpAuthRepository, argonParams utils.ArgonParams, jwtAccessSecret, jwtRefreshSecret, appName string, jwtAccessExpiration, jwtRefreshExpiration int) *SignUpService {
	return &SignUpService{
		UserRepo:             userRepo,
		AuthRepo:             authRepo,
		ArgonParams:          argonParams,
		JwtAccessSecret:      jwtAccessSecret,
		JwtRefreshSecret:     jwtRefreshSecret,
		JwtAccessExpiration:  jwtAccessExpiration,
		JwtRefreshExpiration: jwtRefreshExpiration,
		AppName:              appName,
	}
}

func (s *SignUpService) SignUp(ctx context.Context, data dto.SignUpRequest) (*SignUpResult, error) {
	if data.Password != data.ConfirmPassword {
		return nil, models.ErrPasswordsDontMatch
	}

	if !utils.IsPasswordValid(data.Password) {
		return nil, models.ErrPasswordIsNotValid
	}

	sanitizedData := dto.SignUpRequest{
		Name:     utils.SanitizeHumanName(data.Name),
		Username: utils.SanitizeUsername(data.Username),
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Password: data.Password,
	}

	if sanitizedData.Name == "" || utf8.RuneCountInString(sanitizedData.Name) < 6 {
		return nil, models.ErrInvalidRegistrationName
	}

	if sanitizedData.Username == "" || utf8.RuneCountInString(sanitizedData.Username) < 6 {
		return nil, models.ErrInvalidRegistrationUsername
	}

	_, err := mail.ParseAddress(sanitizedData.Email)
	if err != nil || !strings.Contains(sanitizedData.Email, "@") || strings.LastIndex(sanitizedData.Email, ".") < strings.LastIndex(sanitizedData.Email, "@") {
		return nil, models.ErrInvalidEmailFormat
	}

	emailTaken, usernameTaken, err := s.UserRepo.CheckAvailability(ctx, sanitizedData.Email, sanitizedData.Username)
	if err != nil {
		slog.Error("[SignUp] failed to verify user availability", "error", err)
		return nil, models.ErrFailQueryUser
	}

	if emailTaken {
		return nil, models.ErrEmailAlreadyUsed
	}
	if usernameTaken {
		return nil, models.ErrUsernameAlreadyUsed
	}

	passwordHash, err := utils.HashPassword(sanitizedData.Password, s.ArgonParams)
	if err != nil {
		slog.Error("[SignUp] Argon2 hashing failed for new user registration", "error", err)
		return nil, models.ErrCryptTokenFailed
	}

	dataUser := models.NewUser{
		Name:         sanitizedData.Name,
		Username:     sanitizedData.Username,
		Email:        sanitizedData.Email,
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
			slog.Error("[SignUp] failed to register user", "email", utils.MaskEmail(dataUser.Email), "error", err)
			return nil, models.ErrUserRegistrationFailed
		}
	}

	_, accessToken, err := utils.GenerateToken(userId, sanitizedData.Username, sanitizedData.Email, []string{}, s.JwtAccessSecret, s.AppName, false, time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute))
	if err != nil {
		slog.Warn("[SignUp] user registered, but failed to sign access token", "userID", userId, "error", err)
		return nil, models.ErrAccountCreatedButTokenFailed
	}

	idToken, refreshToken, err := utils.GenerateToken(userId, sanitizedData.Username, sanitizedData.Email, []string{}, s.JwtRefreshSecret, s.AppName, false, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		slog.Warn("[SignUp] user registered, but failed to sign refresh token", "userID", userId, "error", err)
		return nil, models.ErrAccountCreatedButTokenFailed
	}

	err = s.AuthRepo.SaveRefreshToken(ctx, idToken, userId, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		slog.Warn("[SignUp] user registered, but failed to persist refresh session in database", "userID", userId, "error", err)
		return nil, models.ErrAccountCreatedButTokenFailed
	}

	response := &dto.SignUpResponse{
		AccessToken: accessToken,
		Success:     true,
		AutoLogin:   true,
	}

	return &SignUpResult{response, refreshToken}, nil
}
