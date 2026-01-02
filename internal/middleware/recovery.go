package middleware

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"github.com/xiaochangtongxue/my-gin/pkg/response"
	"go.uber.org/zap"
)

const (
	// RequestIDKey 请求 ID 在 gin.Context 中的 key
	RequestIDKey = "RequestID"
)

// RecoveryConfig Recovery 中间件配置
type RecoveryConfig struct {
	EnableStacktrace bool // 是否记录堆栈信息
}

var (
	// defaultRecoveryConfig 默认 Recovery 配置（兼容旧代码）
	defaultRecoveryConfig *RecoveryConfig
)

// InitRecovery 初始化 Recovery 配置（兼容全局配置）
func InitRecovery(cfg *config.Config) {
	defaultRecoveryConfig = &RecoveryConfig{
		EnableStacktrace: cfg.Middleware.Recovery.EnableStacktrace,
	}
}

// Recovery 全局异常恢复中间件
// 捕获 panic 并记录日志，返回统一错误响应
// 使用默认配置
func Recovery() gin.HandlerFunc {
	if defaultRecoveryConfig != nil {
		return RecoveryWithConfig(*defaultRecoveryConfig)
	}
	return RecoveryWithConfig(RecoveryConfig{EnableStacktrace: true})
}

// RecoveryWithConfig 带配置的 Recovery 中间件
func RecoveryWithConfig(recoveryCfg RecoveryConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求信息
				method := c.Request.Method
				path := c.Request.URL.Path
				query := c.Request.URL.RawQuery
				requestID := GetRequestID(c)

				// 构建日志字段
				fields := []zap.Field{
					zap.String("request_id", requestID),
					zap.Any("error", err),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("query", query),
				}

				// 根据配置决定是否记录堆栈信息
				if recoveryCfg.EnableStacktrace {
					fields = append(fields, zap.String("stack", string(debug.Stack())))
				}

				// 记录 panic 日志
				logger.Error("服务发生 panic", fields...)

				// 返回统一错误响应
				if !c.IsAborted() {
					response.ServerError(c, "服务器内部错误")
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

// GetRequestID 从 gin.Context 获取请求 ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}