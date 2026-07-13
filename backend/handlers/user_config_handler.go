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
	userIdContext, existsId := c.Get("userId")
	if !existsId {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeMissingUserIdContext,
			dto.MsgMissingDataFromContext,
		))
		return
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
			dto.CodeMissingCookie,
			dto.MsgMissingRefreshCookie,
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

	var requestBody dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.CodeInvalidRequestBody,
			err.Error(),
		))
		return
	}

	err = h.service.ChangePassword(c.Request.Context(), id, refreshToken, requestBody)
	if errors.Is(err, models.ErrPasswordChangeButNotLogout) {
		c.JSON(http.StatusOK, dto.ChangePasswordResponse{
			Code:                   "PASSWORD_CHANGED_WITH_WARNING",
			Message:                "Password changed successfully, but we couldn't terminate other active sessions.",
			OthersDevicesLoggedOut: false,
		})
		return
	}

	if err != nil {
		if appErr, ok := errors.AsType[*models.AppError](err); ok {
			c.JSON(appErr.StatusCode, dto.NewErrorResponse(appErr.Code, appErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeInternalError, dto.MsgUnexpectedError))
		return
	}

	c.JSON(http.StatusOK, dto.ChangePasswordResponse{
		Code:                   "PASSWORD_CHANGED_SUCCESS",
		Message:                "Password changed successfully and all other sessions were terminated.",
		OthersDevicesLoggedOut: true,
	})
}
