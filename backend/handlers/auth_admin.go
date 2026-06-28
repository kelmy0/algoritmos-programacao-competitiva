package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

type AuthAdminHandler struct {
	service   *services.AuthService
	isProduce bool
}

func NewAuthAdminHandler(service *services.AuthService, isProduce bool) *AuthAdminHandler {
	return &AuthAdminHandler{service: service, isProduce: isProduce}
}

func (h *AuthAdminHandler) Auth(c *gin.Context) {
	var requestBody dto.AuthRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loginResponse, refreshToken, refreshDuration, err := h.service.Auth(c.Request.Context(), requestBody, true)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(
		"refresh_token",
		refreshToken,
		60*60*24*refreshDuration, // Age in seconds converted to days
		"/",
		"",
		h.isProduce, // Secure Https only
		true,
	)

	c.JSON(http.StatusOK, loginResponse)
}
