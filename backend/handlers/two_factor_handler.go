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
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeMissingUserIdContext,
			dto.MsgMissingDataFromContext,
		))
		return
	}

	id, okId := userIdContext.(string)
	email, okEmail := userEmailContext.(string)
	if !okId || !okEmail {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeInternalError,
			dto.MsgUnexpectedError,
		))
		return
	}

	response, err := h.service.Generate2FA(c.Request.Context(), id, email)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TwoFactorHandler) Enable2FA(c *gin.Context) {
	userIdContext, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeMissingUserIdContext,
			dto.MsgMissingDataFromContext,
		))
		return
	}

	id, ok := userIdContext.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeInternalError,
			dto.MsgUnexpectedError,
		))
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

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Two-factor authentication enabled successfully",
	})
}

func (h *TwoFactorHandler) Disable2FA(c *gin.Context) {
	userIdContext, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeMissingUserIdContext,
			dto.MsgMissingDataFromContext,
		))
		return
	}

	id, ok := userIdContext.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeInternalError,
			dto.MsgUnexpectedError,
		))
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
