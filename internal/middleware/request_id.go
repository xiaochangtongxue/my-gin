package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"go.uber.org/zap"
)

// RequestIDConfig 请求 ID 中间件配置
type RequestIDConfig struct {
	Header    string   // 请求ID 的 HTTP 头字段名
	SkipPaths []string // 跳过生成 RequestID 的路径
}

// DefaultRequestIDConfig 默认请求 ID 配置
var DefaultRequestIDConfig = RequestIDConfig{
	Header: "X-Request-ID",
}

var (
	// defaultRequestIDConfig 默认请求 ID 配置（兼容旧代码）
	defaultRequestIDConfig *RequestIDConfig
)

// InitRequestID 初始化请求 ID 配置（兼容全局配置）
func InitRequestID(cfg *config.Config) {
	defaultRequestIDConfig = &RequestIDConfig{
		Header:    cfg.Middleware.RequestID.Header,
		SkipPaths: cfg.Middleware.RequestID.SkipPaths,
	}
}

// RequestID 请求 ID 中间件
// 为每个请求生成唯一 ID，用于链路追踪和日志关联
// 使用默认配置
func RequestID() gin.HandlerFunc {
	if defaultRequestIDConfig != nil {
		return RequestIDWithConfig(*defaultRequestIDConfig)
	}
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestIDWithConfig 带配置的请求 ID 中间件
func RequestIDWithConfig(requestIDCfg RequestIDConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过（支持前缀匹配）
		shouldSkip := false
		requestPath := c.Request.URL.Path
		for _, skipPath := range requestIDCfg.SkipPaths {
			if requestPath == skipPath || strings.HasPrefix(requestPath, skipPath+"/") {
				shouldSkip = true
				break
			}
		}

		if shouldSkip {
			c.Next()
			return
		}

		headerName := requestIDCfg.Header

		// 尝试从请求头获取 RequestID
		requestID := c.GetHeader(headerName)

		// 如果请求头中没有，生成新的 UUID
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 设置到 Context 中，供后续中间件和处理器使用
		c.Set(RequestIDKey, requestID)

		// 设置响应头，方便客户端追踪
		c.Header(headerName, requestID)

		// 记录请求开始日志
		logger.Debug("请求开始",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("client_ip", c.ClientIP()),
		)

		c.Next()
	}
}

// generateRequestID 生成请求 ID
// 使用 UUID v4 作为默认实现
func generateRequestID() string {
	return uuid.New().String()
}