package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
)

func GetAuthContext(c *gin.Context) (userId, accessJti string, accessExpiresAt time.Time, ok bool) {
	userId = c.GetString("userId")
	accessJti = c.GetString("accessJti")
	accessExpiresAt = c.GetTime("accessExpiresAt")

	if userId == "" || accessJti == "" || accessExpiresAt.IsZero() {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeMissingUserIdContext,
			dto.MsgMissingDataFromContext,
		))
		c.Abort()
		return "", "", time.Time{}, false
	}

	return userId, accessJti, accessExpiresAt, true
}
