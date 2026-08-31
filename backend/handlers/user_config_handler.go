package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
)

type UserConfigService interface {
	ChangePassword(ctx context.Context, userIdContext, refreshTokenString string, data dto.ChangePasswordRequest) error
	DefinePassword(ctx context.Context, userIdContext, refreshTokenString string, data dto.DefinePasswordRequest) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, data dto.ResetPasswordRequest) error
	GetMyCredentials(ctx context.Context, id string) (*dto.GetMyCredentialsResponse, error)
}

type UserConfigHandler struct {
	service UserConfigService
}

func NewUserConfigHandler(service UserConfigService) *UserConfigHandler {
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
	id, _, _, _, ok := GetAuthContext(c)

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

func (h *UserConfigHandler) getAuthCredentials(c *gin.Context) (string, string, bool) {
	id, _, _, _, ok := GetAuthContext(c)

	if !ok {
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

	return id, refreshToken, true
}
