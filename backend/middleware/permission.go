package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
)

func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxPermissions, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, dto.NewErrorResponse(
				dto.CodeNoPermission,
				dto.MsgNoPermission,
			))
			c.Abort()
			return
		}

		permissions, ok := ctxPermissions.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
				dto.CodeInternalError,
				dto.MsgUnexpectedError,
			))
			c.Abort()
			return
		}

		hasPermission := slices.Contains(permissions, requiredPermission)

		if !hasPermission {
			c.JSON(http.StatusForbidden, dto.NewErrorResponse(
				dto.CodeNoPermission,
				dto.MsgNoPermission,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
