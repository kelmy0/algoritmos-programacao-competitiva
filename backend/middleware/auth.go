package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(secretKey, issuer string, redisClient *redis.Client) gin.HandlerFunc {
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

		if claims.IssuedAt == nil || claims.IssuedAt.Time.IsZero() {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				dto.CodeInvalidAccessToken,
				"Token is missing issued at (iat) claim.",
			))
			c.Abort()
			return
		}

		ctx := c.Request.Context()

		if claims.ID != "" {
			blacklisted, err := redisClient.Exists(ctx, "blacklist:jti:"+claims.ID).Result()
			if err != nil {
				c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
					dto.CodeInternalError,
					"Failed to verify token revocation status.",
				))
				c.Abort()
				return
			}

			if blacklisted > 0 {
				c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
					dto.CodeInvalidAccessToken,
					"Token has been revoked.",
				))
				c.Abort()
				return
			}
		}

		logoutAllKey := "logout_all:" + claims.Subject
		val, err := redisClient.Get(ctx, logoutAllKey).Result()

		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
				dto.CodeInternalError,
				"Failed to verify session status.",
			))
			c.Abort()
			return
		}

		if err == nil && val != "" {
			parts := strings.Split(val, ":")
			logoutTimestamp, _ := strconv.ParseInt(parts[0], 10, 64)

			allowedAccessJti := ""
			if len(parts) > 1 {
				allowedAccessJti = parts[1]
			}

			if claims.IssuedAt.Time.Unix() <= logoutTimestamp && claims.ID != allowedAccessJti {
				c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
					dto.CodeTokenNolongerValid,
					"Session expired due to logout on other devices.",
				))
				c.Abort()
				return
			}
		}

		c.Set("userId", claims.Subject)
		c.Set("permissions", claims.Permissions)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("isEmployee", claims.IsEmployee)
		c.Set("accessJti", claims.ID)
		c.Set("accessExpiresAt", claims.ExpiresAt.Time)
		c.Next()
	}
}
