package auth_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
)

type testTwoFactorUserRepository interface {
	GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error)
}

type testTwoFactorSessionIssuer interface {
	IssueSession(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error)
}

type testTwoFactorRedis interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
}

type MockTwoFactorUserRepository struct {
	GetUserByIdForAuthFunc func(ctx context.Context, id string) (*models.User, error)
}

func (m *MockTwoFactorUserRepository) GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error) {
	if m.GetUserByIdForAuthFunc != nil {
		return m.GetUserByIdForAuthFunc(ctx, id)
	}
	return nil, models.ErrUserNotFound
}

type MockTwoFactorSessionIssuer struct {
	IssueSessionFunc func(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error)
}

func (m *MockTwoFactorSessionIssuer) IssueSession(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error) {
	if m.IssueSessionFunc != nil {
		return m.IssueSessionFunc(ctx, user, deviceHash, hasPassword)
	}
	return auth.AuthResult{
		LoginResponse: dto.LoginResponse{AccessToken: "mocked_access_token"},
		RefreshToken:  "mocked_refresh_token",
	}, nil
}

type MockTwoFactorRedis struct {
	SetFunc    func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
	ExistsFunc func(ctx context.Context, keys ...string) *redis.IntCmd
}

func (m *MockTwoFactorRedis) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, value, expiration)
	}
	cmd := redis.NewStatusCmd(ctx)
	cmd.SetVal("OK")
	return cmd
}

func (m *MockTwoFactorRedis) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, keys...)
	}
	cmd := redis.NewIntCmd(ctx)
	cmd.SetVal(0)
	return cmd
}

const testEncryptSecret = "12345678901234567890123456789012"

func newTestTwoFactorService(
	userRepo testTwoFactorUserRepository,
	sessionIssuer testTwoFactorSessionIssuer,
	redisClient testTwoFactorRedis,
	pubKey ed25519.PublicKey,
	privKey ed25519.PrivateKey,
) *auth.TwoFactorService {
	return auth.NewTwoFactorService(
		userRepo,
		sessionIssuer,
		redisClient,
		pubKey,
		privKey,
		"localhost",
		testEncryptSecret,
	)
}

func generateTokenWithoutSubject(privKey ed25519.PrivateKey, domain, deviceHash string) string {
	claims := jwt.MapClaims{
		"iss": domain,
		"aud": domain,
		"dvh": deviceHash,
		"exp": time.Now().Add(5 * time.Minute).Unix(),
		"jti": "mock_jti_without_sub",
		"sub": "",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	tokenStr, _ := token.SignedString(privKey)
	return tokenStr
}

func TestVerifyLogin2FA(t *testing.T) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate ed25519 keys: %v", err)
	}

	defaultUserID := "user_2fa_123"
	defaultDeviceHash := "device_hash_456"

	totpKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "TestApp",
		AccountName: "user@example.com",
	})
	if err != nil {
		t.Fatalf("failed to generate totp key for setup: %v", err)
	}
	rawTotpSecret := totpKey.Secret()

	encryptedSecret, err := utils.Encrypt(rawTotpSecret, testEncryptSecret)
	if err != nil {
		t.Fatalf("failed to encrypt totp secret for setup: %v", err)
	}

	validUser := &models.User{
		Id:                      defaultUserID,
		Email:                   "user@example.com",
		Enable:                  true,
		TwoFactorAuthentication: true,
		TwoFactorSecret:         &encryptedSecret,
	}

	tests := []struct {
		name                string
		getPreAuthToken     func(t *testing.T) string
		reqDeviceHash       string
		reqCode             func() string
		getByUserIdFn       func(ctx context.Context, id string) (*models.User, error)
		redisExistsFn       func(ctx context.Context, keys ...string) *redis.IntCmd
		redisSetFn          func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
		issueSessionFn      func(ctx context.Context, user *models.User, deviceHash string, hasPassword bool) (auth.AuthResult, error)
		wantErr             error
		wantHasAccessToken  bool
		wantHasRefreshToken bool
	}{

		{
			name: "Success: Valid 2FA code completes authentication and blacklists token",
			getPreAuthToken: func(t *testing.T) string {
				_, token, err := utils.GeneratePreAuthToken(defaultUserID, "localhost", defaultDeviceHash, privKey, time.Now().Add(5*time.Minute))
				if err != nil {
					t.Fatalf("failed to generate valid pre-auth token: %v", err)
				}
				return token
			},
			reqDeviceHash: defaultDeviceHash,
			reqCode: func() string {
				code, _ := totp.GenerateCode(rawTotpSecret, time.Now())
				return code
			},
			getByUserIdFn: func(ctx context.Context, id string) (*models.User, error) {
				return validUser, nil
			},
			wantErr:             nil,
			wantHasAccessToken:  true,
			wantHasRefreshToken: true,
		},
		{
			name: "Error: PreAuthToken is invalid or expired",
			getPreAuthToken: func(t *testing.T) string {
				return "invalid.jwt.token"
			},
			reqDeviceHash: defaultDeviceHash,
			reqCode:       func() string { return "123456" },
			wantErr:       models.ErrSessionExpired,
		},
		{
			name: "Error: PreAuthToken claims missing Subject field",
			getPreAuthToken: func(t *testing.T) string {
				return generateTokenWithoutSubject(privKey, "localhost", defaultDeviceHash)
			},
			reqDeviceHash: defaultDeviceHash,
			reqCode:       func() string { return "123456" },
			wantErr:       models.ErrSessionData,
		},
		{
			name: "Error: Security Alert - Device hash mismatch",
			getPreAuthToken: func(t *testing.T) string {
				_, token, _ := utils.GeneratePreAuthToken(defaultUserID, "localhost", "original_device_hash", privKey, time.Now().Add(5*time.Minute))
				return token
			},
			reqDeviceHash: "different_device_hash",
			reqCode:       func() string { return "123456" },
			wantErr:       models.ErrSessionExpired,
		},
		{
			name: "Error: PreAuthToken was already used and is blacklisted in Redis",
			getPreAuthToken: func(t *testing.T) string {
				_, token, _ := utils.GeneratePreAuthToken(defaultUserID, "localhost", defaultDeviceHash, privKey, time.Now().Add(5*time.Minute))
				return token
			},
			reqDeviceHash: defaultDeviceHash,
			reqCode:       func() string { return "123456" },
			redisExistsFn: func(ctx context.Context, keys ...string) *redis.IntCmd {
				cmd := redis.NewIntCmd(ctx)
				cmd.SetVal(1)
				return cmd
			},
			wantErr: models.ErrSessionExpired,
		},
		{
			name: "Error: User not found in database",
			getPreAuthToken: func(t *testing.T) string {
				_, token, _ := utils.GeneratePreAuthToken(defaultUserID, "localhost", defaultDeviceHash, privKey, time.Now().Add(5*time.Minute))
				return token
			},
			reqDeviceHash: defaultDeviceHash,
			reqCode:       func() string { return "123456" },
			getByUserIdFn: func(ctx context.Context, id string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			wantErr: models.ErrUserNotFound,
		},
		{
			name: "Error: Database error when fetching user",
			getPreAuthToken: func(t *testing.T) string {
				_, token, _ := utils.GeneratePreAuthToken(defaultUserID, "localhost", defaultDeviceHash, privKey, time.Now().Add(5*time.Minute))
				return token
			},
			reqDeviceHash: defaultDeviceHash,
			reqCode:       func() string { return "123456" },
			getByUserIdFn: func(ctx context.Context, id string) (*models.User, error) {
				return nil, errors.New("db connection failure")
			},
			wantErr: models.ErrFailQueryUser,
		},
		{
			name: "Error: User account is disabled (Enable = false)",
			getPreAuthToken: func(t *testing.T) string {
				_, token, _ := utils.GeneratePreAuthToken(defaultUserID, "localhost", defaultDeviceHash, privKey, time.Now().Add(5*time.Minute))
				return token
			},
			reqDeviceHash: defaultDeviceHash,
			reqCode:       func() string { return "123456" },
			getByUserIdFn: func(ctx context.Context, id string) (*models.User, error) {
				return &models.User{Id: defaultUserID, Enable: false}, nil
			},
			wantErr: models.ErrUserNotEnabled,
		},
		{
			name: "Error: User has no 2FA secret configured (nil or empty)",
			getPreAuthToken: func(t *testing.T) string {
				_, token, _ := utils.GeneratePreAuthToken(defaultUserID, "localhost", defaultDeviceHash, privKey, time.Now().Add(5*time.Minute))
				return token
			},
			reqDeviceHash: defaultDeviceHash,
			reqCode:       func() string { return "123456" },
			getByUserIdFn: func(ctx context.Context, id string) (*models.User, error) {
				return &models.User{Id: defaultUserID, Enable: true, TwoFactorSecret: nil}, nil
			},
			wantErr: models.Err2FANotInitiated,
		},
		{
			name: "Error: Decryption of 2FA secret fails (corrupted secret)",
			getPreAuthToken: func(t *testing.T) string {
				_, token, _ := utils.GeneratePreAuthToken(defaultUserID, "localhost", defaultDeviceHash, privKey, time.Now().Add(5*time.Minute))
				return token
			},
			reqDeviceHash: defaultDeviceHash,
			reqCode:       func() string { return "123456" },
			getByUserIdFn: func(ctx context.Context, id string) (*models.User, error) {
				invalidSecret := "not_a_valid_aes_ciphertext"
				return &models.User{
					Id:              defaultUserID,
					Enable:          true,
					TwoFactorSecret: &invalidSecret,
				}, nil
			},
			wantErr: models.ErrUnexpectedLogin,
		},
		{
			name: "Error: Invalid TOTP code provided",
			getPreAuthToken: func(t *testing.T) string {
				_, token, _ := utils.GeneratePreAuthToken(defaultUserID, "localhost", defaultDeviceHash, privKey, time.Now().Add(5*time.Minute))
				return token
			},
			reqDeviceHash: defaultDeviceHash,
			reqCode:       func() string { return "000000" },
			getByUserIdFn: func(ctx context.Context, id string) (*models.User, error) {
				return validUser, nil
			},
			wantErr: models.Err2FAInvalid,
		},
		{
			name: "Error: SessionIssuer fails to issue final session",
			getPreAuthToken: func(t *testing.T) string {
				_, token, _ := utils.GeneratePreAuthToken(defaultUserID, "localhost", defaultDeviceHash, privKey, time.Now().Add(5*time.Minute))
				return token
			},
			reqDeviceHash: defaultDeviceHash,
			reqCode: func() string {
				code, _ := totp.GenerateCode(rawTotpSecret, time.Now())
				return code
			},
			getByUserIdFn: func(ctx context.Context, id string) (*models.User, error) {
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
			mockRepo := &MockTwoFactorUserRepository{
				GetUserByIdForAuthFunc: tt.getByUserIdFn,
			}

			mockIssuer := &MockTwoFactorSessionIssuer{
				IssueSessionFunc: tt.issueSessionFn,
			}

			mockRedis := &MockTwoFactorRedis{
				ExistsFunc: tt.redisExistsFn,
				SetFunc:    tt.redisSetFn,
			}

			service := newTestTwoFactorService(mockRepo, mockIssuer, mockRedis, pubKey, privKey)

			req := dto.Verify2FARequest{
				PreAuthToken: tt.getPreAuthToken(t),
				DeviceHash:   tt.reqDeviceHash,
				Code:         tt.reqCode(),
			}

			result, err := service.VerifyLogin2FA(context.Background(), req)

			if tt.wantErr != nil {
				if err == nil || (!errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error()) {
					t.Fatalf("VerifyLogin2FA() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("VerifyLogin2FA() unexpected error = %v", err)
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
