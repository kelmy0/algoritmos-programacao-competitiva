package handlers

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
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

type MockSignUpService struct {
	SignUpFunc func(ctx context.Context, req dto.SignUpRequest) (services.SignUpResult, error)
}

func (m *MockSignUpService) SignUp(ctx context.Context, req dto.SignUpRequest) (services.SignUpResult, error) {
	if m.SignUpFunc != nil {
		return m.SignUpFunc(ctx, req)
	}
	return services.SignUpResult{}, nil
}

func TestSignUpHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           any
		userAgent      string
		setupService   func() *MockSignUpService
		wantStatus     int
		wantCookieName string
		wantAutoLogin  *bool
	}{
		{
			name: "Error: Malformed JSON Body",
			body: "invalid-json-body",
			setupService: func() *MockSignUpService {
				return &MockSignUpService{}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Error: Service returns domain AppError (Email already used)",
			body: dto.SignUpRequest{
				Name:            "John Smith",
				Username:        "johnsmith",
				Email:           "used@email.com",
				Password:        "ValidPassword123!",
				ConfirmPassword: "ValidPassword123!",
			},
			setupService: func() *MockSignUpService {
				return &MockSignUpService{
					SignUpFunc: func(ctx context.Context, req dto.SignUpRequest) (services.SignUpResult, error) {
						return services.SignUpResult{}, models.ErrEmailAlreadyUsed
					},
				}
			},
			wantStatus: models.ErrEmailAlreadyUsed.StatusCode,
		},
		{
			name: "Error: Service returns unexpected generic error (Falls back to 500)",
			body: dto.SignUpRequest{
				Name:            "John Smith",
				Username:        "johnsmith",
				Email:           "john@email.com",
				Password:        "ValidPassword123!",
				ConfirmPassword: "ValidPassword123!",
			},
			setupService: func() *MockSignUpService {
				return &MockSignUpService{
					SignUpFunc: func(ctx context.Context, req dto.SignUpRequest) (services.SignUpResult, error) {
						return services.SignUpResult{}, errors.New("database connection timeout")
					},
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "Success: User registered without auto-login (No Refresh Cookie)",
			body: dto.SignUpRequest{
				Name:            "John Smith",
				Username:        "johnsmith",
				Email:           "john@email.com",
				Password:        "ValidPassword123!",
				ConfirmPassword: "ValidPassword123!",
			},
			setupService: func() *MockSignUpService {
				return &MockSignUpService{
					SignUpFunc: func(ctx context.Context, req dto.SignUpRequest) (services.SignUpResult, error) {
						return services.SignUpResult{
							RefreshToken: "",
							SignUpResponse: dto.SignUpResponse{
								AutoLogin: false,
							},
						}, nil
					},
				}
			},
			wantStatus:     http.StatusOK,
			wantCookieName: "",
			wantAutoLogin:  boolPtr(false),
		},
		{
			name:      "Success: User registered with User-Agent and Refresh Cookie set",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/152.0.0.0 Safari/537.36",
			body: dto.SignUpRequest{
				Name:            "John Smith",
				Username:        "johnsmith",
				Email:           "john@email.com",
				Password:        "ValidPassword123!",
				ConfirmPassword: "ValidPassword123!",
			},
			setupService: func() *MockSignUpService {
				return &MockSignUpService{
					SignUpFunc: func(ctx context.Context, req dto.SignUpRequest) (services.SignUpResult, error) {
						return services.SignUpResult{
							RefreshToken: "mocked-refresh-token",
							SignUpResponse: dto.SignUpResponse{
								AutoLogin: true,
							},
						}, nil
					},
				}
			},
			wantStatus:     http.StatusOK,
			wantCookieName: "refresh_token",
			wantAutoLogin:  boolPtr(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupService()
			handler := NewSignUpHandler(mockService, 7, "localhost", false)

			r := gin.New()
			r.POST("/signup", handler.SignUp)

			jsonBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", "application/json")
			if tt.userAgent != "" {
				req.Header.Set("User-Agent", tt.userAgent)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("SignUpHandler status = %d, want %d", w.Code, tt.wantStatus)
			}

			cookies := w.Result().Cookies()
			if tt.wantCookieName != "" {
				found := false
				for _, c := range cookies {
					if c.Name == tt.wantCookieName {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("SignUpHandler expected cookie %s not found in response", tt.wantCookieName)
				}
			} else if len(cookies) > 0 {
				for _, c := range cookies {
					if c.Name == "refresh_token" && c.Value != "" {
						t.Errorf("SignUpHandler expected no refresh_token cookie, but got: %s", c.Value)
					}
				}
			}

			if tt.wantAutoLogin != nil {
				var resp dto.SignUpResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}
				if resp.AutoLogin != *tt.wantAutoLogin {
					t.Errorf("SignUpHandler response AutoLogin = %v, want %v", resp.AutoLogin, *tt.wantAutoLogin)
				}
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return new(b)
}
