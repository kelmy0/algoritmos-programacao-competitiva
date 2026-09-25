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

type MockSignUpService struct {
	SignUpFunc func(ctx context.Context, data dto.SignUpRequest) (auth.SignUpResult, error)
}

func (m *MockSignUpService) SignUp(ctx context.Context, data dto.SignUpRequest) (auth.SignUpResult, error) {
	if m.SignUpFunc != nil {
		return m.SignUpFunc(ctx, data)
	}
	return auth.SignUpResult{}, nil
}
func TestSignUpHandler_SignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const (
		testAppDomain           = "localhost"
		testRefreshDurationDays = 7
		testIsProduction        = false
	)

	validPayload := dto.SignUpRequest{
		Name:            "New User",
		Username:        "newuser",
		Email:           "new.user@exemplo.com",
		Password:        "Password123!",
		ConfirmPassword: "Password123!",
	}

	validResultWithCookie := auth.SignUpResult{
		SignUpResponse: dto.SignUpResponse{
			AccessToken: "access_token_123",
			Success:     true,
			AutoLogin:   true,
		},
		RefreshToken: "refresh_token_valid_456",
	}

	tests := []struct {
		name             string
		requestBody      any
		signUpFn         func(ctx context.Context, data dto.SignUpRequest) (auth.SignUpResult, error)
		wantStatusCode   int
		wantRefreshToken bool
	}{
		{
			name:             "Error: Invalid JSON body returns 400 Bad Request",
			requestBody:      `{ json_invalido: `,
			signUpFn:         nil,
			wantStatusCode:   http.StatusBadRequest,
			wantRefreshToken: false,
		},
		{
			name:        "Error: Service returns AppError mapped status code (Conflict 409)",
			requestBody: validPayload,
			signUpFn: func(ctx context.Context, data dto.SignUpRequest) (auth.SignUpResult, error) {
				return auth.SignUpResult{}, &models.AppError{
					StatusCode: http.StatusConflict,
					Code:       "EMAIL_EXISTS",
					Message:    "Email already registered",
				}
			},
			wantStatusCode:   http.StatusConflict,
			wantRefreshToken: false,
		},
		{
			name:        "Error: Service returns unexpected error mapped to 500 Internal Server Error",
			requestBody: validPayload,
			signUpFn: func(ctx context.Context, data dto.SignUpRequest) (auth.SignUpResult, error) {
				return auth.SignUpResult{}, errors.New("database failure")
			},
			wantStatusCode:   http.StatusInternalServerError,
			wantRefreshToken: false,
		},
		{
			name:        "Success: Registration with autoLogin sets cookie and returns 200 OK",
			requestBody: validPayload,
			signUpFn: func(ctx context.Context, data dto.SignUpRequest) (auth.SignUpResult, error) {
				return validResultWithCookie, nil
			},
			wantStatusCode:   http.StatusOK,
			wantRefreshToken: true,
		},
		{
			name:        "Success: Registration without autoLogin omits cookie and returns 200 OK",
			requestBody: validPayload,
			signUpFn: func(ctx context.Context, data dto.SignUpRequest) (auth.SignUpResult, error) {
				return auth.SignUpResult{
					SignUpResponse: dto.SignUpResponse{
						AccessToken: "",
						Success:     true,
						AutoLogin:   false,
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
			mockService := &MockSignUpService{
				SignUpFunc: tt.signUpFn,
			}

			handler := authhandler.NewSignUpHandler(
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

			req, err := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler.SignUp(c)

			if w.Code != tt.wantStatusCode {
				t.Errorf("SignUp() statusCode = %d, want %d. Body: %s", w.Code, tt.wantStatusCode, w.Body.String())
			}

			cookies := w.Result().Cookies()
			hasCookie := len(cookies) > 0

			if hasCookie != tt.wantRefreshToken {
				t.Errorf("SignUp() cookie presence = %v, want %v", hasCookie, tt.wantRefreshToken)
			}

			if w.Code == http.StatusOK {
				var gotResponse dto.SignUpResponse
				if err := json.Unmarshal(w.Body.Bytes(), &gotResponse); err != nil {
					t.Fatalf("failed to unmarshal response body into dto.SignUpResponse: %v", err)
				}

				if !gotResponse.Success {
					t.Errorf("SignUp() response Success = false, want true")
				}
			}
		})
	}
}
