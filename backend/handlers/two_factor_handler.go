package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TwoFactorHandler) Enable2FA(c *gin.Context) {}

func (h *TwoFactorHandler) Disable2FA(c *gin.Context) {}
