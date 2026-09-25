package auth_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
)

type testIssuerRepository interface {
	SaveRefreshToken(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error
}

type MockSessionIssuerRepository struct {
	SaveRefreshTokenFunc func(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error
}

func (m *MockSessionIssuerRepository) SaveRefreshToken(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error {
	if m.SaveRefreshTokenFunc != nil {
		return m.SaveRefreshTokenFunc(ctx, tokenId, userId, familyId, expiresAt)
	}
	return nil
}

func NewTestSessionIssuer(authRepo testIssuerRepository, accessKey, refreshKey ed25519.PrivateKey) *auth.SessionIssuer {
	if len(accessKey) == 0 {
		_, accessKey, _ = ed25519.GenerateKey(rand.Reader)
	}
	if len(refreshKey) == 0 {
		_, refreshKey, _ = ed25519.GenerateKey(rand.Reader)
	}

	return auth.NewSessionIssuer(
		authRepo,
		accessKey,
		refreshKey,
		15,
		7,
		"localhost",
	)
}

func TestIssueSession(t *testing.T) {
	_, validAccessKey, _ := ed25519.GenerateKey(rand.Reader)
	_, validRefreshKey, _ := ed25519.GenerateKey(rand.Reader)

	validUser := &models.User{
		Id:                      "user_123",
		Name:                    "John Doe",
		Username:                "johndoe",
		Email:                   "john@example.com",
		Permissions:             []string{"read:users", "write:users"},
		Role:                    models.Role{IsEmployee: true},
		TwoFactorAuthentication: false,
	}

	tests := []struct {
		name                string
		user                *models.User
		deviceHash          string
		hasPassword         bool
		accessKeyOverride   ed25519.PrivateKey
		refreshKeyOverride  ed25519.PrivateKey
		saveRefreshTokenFn  func(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error
		wantErr             error
		wantAccessToken     bool
		wantRefreshToken    bool
		wantSaveTokenCalled bool
	}{
		{
			name:                "Success: Issue session successfully for standard user",
			user:                validUser,
			deviceHash:          "device_hash_xyz",
			hasPassword:         true,
			wantErr:             nil,
			wantAccessToken:     true,
			wantRefreshToken:    true,
			wantSaveTokenCalled: true,
		},
		{
			name: "Success: Issue session for non-employee user with 2FA",
			user: &models.User{
				Id:                      "user_456",
				Name:                    "Jane Doe",
				Username:                "janedoe",
				Email:                   "jane@example.com",
				Permissions:             []string{},
				Role:                    models.Role{IsEmployee: false},
				TwoFactorAuthentication: true,
			},
			deviceHash:          "device_hash_abc",
			hasPassword:         false,
			wantErr:             nil,
			wantAccessToken:     true,
			wantRefreshToken:    true,
			wantSaveTokenCalled: true,
		},
		{
			name:              "Error: Failed to generate access token due to invalid private key",
			user:              validUser,
			deviceHash:        "device_hash_xyz",
			hasPassword:       true,
			accessKeyOverride: ed25519.PrivateKey([]byte("invalid_key")),
			wantErr:           models.ErrGeneratingToken,
		},
		{
			name:               "Error: Failed to generate refresh token due to invalid private key",
			user:               validUser,
			deviceHash:         "device_hash_xyz",
			hasPassword:        true,
			refreshKeyOverride: ed25519.PrivateKey([]byte("invalid_key")),
			wantErr:            models.ErrGeneratingToken,
		},
		{
			name:        "Error: Database fails when persisting refresh token",
			user:        validUser,
			deviceHash:  "device_hash_xyz",
			hasPassword: true,
			saveRefreshTokenFn: func(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error {
				return errors.New("database failure: connection connection timeout")
			},
			wantErr: models.ErrGeneratingToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessKey := validAccessKey
			if len(tt.accessKeyOverride) > 0 {
				accessKey = tt.accessKeyOverride
			}

			refreshKey := validRefreshKey
			if len(tt.refreshKeyOverride) > 0 {
				refreshKey = tt.refreshKeyOverride
			}

			saveCalled := false
			mockRepo := &MockSessionIssuerRepository{
				SaveRefreshTokenFunc: func(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error {
					saveCalled = true

					if userId != tt.user.Id {
						t.Errorf("SaveRefreshToken() userId = %v, want %v", userId, tt.user.Id)
					}
					if tokenId == "" {
						t.Errorf("SaveRefreshToken() tokenId should not be empty")
					}
					if familyId == "" {
						t.Errorf("SaveRefreshToken() familyId should not be empty")
					}
					if expiresAt.Before(time.Now()) {
						t.Errorf("SaveRefreshToken() expiresAt should be in the future")
					}

					if tt.saveRefreshTokenFn != nil {
						return tt.saveRefreshTokenFn(ctx, tokenId, userId, familyId, expiresAt)
					}
					return nil
				},
			}

			issuer := NewTestSessionIssuer(mockRepo, accessKey, refreshKey)

			result, err := issuer.IssueSession(context.Background(), tt.user, tt.deviceHash, tt.hasPassword)

			if tt.wantErr != nil {
				if err == nil || (!errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error()) {
					t.Fatalf("IssueSession() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("IssueSession() unexpected error = %v", err)
			}

			if tt.wantAccessToken && result.LoginResponse.AccessToken == "" {
				t.Errorf("IssueSession() AccessToken was expected, got empty")
			}

			if tt.wantRefreshToken && result.RefreshToken == "" {
				t.Errorf("IssueSession() RefreshToken was expected, got empty")
			}

			if result.LoginResponse.Requires2FA != false {
				t.Errorf("IssueSession() Requires2FA = %v, want false", result.LoginResponse.Requires2FA)
			}

			if saveCalled != tt.wantSaveTokenCalled {
				t.Errorf("SaveRefreshToken() execution called = %v, want %v", saveCalled, tt.wantSaveTokenCalled)
			}
		})
	}
}
