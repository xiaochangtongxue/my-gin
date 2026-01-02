package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/jwt"
	"github.com/xiaochangtongxue/my-gin/pkg/metrics"
)

// DefaultMiddleware 默认中间件配置
type DefaultMiddleware struct {
	// EnableRecovery 是否启用 Recovery 中间件
	EnableRecovery bool
	// EnableRequestID 是否启用 RequestID 中间件
	EnableRequestID bool
	// EnableLogger 是否启用 Logger 中间件
	EnableLogger bool
	// EnableCORS 是否启用 CORS 中间件
	EnableCORS bool
	// EnableSecurity 是否启用安全响应头中间件
	EnableSecurity bool
	// EnableXSS 是否启用 XSS 防护中间件
	EnableXSS bool
	// EnableCSRF 是否启用 CSRF 中间件
	EnableCSRF bool
	// EnableAuth 是否启用 JWT 认证中间件
	EnableAuth bool
	// EnableMetrics 是否启用 Metrics 中间件
	EnableMetrics bool

	// AuthConfig 认证中间件配置（EnableAuth=true 时需要）
	AuthSkipPaths []string     // 跳过认证的路径
	AuthJWTMgr    *jwt.Manager // JWT 管理器
	AuthCache     cache.Cache  // 黑名单缓存
}

// ApplyWithConfig 使用自定义配置应用中间件
// cfg 为全局配置，用于初始化各中间件的配置
func ApplyWithConfig(engine *gin.Engine, cfg *config.Config, middlewareConfig DefaultMiddleware) {
	// 初始化各中间件配置（依赖注入方式）
	InitRecovery(cfg)
	InitRequestID(cfg)
	InitLogger(cfg)
	InitCORS(cfg)
	InitSecurity(cfg)
	InitCSRF(cfg)

	// Recovery 必须在最前面，用于捕获后续中间件的 panic
	if middlewareConfig.EnableRecovery {
		engine.Use(Recovery())
	}

	// RequestID 应该尽早执行，确保后续中间件都能获取到请求 ID
	if middlewareConfig.EnableRequestID {
		engine.Use(RequestID())
	}

	// 安全响应头
	if middlewareConfig.EnableSecurity {
		engine.Use(Security())
	}

	// CORS 跨域
	if middlewareConfig.EnableCORS {
		engine.Use(CORS())
	}

	// XSS 防护
	if middlewareConfig.EnableXSS {
		engine.Use(XSS())
	}

	// CSRF 保护
	if middlewareConfig.EnableCSRF {
		engine.Use(CSRF())
	}

	// JWT 认证
	if middlewareConfig.EnableAuth {
		// 初始化认证中间件配置
		if middlewareConfig.AuthJWTMgr != nil {
			InitAuth(middlewareConfig.AuthSkipPaths, middlewareConfig.AuthJWTMgr, middlewareConfig.AuthCache)
		}
		engine.Use(Auth())
	}

	// Prometheus Metrics 指标收集
	if middlewareConfig.EnableMetrics {
		engine.Use(MetricsWithConfig(metrics.Config{
			SkipPaths: cfg.Middleware.Metrics.SkipPaths,
		}))
	}

	// 请求日志（放在最后，记录所有中间件处理完成后的信息）
	if middlewareConfig.EnableLogger {
		engine.Use(Logger())
	}
}
