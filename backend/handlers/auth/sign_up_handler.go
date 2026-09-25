package authhandler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/handlers"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
)

type signUpService interface {
	SignUp(ctx context.Context, data dto.SignUpRequest) (result auth.SignUpResult, err error)
}

type SignUpHandler struct {
	service             signUpService
	refreshDurationDays int
	isProduction        bool
	appDomain           string
}

func NewSignUpHandler(service signUpService, refreshDurationDays int, isProduction bool, appDomain string) *SignUpHandler {
	return &SignUpHandler{
		service:             service,
		refreshDurationDays: refreshDurationDays,
		appDomain:           appDomain,
		isProduction:        isProduction,
	}
}

func (h *SignUpHandler) SignUp(c *gin.Context) {
	var requestBody dto.SignUpRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInternalError,
			err.Error(),
		))
		return
	}

	requestBody.DeviceHash = ExtractDeviceHash(c.Request)
	result, err := h.service.SignUp(c.Request.Context(), requestBody)

	if err != nil {
		handlers.HandleAPIError(c, err)
		return
	}

	if result.RefreshToken != "" {
		SetRefreshCookie(c, result.RefreshToken, h.appDomain, h.refreshDurationDays, h.isProduction)
	}

	c.JSON(http.StatusOK, result.SignUpResponse)
}
