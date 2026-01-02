package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"go.uber.org/zap"
)

const (
	// CSRFTokenKey CSRF Token 在 Context 中的 key
	CSRFTokenKey = "csrf_token"
)

// CSRFConfig CSRF 中间件配置
type CSRFConfig struct {
	Enable         bool     // 是否启用 CSRF 保护
	TokenLength    int      // Token 长度
	TokenExpiry    int      // Token 过期时间（秒）
	CookieName     string   // Cookie 名称
	HeaderName     string   // Header 字段名
	TrustedOrigins []string // 信任的源
}

// DefaultCSRFConfig 默认 CSRF 配置
var DefaultCSRFConfig = CSRFConfig{
	TokenLength:    32,
	TokenExpiry:    int(24 * time.Hour.Seconds()),
	CookieName:     "_csrf",
	HeaderName:     "X-CSRF-Token",
	TrustedOrigins: []string{},
}

var (
	// defaultCSRFConfig 默认 CSRF 配置（兼容旧代码）
	defaultCSRFConfig *CSRFConfig
)

// InitCSRF 初始化 CSRF 配置（兼容全局配置）
func InitCSRF(cfg *config.Config) {
	defaultCSRFConfig = &CSRFConfig{
		Enable:         cfg.Middleware.CSRF.Enable,
		TokenLength:    cfg.Middleware.CSRF.TokenLength,
		TokenExpiry:    int(cfg.Middleware.CSRF.TokenExpiry.Seconds()),
		CookieName:     cfg.Middleware.CSRF.CookieName,
		HeaderName:     cfg.Middleware.CSRF.HeaderName,
		TrustedOrigins: cfg.Middleware.CSRF.TrustedOrigins,
	}
}

// CSRF CSRF 保护中间件
// 使用默认配置
func CSRF() gin.HandlerFunc {
	if defaultCSRFConfig != nil {
		return CSRFWithConfig(*defaultCSRFConfig)
	}
	return CSRFWithConfig(DefaultCSRFConfig)
}

// CSRFWithConfig 带配置的 CSRF 中间件
func CSRFWithConfig(csrfCfg CSRFConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否启用 CSRF 保护
		if !csrfCfg.Enable {
			c.Next()
			return
		}

		// 检查是否是信任的源
		origin := c.Request.Header.Get("Origin")
		referer := c.Request.Header.Get("Referer")
		if isTrustedOrigin(origin, csrfCfg.TrustedOrigins) || isTrustedOrigin(referer, csrfCfg.TrustedOrigins) {
			c.Next()
			return
		}

		// 对于安全方法（GET、HEAD、OPTIONS、TRACE），生成新 Token
		if isSafeMethod(c.Request.Method) {
			token := generateCSRFToken(csrfCfg.TokenLength)
			c.Set(CSRFTokenKey, token)

			// 设置 Cookie
			setCSRFCookie(c, token, csrfCfg)

			// 设置响应头（方便前端获取）
			c.Header(csrfCfg.HeaderName, token)
			c.Next()
			return
		}

		// 对于非安全方法（POST、PUT、DELETE、PATCH），验证 Token
		if !validateCSRFToken(c, csrfCfg) {
			logger.Warn("CSRF Token 验证失败",
				zap.String("request_id", GetRequestID(c)),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("client_ip", c.ClientIP()),
			)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

// isSafeMethod 检查是否是安全方法
func isSafeMethod(method string) bool {
	return method == "GET" || method == "HEAD" || method == "OPTIONS" || method == "TRACE"
}

// isTrustedOrigin 检查是否是信任的源
func isTrustedOrigin(origin string, trustedOrigins []string) bool {
	if origin == "" {
		return false
	}

	for _, trusted := range trustedOrigins {
		if strings.HasSuffix(origin, trusted) || origin == trusted {
			return true
		}
	}
	return false
}

// generateCSRFToken 生成 CSRF Token
func generateCSRFToken(length int) string {
	if length <= 0 {
		length = 32
	}

	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用备用方案
		return generateFallbackToken(length)
	}

	return hex.EncodeToString(bytes)
}

// generateFallbackToken 生成备用 Token
func generateFallbackToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}

// setCSRFCookie 设置 CSRF Cookie
func setCSRFCookie(c *gin.Context, token string, csrfCfg CSRFConfig) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		csrfCfg.CookieName,
		token,
		csrfCfg.TokenExpiry,
		"",
		"",
		true, // Secure (HTTPS only)
		true, // HttpOnly
	)
}

// validateCSRFToken 验证 CSRF Token
func validateCSRFToken(c *gin.Context, csrfCfg CSRFConfig) bool {
	// 从 Header 获取 Token
	token := c.GetHeader(csrfCfg.HeaderName)

	// 如果 Header 中没有，尝试从 Form 获取
	if token == "" {
		token = c.PostForm(csrfCfg.CookieName)
	}

	// 如果还没有，尝试从 Query 获取
	if token == "" {
		token = c.Query(csrfCfg.CookieName)
	}

	if token == "" {
		return false
	}

	// 从 Cookie 获取 Token 进行比对
	cookieToken, err := c.Cookie(csrfCfg.CookieName)
	if err != nil {
		return false
	}

	return token == cookieToken
}

// GetCSRFToken 从 Context 获取 CSRF Token
func GetCSRFToken(c *gin.Context) string {
	if token, exists := c.Get(CSRFTokenKey); exists {
		if t, ok := token.(string); ok {
			return t
		}
	}
	return ""
}