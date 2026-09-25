package authhandler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
)

type loginService interface {
	Login(ctx context.Context, data dto.AuthRequest) (auth.AuthResult, error)
}

type LoginHandler struct {
	service             loginService
	refreshDurationDays int
	isProduction        bool
	appDomain           string
}

func NewLoginHandler(loginService loginService, refreshDurationDays int,
	isProduction bool, appDomain string) *LoginHandler {
	return &LoginHandler{
		service:             loginService,
		refreshDurationDays: refreshDurationDays,
		isProduction:        isProduction,
		appDomain:           appDomain,
	}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var requestBody dto.AuthRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	requestBody.DeviceHash = ExtractDeviceHash(c.Request)
	result, err := h.service.Login(c.Request.Context(), requestBody)
	if err != nil {
		handlers.HandleAPIError(c, err)
		return
	}

	if result.LoginResponse.Requires2FA {
		c.JSON(http.StatusOK, result.LoginResponse)
		return
	}

	if result.RefreshToken != "" {
		SetRefreshCookie(c, result.RefreshToken, h.appDomain, h.refreshDurationDays, h.isProduction)
	}
	c.JSON(http.StatusOK, result.LoginResponse)
}
