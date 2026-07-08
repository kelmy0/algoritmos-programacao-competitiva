package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

type TwoFactorHandler struct {
	service *services.TwoFactorService
}

func NewTwoFactorHandler(service *services.TwoFactorService) *TwoFactorHandler {
	return &TwoFactorHandler{service: service}
}

func (h *TwoFactorHandler) Generate2FA(c *gin.Context) {}

func (h *TwoFactorHandler) Enable2FA(c *gin.Context) {}

func (h *TwoFactorHandler) Disable2FA(c *gin.Context) {}
