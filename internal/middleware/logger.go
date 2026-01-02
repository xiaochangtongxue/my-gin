package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"go.uber.org/zap"
)

// LoggerConfig 日志中间件配置
type LoggerConfig struct {
	SkipPaths    []string        // 跳过日志记录的路径
	SlowThreshold time.Duration // 慢请求阈值
	LogQuery     bool           // 是否记录查询参数
}

var (
	// defaultLoggerConfig 默认日志配置（兼容旧代码）
	defaultLoggerConfig *LoggerConfig
)

// InitLogger 初始化日志配置（兼容全局配置）
func InitLogger(cfg *config.Config) {
	defaultLoggerConfig = &LoggerConfig{
		SkipPaths:    cfg.Middleware.Logger.SkipPaths,
		SlowThreshold: cfg.Middleware.Logger.SlowThreshold,
		LogQuery:     cfg.Middleware.Logger.LogQuery,
	}
}

// Logger 请求日志中间件
// 使用默认配置
func Logger() gin.HandlerFunc {
	if defaultLoggerConfig != nil {
		return LoggerWithConfig(*defaultLoggerConfig)
	}
	return LoggerWithConfig(LoggerConfig{})
}

// LoggerWithConfig 带配置的日志中间件
func LoggerWithConfig(loggerCfg LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过（支持前缀匹配）
		shouldSkip := false
		requestPath := c.Request.URL.Path
		for _, skipPath := range loggerCfg.SkipPaths {
			if requestPath == skipPath || strings.HasPrefix(requestPath, skipPath+"/") {
				shouldSkip = true
				break
			}
		}

		if shouldSkip {
			c.Next()
			return
		}

		// 记录开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)

		// 获取请求 ID
		requestID := GetRequestID(c)

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 获取客户端 IP
		clientIP := c.ClientIP()

		// 获取请求方法
		method := c.Request.Method

		// 构建日志字段
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// 记录查询参数（根据配置）
		if loggerCfg.LogQuery && query != "" {
			fields = append(fields, zap.String("query", query))
		}

		// 记录错误信息
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// 根据状态码和耗时选择日志级别
		switch {
		case statusCode >= 500:
			// 服务器错误
			logger.Error("请求完成（服务器错误）", fields...)
		case statusCode >= 400:
			// 客户端错误
			logger.Warn("请求完成（客户端错误）", fields...)
		case latency > loggerCfg.SlowThreshold:
			// 慢请求
			logger.Warn("请求完成（慢请求）",
				append(fields, zap.Duration("slow_threshold", loggerCfg.SlowThreshold))...,
			)
		default:
			// 正常请求
			logger.Debug("请求完成", fields...)
		}
	}
}