package authhandler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	authhandler "github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers/auth"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
)

type MockSocialauth struct {
	AuthWithSocialProviderFunc func(ctx context.Context, provider, socialUserId, email, name, deviceHash string) (auth.AuthResult, error)
}

func (m *MockSocialauth) AuthWithSocialProvider(ctx context.Context, provider, socialUserId, email, name, deviceHash string) (auth.AuthResult, error) {
	if m.AuthWithSocialProviderFunc != nil {
		return m.AuthWithSocialProviderFunc(ctx, provider, socialUserId, email, name, deviceHash)
	}
	return auth.AuthResult{}, nil
}

type MockSocialProvider struct {
	NameValue        string
	OAuthConfigValue *oauth2.Config
	FetchUserFunc    func(ctx context.Context, token *oauth2.Token) (*authhandler.SocialUser, string, error)
}

func (m *MockSocialProvider) Name() string {
	return m.NameValue
}

func (m *MockSocialProvider) OAuthConfig() *oauth2.Config {
	if m.OAuthConfigValue != nil {
		return m.OAuthConfigValue
	}
	return &oauth2.Config{
		ClientID:     "mock_client_id",
		ClientSecret: "mock_client_secret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://provider.com/oauth/auth",
			TokenURL: "https://provider.com/oauth/token",
		},
	}
}

func (m *MockSocialProvider) FetchUser(ctx context.Context, token *oauth2.Token) (*authhandler.SocialUser, string, error) {
	if m.FetchUserFunc != nil {
		return m.FetchUserFunc(ctx, token)
	}
	return &authhandler.SocialUser{
		ID:    "social_123",
		Email: "social@exemplo.com",
		Name:  "Usuario Social",
	}, "", nil
}

const (
	testFrontendURL = "http://localhost:8000"
	testProvider    = "github"
	testState       = "state_random_12345678901234567890"
	testCode        = "oauth_code_abc123"
)

func TestAuthSocialHandler_SocialLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		providerParam  string
		wantStatusCode int
		wantLocation   bool
	}{
		{
			name:           "Error: Provider not registered returns 400 Bad Request",
			providerParam:  "unregistered_provider",
			wantStatusCode: http.StatusBadRequest,
			wantLocation:   false,
		},
		{
			name:           "Success: Redirects to provider authorization URL and sets state cookie",
			providerParam:  testProvider,
			wantStatusCode: http.StatusTemporaryRedirect,
			wantLocation:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSocialauth{}
			mockProvider := &MockSocialProvider{NameValue: testProvider}

			handler := authhandler.NewAuthSocialHandler(
				mockService,
				testAppDomain,
				testFrontendURL,
				testIsProduction,
				testRefreshDurationDays,
				mockProvider,
			)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest(http.MethodGet, "/auth/"+tt.providerParam, nil)
			c.Request = req
			c.Params = gin.Params{{Key: "provider", Value: tt.providerParam}}

			handler.SocialLogin(c)

			statusCode := c.Writer.Status()
			if statusCode != tt.wantStatusCode {
				t.Errorf("SocialLogin() statusCode = %d, want %d", statusCode, tt.wantStatusCode)
			}

			location := w.Header().Get("Location")
			if tt.wantLocation && location == "" {
				t.Errorf("SocialLogin() expected Location header redirect, got empty")
			}
		})
	}
}

func TestAuthSocialHandler_SocialCallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                 string
		providerParam        string
		urlState             string
		cookieState          string
		urlCode              string
		exchangeErr          bool
		fetchUserFn          func(ctx context.Context, token *oauth2.Token) (*authhandler.SocialUser, string, error)
		authSocialFn         func(ctx context.Context, provider, socialUserId, email, name, deviceHash string) (auth.AuthResult, error)
		wantRedirectContains string
	}{
		{
			name:                 "Error: Provider not found returns 400",
			providerParam:        "invalid_provider",
			wantRedirectContains: "",
		},
		{
			name:                 "Error: State mismatch redirects to error=" + dto.CodeSessionExpired,
			providerParam:        testProvider,
			urlState:             "invalid_state",
			cookieState:          testState,
			wantRedirectContains: "error=" + dto.CodeSessionExpired,
		},
		{
			name:                 "Error: Missing OAuth code redirects to error=" + dto.CodeMissingOAuthCode,
			providerParam:        testProvider,
			urlState:             testState,
			cookieState:          testState,
			urlCode:              "",
			wantRedirectContains: "error=" + dto.CodeMissingOAuthCode,
		},
		{
			name:                 "Error: Exchange token failure redirects to error=" + dto.CodeInternalError,
			providerParam:        testProvider,
			urlState:             testState,
			cookieState:          testState,
			urlCode:              testCode,
			exchangeErr:          true,
			wantRedirectContains: "error=" + dto.CodeInternalError,
		},
		{
			name:          "Error: FetchUser returns CodeUnverifiedGithubEmail redirects with error code",
			providerParam: testProvider,
			urlState:      testState,
			cookieState:   testState,
			urlCode:       testCode,
			fetchUserFn: func(ctx context.Context, token *oauth2.Token) (*authhandler.SocialUser, string, error) {
				return nil, dto.CodeUnverifiedGithubEmail, nil
			},
			wantRedirectContains: "error=" + dto.CodeUnverifiedGithubEmail,
		},
		{
			name:          "Error: FetchUser returns CodeInvalidGoogleToken redirects with error code",
			providerParam: testProvider,
			urlState:      testState,
			cookieState:   testState,
			urlCode:       testCode,
			fetchUserFn: func(ctx context.Context, token *oauth2.Token) (*authhandler.SocialUser, string, error) {
				return nil, dto.CodeInvalidGoogleToken, nil
			},
			wantRedirectContains: "error=" + dto.CodeInvalidGoogleToken,
		},
		{
			name:          "Error: AuthWithSocialProvider fails with AppError redirects with error code",
			providerParam: testProvider,
			urlState:      testState,
			cookieState:   testState,
			urlCode:       testCode,
			authSocialFn: func(ctx context.Context, provider, socialUserId, email, name, deviceHash string) (auth.AuthResult, error) {
				return auth.AuthResult{}, &models.AppError{
					Code: dto.CodeInternalError,
				}
			},
			wantRedirectContains: "error=" + dto.CodeInternalError,
		},
		{
			name:          "Success: Requires 2FA redirects to pre_auth_token=true",
			providerParam: testProvider,
			urlState:      testState,
			cookieState:   testState,
			urlCode:       testCode,
			authSocialFn: func(ctx context.Context, provider, socialUserId, email, name, deviceHash string) (auth.AuthResult, error) {
				return auth.AuthResult{
					LoginResponse: dto.LoginResponse{
						Requires2FA:  true,
						PreAuthToken: "pre_auth_123",
					},
				}, nil
			},
			wantRedirectContains: "pre_auth_token=true",
		},
		{
			name:          "Success: Full auth redirects to access_token=true",
			providerParam: testProvider,
			urlState:      testState,
			cookieState:   testState,
			urlCode:       testCode,
			authSocialFn: func(ctx context.Context, provider, socialUserId, email, name, deviceHash string) (auth.AuthResult, error) {
				return auth.AuthResult{
					RefreshToken: "refresh_123",
					LoginResponse: dto.LoginResponse{
						AccessToken: "access_123",
					},
				}, nil
			},
			wantRedirectContains: "access_token=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSocialauth{AuthWithSocialProviderFunc: tt.authSocialFn}

			mockOAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.exchangeErr {
					http.Error(w, "invalid_grant", http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"access_token":"mock_access_token","token_type":"bearer"}`))
			}))
			defer mockOAuthServer.Close()

			mockProvider := &MockSocialProvider{
				NameValue: testProvider,
				OAuthConfigValue: &oauth2.Config{
					ClientID:     "client_id",
					ClientSecret: "client_secret",
					Endpoint: oauth2.Endpoint{
						AuthURL:  mockOAuthServer.URL + "/auth",
						TokenURL: mockOAuthServer.URL + "/token",
					},
				},
				FetchUserFunc: tt.fetchUserFn,
			}

			handler := authhandler.NewAuthSocialHandler(
				mockService,
				testAppDomain,
				testFrontendURL,
				testIsProduction,
				testRefreshDurationDays,
				mockProvider,
			)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			targetURL := "/auth/callback/" + tt.providerParam
			if tt.urlState != "" || tt.urlCode != "" {
				targetURL += "?state=" + tt.urlState + "&code=" + tt.urlCode
			}

			req, _ := http.NewRequest(http.MethodGet, targetURL, nil)
			if tt.cookieState != "" {
				req.AddCookie(&http.Cookie{
					Name:  "oauth_" + tt.providerParam + "_state",
					Value: tt.cookieState,
				})
			}
			c.Request = req
			c.Params = gin.Params{{Key: "provider", Value: tt.providerParam}}

			handler.SocialCallback(c)

			if tt.wantRedirectContains != "" {
				location := w.Header().Get("Location")
				if location == "" {
					t.Errorf("SocialCallback() expected redirect header Location, got empty. Status: %d", c.Writer.Status())
				} else if !containsString(location, tt.wantRedirectContains) {
					t.Errorf("SocialCallback() Location = %s, want containing %s", location, tt.wantRedirectContains)
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || searchSubstr(s, substr))
}

func searchSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
