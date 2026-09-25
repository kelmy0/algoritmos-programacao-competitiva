package authhandler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
)

type twoFactorService interface {
	VerifyLogin2FA(ctx context.Context, data dto.Verify2FARequest) (auth.AuthResult, error)
}

type TwoFactorHandler struct {
	service             twoFactorService
	refreshDurationDays int
	isProduction        bool
	appDomain           string
}

func NewTwoFactorHandler(twoFactorService twoFactorService, refreshDurationDays int,
	isProduction bool, appDomain string) *TwoFactorHandler {
	return &TwoFactorHandler{
		service:             twoFactorService,
		refreshDurationDays: refreshDurationDays,
		isProduction:        isProduction,
		appDomain:           appDomain,
	}
}

func (h *TwoFactorHandler) Verify2FA(c *gin.Context) {
	preAuthToken, err := c.Cookie("pre_auth_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
			dto.CodeMissingCookie,
			dto.MsgMissingRefreshCookie,
		))
		return
	}

	var requestBody dto.Verify2FARequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	requestBody.PreAuthToken = preAuthToken
	requestBody.DeviceHash = ExtractDeviceHash(c.Request)
	result, err := h.service.VerifyLogin2FA(c.Request.Context(), requestBody)
	if err != nil {
		handlers.HandleAPIError(c, err)
		return
	}

	SetRefreshCookie(c, result.RefreshToken, h.appDomain, h.refreshDurationDays, h.isProduction)
	c.JSON(http.StatusOK, result.LoginResponse)
}
