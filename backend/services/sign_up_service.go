package services

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type SignUpUserRepository interface {
	CheckUserExists(ctx context.Context, email string) (bool, error)
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

func NewSignUpService(userRepo SignUpUserRepository, authRepo SignUpAuthRepository, parallelism uint8, memory, iterarions, saltLength, keyLength uint32, jwtAccessSecret, jwtRefreshSecret, appName string, jwtAccessExpiration int, jwtRefreshExpiration int) *SignUpService {
	return &SignUpService{
		UserRepo: userRepo,
		AuthRepo: authRepo,
		ArgonParams: utils.ArgonParams{
			Memory:      memory,
			Iterations:  iterarions,
			Parallelism: parallelism,
			SaltLength:  saltLength,
			KeyLength:   keyLength,
		},
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

	sanitizedData := dto.SignUpRequest{
		Name:     utils.SanitizeHumanName(data.Name),
		Username: utils.SanitizeHumanName(data.Username),
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Password: data.Password,
	}

	if sanitizedData.Name == "" || sanitizedData.Username == "" || utf8.RuneCountInString(sanitizedData.Name) < 6 || utf8.RuneCountInString(sanitizedData.Username) < 6 {
		return nil, models.ErrInvalidRegistrationFields
	}

	_, err := mail.ParseAddress(sanitizedData.Email)
	if err != nil || !strings.Contains(sanitizedData.Email, "@") || strings.LastIndex(sanitizedData.Email, ".") < strings.LastIndex(sanitizedData.Email, "@") {
		return nil, models.ErrInvalidEmailFormat
	}

	userExists, err := s.UserRepo.CheckUserExists(ctx, sanitizedData.Email)
	if err != nil {
		return nil, models.ErrFailQueryUser
	}

	if userExists {
		return nil, models.ErrUserAlreadyExists
	}

	passwordHash, err := utils.HashPassword(sanitizedData.Password, s.ArgonParams)
	if err != nil {
		return nil, models.ErrUserRegistrationFailed
	}

	dataUser := models.NewUser{
		Name:         sanitizedData.Name,
		Username:     sanitizedData.Username,
		Email:        sanitizedData.Email,
		PasswordHash: passwordHash,
	}

	userId, err := s.UserRepo.CreateUser(ctx, dataUser)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			return nil, models.ErrUserAlreadyExists
		}
		return nil, models.ErrUserRegistrationFailed
	}

	_, accessToken, err := utils.GenerateToken(userId, sanitizedData.Username, sanitizedData.Email, []string{}, s.JwtAccessSecret, s.AppName, false, time.Now().Add(time.Duration(s.JwtAccessExpiration)*time.Minute))
	if err != nil {
		return nil, models.ErrAccountCreatedButTokenFailed
	}

	idToken, refreshToken, err := utils.GenerateToken(userId, sanitizedData.Username, sanitizedData.Email, []string{}, s.JwtRefreshSecret, s.AppName, false, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		return nil, models.ErrAccountCreatedButTokenFailed
	}

	err = s.AuthRepo.SaveRefreshToken(ctx, idToken, userId, time.Now().AddDate(0, 0, s.JwtRefreshExpiration))
	if err != nil {
		return nil, models.ErrAccountCreatedButTokenFailed
	}

	response := &dto.SignUpResponse{
		AccessToken: accessToken,
		Success:     true,
		AutoLogin:   true,
	}

	return &SignUpResult{response, refreshToken}, nil
}
