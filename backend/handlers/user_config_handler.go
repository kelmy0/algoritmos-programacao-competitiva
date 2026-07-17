package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
)

type UserConfigHandler struct {
	service *services.UserConfigService
}

func NewUserConfigHandler(service *services.UserConfigService) *UserConfigHandler {
	return &UserConfigHandler{service: service}
}

func (h *UserConfigHandler) ChangePassword(c *gin.Context) {
	id, refreshToken, ok := h.getAuthCredentials(c)
	if !ok {
		return
	}

	var requestBody dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	err := h.service.ChangePassword(c.Request.Context(), id, refreshToken, requestBody)
	if errors.Is(err, models.ErrPasswordChangeButNotLogout) {
		c.JSON(http.StatusOK, dto.ChangePasswordResponse{
			Code:                   "PASSWORD_CHANGED_WITH_WARNING",
			Message:                "Password changed successfully, but we couldn't terminate other active sessions.",
			OthersDevicesLoggedOut: false,
		})
		return
	}

	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ChangePasswordResponse{
		Code:                   "PASSWORD_CHANGED_SUCCESS",
		Message:                "Password changed successfully and all other sessions were terminated.",
		OthersDevicesLoggedOut: true,
	})
}

func (h *UserConfigHandler) DefinePassword(c *gin.Context) {
	id, refreshToken, ok := h.getAuthCredentials(c)
	if !ok {
		return
	}

	var requestBody dto.DefinePasswordRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	err := h.service.DefinePassword(c.Request.Context(), id, refreshToken, requestBody)
	if errors.Is(err, models.ErrPasswordSetButNotLogout) {
		c.JSON(http.StatusOK, dto.ChangePasswordResponse{
			Code:                   "PASSWORD_SETTED_WITH_WARNING",
			Message:                "Password set successfully, but we couldn't terminate other active sessions.",
			OthersDevicesLoggedOut: false,
		})
		return
	}

	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ChangePasswordResponse{
		Code:                   "PASSWORD_SETTED_SUCCESS",
		Message:                "Password set successfully and all other sessions were terminated.",
		OthersDevicesLoggedOut: true,
	})
}

func (h *UserConfigHandler) ForgotPassword(c *gin.Context) {
	var requestBody dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	err := h.service.ForgotPassword(c.Request.Context(), requestBody.Email)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Recovery link sent to email.",
	})
}

func (h *UserConfigHandler) ResetPassword(c *gin.Context) {
	var requestBody dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), requestBody)
	if err != nil {
		HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Password recovered!",
	})
}

func (h *UserConfigHandler) GetMyCredentials(c *gin.Context) {
	id, ok := h.getCredentials(c)
	if !ok {
		return
	}

	user, err := h.service.GetMyCredentials(c.Request.Context(), id)
	if err != nil {
		HandleAPIError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserConfigHandler) getCredentials(c *gin.Context) (string, bool) {
	userIdContext, existsId := c.Get("userId")
	if !existsId {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeMissingUserIdContext,
			dto.MsgMissingDataFromContext,
		))
		return "", false
	}

	id, ok := userIdContext.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeInternalError,
			dto.MsgUnexpectedError,
		))
		return "", false
	}

	return id, true
}

func (h *UserConfigHandler) getAuthCredentials(c *gin.Context) (string, string, bool) {
	id, success := h.getCredentials(c)
	if !success {
		return "", "", false
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
			dto.CodeMissingCookie,
			dto.MsgMissingRefreshCookie,
		))
		return "", "", false
	}

	return id, refreshToken, success
}
