package middleware

import (
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
)

func LimitBodySize(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

func LimitQueryParamsSize(maxChars int) gin.HandlerFunc {
	return func(c *gin.Context) {
		fullURL := c.Request.URL.RequestURI()

		if utf8.RuneCountInString(fullURL) > maxChars {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
				dto.CodeUrlTooLarge,
				dto.MsgUrlTooLarge,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
