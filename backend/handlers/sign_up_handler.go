package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

type SignUpHandler struct {
	Service             *services.SignUpService
	RefreshDurationDays int
	AppDomain           string
	IsProduce           bool
}

func NewSignUpHandler(service *services.SignUpService, refreshDurationDays int, appDomain string, isProduce bool) *SignUpHandler {
	return &SignUpHandler{
		Service:             service,
		RefreshDurationDays: refreshDurationDays,
		AppDomain:           appDomain,
		IsProduce:           isProduce,
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
	result, err := h.Service.SignUp(c.Request.Context(), requestBody)

	if err != nil {
		HandleAPIError(c, err)
		return
	}

	SetRefreshCookie(c, result.RefreshToken, h.AppDomain, h.RefreshDurationDays, h.IsProduce)
	c.JSON(http.StatusOK, result.SignUpResponse)
}
