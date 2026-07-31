package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
)

func SitemapMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := c.GetHeader("x-sitemap-secret")

		if subtle.ConstantTimeCompare([]byte(secret), []byte(secretKey)) != 1 {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				dto.CodeInvalidSecret,
				"Secrets don't matches.",
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
