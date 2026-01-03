package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/limiter"
	"github.com/xiaochangtongxue/my-gin/pkg/response"
	"strconv"
)

var (
	globalLimiter   *limiter.TokenBucketLimiter
	ipLimiter       *limiter.TokenBucketLimiter
	userLimiter     *limiter.TokenBucketLimiter
	rateLimitConfig *limiter.Config
)

// InitRateLimit 初始化限流器配置
func InitRateLimit(cfg *config.Config, cache cache.Cache) {
	rateLimitConfig = &limiter.Config{
		Enable: cfg.Middleware.RateLimit.Enable,
	}

	if cfg.Middleware.RateLimit.Global.Rate > 0 {
		globalLimiter = limiter.NewTokenBucketLimiter(
			cache,
			cfg.Middleware.RateLimit.Global.Rate,
			int64(cfg.Middleware.RateLimit.Global.Burst),
		)
	}

	if cfg.Middleware.RateLimit.IP.Rate > 0 {
		ipLimiter = limiter.NewTokenBucketLimiter(
			cache,
			cfg.Middleware.RateLimit.IP.Rate,
			int64(cfg.Middleware.RateLimit.IP.Burst),
		)
	}

	if cfg.Middleware.RateLimit.User.Rate > 0 {
		userLimiter = limiter.NewTokenBucketLimiter(
			cache,
			cfg.Middleware.RateLimit.User.Rate,
			int64(cfg.Middleware.RateLimit.User.Burst),
		)
	}
}

// RateLimit 限流中间件
// 使用默认配置
func RateLimit() gin.HandlerFunc {
	return RateLimitWithConfig(rateLimitConfig)
}

// RateLimitWithConfig 带配置的限流中间件
func RateLimitWithConfig(cfg *limiter.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.Enable {
			c.Next()
			return
		}

		// 检查全局限流
		if globalLimiter != nil {
			if !globalLimiter.Allow(c.Request.Context(), "global") {
				response.RateLimitError(c)
				c.Abort()
				return
			}
		}

		// 检查 IP 限流
		if ipLimiter != nil {
			clientIP := c.ClientIP()
			if !ipLimiter.Allow(c.Request.Context(), "ip:"+clientIP) {
				response.RateLimitError(c)
				c.Abort()
				return
			}
		}

		// 检查用户限流（需要认证）
		if userLimiter != nil {
			if userID := getUserID(c); userID != 0 {
				key := "user:" + strconv.Itoa(int(userID))
				if !userLimiter.Allow(c.Request.Context(), key) {
					response.RateLimitError(c)
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}

// getUserID 从上下文获取用户 ID
func getUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}
