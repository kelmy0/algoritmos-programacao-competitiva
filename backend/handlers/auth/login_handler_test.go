package authhandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	authhandler "github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers/auth"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
)

type MockLoginService struct {
	LoginFunc func(ctx context.Context, data dto.AuthRequest) (auth.AuthResult, error)
}

func (m *MockLoginService) Login(ctx context.Context, data dto.AuthRequest) (auth.AuthResult, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, data)
	}
	return auth.AuthResult{}, nil
}

func TestLoginHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const (
		testAppDomain           = "localhost"
		testRefreshDurationDays = 7
		testIsProduction        = false
	)

	validPayload := dto.AuthRequest{
		Email:    "usuario@exemplo.com",
		Password: "password123",
	}

	validResultDirectLogin := auth.AuthResult{
		LoginResponse: dto.LoginResponse{
			AccessToken: "access_token_123",
			Requires2FA: false,
		},
		RefreshToken: "refresh_token_xyz_789",
	}

	validResultRequires2FA := auth.AuthResult{
		LoginResponse: dto.LoginResponse{
			PreAuthToken: "pre_auth_token_456",
			Requires2FA:  true,
		},
		RefreshToken: "refresh_token_nao_deve_ser_gravado",
	}

	tests := []struct {
		name             string
		requestBody      any
		loginFn          func(ctx context.Context, data dto.AuthRequest) (auth.AuthResult, error)
		wantStatusCode   int
		wantRefreshToken bool
	}{
		{
			name:             "Error: Invalid JSON body returns 400 Bad Request",
			requestBody:      `{ json_invalido: `,
			loginFn:          nil,
			wantStatusCode:   http.StatusBadRequest,
			wantRefreshToken: false,
		},
		{
			name: "Error: Payload failing binding rules (password < 8 chars) returns 400",
			requestBody: dto.AuthRequest{
				Email:    "valid@email.com",
				Password: "123",
			},
			loginFn:          nil,
			wantStatusCode:   http.StatusBadRequest,
			wantRefreshToken: false,
		},
		{
			name:        "Error: Service returns AppError (Unauthorized 401)",
			requestBody: validPayload,
			loginFn: func(ctx context.Context, data dto.AuthRequest) (auth.AuthResult, error) {
				return auth.AuthResult{}, &models.AppError{
					StatusCode: http.StatusUnauthorized,
					Code:       "INVALID_CREDENTIALS",
					Message:    "Invalid email or password",
				}
			},
			wantStatusCode:   http.StatusUnauthorized,
			wantRefreshToken: false,
		},
		{
			name:        "Error: Service returns unexpected error mapped to 500",
			requestBody: validPayload,
			loginFn: func(ctx context.Context, data dto.AuthRequest) (auth.AuthResult, error) {
				return auth.AuthResult{}, errors.New("db connection timeout")
			},
			wantStatusCode:   http.StatusInternalServerError,
			wantRefreshToken: false,
		},
		{
			name:        "Success: Requires 2FA returns 200 OK and skips cookie even if RefreshToken is present",
			requestBody: validPayload,
			loginFn: func(ctx context.Context, data dto.AuthRequest) (auth.AuthResult, error) {
				return validResultRequires2FA, nil
			},
			wantStatusCode:   http.StatusOK,
			wantRefreshToken: false,
		},
		{
			name:        "Success: Direct login sets refresh cookie and returns 200 OK",
			requestBody: validPayload,
			loginFn: func(ctx context.Context, data dto.AuthRequest) (auth.AuthResult, error) {
				return validResultDirectLogin, nil
			},
			wantStatusCode:   http.StatusOK,
			wantRefreshToken: true,
		},
		{
			name:        "Success: Direct login without RefreshToken omits cookie and returns 200 OK",
			requestBody: validPayload,
			loginFn: func(ctx context.Context, data dto.AuthRequest) (auth.AuthResult, error) {
				return auth.AuthResult{
					LoginResponse: dto.LoginResponse{
						AccessToken: "access_token_123",
						Requires2FA: false,
					},
					RefreshToken: "",
				}, nil
			},
			wantStatusCode:   http.StatusOK,
			wantRefreshToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockLoginService{
				LoginFunc: tt.loginFn,
			}

			handler := authhandler.NewLoginHandler(
				mockService,
				testRefreshDurationDays,
				testIsProduction,
				testAppDomain,
			)

			var bodyBytes []byte
			switch v := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, _ = json.Marshal(v)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "Go-Test-Agent")
			c.Request = req

			handler.Login(c)

			if w.Code != tt.wantStatusCode {
				t.Errorf("Login() statusCode = %d, want %d. Body: %s", w.Code, tt.wantStatusCode, w.Body.String())
			}

			hasCookie := false
			for _, cookie := range w.Result().Cookies() {
				if cookie.Name == "refresh_token" {
					hasCookie = true
					break
				}
			}

			if hasCookie != tt.wantRefreshToken {
				t.Errorf("Login() refresh cookie presence = %v, want %v", hasCookie, tt.wantRefreshToken)
			}
		})
	}
}
