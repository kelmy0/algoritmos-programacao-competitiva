package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	clients map[string]*client
	mu      sync.RWMutex
	r       rate.Limit
	b       int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	i := &RateLimiter{
		clients: make(map[string]*client),
		r:       r,
		b:       b,
	}

	go i.cleanupClients()

	return i
}

func (i *RateLimiter) getLimiter(key string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	v, exists := i.clients[key]
	if !exists {
		limiter := rate.NewLimiter(i.r, i.b)
		i.clients[key] = &client{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (i *RateLimiter) cleanupClients() {
	for {
		time.Sleep(10 * time.Minute)
		i.mu.Lock()
		for key, client := range i.clients {
			if time.Since(client.lastSeen) > 30*time.Minute {
				delete(i.clients, key)
			}
		}
		i.mu.Unlock()
	}
}

func RateLimitMiddleware(limiterManager *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		if userId, exists := c.Get("userId"); exists {
			if idStr, ok := userId.(string); ok && idStr != "" {
				key = "usr_" + idStr
			}
		}

		if key == "" {
			clientIP := c.GetHeader("X-Forwarded-For")
			if clientIP == "" {
				clientIP = c.GetHeader("X-Real-IP")
			}
			if clientIP == "" {
				clientIP = c.ClientIP()
			}

			key = "ip_" + clientIP
		}

		limiter := limiterManager.getLimiter(key)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, dto.NewErrorResponse(dto.CodeTooManyRequests, dto.MsgTooManyRequests))
			c.Abort()
			return
		}

		c.Next()
	}
}
