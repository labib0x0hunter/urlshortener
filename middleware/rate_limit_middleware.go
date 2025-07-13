package middleware

import (
	"crypto/sha1"
	"encoding/hex"
	"time"
	"urlshortener/cache"

	"github.com/gin-gonic/gin"
)

// RateLimitByUserAgentMiddleware is a Gin middleware that limits the number of requests
// from the same User-Agent within a specified time window using Redis as a backend.
// If a User-Agent exceeds the allowed number of requests, it returns HTTP 429 (Too Many Requests).
//
// Usage:
//
//	router.Use(RateLimitByUserAgentMiddleware(redisCache))
//
// Arguments:
//
//	redis cache.Cache: The cache implementation (e.g., Redis) used for rate limiting.
//
// Rate limit: 5 requests per User-Agent per hour.
func RateLimitByUserAgentMiddleware(redis cache.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userAgent := ctx.Request.UserAgent() // Get the User-Agent string from the request

		// Hash the User-Agent to create a consistent and safe Redis key
		hash := sha1.New()
		hash.Write([]byte(userAgent))
		key := "rate:user_agent:" + hex.EncodeToString(hash.Sum(nil))

		// Increment the request count for this User-Agent in Redis
		count, err := redis.Incr(key)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": "Internal server error",
			})
			ctx.Abort()
			return
		}

		// If the request count exceeds the limit, block the request
		if count > 5 {
			ctx.JSON(429, gin.H{
				"error": "Rate limit exceeded",
			})
			ctx.Abort()
			return
		}

		// Set the expiration for the rate limit window if this is the first request
		if count == 1 {
			redis.Expire(key, 1*time.Hour)
		}

		ctx.Next() // Continue to the next handler if not rate limited
	}
}
