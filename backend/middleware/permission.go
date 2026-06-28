package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxPermissions, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: no permissions assigned to this user"})
			c.Abort()
			return
		}

		permissions, ok := ctxPermissions.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error processing permissions"})
			c.Abort()
			return
		}

		hasPermission := false
		for _, p := range permissions {
			if p == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: you do not have the required permission"})
			c.Abort()
			return
		}

		c.Next()
	}
}
