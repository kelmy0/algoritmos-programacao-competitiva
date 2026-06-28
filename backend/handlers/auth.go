package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

type AuthHandler struct {
	service   *services.AuthService
	isProduce bool
}

func NewAuthHandler(service *services.AuthService, isProduce bool) *AuthHandler {
	return &AuthHandler{service: service, isProduce: isProduce}
}

func (h *AuthHandler) Auth(c *gin.Context) {
	var requestBody dto.AuthRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loginResponse, refreshToken, refreshDuration, err := h.service.Auth(c.Request.Context(), requestBody)

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

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token cookie is required"})
		c.Abort()
		return
	}

	accessToken, err := h.service.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
