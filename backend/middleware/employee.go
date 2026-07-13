package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
)

func EmployeeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxIsEmployee, exists := c.Get("isEmployee")
		if !exists {
			c.JSON(http.StatusForbidden, dto.NewErrorResponse(
				dto.CodeRestrictedArea,
				dto.MsgRestrictedArea,
			))
			c.Abort()
			return
		}

		isEmployee, ok := ctxIsEmployee.(bool)
		if !ok || !isEmployee {
			c.JSON(http.StatusForbidden, dto.NewErrorResponse(
				dto.CodeRestrictedArea,
				dto.MsgRestrictedArea,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
