package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

func RequireCaptcha(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if secretKey == "" {
			c.Next()
			return
		}

		token := c.GetHeader("X-CF-Turnstile-Response")

		if token == "" {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
				dto.CodeMissingHeader,
				"Captcha header is missing.",
			))
			c.Abort()
			return
		}

		clientIP := c.ClientIP()
		valid, err := utils.VerifyTurnstile(secretKey, token, clientIP)

		if err != nil || !valid {
			c.JSON(http.StatusForbidden, dto.NewErrorResponse(
				dto.CodeInvalidCaptcha,
				"Captcha is not valid. Try again!",
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
