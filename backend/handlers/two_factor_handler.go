package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

type TwoFactorHandler struct {
	service *services.TwoFactorService
}

func NewTwoFactorHandler(service *services.TwoFactorService) *TwoFactorHandler {
	return &TwoFactorHandler{service: service}
}

func (h *TwoFactorHandler) Generate2FA(c *gin.Context) {
	userIdContext, existsId := c.Get("userId")
	userEmailContext, existsEmail := c.Get("email")

	if !existsId || !existsEmail {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication data missing from the context"})
		return
	}

	id, okId := userIdContext.(string)
	email, okEmail := userEmailContext.(string)
	if !okId || !okEmail {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing user data"})
		return
	}

	response, err := h.service.Generate2FA(c.Request.Context(), id, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TwoFactorHandler) Enable2FA(c *gin.Context) {
	userIdContext, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication data missing"})
		return
	}

	id, ok := userIdContext.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error processing user ID"})
		return
	}

	var req dto.TwoFactorEnableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload. Code must be 6 digits."})
		return
	}
	err := h.service.Enable2FA(c.Request.Context(), id, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Two-factor authentication enabled successfully"})
}

func (h *TwoFactorHandler) Disable2FA(c *gin.Context) {}
