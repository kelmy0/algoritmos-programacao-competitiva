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
	id, email, _, _, ok := GetAuthContext(c)

	if !ok {
		return
	}

	var requestBody dto.TwoFactorGenerateRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	response, err := h.service.Generate2FA(c.Request.Context(), id, email, requestBody.Password)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TwoFactorHandler) Enable2FA(c *gin.Context) {
	id, _, _, _, ok := GetAuthContext(c)

	if !ok {
		return
	}

	var req dto.TwoFactorEnableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			"Code must be 6 digits.",
		))
		return
	}
	err := h.service.Enable2FA(c.Request.Context(), id, req.Code)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TwoFactorHandler) Disable2FA(c *gin.Context) {
	id, _, _, _, ok := GetAuthContext(c)

	if !ok {
		return
	}

	var req dto.TwoFactorDisableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}
	err := h.service.Disable2FA(c.Request.Context(), id, req.Password)
	if err != nil {
		HandleAPIError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Two-factor authentication disabled successfully",
	})
}
