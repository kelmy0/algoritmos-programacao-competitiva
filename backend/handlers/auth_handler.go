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
		c.SetCookie("refresh_token", result.RefreshToken, 60*60*24*h.refreshDurationDays, "/", h.appDomain, h.isProduce, true)
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

	result, err := h.service.VerifyLogin2FA(c.Request.Context(), requestBody)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.SetCookie("refresh_token", result.RefreshToken, 60*60*24*h.refreshDurationDays, "/", h.appDomain, h.isProduce, true)
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

	accessToken, err := h.service.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, &dto.RefreshResponse{
		AccessToken: accessToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeMissingUserIdContext,
			dto.MsgMissingDataFromContext,
		))
		return
	}

	id, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeInternalError,
			dto.MsgUnexpectedError,
		))
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

	err = h.service.Logout(c.Request.Context(), id, refreshToken)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", h.appDomain, h.isProduce, true)
	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully logged out.",
	})
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeMissingUserIdContext,
			dto.MsgMissingDataFromContext,
		))
		return
	}

	id, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeInternalError,
			dto.MsgUnexpectedError,
		))
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

	err = h.service.LogoutAll(c.Request.Context(), id, refreshToken)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully logged out from all devices.",
	})
}
