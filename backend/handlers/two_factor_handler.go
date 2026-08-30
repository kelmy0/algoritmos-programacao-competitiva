package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

type TwoFactorService interface {
	Generate2FA(ctx context.Context, data dto.TwoFactorGenerateRequest) (dto.TwoFactorGenerateResponse, error)
	Enable2FA(ctx context.Context, data dto.TwoFactorEnableRequest) (services.Enable2FAResult, error)
	Disable2FA(ctx context.Context, data dto.TwoFactorDisableRequest) error
}

type TwoFactorHandler struct {
	service             TwoFactorService
	isProduction        bool
	appDomain           string
	refreshDurationDays int
}

func NewTwoFactorHandler(service TwoFactorService, isProduction bool, appDomain string, refreshDuration int) *TwoFactorHandler {
	return &TwoFactorHandler{
		service:             service,
		isProduction:        isProduction,
		appDomain:           appDomain,
		refreshDurationDays: refreshDuration,
	}
}

func (h *TwoFactorHandler) Generate2FA(c *gin.Context) {
	id, email, _, _, ok := GetAuthContext(c)

	if !ok {
		return
	}

	var requestBody dto.TwoFactorGenerateRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	requestBody.Email = email
	requestBody.UserId = id

	response, err := h.service.Generate2FA(c.Request.Context(), requestBody)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TwoFactorHandler) Enable2FA(c *gin.Context) {
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

	var requestBody dto.TwoFactorEnableRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			"Code must be 6 digits.",
		))
		return
	}

	requestBody.DeviceHash = ExtractDeviceHash(c.Request)
	requestBody.RefreshToken = refreshToken
	requestBody.UserId = id

	result, err := h.service.Enable2FA(c.Request.Context(), requestBody)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	SetRefreshCookie(c, result.RefreshToken, h.appDomain, h.refreshDurationDays, h.isProduction)
	c.JSON(http.StatusOK, &dto.TwoFactorEnableResponse{
		AccessToken: result.AccessToken,
	})
}

func (h *TwoFactorHandler) Disable2FA(c *gin.Context) {
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

	var requestBody dto.TwoFactorDisableRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	requestBody.DeviceHash = ExtractDeviceHash(c.Request)
	requestBody.RefreshToken = refreshToken
	requestBody.UserId = id

	err = h.service.Disable2FA(c.Request.Context(), requestBody)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	SetRefreshCookie(c, "", h.appDomain, -1, h.isProduction)
	c.Status(http.StatusNoContent)
}
