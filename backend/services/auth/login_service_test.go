package auth_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"testing"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type testUserRepo interface {
	GetUserByEmailForAuth(ctx context.Context, email string) (*models.User, error)
}

type testSessionIssuer interface {
	IssueSession(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error)
}

type MockLoginUserRepo struct {
	GetUserByEmailForAuthFunc func(ctx context.Context, email string) (*models.User, error)
}

func (m *MockLoginUserRepo) GetUserByEmailForAuth(ctx context.Context, email string) (*models.User, error) {
	if m.GetUserByEmailForAuthFunc != nil {
		return m.GetUserByEmailForAuthFunc(ctx, email)
	}

	hashedPassword, _ := utils.HashPassword("ValidPassword123!", utils.ArgonParams{
		Memory: 8 * 1024, Iterations: 1, Parallelism: 1, SaltLength: 16, KeyLength: 32,
	})
	return &models.User{
		Id:                      "user_id_123",
		Email:                   email,
		Enable:                  true,
		PasswordHash:            &hashedPassword,
		TwoFactorAuthentication: false,
	}, nil
}

type MockLoginSessionIssuer struct {
	IssueSessionFunc func(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error)
}

func (m *MockLoginSessionIssuer) IssueSession(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error) {
	if m.IssueSessionFunc != nil {
		return m.IssueSessionFunc(ctx, user, deviceHash, hasPassword)
	}

	return auth.AuthResult{
		LoginResponse: dto.LoginResponse{AccessToken: "mocked_access_token"},
		RefreshToken:  "mocked_refresh_token",
	}, nil
}

func newTestLoginService(userRepo testUserRepo, sessionIssuer testSessionIssuer) *auth.LoginService {
	_, accessPrivKey, _ := ed25519.GenerateKey(rand.Reader)

	return auth.NewLoginService(
		userRepo,
		sessionIssuer,
		accessPrivKey,
		"localhost",
	)
}

func TestLogin(t *testing.T) {
	fastArgonParams := utils.ArgonParams{
		Memory: 8 * 1024, Iterations: 1, Parallelism: 1, SaltLength: 16, KeyLength: 32,
	}
	validHashedPassword, _ := utils.HashPassword("ValidPassword123!", fastArgonParams)

	validReq := dto.AuthRequest{
		Email:      "john@email.com",
		Password:   "ValidPassword123!",
		DeviceHash: "device-hash-123",
	}

	tests := []struct {
		name             string
		req              dto.AuthRequest
		getUserFn        func(ctx context.Context, email string) (*models.User, error)
		issueSessionFn   func(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error)
		wantErr          error
		want2FA          bool
		wantPreAuthToken bool
	}{
		{
			name:    "Success: User authenticated without 2FA",
			req:     validReq,
			wantErr: nil,
		},
		{
			name: "Success: User requires 2FA",
			req:  validReq,
			getUserFn: func(ctx context.Context, email string) (*models.User, error) {
				return &models.User{
					Id:                      "user_id_123",
					Email:                   email,
					Enable:                  true,
					PasswordHash:            &validHashedPassword,
					TwoFactorAuthentication: true,
				}, nil
			},
			wantErr:          nil,
			want2FA:          true,
			wantPreAuthToken: true,
		},
		{
			name: "Error: User not found",
			req:  validReq,
			getUserFn: func(ctx context.Context, email string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			wantErr: models.ErrInvalidEmailOrPassword,
		},
		{
			name: "Error: Database unexpected failure during GetUser",
			req:  validReq,
			getUserFn: func(ctx context.Context, email string) (*models.User, error) {
				return nil, errors.New("sql: connection refused")
			},
			wantErr: models.ErrInvalidEmailOrPassword,
		},
		{
			name: "Error: Account is disabled",
			req:  validReq,
			getUserFn: func(ctx context.Context, email string) (*models.User, error) {
				return &models.User{
					Id:           "user_id_123",
					Enable:       false,
					PasswordHash: &validHashedPassword,
				}, nil
			},
			wantErr: models.ErrInvalidEmailOrPassword,
		},
		{
			name: "Error: User has no local password configured (OAuth only user)",
			req:  validReq,
			getUserFn: func(ctx context.Context, email string) (*models.User, error) {
				return &models.User{
					Id:           "user_id_123",
					Enable:       true,
					PasswordHash: nil,
				}, nil
			},
			wantErr: models.ErrInvalidEmailOrPassword,
		},
		{
			name: "Error: Incorrect password provided",
			req: dto.AuthRequest{
				Email:      "john@email.com",
				Password:   "WrongPassword123!",
				DeviceHash: "device-hash-123",
			},
			wantErr: models.ErrInvalidEmailOrPassword,
		},
		{
			name: "Error: Session issuance fails",
			req:  validReq,
			issueSessionFn: func(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error) {
				return auth.AuthResult{}, errors.New("failed to issue refresh token")
			},
			wantErr: errors.New("failed to issue refresh token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &MockLoginUserRepo{
				GetUserByEmailForAuthFunc: tt.getUserFn,
			}
			sessionIssuer := &MockLoginSessionIssuer{
				IssueSessionFunc: tt.issueSessionFn,
			}

			service := newTestLoginService(userRepo, sessionIssuer)

			result, err := service.Login(context.Background(), tt.req)

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Fatalf("Login() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Login() unexpected error = %v", err)
			}

			if result.LoginResponse.Requires2FA != tt.want2FA {
				t.Errorf("Login() Requires2FA = %v, want2FA %v", result.LoginResponse.Requires2FA, tt.want2FA)
			}

			hasPreAuthToken := result.LoginResponse.PreAuthToken != ""
			if hasPreAuthToken != tt.wantPreAuthToken {
				t.Errorf("Login() PreAuthToken presence = %v, wantPreAuthToken %v", hasPreAuthToken, tt.wantPreAuthToken)
			}
		})
	}
}
