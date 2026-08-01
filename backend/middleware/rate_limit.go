package middleware

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

var tokenBucketScript = redis.NewScript(`
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local fill_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local ttl = math.ceil(capacity / fill_rate)

local data = redis.call("HMGET", key, "tokens", "last_updated")
local tokens = tonumber(data[1])
local last_updated = tonumber(data[2])

if not tokens then
    tokens = capacity
    last_updated = now
else
    local delta = math.max(0, now - last_updated)
    tokens = math.min(capacity, tokens + (delta * fill_rate))
    last_updated = now
end

if tokens >= 1 then
    tokens = tokens - 1
    redis.call("HMSET", key, "tokens", tokens, "last_updated", last_updated)
    redis.call("EXPIRE", key, ttl)
    return {1, math.floor(tokens)}
else
    redis.call("EXPIRE", key, ttl)
    return {0, 0}
end
`)

type RateLimiter struct {
	rdb      *redis.Client
	fillRate float64
	capacity int
}

func NewRateLimiter(rdb *redis.Client, r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		rdb:      rdb,
		fillRate: float64(r),
		capacity: b,
	}
}

func RateLimitMiddleware(limiterManager *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiterManager == nil || limiterManager.rdb == nil {
			c.Next()
			return
		}

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

		redisKey := fmt.Sprintf("ratelimit:tb:%s:%s", c.FullPath(), key)
		now := time.Now().Unix()

		ctx := c.Request.Context()

		res, err := tokenBucketScript.Run(
			ctx,
			limiterManager.rdb,
			[]string{redisKey},
			fmt.Sprintf("%d", limiterManager.capacity),
			fmt.Sprintf("%f", limiterManager.fillRate),
			fmt.Sprintf("%d", now),
		).Result()

		if err != nil {
			log.Printf("❌ [RateLimit Redis Error Go]: %v", err)
			c.Next()
			return
		}

		results, ok := res.([]any)
		if !ok || len(results) < 2 {
			log.Printf("❌ [RateLimit Parse Error Go]: %v", res)
			c.Next()
			return
		}

		allowed := results[0].(int64)
		remainingTokens := results[1].(int64)

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiterManager.capacity))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remainingTokens))

		if allowed == 0 {
			retryAfter := math.Ceil(1.0 / limiterManager.fillRate)
			c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter))

			c.JSON(http.StatusTooManyRequests, dto.NewErrorResponse(dto.CodeTooManyRequests, dto.MsgTooManyRequests))
			c.Abort()
			return
		}

		c.Next()
	}
}
