package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

type AuthHandler struct {
	service             *services.AuthService
	isProduce           bool
	appDomain           string
	refreshDurationDays int
}

func NewAuthHandler(service *services.AuthService, isProduce bool, appDomain string, refreshDurationDays int) *AuthHandler {
	return &AuthHandler{service: service, isProduce: isProduce, appDomain: appDomain, refreshDurationDays: refreshDurationDays}
}

func (h *AuthHandler) Auth(c *gin.Context) {
	var requestBody dto.AuthRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	requestBody.DeviceHash = ExtractDeviceHash(c.Request)
	result, err := h.service.Auth(c.Request.Context(), requestBody)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	if result.LoginResponse.Requires2FA {
		c.JSON(http.StatusOK, result.LoginResponse)
		return
	}

	if result.RefreshToken != "" {
		SetRefreshCookie(c, result.RefreshToken, h.appDomain, h.refreshDurationDays, h.isProduce)
	}
	c.JSON(http.StatusOK, result.LoginResponse)
}

func (h *AuthHandler) Verify2FA(c *gin.Context) {
	var requestBody dto.Verify2FARequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	requestBody.DeviceHash = ExtractDeviceHash(c.Request)
	result, err := h.service.VerifyLogin2FA(c.Request.Context(), requestBody)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	SetRefreshCookie(c, result.RefreshToken, h.appDomain, h.refreshDurationDays, h.isProduce)
	c.JSON(http.StatusOK, result.LoginResponse)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
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
		HandleAPIError(c, err)
		return
	}

	SetRefreshCookie(c, result.RefreshToken, h.appDomain, h.refreshDurationDays, h.isProduce)
	c.JSON(http.StatusOK, &dto.RefreshResponse{
		AccessToken: result.AccessToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	id, _, accessJti, accessExpiresAt, ok := GetAuthContext(c)
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
		HandleAPIError(c, err)
		return
	}

	ClearCookie(c, "refresh_token", h.appDomain, h.isProduce)
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) LogoutOtherDevices(c *gin.Context) {
	id, _, accessJti, _, ok := GetAuthContext(c)
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
		HandleAPIError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) LogoutAllDevices(c *gin.Context) {
	id, _, _, _, ok := GetAuthContext(c)
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
		HandleAPIError(c, err)
		return
	}

	ClearCookie(c, "refresh_token", h.appDomain, h.isProduce)
	c.Status(http.StatusNoContent)
}
