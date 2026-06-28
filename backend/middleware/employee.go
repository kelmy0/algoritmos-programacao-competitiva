package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func EmployeeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxIsEmployee, exists := c.Get("isEmployee")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: access denied"})
			c.Abort()
			return
		}

		isEmployee, ok := ctxIsEmployee.(bool)
		if !ok || !isEmployee {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: this area is restricted to employees only"})
			c.Abort()
			return
		}

		c.Next()
	}
}
