package handlers

import (
	"errors"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.Service.SignUp(c.Request.Context(), requestBody)
	if err != nil {
		if errors.Is(err, services.ErrAccountCreatedButTokenFailed) {
			c.JSON(http.StatusCreated, dto.SignUpResponse{
				Success:   true,
				AutoLogin: false,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refresh_token", result.RefreshToken, 60*60*24*h.RefreshDurationDays, "/", h.AppDomain, h.IsProduce, true)
	c.JSON(http.StatusOK, result.SignUpResponse)
}
