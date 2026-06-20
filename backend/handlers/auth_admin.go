package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

type AuthAdminHandler struct {
	service *services.AuthAdminService
}

func NewAuthAdminHandler(service *services.AuthAdminService) *AuthAdminHandler {
	return &AuthAdminHandler{service: service}
}

func (h *AuthAdminHandler) Auth(c *gin.Context) {
	var requestBody dto.AuthAdminRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	admin, err := h.service.AuthAdmin(c.Request.Context(), requestBody)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": admin,
	})
}
