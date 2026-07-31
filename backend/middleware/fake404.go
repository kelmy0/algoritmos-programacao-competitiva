package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
)

func Fake404Middleware(expectedHash string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientHash := c.GetHeader("x-admin-secret")
		if subtle.ConstantTimeCompare([]byte(clientHash), []byte(expectedHash)) != 1 {
			c.JSON(http.StatusNotFound, dto.NewErrorResponse(
				dto.CodePageNotFound,
				dto.MsgPageNotFound,
			))
			c.Abort()
			return
		}
		c.Next()
	}
}
