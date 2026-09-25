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
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"github.com/redis/go-redis/v9"
)

type testSessionAuthRepository interface {
	RotateRefreshToken(ctx context.Context, oldTokenId, newTokenId, userId, familyId string, newExpiresAt time.Time) error
	RevokeFamily(ctx context.Context, familyId string) error
	GetRefreshTokenById(ctx context.Context, id string) (*models.RefreshToken, error)
	DeleteFamily(ctx context.Context, familyId string) error
	DeleteOtherFamilies(ctx context.Context, userId, currentFamilyId string) error
	DeleteAllUserRefreshTokens(ctx context.Context, userId string) error
}

type testSessionUserRepository interface {
	GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error)
}

type testSessionRedis interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
}

type MockSessionAuthRepository struct {
	RotateRefreshTokenFunc         func(ctx context.Context, oldTokenId, newTokenId, userId, familyId string, newExpiresAt time.Time) error
	RevokeFamilyFunc               func(ctx context.Context, familyId string) error
	GetRefreshTokenByIdFunc        func(ctx context.Context, id string) (*models.RefreshToken, error)
	DeleteFamilyFunc               func(ctx context.Context, familyId string) error
	DeleteOtherFamiliesFunc        func(ctx context.Context, userId, currentFamilyId string) error
	DeleteAllUserRefreshTokensFunc func(ctx context.Context, userId string) error
}

func (m *MockSessionAuthRepository) RotateRefreshToken(ctx context.Context, oldTokenId, newTokenId, userId, familyId string, newExpiresAt time.Time) error {
	if m.RotateRefreshTokenFunc != nil {
		return m.RotateRefreshTokenFunc(ctx, oldTokenId, newTokenId, userId, familyId, newExpiresAt)
	}
	return nil
}

func (m *MockSessionAuthRepository) RevokeFamily(ctx context.Context, familyId string) error {
	if m.RevokeFamilyFunc != nil {
		return m.RevokeFamilyFunc(ctx, familyId)
	}
	return nil
}

func (m *MockSessionAuthRepository) GetRefreshTokenById(ctx context.Context, id string) (*models.RefreshToken, error) {
	if m.GetRefreshTokenByIdFunc != nil {
		return m.GetRefreshTokenByIdFunc(ctx, id)
	}
	return &models.RefreshToken{
		Id:        id,
		UserId:    "user_123",
		FamilyId:  "family_123",
		IsRevoked: false,
	}, nil
}

func (m *MockSessionAuthRepository) DeleteFamily(ctx context.Context, familyId string) error {
	if m.DeleteFamilyFunc != nil {
		return m.DeleteFamilyFunc(ctx, familyId)
	}
	return nil
}

func (m *MockSessionAuthRepository) DeleteOtherFamilies(ctx context.Context, userId, currentFamilyId string) error {
	if m.DeleteOtherFamiliesFunc != nil {
		return m.DeleteOtherFamiliesFunc(ctx, userId, currentFamilyId)
	}
	return nil
}

func (m *MockSessionAuthRepository) DeleteAllUserRefreshTokens(ctx context.Context, userId string) error {
	if m.DeleteAllUserRefreshTokensFunc != nil {
		return m.DeleteAllUserRefreshTokensFunc(ctx, userId)
	}
	return nil
}

type MockSessionUserRepository struct {
	GetUserByIdForAuthFunc func(ctx context.Context, id string) (*models.User, error)
}

func (m *MockSessionUserRepository) GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error) {
	if m.GetUserByIdForAuthFunc != nil {
		return m.GetUserByIdForAuthFunc(ctx, id)
	}
	return &models.User{
		Id:          id,
		Email:       "user@example.com",
		Username:    "user123",
		Name:        "User Test",
		Enable:      true,
		Permissions: []string{"read"},
		Role:        models.Role{IsEmployee: false},
	}, nil
}

type MockSessionRedis struct {
	SetFunc func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
}

func (m *MockSessionRedis) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, value, expiration)
	}
	cmd := redis.NewStatusCmd(ctx)
	cmd.SetVal("OK")
	return cmd
}

type testKeys struct {
	pubKey     ed25519.PublicKey
	accessPriv ed25519.PrivateKey
	refPriv    ed25519.PrivateKey
}

func generateTestKeys() testKeys {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	_, accessPriv, _ := ed25519.GenerateKey(rand.Reader)
	return testKeys{
		pubKey:     pub,
		accessPriv: accessPriv,
		refPriv:    priv,
	}
}

func newTestSessionService(authRepo testSessionAuthRepository, userRepo testSessionUserRepository, redisClient testSessionRedis, keys testKeys) *auth.SessionService {
	return auth.NewSessionService(
		authRepo,
		userRepo,
		redisClient,
		keys.pubKey,
		keys.accessPriv,
		keys.refPriv,
		15,
		7,
		"localhost",
	)
}

func TestRefreshToken(t *testing.T) {
	keys := generateTestKeys()
	appDomain := "localhost"
	deviceHash := "device_123"
	userId := "user_123"

	makeValidToken := func(uId, familyId, devHash string, expDays int) string {
		exp := time.Now().AddDate(0, 0, expDays)
		_, _, token, _ := utils.GenerateRefreshToken(uId, appDomain, familyId, devHash, keys.refPriv, exp)
		return token
	}

	validToken := makeValidToken(userId, "family_123", deviceHash, 7)

	tests := []struct {
		name               string
		refreshTokenString string
		deviceHash         string
		getDbTokenFn       func(ctx context.Context, id string) (*models.RefreshToken, error)
		getUserFn          func(ctx context.Context, id string) (*models.User, error)
		rotateTokenFn      func(ctx context.Context, oldTokenId, newTokenId, userId, familyId string, newExpiresAt time.Time) error
		accessKeyOverride  ed25519.PrivateKey
		refreshKeyOverride ed25519.PrivateKey
		wantRevokeCalled   bool
		wantErr            error
	}{
		{
			name:               "Success: Refresh token successfully rotated",
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			wantErr:            nil,
		},
		{
			name:               "Error: JWT validation failure (invalid/malformed token)",
			refreshTokenString: "invalid.jwt.token",
			deviceHash:         deviceHash,
			wantErr:            models.ErrInvalidOrExpiredRefresh,
		},
		{
			name:               "Error: Security Alert - Device hash mismatch (triggers RevokeFamily)",
			refreshTokenString: validToken,
			deviceHash:         "different_device_hash",
			wantRevokeCalled:   true,
			wantErr:            models.ErrInvalidOrExpiredRefresh,
		},
		{
			name:               "Error: Repository error during GetRefreshTokenById query",
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			getDbTokenFn: func(ctx context.Context, id string) (*models.RefreshToken, error) {
				return nil, errors.New("db connection failure")
			},
			wantErr: models.ErrInvalidOrExpiredRefresh,
		},
		{
			name:               "Error: Refresh token not found in DB (dbToken == nil)",
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			getDbTokenFn: func(ctx context.Context, id string) (*models.RefreshToken, error) {
				return nil, nil
			},
			wantErr: models.ErrInvalidOrExpiredRefresh,
		},
		{
			name:               "Error: Security Alert - Token already revoked in DB (triggers RevokeFamily)",
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			getDbTokenFn: func(ctx context.Context, id string) (*models.RefreshToken, error) {
				return &models.RefreshToken{
					Id:        id,
					UserId:    userId,
					FamilyId:  "family_123",
					IsRevoked: true,
				}, nil
			},
			wantRevokeCalled: true,
			wantErr:          models.ErrInvalidOrExpiredRefresh,
		},
		{
			name:               "Error: Token integrity mismatch (DB UserId != JWT Subject)",
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			getDbTokenFn: func(ctx context.Context, id string) (*models.RefreshToken, error) {
				return &models.RefreshToken{
					Id:        id,
					UserId:    "different_user_id",
					FamilyId:  "family_123",
					IsRevoked: false,
				}, nil
			},
			wantErr: models.ErrTokenMetadataMisMatch,
		},
		{
			name:               "Error: User not found in database (models.ErrUserNotFound)",
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			getUserFn: func(ctx context.Context, id string) (*models.User, error) {
				return nil, models.ErrUserNotFound
			},
			wantErr: models.ErrUserNotFound,
		},
		{
			name:               "Error: Database failure during GetUserByIdForAuth",
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			getUserFn: func(ctx context.Context, id string) (*models.User, error) {
				return nil, errors.New("sql select error")
			},
			wantErr: models.ErrUserNotFound,
		},
		{
			name:               "Error: User account is disabled",
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			getUserFn: func(ctx context.Context, id string) (*models.User, error) {
				return &models.User{Id: id, Enable: false}, nil
			},
			wantErr: models.ErrUserNotEnabled,
		},
		{
			name:               "Error: GenerateAccessToken fails due to invalid access key",
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			accessKeyOverride:  ed25519.PrivateKey([]byte("invalid_key")),
			wantErr:            models.ErrGeneratingToken,
		},
		{
			name:               "Error: GenerateRefreshToken fails due to invalid refresh key",
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			refreshKeyOverride: ed25519.PrivateKey([]byte("invalid_key")),
			wantErr:            models.ErrGeneratingToken,
		},
		{
			name:               "Error: Database failure during RotateRefreshToken",
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			rotateTokenFn: func(ctx context.Context, oldTokenId, newTokenId, userId, familyId string, newExpiresAt time.Time) error {
				return errors.New("db transaction failed")
			},
			wantErr: models.ErrGeneratingToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			revokeCalled := false
			mockAuthRepo := &MockSessionAuthRepository{
				GetRefreshTokenByIdFunc: tt.getDbTokenFn,
				RotateRefreshTokenFunc:  tt.rotateTokenFn,
				RevokeFamilyFunc: func(ctx context.Context, familyId string) error {
					revokeCalled = true
					return nil
				},
			}

			mockUserRepo := &MockSessionUserRepository{
				GetUserByIdForAuthFunc: tt.getUserFn,
			}

			currentKeys := keys
			if len(tt.accessKeyOverride) > 0 {
				currentKeys.accessPriv = tt.accessKeyOverride
			}
			if len(tt.refreshKeyOverride) > 0 {
				currentKeys.refPriv = tt.refreshKeyOverride
			}

			service := newTestSessionService(mockAuthRepo, mockUserRepo, &MockSessionRedis{}, currentKeys)

			res, err := service.RefreshToken(context.Background(), tt.refreshTokenString, tt.deviceHash)

			if tt.wantErr != nil {
				if err == nil || (!errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error()) {
					t.Fatalf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Fatalf("RefreshToken() unexpected error = %v", err)
				}
				if res.AccessToken == "" || res.RefreshToken == "" {
					t.Errorf("RefreshToken() expected tokens in response, got empty")
				}
			}

			if revokeCalled != tt.wantRevokeCalled {
				t.Errorf("RevokeFamily() called = %v, wantRevokeCalled %v", revokeCalled, tt.wantRevokeCalled)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	keys := generateTestKeys()
	appDomain := "localhost"
	userId := "user_123"

	exp := time.Now().AddDate(0, 0, 7)
	_, _, validToken, _ := utils.GenerateRefreshToken(userId, appDomain, "family_123", "device_123", keys.refPriv, exp)

	tests := []struct {
		name               string
		userId             string
		refreshTokenString string
		accessJti          string
		accessExpiresAt    time.Time
		deleteFamilyFn     func(ctx context.Context, familyId string) error
		redisSetFn         func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
		wantErr            error
		wantRedisCalled    bool
	}{
		{
			name:               "Success: Logout with access token active (added to Redis blacklist)",
			userId:             userId,
			refreshTokenString: validToken,
			accessJti:          "jti_123",
			accessExpiresAt:    time.Now().Add(10 * time.Minute),
			wantErr:            nil,
			wantRedisCalled:    true,
		},
		{
			name:               "Success: Logout with access token already expired (skips Redis)",
			userId:             userId,
			refreshTokenString: validToken,
			accessJti:          "jti_123",
			accessExpiresAt:    time.Now().Add(-1 * time.Minute),
			wantErr:            nil,
			wantRedisCalled:    false,
		},
		{
			name:               "Success: Logout gracefully handles Redis error",
			userId:             userId,
			refreshTokenString: validToken,
			accessJti:          "jti_123",
			accessExpiresAt:    time.Now().Add(10 * time.Minute),
			redisSetFn: func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
				cmd := redis.NewStatusCmd(ctx)
				cmd.SetErr(errors.New("redis connection down"))
				return cmd
			},
			wantErr:         nil,
			wantRedisCalled: true,
		},
		{
			name:               "Error: Invalid refresh token",
			userId:             userId,
			refreshTokenString: "invalid_token",
			accessJti:          "jti_123",
			accessExpiresAt:    time.Now().Add(10 * time.Minute),
			wantErr:            models.ErrInvalidOrExpiredRefresh,
		},
		{
			name:               "Error: Subject in token does not match userId parameter",
			userId:             "wrong_user_id",
			refreshTokenString: validToken,
			accessJti:          "jti_123",
			accessExpiresAt:    time.Now().Add(10 * time.Minute),
			wantErr:            models.ErrTokenMetadataMisMatch,
		},
		{
			name:               "Error: Database failure during DeleteFamily",
			userId:             userId,
			refreshTokenString: validToken,
			accessJti:          "jti_123",
			accessExpiresAt:    time.Now().Add(10 * time.Minute),
			deleteFamilyFn: func(ctx context.Context, familyId string) error {
				return errors.New("database error")
			},
			wantErr: models.ErrUnexpectedLogout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisCalled := false
			mockAuthRepo := &MockSessionAuthRepository{
				DeleteFamilyFunc: tt.deleteFamilyFn,
			}
			mockRedis := &MockSessionRedis{
				SetFunc: func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
					redisCalled = true
					if tt.redisSetFn != nil {
						return tt.redisSetFn(ctx, key, value, expiration)
					}
					cmd := redis.NewStatusCmd(ctx)
					cmd.SetVal("OK")
					return cmd
				},
			}

			service := newTestSessionService(mockAuthRepo, &MockSessionUserRepository{}, mockRedis, keys)

			err := service.Logout(context.Background(), tt.userId, tt.refreshTokenString, tt.accessJti, tt.accessExpiresAt)

			if tt.wantErr != nil {
				if err == nil || (!errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error()) {
					t.Fatalf("Logout() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Logout() unexpected error = %v", err)
			}

			if redisCalled != tt.wantRedisCalled {
				t.Errorf("Redis Set called = %v, wantRedisCalled %v", redisCalled, tt.wantRedisCalled)
			}
		})
	}
}

func TestLogoutOtherDevices(t *testing.T) {
	keys := generateTestKeys()
	appDomain := "localhost"
	userId := "user_123"
	deviceHash := "device_123"

	exp := time.Now().AddDate(0, 0, 7)
	_, _, validToken, _ := utils.GenerateRefreshToken(userId, appDomain, "family_123", deviceHash, keys.refPriv, exp)

	tests := []struct {
		name               string
		userId             string
		refreshTokenString string
		accessJti          string
		deviceHash         string
		deleteOtherFn      func(ctx context.Context, userId, currentFamilyId string) error
		redisSetFn         func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
		wantRevokeCalled   bool
		wantErr            error
	}{
		{
			name:               "Success: Logout other devices successfully",
			userId:             userId,
			refreshTokenString: validToken,
			accessJti:          "jti_123",
			deviceHash:         deviceHash,
			wantErr:            nil,
		},
		{
			name:               "Error: Invalid refresh token",
			userId:             userId,
			refreshTokenString: "invalid_token",
			accessJti:          "jti_123",
			deviceHash:         deviceHash,
			wantErr:            models.ErrInvalidOrExpiredRefresh,
		},
		{
			name:               "Error: Subject in token does not match userId parameter",
			userId:             "different_user",
			refreshTokenString: validToken,
			accessJti:          "jti_123",
			deviceHash:         deviceHash,
			wantErr:            models.ErrTokenMetadataMisMatch,
		},
		{
			name:               "Error: Device hash mismatch (triggers RevokeFamily)",
			userId:             userId,
			refreshTokenString: validToken,
			accessJti:          "jti_123",
			deviceHash:         "wrong_device",
			wantRevokeCalled:   true,
			wantErr:            models.ErrInvalidOrExpiredRefresh,
		},
		{
			name:               "Error: Database failure during DeleteOtherFamilies",
			userId:             userId,
			refreshTokenString: validToken,
			accessJti:          "jti_123",
			deviceHash:         deviceHash,
			deleteOtherFn: func(ctx context.Context, userId, currentFamilyId string) error {
				return errors.New("db error")
			},
			wantErr: models.ErrUnexpectedLogout,
		},
		{
			name:               "Error: Redis set failure",
			userId:             userId,
			refreshTokenString: validToken,
			accessJti:          "jti_123",
			deviceHash:         deviceHash,
			redisSetFn: func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
				cmd := redis.NewStatusCmd(ctx)
				cmd.SetErr(errors.New("redis timeout"))
				return cmd
			},
			wantErr: models.ErrUnexpectedLogout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			revokeCalled := false
			mockAuthRepo := &MockSessionAuthRepository{
				DeleteOtherFamiliesFunc: tt.deleteOtherFn,
				RevokeFamilyFunc: func(ctx context.Context, familyId string) error {
					revokeCalled = true
					return nil
				},
			}

			mockRedis := &MockSessionRedis{
				SetFunc: tt.redisSetFn,
			}

			service := newTestSessionService(mockAuthRepo, &MockSessionUserRepository{}, mockRedis, keys)

			err := service.LogoutOtherDevices(context.Background(), tt.userId, tt.refreshTokenString, tt.accessJti, tt.deviceHash)

			if tt.wantErr != nil {
				if err == nil || (!errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error()) {
					t.Fatalf("LogoutOtherDevices() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Fatalf("LogoutOtherDevices() unexpected error = %v", err)
			}

			if revokeCalled != tt.wantRevokeCalled {
				t.Errorf("RevokeFamily() called = %v, wantRevokeCalled %v", revokeCalled, tt.wantRevokeCalled)
			}
		})
	}
}

func TestLogoutAllDevices(t *testing.T) {
	keys := generateTestKeys()
	appDomain := "localhost"
	userId := "user_123"
	deviceHash := "device_123"

	exp := time.Now().AddDate(0, 0, 7)
	_, _, validToken, _ := utils.GenerateRefreshToken(userId, appDomain, "family_123", deviceHash, keys.refPriv, exp)

	tests := []struct {
		name               string
		userId             string
		refreshTokenString string
		deviceHash         string
		deleteAllFn        func(ctx context.Context, userId string) error
		redisSetFn         func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
		wantRevokeCalled   bool
		wantErr            error
	}{
		{
			name:               "Success: Logout all devices successfully",
			userId:             userId,
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			wantErr:            nil,
		},
		{
			name:               "Error: Invalid refresh token",
			userId:             userId,
			refreshTokenString: "invalid_token",
			deviceHash:         deviceHash,
			wantErr:            models.ErrInvalidOrExpiredRefresh,
		},
		{
			name:               "Error: Subject in token does not match userId parameter",
			userId:             "different_user",
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			wantErr:            models.ErrTokenMetadataMisMatch,
		},
		{
			name:               "Error: Device hash mismatch (triggers RevokeFamily)",
			userId:             userId,
			refreshTokenString: validToken,
			deviceHash:         "wrong_device",
			wantRevokeCalled:   true,
			wantErr:            models.ErrInvalidOrExpiredRefresh,
		},
		{
			name:               "Error: Database failure during DeleteAllUserRefreshTokens",
			userId:             userId,
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			deleteAllFn: func(ctx context.Context, userId string) error {
				return errors.New("db delete error")
			},
			wantErr: models.ErrUnexpectedLogout,
		},
		{
			name:               "Error: Redis set failure",
			userId:             userId,
			refreshTokenString: validToken,
			deviceHash:         deviceHash,
			redisSetFn: func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
				cmd := redis.NewStatusCmd(ctx)
				cmd.SetErr(errors.New("redis connection refused"))
				return cmd
			},
			wantErr: models.ErrUnexpectedLogout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			revokeCalled := false
			mockAuthRepo := &MockSessionAuthRepository{
				DeleteAllUserRefreshTokensFunc: tt.deleteAllFn,
				RevokeFamilyFunc: func(ctx context.Context, familyId string) error {
					revokeCalled = true
					return nil
				},
			}

			mockRedis := &MockSessionRedis{
				SetFunc: tt.redisSetFn,
			}

			service := newTestSessionService(mockAuthRepo, &MockSessionUserRepository{}, mockRedis, keys)

			err := service.LogoutAllDevices(context.Background(), tt.userId, tt.refreshTokenString, tt.deviceHash)

			if tt.wantErr != nil {
				if err == nil || (!errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error()) {
					t.Fatalf("LogoutAllDevices() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Fatalf("LogoutAllDevices() unexpected error = %v", err)
			}

			if revokeCalled != tt.wantRevokeCalled {
				t.Errorf("RevokeFamily() called = %v, wantRevokeCalled %v", revokeCalled, tt.wantRevokeCalled)
			}
		})
	}
}
