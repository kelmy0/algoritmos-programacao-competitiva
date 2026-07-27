package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func CacheControl(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d, stale-while-revalidate=600", int(duration.Seconds())))
		c.Next()
	}
}
