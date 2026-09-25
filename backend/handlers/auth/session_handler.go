package authhandler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
)

type sessionService interface {
	RefreshToken(ctx context.Context, refreshTokenString, deviceHash string) (auth.RefreshTokenResult, error)
	Logout(ctx context.Context, userId, refreshTokenString, accessJti string, accessExpiresAt time.Time) error
	LogoutOtherDevices(ctx context.Context, userId, refreshTokenString, accessJti, deviceHash string) error
	LogoutAllDevices(ctx context.Context, userId, refreshTokenString, deviceHash string) error
}

type SessionHandler struct {
	service             sessionService
	refreshDurationDays int
	isProduction        bool
	appDomain           string
}

func NewSessionHandler(sessionService sessionService, refreshDurationDays int, isProduction bool, appDomain string) *SessionHandler {
	return &SessionHandler{
		service:             sessionService,
		refreshDurationDays: refreshDurationDays,
		isProduction:        isProduction,
		appDomain:           appDomain,
	}
}

func (h *SessionHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
			dto.CodeMissingCookie,
			dto.MsgMissingRefreshCookie,
		))
		return
	}

	deviceHash := ExtractDeviceHash(c.Request)
	result, err := h.service.RefreshToken(c.Request.Context(), refreshToken, deviceHash)
	if err != nil {
		handlers.HandleAPIError(c, err)
		return
	}

	SetRefreshCookie(c, result.RefreshToken, h.appDomain, h.refreshDurationDays, h.isProduction)
	c.JSON(http.StatusOK, &dto.RefreshResponse{
		AccessToken: result.AccessToken,
	})
}

func (h *SessionHandler) Logout(c *gin.Context) {
	id, _, accessJti, accessExpiresAt, ok := handlers.GetAuthContext(c)
	if !ok {
		return
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
			dto.CodeMissingCookie,
			dto.MsgMissingRefreshCookie,
		))
		return
	}

	err = h.service.Logout(c.Request.Context(), id, refreshToken, accessJti, accessExpiresAt)
	if err != nil {
		handlers.HandleAPIError(c, err)
		return
	}

	ClearCookie(c, "refresh_token", h.appDomain, h.isProduction)
	c.Status(http.StatusNoContent)
}

func (h *SessionHandler) LogoutOtherDevices(c *gin.Context) {
	id, _, accessJti, _, ok := handlers.GetAuthContext(c)
	if !ok {
		return
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
			dto.CodeMissingCookie,
			dto.MsgMissingRefreshCookie,
		))
		return
	}

	dvh := ExtractDeviceHash(c.Request)
	err = h.service.LogoutOtherDevices(c.Request.Context(), id, refreshToken, accessJti, dvh)
	if err != nil {
		handlers.HandleAPIError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SessionHandler) LogoutAllDevices(c *gin.Context) {
	id, _, _, _, ok := handlers.GetAuthContext(c)
	if !ok {
		return
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
			dto.CodeMissingCookie,
			dto.MsgMissingRefreshCookie,
		))
		return
	}

	dvh := ExtractDeviceHash(c.Request)
	err = h.service.LogoutAllDevices(c.Request.Context(), id, refreshToken, dvh)
	if err != nil {
		handlers.HandleAPIError(c, err)
		return
	}

	ClearCookie(c, "refresh_token", h.appDomain, h.isProduction)
	c.Status(http.StatusNoContent)
}
