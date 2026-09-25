package auth_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type testSignUpUserRepository interface {
	CheckAvailability(ctx context.Context, email, username string) (emailTaken bool, usernameTaken bool, err error)
	CreateUser(ctx context.Context, data models.NewUser) (string, error)
}

type testSignUpAuthRepository interface {
	SaveRefreshToken(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error
}

type MockUserRepo struct {
	CheckAvailabilityFunc func(ctx context.Context, email, username string) (bool, bool, error)
	CreateUserFunc        func(ctx context.Context, data models.NewUser) (string, error)
}

func (m *MockUserRepo) CheckAvailability(ctx context.Context, email, username string) (bool, bool, error) {
	if m.CheckAvailabilityFunc != nil {
		return m.CheckAvailabilityFunc(ctx, email, username)
	}
	return false, false, nil
}

func (m *MockUserRepo) CreateUser(ctx context.Context, data models.NewUser) (string, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, data)
	}
	return "user_id_123", nil
}

type MockAuthRepo struct {
	SaveRefreshTokenFunc func(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error
}

func (m *MockAuthRepo) SaveRefreshToken(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error {
	if m.SaveRefreshTokenFunc != nil {
		return m.SaveRefreshTokenFunc(ctx, tokenId, userId, familyId, expiresAt)
	}
	return nil
}

func newTestSignUpService(userRepo testSignUpUserRepository, authRepo testSignUpAuthRepository) *auth.SignUpService {
	_, accessPrivKey, _ := ed25519.GenerateKey(rand.Reader)
	_, refreshPrivKey, _ := ed25519.GenerateKey(rand.Reader)

	fastArgonParams := utils.ArgonParams{
		Memory:      8 * 1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}

	return auth.NewSignUpService(
		userRepo,
		authRepo,
		fastArgonParams,
		accessPrivKey,
		refreshPrivKey,
		"localhost",
		15,
		7,
	)
}

func TestSignUp(t *testing.T) {
	validReq := dto.SignUpRequest{
		Name:            "John Smith",
		Username:        "johnsmith",
		Email:           "john@email.com",
		Password:        "ValidPassword123!",
		ConfirmPassword: "ValidPassword123!",
		DeviceHash:      "device-123",
	}

	tests := []struct {
		name                string
		req                 dto.SignUpRequest
		checkAvailabilityFn func(ctx context.Context, email, username string) (bool, bool, error)
		createUserFn        func(ctx context.Context, data models.NewUser) (string, error)
		saveRefreshTokenFn  func(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error
		wantErr             error
		wantAutoLogin       bool
		wantRefreshToken    bool
	}{
		{
			name:             "Success: User created with auto-login",
			req:              validReq,
			wantErr:          nil,
			wantAutoLogin:    true,
			wantRefreshToken: true,
		},
		{
			name: "Error: Passwords do not match",
			req: dto.SignUpRequest{
				Password:        "ValidPassword123!",
				ConfirmPassword: "OtherPassword123!",
			},
			wantErr: models.ErrPasswordsDontMatch,
		},
		{
			name: "Error: Name too short",
			req: func() dto.SignUpRequest {
				r := validReq
				r.Name = "John"
				return r
			}(),
			wantErr: models.ErrInvalidRegistrationName,
		},
		{
			name: "Error: Email already used",
			req:  validReq,
			checkAvailabilityFn: func(ctx context.Context, email, username string) (bool, bool, error) {
				return true, false, nil
			},
			wantErr: models.ErrEmailAlreadyUsed,
		},
		{
			name: "Error: Username already used",
			req:  validReq,
			checkAvailabilityFn: func(ctx context.Context, email, username string) (bool, bool, error) {
				return false, true, nil
			},
			wantErr: models.ErrUsernameAlreadyUsed,
		},
		{
			name: "Error: Database failure during availability check",
			req:  validReq,
			checkAvailabilityFn: func(ctx context.Context, email, username string) (bool, bool, error) {
				return false, false, errors.New("sql: connection timeout")
			},
			wantErr: models.ErrFailQueryUser,
		},
		{
			name: "Error: Database failure during user creation",
			req:  validReq,
			createUserFn: func(ctx context.Context, data models.NewUser) (string, error) {
				return "", errors.New("pq: deadlock detected")
			},
			wantErr: models.ErrUserRegistrationFailed,
		},
		{
			name: "Success: Registration succeeds but auto-login fails gracefully",
			req:  validReq,
			saveRefreshTokenFn: func(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error {
				return errors.New("redis/db down")
			},
			wantErr:          nil,
			wantAutoLogin:    false,
			wantRefreshToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			userRepo := &MockUserRepo{
				CheckAvailabilityFunc: tt.checkAvailabilityFn,
				CreateUserFunc:        tt.createUserFn,
			}
			authRepo := &MockAuthRepo{
				SaveRefreshTokenFunc: tt.saveRefreshTokenFn,
			}

			service := newTestSignUpService(userRepo, authRepo)

			result, err := service.SignUp(context.Background(), tt.req)

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Fatalf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("SignUp() unexpected error = %v", err)
			}

			if result.SignUpResponse.AutoLogin != tt.wantAutoLogin {
				t.Errorf("SignUp() AutoLogin = %v, wantAutoLogin %v", result.SignUpResponse.AutoLogin, tt.wantAutoLogin)
			}

			hasRefreshToken := result.RefreshToken != ""
			if hasRefreshToken != tt.wantRefreshToken {
				t.Errorf("SignUp() RefreshToken presence = %v, wantRefreshToken %v", hasRefreshToken, tt.wantRefreshToken)
			}
		})
	}
}
