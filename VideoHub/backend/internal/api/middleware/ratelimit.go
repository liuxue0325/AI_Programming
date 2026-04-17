package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter  *rate.Limiter
	bursts   int
}

func NewRateLimiter(rps float64, bursts int) *RateLimiter {
	return &RateLimiter{
		limiter:  rate.NewLimiter(rate.Limit(rps), bursts),
		bursts:   bursts,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    1005,
				"message": "rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
