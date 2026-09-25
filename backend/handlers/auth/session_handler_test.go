package authhandler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	authhandler "github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers/auth"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
)

type MockSessionService struct {
	RefreshTokenFunc       func(ctx context.Context, refreshTokenString, deviceHash string) (auth.RefreshTokenResult, error)
	LogoutFunc             func(ctx context.Context, userId, refreshTokenString, accessJti string, accessExpiresAt time.Time) error
	LogoutOtherDevicesFunc func(ctx context.Context, userId, refreshTokenString, accessJti, deviceHash string) error
	LogoutAllDevicesFunc   func(ctx context.Context, userId, refreshTokenString, deviceHash string) error
}

func (m *MockSessionService) RefreshToken(ctx context.Context, refreshTokenString, deviceHash string) (auth.RefreshTokenResult, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(ctx, refreshTokenString, deviceHash)
	}
	return auth.RefreshTokenResult{}, nil
}

func (m *MockSessionService) Logout(ctx context.Context, userId, refreshTokenString, accessJti string, accessExpiresAt time.Time) error {
	if m.LogoutFunc != nil {
		return m.LogoutFunc(ctx, userId, refreshTokenString, accessJti, accessExpiresAt)
	}
	return nil
}

func (m *MockSessionService) LogoutOtherDevices(ctx context.Context, userId, refreshTokenString, accessJti, deviceHash string) error {
	if m.LogoutOtherDevicesFunc != nil {
		return m.LogoutOtherDevicesFunc(ctx, userId, refreshTokenString, accessJti, deviceHash)
	}
	return nil
}

func (m *MockSessionService) LogoutAllDevices(ctx context.Context, userId, refreshTokenString, deviceHash string) error {
	if m.LogoutAllDevicesFunc != nil {
		return m.LogoutAllDevicesFunc(ctx, userId, refreshTokenString, deviceHash)
	}
	return nil
}

func setAuthContext(c *gin.Context, userID, email, accessID string, expiresAt time.Time) {
	c.Set("userId", userID)
	c.Set("email", email)
	c.Set("accessId", accessID)
	c.Set("accessExpiresAt", expiresAt)
}

func isCookieCleared(cookies []*http.Cookie, cookieName string) bool {
	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			return cookie.MaxAge < 0 || cookie.Value == ""
		}
	}
	return false
}

const (
	testAppDomain           = "localhost"
	testRefreshDurationDays = 7
	testIsProduction        = false
	testRefreshToken        = "valid_refresh_token_123"
	testUserID              = "usr_999"
	testEmail               = "usuario@exemplo.com"
	testAccessID            = "jti_abc_123"
)

func TestSessionHandler_Refresh(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		hasRefreshCookie bool
		refreshTokenFn   func(ctx context.Context, refreshTokenString, deviceHash string) (auth.RefreshTokenResult, error)
		wantStatusCode   int
		wantNewCookie    bool
	}{
		{
			name:             "Error: Missing refresh_token cookie returns 401",
			hasRefreshCookie: false,
			refreshTokenFn:   nil,
			wantStatusCode:   http.StatusUnauthorized,
			wantNewCookie:    false,
		},
		{
			name:             "Error: Service returns AppError mapped status code",
			hasRefreshCookie: true,
			refreshTokenFn: func(ctx context.Context, refreshTokenString, deviceHash string) (auth.RefreshTokenResult, error) {
				return auth.RefreshTokenResult{}, &models.AppError{
					StatusCode: http.StatusUnauthorized,
					Code:       "TOKEN_EXPIRED",
					Message:    "Refresh token expired",
				}
			},
			wantStatusCode: http.StatusUnauthorized,
			wantNewCookie:  false,
		},
		{
			name:             "Success: Refresh token rotated successfully returns 200 OK and sets new cookie",
			hasRefreshCookie: true,
			refreshTokenFn: func(ctx context.Context, refreshTokenString, deviceHash string) (auth.RefreshTokenResult, error) {
				return auth.RefreshTokenResult{
					AccessToken:  "new_access_token_456",
					RefreshToken: "new_refresh_token_789",
				}, nil
			},
			wantStatusCode: http.StatusOK,
			wantNewCookie:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSessionService{RefreshTokenFunc: tt.refreshTokenFn}
			handler := authhandler.NewSessionHandler(mockService, testRefreshDurationDays, testIsProduction, testAppDomain)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
			if tt.hasRefreshCookie {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: testRefreshToken})
			}
			c.Request = req

			handler.Refresh(c)

			if w.Code != tt.wantStatusCode {
				t.Errorf("Refresh() statusCode = %d, want %d. Body: %s", w.Code, tt.wantStatusCode, w.Body.String())
			}

			hasCookie := false
			for _, ck := range w.Result().Cookies() {
				if ck.Name == "refresh_token" && ck.Value != "" && ck.MaxAge >= 0 {
					hasCookie = true
				}
			}

			if hasCookie != tt.wantNewCookie {
				t.Errorf("Refresh() new cookie presence = %v, want %v", hasCookie, tt.wantNewCookie)
			}
		})
	}
}

func TestSessionHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now().Add(1 * time.Hour)

	tests := []struct {
		name             string
		setAuthCtx       bool
		hasRefreshCookie bool
		logoutFn         func(ctx context.Context, userId, refreshTokenString, accessJti string, accessExpiresAt time.Time) error
		wantStatusCode   int
		wantCookieClear  bool
	}{
		{
			name:             "Error: Missing auth context aborts request with 500",
			setAuthCtx:       false,
			hasRefreshCookie: true,
			logoutFn:         nil,
			wantStatusCode:   http.StatusInternalServerError,
			wantCookieClear:  false,
		},
		{
			name:             "Error: Missing refresh_token cookie returns 401",
			setAuthCtx:       true,
			hasRefreshCookie: false,
			logoutFn:         nil,
			wantStatusCode:   http.StatusUnauthorized,
			wantCookieClear:  false,
		},
		{
			name:             "Error: Service error mapped via HandleAPIError",
			setAuthCtx:       true,
			hasRefreshCookie: true,
			logoutFn: func(ctx context.Context, userId, refreshTokenString, accessJti string, accessExpiresAt time.Time) error {
				return errors.New("redis error")
			},
			wantStatusCode:  http.StatusInternalServerError,
			wantCookieClear: false,
		},
		{
			name:             "Success: Logout revokes session, clears cookie and returns 204 No Content",
			setAuthCtx:       true,
			hasRefreshCookie: true,
			logoutFn: func(ctx context.Context, userId, refreshTokenString, accessJti string, accessExpiresAt time.Time) error {
				if userId != testUserID || accessJti != testAccessID {
					return errors.New("invalid arguments passed to service")
				}
				return nil
			},
			wantStatusCode:  http.StatusNoContent,
			wantCookieClear: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSessionService{LogoutFunc: tt.logoutFn}
			handler := authhandler.NewSessionHandler(mockService, testRefreshDurationDays, testIsProduction, testAppDomain)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setAuthCtx {
				setAuthContext(c, testUserID, testEmail, testAccessID, now)
			}

			req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
			if tt.hasRefreshCookie {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: testRefreshToken})
			}
			c.Request = req

			handler.Logout(c)

			statusCode := c.Writer.Status()

			if statusCode != tt.wantStatusCode {
				t.Errorf("Logout() statusCode = %d, want %d. Body: %s", statusCode, tt.wantStatusCode, w.Body.String())
			}

			cleared := isCookieCleared(w.Result().Cookies(), "refresh_token")
			if cleared != tt.wantCookieClear {
				t.Errorf("Logout() cookie cleared = %v, want %v", cleared, tt.wantCookieClear)
			}
		})
	}
}

func TestSessionHandler_LogoutOtherDevices(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                 string
		setAuthCtx           bool
		hasRefreshCookie     bool
		logoutOtherDevicesFn func(ctx context.Context, userId, refreshTokenString, accessJti, deviceHash string) error
		wantStatusCode       int
	}{
		{
			name:                 "Error: Missing auth context aborts request with 500",
			setAuthCtx:           false,
			hasRefreshCookie:     true,
			logoutOtherDevicesFn: nil,
			wantStatusCode:       http.StatusInternalServerError,
		},
		{
			name:                 "Error: Missing refresh_token cookie returns 401",
			setAuthCtx:           true,
			hasRefreshCookie:     false,
			logoutOtherDevicesFn: nil,
			wantStatusCode:       http.StatusUnauthorized,
		},
		{
			name:             "Error: Service error mapped via HandleAPIError",
			setAuthCtx:       true,
			hasRefreshCookie: true,
			logoutOtherDevicesFn: func(ctx context.Context, userId, refreshTokenString, accessJti, deviceHash string) error {
				return errors.New("db error")
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:             "Success: Revokes other sessions and returns 204 No Content",
			setAuthCtx:       true,
			hasRefreshCookie: true,
			logoutOtherDevicesFn: func(ctx context.Context, userId, refreshTokenString, accessJti, deviceHash string) error {
				return nil
			},
			wantStatusCode: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSessionService{LogoutOtherDevicesFunc: tt.logoutOtherDevicesFn}
			handler := authhandler.NewSessionHandler(mockService, testRefreshDurationDays, testIsProduction, testAppDomain)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setAuthCtx {
				setAuthContext(c, testUserID, testEmail, testAccessID, time.Now())
			}

			req, _ := http.NewRequest(http.MethodPost, "/auth/logout-others", nil)
			if tt.hasRefreshCookie {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: testRefreshToken})
			}
			c.Request = req

			handler.LogoutOtherDevices(c)

			statusCode := c.Writer.Status()

			if statusCode != tt.wantStatusCode {
				t.Errorf("LogoutOtherDevices() statusCode = %d, want %d. Body: %s", statusCode, tt.wantStatusCode, w.Body.String())
			}
		})
	}
}

func TestSessionHandler_LogoutAllDevices(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		setAuthCtx         bool
		hasRefreshCookie   bool
		logoutAllDevicesFn func(ctx context.Context, userId, refreshTokenString, deviceHash string) error
		wantStatusCode     int
		wantCookieClear    bool
	}{
		{
			name:               "Error: Missing auth context aborts request with 500",
			setAuthCtx:         false,
			hasRefreshCookie:   true,
			logoutAllDevicesFn: nil,
			wantStatusCode:     http.StatusInternalServerError,
			wantCookieClear:    false,
		},
		{
			name:               "Error: Missing refresh_token cookie returns 401",
			setAuthCtx:         true,
			hasRefreshCookie:   false,
			logoutAllDevicesFn: nil,
			wantStatusCode:     http.StatusUnauthorized,
			wantCookieClear:    false,
		},
		{
			name:             "Error: Service error mapped via HandleAPIError",
			setAuthCtx:       true,
			hasRefreshCookie: true,
			logoutAllDevicesFn: func(ctx context.Context, userId, refreshTokenString, deviceHash string) error {
				return errors.New("unexpected error")
			},
			wantStatusCode:  http.StatusInternalServerError,
			wantCookieClear: false,
		},
		{
			name:             "Success: Logout all devices revokes sessions, clears cookie and returns 204 No Content",
			setAuthCtx:       true,
			hasRefreshCookie: true,
			logoutAllDevicesFn: func(ctx context.Context, userId, refreshTokenString, deviceHash string) error {
				return nil
			},
			wantStatusCode:  http.StatusNoContent,
			wantCookieClear: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSessionService{LogoutAllDevicesFunc: tt.logoutAllDevicesFn}
			handler := authhandler.NewSessionHandler(mockService, testRefreshDurationDays, testIsProduction, testAppDomain)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setAuthCtx {
				setAuthContext(c, testUserID, testEmail, testAccessID, time.Now())
			}

			req, _ := http.NewRequest(http.MethodPost, "/auth/logout-all", nil)
			if tt.hasRefreshCookie {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: testRefreshToken})
			}
			c.Request = req

			handler.LogoutAllDevices(c)

			statusCode := c.Writer.Status()

			if statusCode != tt.wantStatusCode {
				t.Errorf("LogoutAllDevices() statusCode = %d, want %d. Body: %s", statusCode, tt.wantStatusCode, w.Body.String())
			}

			cleared := isCookieCleared(w.Result().Cookies(), "refresh_token")
			if cleared != tt.wantCookieClear {
				t.Errorf("LogoutAllDevices() cookie cleared = %v, want %v", cleared, tt.wantCookieClear)
			}
		})
	}
}
