package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
)

func HandleAPIError(c *gin.Context, err error) {
	if appErr, ok := errors.AsType[*models.AppError](err); ok {
		c.JSON(appErr.StatusCode, dto.NewErrorResponse(appErr.Code, appErr.Message))
		return
	}

	c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
		dto.CodeInternalError,
		dto.MsgUnexpectedError,
	))
}

func GetAuthContext(c *gin.Context) (userId, email, accessId string, accessExpiresAt time.Time, ok bool) {
	userId = c.GetString("userId")
	accessId = c.GetString("accessId")
	accessExpiresAt = c.GetTime("accessExpiresAt")
	email = c.GetString("email")

	if userId == "" || email == "" || accessId == "" || accessExpiresAt.IsZero() {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeMissingUserIdContext,
			dto.MsgMissingDataFromContext,
		))
		c.Abort()
		return "", "", "", time.Time{}, false
	}

	return userId, email, accessId, accessExpiresAt, true
}

func SetRefreshCookie(c *gin.Context, value, domain string, refreshDurationDays int, isProduction bool) {
	maxAge := 60 * 60 * 24 * refreshDurationDays

	c.SetCookieData(&http.Cookie{
		Name:     "refresh_token",
		Value:    value,
		MaxAge:   maxAge,
		Expires:  time.Now().Add(time.Duration(maxAge) * time.Second),
		Path:     "",
		Domain:   domain,
		Secure:   isProduction,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
