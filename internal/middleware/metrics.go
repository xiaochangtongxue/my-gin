package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/pkg/metrics"
)

// InitMetrics 初始化 Metrics 中间件配置
func InitMetrics(cfg *metrics.Config) {
	// Metrics 中间件使用传入的配置，无需存储全局状态
}

// Metrics Prometheus 指标收集中间件
// 使用默认配置
func Metrics() gin.HandlerFunc {
	return MetricsWithConfig(metrics.Config{})
}

// MetricsWithConfig 带配置的 Metrics 中间件
func MetricsWithConfig(cfg metrics.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否在白名单中（支持前缀匹配）
		for _, path := range cfg.SkipPaths {
			if c.Request.URL.Path == path || strings.HasPrefix(c.Request.URL.Path, path+"/") {
				c.Next()
				return
			}
		}

		// 记录开始时间
		start := time.Now()

		// 增加处理中请求数
		metrics.HTTPRequestsInFlight.Inc()

		// 处理请求
		c.Next()

		// 减少处理中请求数
		metrics.HTTPRequestsInFlight.Dec()

		// 计算耗时（秒）
		duration := time.Since(start).Seconds()

		// 获取路径
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// 记录请求数
		metrics.HTTPRequestTotal.WithLabelValues(
			c.Request.Method,
			path,
			strconv.Itoa(c.Writer.Status()),
		).Inc()

		// 记录请求耗时
		metrics.HTTPRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)
	}
}