package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

func AuthMiddleware(secretKey, issuer string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				dto.CodeMissingHeader,
				"Authorization header is required.",
			))
			c.Abort()
			return
		}

		// Verify bearer <token> format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				dto.CodeInvalidHeaderFormat,
				"Authorization header format must be Bearer {token}.",
			))
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString, secretKey, issuer)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				dto.CodeInvalidAccessToken,
				"Invalid or expired access token.",
			))
			c.Abort()
			return
		}

		c.Set("userId", claims.Subject)
		c.Set("permissions", claims.Permissions)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("isEmployee", claims.IsEmployee)
		c.Next()
	}
}
