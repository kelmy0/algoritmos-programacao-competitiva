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
)

type testSocialUserRepository interface {
	GetUserByEmailForAuth(ctx context.Context, email string) (*models.User, error)
	GetUserBySocialID(ctx context.Context, provider, socialId string) (*models.User, error)
	CreateSocialUser(ctx context.Context, newUser models.NewUserGoogle, provider, socialId string) (*models.User, error)
}

type testSocialSessionIssuer interface {
	IssueSession(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error)
}

type MockSocialUserRepository struct {
	GetUserBySocialIDFunc     func(ctx context.Context, provider, socialId string) (*models.User, error)
	GetUserByEmailForAuthFunc func(ctx context.Context, email string) (*models.User, error)
	CreateSocialUserFunc      func(ctx context.Context, newUser models.NewUserGoogle, provider, socialId string) (*models.User, error)
}

func (m *MockSocialUserRepository) GetUserBySocialID(ctx context.Context, provider, socialId string) (*models.User, error) {
	if m.GetUserBySocialIDFunc != nil {
		return m.GetUserBySocialIDFunc(ctx, provider, socialId)
	}
	return nil, models.ErrUserNotFound
}

func (m *MockSocialUserRepository) GetUserByEmailForAuth(ctx context.Context, email string) (*models.User, error) {
	if m.GetUserByEmailForAuthFunc != nil {
		return m.GetUserByEmailForAuthFunc(ctx, email)
	}
	return nil, models.ErrUserNotFound
}

func (m *MockSocialUserRepository) CreateSocialUser(ctx context.Context, newUser models.NewUserGoogle, provider, socialId string) (*models.User, error) {
	if m.CreateSocialUserFunc != nil {
		return m.CreateSocialUserFunc(ctx, newUser, provider, socialId)
	}
	return &models.User{
		Id:                      "new_social_user_id",
		Email:                   newUser.Email,
		Name:                    newUser.Name,
		Username:                newUser.Username,
		Enable:                  true,
		TwoFactorAuthentication: false,
	}, nil
}

type MockSocialSessionIssuer struct {
	IssueSessionFunc func(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error)
}

func (m *MockSocialSessionIssuer) IssueSession(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error) {
	if m.IssueSessionFunc != nil {
		return m.IssueSessionFunc(ctx, user, deviceHash, hasPassword)
	}
	return auth.AuthResult{
		LoginResponse: dto.LoginResponse{AccessToken: "mocked_access_token"},
		RefreshToken:  "mocked_refresh_token",
	}, nil
}

func newTestSocialService(userRepo testSocialUserRepository, sessionIssuer testSocialSessionIssuer, privKey ed25519.PrivateKey) *auth.SocialService {
	if len(privKey) == 0 {
		_, privKey, _ = ed25519.GenerateKey(rand.Reader)
	}
	return auth.NewSocialService(userRepo, sessionIssuer, privKey, "localhost")
}

func TestAuthWithSocialProvider(t *testing.T) {
	_, validAccessPrivKey, _ := ed25519.GenerateKey(rand.Reader)

	defaultProvider := "google"
	defaultSocialID := "google_123456"
	defaultEmail := "john.doe@example.com"
	defaultName := "John Doe"
	defaultDeviceHash := "device_hash_123"

	validUser := &models.User{
		Id:                      "user_123",
		Email:                   defaultEmail,
		Name:                    defaultName,
		Username:                "johndoe",
		Enable:                  true,
		TwoFactorAuthentication: false,
	}

	tests := []struct {
		name                string
		provider            string
		socialUserId        string
		email               string
		userName            string
		deviceHash          string
		accessKeyOverride   ed25519.PrivateKey
		getBySocialIdFn     func(ctx context.Context, provider, socialId string) (*models.User, error)
		getByEmailFn        func(ctx context.Context, email string) (*models.User, error)
		createSocialUserFn  func(ctx context.Context, newUser models.NewUserGoogle, provider, socialId string) (*models.User, error)
		issueSessionFn      func(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error)
		wantErr             error
		wantRequires2FA     bool
		wantHasPreAuthToken bool
		wantHasAccessToken  bool
		wantHasRefreshToken bool
	}{
		{
			name:         "Success: Existing social user logs in",
			provider:     defaultProvider,
			socialUserId: defaultSocialID,
			email:        defaultEmail,
			userName:     defaultName,
			deviceHash:   defaultDeviceHash,
			getBySocialIdFn: func(ctx context.Context, provider, socialId string) (*models.User, error) {
				return validUser, nil
			},
			wantErr:             nil,
			wantHasAccessToken:  true,
			wantHasRefreshToken: true,
		},
		{
			name:         "Success: New social user registered and session issued",
			provider:     defaultProvider,
			socialUserId: defaultSocialID,
			email:        defaultEmail,
			userName:     defaultName,
			deviceHash:   defaultDeviceHash,
			getBySocialIdFn: func(ctx context.Context, provider, socialId string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			getByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			wantErr:             nil,
			wantHasAccessToken:  true,
			wantHasRefreshToken: true,
		},
		{
			name:         "Success: Social user with 2FA requires pre-auth token",
			provider:     defaultProvider,
			socialUserId: defaultSocialID,
			email:        defaultEmail,
			userName:     defaultName,
			deviceHash:   defaultDeviceHash,
			getBySocialIdFn: func(ctx context.Context, provider, socialId string) (*models.User, error) {
				return &models.User{
					Id:                      "user_2fa",
					Email:                   defaultEmail,
					Enable:                  true,
					TwoFactorAuthentication: true,
				}, nil
			},
			wantErr:             nil,
			wantRequires2FA:     true,
			wantHasPreAuthToken: true,
			wantHasAccessToken:  false,
			wantHasRefreshToken: false,
		},

		{
			name:         "Error: Email already exists with traditional password login (blocked)",
			provider:     defaultProvider,
			socialUserId: defaultSocialID,
			email:        defaultEmail,
			userName:     defaultName,
			deviceHash:   defaultDeviceHash,
			getBySocialIdFn: func(ctx context.Context, provider, socialId string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			getByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return &models.User{Id: "existing_user", Email: email}, nil
			},
			wantErr: models.ErrUserAlreadyExists,
		},
		{
			name:         "Error: Social user is disabled (Enable = false)",
			provider:     defaultProvider,
			socialUserId: defaultSocialID,
			email:        defaultEmail,
			userName:     defaultName,
			deviceHash:   defaultDeviceHash,
			getBySocialIdFn: func(ctx context.Context, provider, socialId string) (*models.User, error) {
				return &models.User{Id: "user_123", Enable: false}, nil
			},
			wantErr: models.ErrUserNotEnabled,
		},
		{
			name:         "Error: Database failure during GetUserBySocialID",
			provider:     defaultProvider,
			socialUserId: defaultSocialID,
			email:        defaultEmail,
			userName:     defaultName,
			deviceHash:   defaultDeviceHash,
			getBySocialIdFn: func(ctx context.Context, provider, socialId string) (*models.User, error) {
				return nil, errors.New("sql: database connection timeout")
			},
			wantErr: models.ErrFailQueryUser,
		},
		{
			name:         "Error: Database failure during GetUserByEmailForAuth",
			provider:     defaultProvider,
			socialUserId: defaultSocialID,
			email:        defaultEmail,
			userName:     defaultName,
			deviceHash:   defaultDeviceHash,
			getBySocialIdFn: func(ctx context.Context, provider, socialId string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			getByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return nil, errors.New("sql: syntax error")
			},
			wantErr: models.ErrFailQueryUser,
		},
		{
			name:         "Error: CreateSocialUser returns ErrEmailAlreadyUsed",
			provider:     defaultProvider,
			socialUserId: defaultSocialID,
			email:        defaultEmail,
			userName:     defaultName,
			deviceHash:   defaultDeviceHash,
			getBySocialIdFn: func(ctx context.Context, provider, socialId string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			getByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			createSocialUserFn: func(ctx context.Context, newUser models.NewUserGoogle, provider, socialId string) (*models.User, error) {
				return nil, models.ErrEmailAlreadyUsed
			},
			wantErr: models.ErrEmailAlreadyUsed,
		},
		{
			name:         "Error: CreateSocialUser returns ErrUsernameAlreadyUsed",
			provider:     defaultProvider,
			socialUserId: defaultSocialID,
			email:        defaultEmail,
			userName:     defaultName,
			deviceHash:   defaultDeviceHash,
			getBySocialIdFn: func(ctx context.Context, provider, socialId string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			getByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			createSocialUserFn: func(ctx context.Context, newUser models.NewUserGoogle, provider, socialId string) (*models.User, error) {
				return nil, models.ErrUsernameAlreadyUsed
			},
			wantErr: models.ErrUsernameAlreadyUsed,
		},
		{
			name:         "Error: CreateSocialUser returns unexpected internal error",
			provider:     defaultProvider,
			socialUserId: defaultSocialID,
			email:        defaultEmail,
			userName:     defaultName,
			deviceHash:   defaultDeviceHash,
			getBySocialIdFn: func(ctx context.Context, provider, socialId string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			getByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			createSocialUserFn: func(ctx context.Context, newUser models.NewUserGoogle, provider, socialId string) (*models.User, error) {
				return nil, errors.New("unhandled insert failure")
			},
			wantErr: models.ErrRegisterSocialUser,
		},

		{
			name:         "Error: Failed to generate 2FA PreAuthToken due to corrupt private key",
			provider:     defaultProvider,
			socialUserId: defaultSocialID,
			email:        defaultEmail,
			userName:     defaultName,
			deviceHash:   defaultDeviceHash,
			getBySocialIdFn: func(ctx context.Context, provider, socialId string) (*models.User, error) {
				return &models.User{
					Id:                      "user_2fa",
					Enable:                  true,
					TwoFactorAuthentication: true,
				}, nil
			},
			accessKeyOverride: ed25519.PrivateKey([]byte("invalid_key")),
			wantErr:           models.ErrUnexpectedLogin,
		},
		{
			name:         "Error: SessionIssuer.IssueSession returns error",
			provider:     defaultProvider,
			socialUserId: defaultSocialID,
			email:        defaultEmail,
			userName:     defaultName,
			deviceHash:   defaultDeviceHash,
			getBySocialIdFn: func(ctx context.Context, provider, socialId string) (*models.User, error) {
				return validUser, nil
			},
			issueSessionFn: func(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error) {
				return auth.AuthResult{}, models.ErrGeneratingToken
			},
			wantErr: models.ErrGeneratingToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockSocialUserRepository{
				GetUserBySocialIDFunc:     tt.getBySocialIdFn,
				GetUserByEmailForAuthFunc: tt.getByEmailFn,
				CreateSocialUserFunc:      tt.createSocialUserFn,
			}

			mockIssuer := &MockSocialSessionIssuer{
				IssueSessionFunc: tt.issueSessionFn,
			}

			keyToUse := validAccessPrivKey
			if len(tt.accessKeyOverride) > 0 {
				keyToUse = tt.accessKeyOverride
			}

			service := newTestSocialService(mockRepo, mockIssuer, keyToUse)

			result, err := service.AuthWithSocialProvider(
				context.Background(),
				tt.provider,
				tt.socialUserId,
				tt.email,
				tt.userName,
				tt.deviceHash,
			)

			if tt.wantErr != nil {
				if err == nil || (!errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error()) {
					t.Fatalf("AuthWithSocialProvider() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("AuthWithSocialProvider() unexpected error = %v", err)
			}

			if result.LoginResponse.Requires2FA != tt.wantRequires2FA {
				t.Errorf("Requires2FA = %v, want %v", result.LoginResponse.Requires2FA, tt.wantRequires2FA)
			}

			hasPreAuthToken := result.LoginResponse.PreAuthToken != ""
			if hasPreAuthToken != tt.wantHasPreAuthToken {
				t.Errorf("PreAuthToken presence = %v, want %v", hasPreAuthToken, tt.wantHasPreAuthToken)
			}

			hasAccessToken := result.LoginResponse.AccessToken != ""
			if hasAccessToken != tt.wantHasAccessToken {
				t.Errorf("AccessToken presence = %v, want %v", hasAccessToken, tt.wantHasAccessToken)
			}

			hasRefreshToken := result.RefreshToken != ""
			if hasRefreshToken != tt.wantHasRefreshToken {
				t.Errorf("RefreshToken presence = %v, want %v", hasRefreshToken, tt.wantHasRefreshToken)
			}
		})
	}
}
