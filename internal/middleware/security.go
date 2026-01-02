package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
)

// SecurityConfig 安全响应头配置
type SecurityConfig struct {
	Enable                bool
	XFrameOptions          string
	XContentTypeOptions    string
	XSSProtection          string
	HSTSEnable             bool
	HSTSMaxAge             int
	HSTSIncludeSubDomains  bool
	ContentSecurityPolicy  string
	ReferrerPolicy         string
}

var (
	// defaultSecurityConfig 默认安全配置（兼容旧代码）
	defaultSecurityConfig *SecurityConfig
)

// InitSecurity 初始化安全配置（兼容全局配置）
func InitSecurity(cfg *config.Config) {
	defaultSecurityConfig = &SecurityConfig{
		Enable:               cfg.Middleware.Security.Enable,
		XFrameOptions:         cfg.Middleware.Security.XFrameOptions,
		XContentTypeOptions:   cfg.Middleware.Security.XContentTypeOptions,
		XSSProtection:         cfg.Middleware.Security.XSSProtection,
		HSTSEnable:            cfg.Middleware.Security.HSTS.Enable,
		HSTSMaxAge:            int(cfg.Middleware.Security.HSTS.MaxAge.Seconds()),
		HSTSIncludeSubDomains: cfg.Middleware.Security.HSTS.IncludeSubDomains,
		ContentSecurityPolicy: cfg.Middleware.Security.ContentSecurityPolicy,
		ReferrerPolicy:        cfg.Middleware.Security.ReferrerPolicy,
	}
}

// Security 安全响应头中间件
// 使用默认配置
func Security() gin.HandlerFunc {
	if defaultSecurityConfig != nil {
		return SecurityWithConfig(*defaultSecurityConfig)
	}
	return SecurityWithConfig(SecurityConfig{Enable: true})
}

// SecurityWithConfig 带配置的安全响应头中间件
func SecurityWithConfig(securityCfg SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否启用安全响应头
		if !securityCfg.Enable {
			c.Next()
			return
		}

		// X-Frame-Options: 防止点击劫持
		if securityCfg.XFrameOptions != "" {
			c.Header("X-Frame-Options", securityCfg.XFrameOptions)
		}

		// X-Content-Type-Options: 防止 MIME 类型嗅探
		if securityCfg.XContentTypeOptions != "" {
			c.Header("X-Content-Type-Options", securityCfg.XContentTypeOptions)
		}

		// X-XSS-Protection: XSS 保护
		if securityCfg.XSSProtection != "" {
			c.Header("X-XSS-Protection", securityCfg.XSSProtection)
		}

		// Strict-Transport-Security: HSTS
		if securityCfg.HSTSEnable {
			hstsValue := "max-age=" + strconv.Itoa(securityCfg.HSTSMaxAge)
			if securityCfg.HSTSIncludeSubDomains {
				hstsValue += "; includeSubDomains"
			}
			c.Header("Strict-Transport-Security", hstsValue)
		}

		// Content-Security-Policy: CSP
		if securityCfg.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", securityCfg.ContentSecurityPolicy)
		}

		// Referrer-Policy: Referrer 策略
		if securityCfg.ReferrerPolicy != "" {
			c.Header("Referrer-Policy", securityCfg.ReferrerPolicy)
		}

		// X-Permitted-Cross-Domain-Policies: 限制跨域策略
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		// Permissions-Policy: 功能策略（原 Feature-Policy）
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Cross-Origin-Opener-Policy: 控制跨域窗口行为
		c.Header("Cross-Origin-Opener-Policy", "same-origin")

		// Cross-Origin-Resource-Policy: 控制跨域资源访问
		c.Header("Cross-Origin-Resource-Policy", "same-origin")

		c.Next()
	}
}