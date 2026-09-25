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

type MockTwoFactorService struct {
	VerifyLogin2FAFunc func(ctx context.Context, data dto.Verify2FARequest) (auth.AuthResult, error)
}

func (m *MockTwoFactorService) VerifyLogin2FA(ctx context.Context, data dto.Verify2FARequest) (auth.AuthResult, error) {
	if m.VerifyLogin2FAFunc != nil {
		return m.VerifyLogin2FAFunc(ctx, data)
	}
	return auth.AuthResult{}, nil
}

func TestTwoFactorHandler_Verify2FA(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const (
		testAppDomain           = "localhost"
		testRefreshDurationDays = 7
		testIsProduction        = false
		testPreAuthToken        = "pre_auth_token_valid_123"
	)

	validPayload := dto.Verify2FARequest{
		Code: "123456",
	}

	validResult := auth.AuthResult{
		LoginResponse: dto.LoginResponse{
			AccessToken: "access_token_final_789",
		},
		RefreshToken: "refresh_token_final_xyz",
	}

	tests := []struct {
		name             string
		hasPreAuthCookie bool
		requestBody      any
		verify2FAFn      func(ctx context.Context, data dto.Verify2FARequest) (auth.AuthResult, error)
		wantStatusCode   int
		wantRefreshToken bool
	}{
		{
			name:             "Error: Missing pre_auth_token cookie returns 401 Unauthorized",
			hasPreAuthCookie: false,
			requestBody:      validPayload,
			verify2FAFn:      nil,
			wantStatusCode:   http.StatusUnauthorized,
			wantRefreshToken: false,
		},
		{
			name:             "Error: Invalid JSON body returns 400 Bad Request",
			hasPreAuthCookie: true,
			requestBody:      `{ json_invalido: `,
			verify2FAFn:      nil,
			wantStatusCode:   http.StatusBadRequest,
			wantRefreshToken: false,
		},
		{
			name:             "Error: Code with invalid length (len != 6) returns 400 Bad Request",
			hasPreAuthCookie: true,
			requestBody: dto.Verify2FARequest{
				Code: "123",
			},
			verify2FAFn:      nil,
			wantStatusCode:   http.StatusBadRequest,
			wantRefreshToken: false,
		},
		{
			name:             "Error: Service returns AppError for invalid 2FA code (401 Unauthorized)",
			hasPreAuthCookie: true,
			requestBody:      validPayload,
			verify2FAFn: func(ctx context.Context, data dto.Verify2FARequest) (auth.AuthResult, error) {
				return auth.AuthResult{}, &models.AppError{
					StatusCode: http.StatusUnauthorized,
					Code:       "INVALID_2FA_CODE",
					Message:    "Invalid 2FA code",
				}
			},
			wantStatusCode:   http.StatusUnauthorized,
			wantRefreshToken: false,
		},
		{
			name:             "Error: Service returns unexpected error mapped to 500",
			hasPreAuthCookie: true,
			requestBody:      validPayload,
			verify2FAFn: func(ctx context.Context, data dto.Verify2FARequest) (auth.AuthResult, error) {
				return auth.AuthResult{}, errors.New("redis session store unreachable")
			},
			wantStatusCode:   http.StatusInternalServerError,
			wantRefreshToken: false,
		},

		{
			name:             "Success: Valid 2FA code sets refresh cookie and returns 200 OK",
			hasPreAuthCookie: true,
			requestBody:      validPayload,
			verify2FAFn: func(ctx context.Context, data dto.Verify2FARequest) (auth.AuthResult, error) {
				if data.PreAuthToken != testPreAuthToken {
					return auth.AuthResult{}, errors.New("pre_auth_token mismatch")
				}
				return validResult, nil
			},
			wantStatusCode:   http.StatusOK,
			wantRefreshToken: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockTwoFactorService{
				VerifyLogin2FAFunc: tt.verify2FAFn,
			}

			handler := authhandler.NewTwoFactorHandler(
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

			req, err := http.NewRequest(http.MethodPost, "/auth/2fa/verify", bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "Go-Test-Agent")

			if tt.hasPreAuthCookie {
				req.AddCookie(&http.Cookie{
					Name:  "pre_auth_token",
					Value: testPreAuthToken,
				})
			}

			c.Request = req

			handler.Verify2FA(c)

			if w.Code != tt.wantStatusCode {
				t.Errorf("Verify2FA() statusCode = %d, want %d. Body: %s", w.Code, tt.wantStatusCode, w.Body.String())
			}

			hasCookie := false
			for _, cookie := range w.Result().Cookies() {
				if cookie.Name == "refresh_token" {
					hasCookie = true
					break
				}
			}

			if hasCookie != tt.wantRefreshToken {
				t.Errorf("Verify2FA() refresh cookie presence = %v, want %v", hasCookie, tt.wantRefreshToken)
			}
		})
	}
}
