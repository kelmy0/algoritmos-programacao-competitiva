package services

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

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
	return "ABCabc123", nil
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

func TestSignUp(t *testing.T) {
	_, accessPrivKey, _ := ed25519.GenerateKey(rand.Reader)
	_, refreshPrivKey, _ := ed25519.GenerateKey(rand.Reader)

	argonParams := utils.ArgonParams{
		Memory:      64 * 1024,
		Iterations:  1,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	tests := []struct {
		name             string
		req              dto.SignUpRequest
		setupUserRepo    func() *MockUserRepo
		setupAuthRepo    func() *MockAuthRepo
		wantErr          error
		wantAutoLogin    bool
		wantRefreshToken bool
	}{
		{
			name: "Error: Password and Confirm Password dont match",
			req: dto.SignUpRequest{
				Password:        "ValidPassword123!",
				ConfirmPassword: "OtherPassword123!",
			},
			setupUserRepo: func() *MockUserRepo { return &MockUserRepo{} },
			setupAuthRepo: func() *MockAuthRepo { return &MockAuthRepo{} },
			wantErr:       models.ErrPasswordsDontMatch,
		},
		{
			name: "Error: Invalid name",
			req: dto.SignUpRequest{
				Name:            "John",
				Username:        "johnsmith",
				Email:           "john@email.com",
				Password:        "ValidPassword123!",
				ConfirmPassword: "ValidPassword123!",
			},
			setupUserRepo: func() *MockUserRepo { return &MockUserRepo{} },
			setupAuthRepo: func() *MockAuthRepo { return &MockAuthRepo{} },
			wantErr:       models.ErrInvalidRegistrationName,
		},
		{
			name: "Error: Email already used!",
			req: dto.SignUpRequest{
				Name:            "John Smith",
				Username:        "johnsmith",
				Email:           "used@email.com",
				Password:        "ValidPassword123!",
				ConfirmPassword: "ValidPassword123!",
			},
			setupUserRepo: func() *MockUserRepo {
				return &MockUserRepo{
					CheckAvailabilityFunc: func(ctx context.Context, email, username string) (bool, bool, error) {
						return true, false, nil
					},
				}
			},
			setupAuthRepo: func() *MockAuthRepo { return &MockAuthRepo{} },
			wantErr:       models.ErrEmailAlreadyUsed,
		},
		{
			name: "Error: Username already used!",
			req: dto.SignUpRequest{
				Name:            "John Smith",
				Username:        "johnsmith",
				Email:           "john@email.com",
				Password:        "ValidPassword123!",
				ConfirmPassword: "ValidPassword123!",
			},
			setupUserRepo: func() *MockUserRepo {
				return &MockUserRepo{
					CheckAvailabilityFunc: func(ctx context.Context, email, username string) (bool, bool, error) {
						return false, true, nil
					},
				}
			},
			setupAuthRepo: func() *MockAuthRepo { return &MockAuthRepo{} },
			wantErr:       models.ErrUsernameAlreadyUsed,
		},
		{
			name: "Error: Database fails during CheckAvailability",
			req: dto.SignUpRequest{
				Name:            "John Smith",
				Username:        "johnsmith",
				Email:           "john@email.com",
				Password:        "ValidPassword123!",
				ConfirmPassword: "ValidPassword123!",
			},
			setupUserRepo: func() *MockUserRepo {
				return &MockUserRepo{
					CheckAvailabilityFunc: func(ctx context.Context, email, username string) (bool, bool, error) {
						return false, false, errors.New("sql: availability error")
					},
				}
			},
			setupAuthRepo: func() *MockAuthRepo { return &MockAuthRepo{} },
			wantErr:       models.ErrFailQueryUser,
		},
		{
			name: "Error: Database fails unexpectedly during CreateUser",
			req: dto.SignUpRequest{
				Name:            "John Smith",
				Username:        "johnsmith",
				Email:           "john@email.com",
				Password:        "ValidPassword123!",
				ConfirmPassword: "ValidPassword123!",
			},
			setupUserRepo: func() *MockUserRepo {
				return &MockUserRepo{
					CreateUserFunc: func(ctx context.Context, data models.NewUser) (string, error) {
						return "", errors.New("pq: deadlock detected")
					},
				}
			},
			setupAuthRepo: func() *MockAuthRepo { return &MockAuthRepo{} },
			wantErr:       models.ErrUserRegistrationFailed,
		},
		{
			name: "Error: Database returns known conflict error during CreateUser",
			req: dto.SignUpRequest{
				Name:            "John Smith",
				Username:        "johnsmith",
				Email:           "john@email.com",
				Password:        "ValidPassword123!",
				ConfirmPassword: "ValidPassword123!",
			},
			setupUserRepo: func() *MockUserRepo {
				return &MockUserRepo{
					CreateUserFunc: func(ctx context.Context, data models.NewUser) (string, error) {
						return "", models.ErrEmailAlreadyUsed
					},
				}
			},
			setupAuthRepo: func() *MockAuthRepo { return &MockAuthRepo{} },
			wantErr:       models.ErrEmailAlreadyUsed,
		},
		{
			name: "Success: User registered but auto-login fails due to SaveRefreshToken error",
			req: dto.SignUpRequest{
				Name:            "John Smith",
				Username:        "johnsmith",
				Email:           "john@email.com",
				Password:        "ValidPassword123!",
				ConfirmPassword: "ValidPassword123!",
				DeviceHash:      "device-hash-123",
			},
			setupUserRepo: func() *MockUserRepo { return &MockUserRepo{} },
			setupAuthRepo: func() *MockAuthRepo {
				return &MockAuthRepo{
					SaveRefreshTokenFunc: func(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error {
						return errors.New("sql: connection refused")
					},
				}
			},
			wantErr:          nil,
			wantAutoLogin:    false,
			wantRefreshToken: false,
		},
		{
			name: "Success: User created with autologin!",
			req: dto.SignUpRequest{
				Name:            "John Smith",
				Username:        "johnsmith",
				Email:           "john@email.com",
				Password:        "ValidPassword123!",
				ConfirmPassword: "ValidPassword123!",
				DeviceHash:      "device-hash-123",
			},
			setupUserRepo:    func() *MockUserRepo { return &MockUserRepo{} },
			setupAuthRepo:    func() *MockAuthRepo { return &MockAuthRepo{} },
			wantErr:          nil,
			wantAutoLogin:    true,
			wantRefreshToken: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := tt.setupUserRepo()
			authRepo := tt.setupAuthRepo()

			service := NewSignUpService(
				userRepo,
				authRepo,
				argonParams,
				accessPrivKey,
				refreshPrivKey,
				"localhost",
				15,
				7,
			)

			result, err := service.SignUp(context.Background(), tt.req)

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Errorf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
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
