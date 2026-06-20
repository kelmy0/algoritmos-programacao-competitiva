package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Fake404Middleware(expectedHash string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientHash := c.Query("secret")
		if clientHash != expectedHash {
			c.JSON(http.StatusNotFound, gin.H{"error": "page not found."})
			c.Abort()
			return
		}
		c.Next()
	}
}
