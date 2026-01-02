package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
)

// CORSConfig CORS 中间件配置
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig 默认 CORS 配置
var DefaultCORSConfig = CORSConfig{
	AllowOrigins:     []string{"*"},
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
	ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
	AllowCredentials: false,
	MaxAge:           86400,
}

var (
	// defaultCORSConfig 默认 CORS 配置（兼容旧代码）
	defaultCORSConfig *CORSConfig
)

// InitCORS 初始化 CORS 配置（兼容全局配置）
func InitCORS(cfg *config.Config) {
	defaultCORSConfig = &CORSConfig{
		AllowOrigins:     cfg.Middleware.CORS.AllowOrigins,
		AllowMethods:     cfg.Middleware.CORS.AllowMethods,
		AllowHeaders:     cfg.Middleware.CORS.AllowHeaders,
		ExposeHeaders:    cfg.Middleware.CORS.ExposeHeaders,
		AllowCredentials: cfg.Middleware.CORS.AllowCredentials,
		MaxAge:           cfg.Middleware.CORS.MaxAge,
	}
}

// CORS 跨域资源共享中间件
// 使用默认配置
func CORS() gin.HandlerFunc {
	if defaultCORSConfig != nil {
		return CORSWithConfig(*defaultCORSConfig)
	}
	return CORSWithConfig(DefaultCORSConfig)
}

// CORSWithConfig 带配置的 CORS 中间件
func CORSWithConfig(corsCfg CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查是否允许该源
		allowed := isOriginAllowed(origin, corsCfg.AllowOrigins)

		if allowed && origin != "" {
			// 如果允许所有源且不允许凭证，使用 *
			if len(corsCfg.AllowOrigins) == 1 && corsCfg.AllowOrigins[0] == "*" && !corsCfg.AllowCredentials {
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}

		c.Header("Access-Control-Allow-Methods", stringsJoin(corsCfg.AllowMethods, ", "))
		c.Header("Access-Control-Allow-Headers", stringsJoin(corsCfg.AllowHeaders, ", "))
		c.Header("Access-Control-Expose-Headers", stringsJoin(corsCfg.ExposeHeaders, ", "))
		c.Header("Access-Control-Max-Age", strconv.Itoa(corsCfg.MaxAge))

		if corsCfg.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// isOriginAllowed 检查源是否被允许
func isOriginAllowed(origin string, allowOrigins []string) bool {
	for _, allowOrigin := range allowOrigins {
		if allowOrigin == "*" || allowOrigin == origin {
			return true
		}
	}
	return false
}

// stringsJoin 连接字符串切片（替代 strings.Join 以减少导入）
func stringsJoin(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	var b strings.Builder
	b.WriteString(strs[0])
	for _, s := range strs[1:] {
		b.WriteString(sep)
		b.WriteString(s)
	}
	return b.String()
}
