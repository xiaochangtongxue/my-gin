package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RegisterMetricsRoutes 注册 Metrics 路由
func RegisterMetricsRoutes(r *gin.Engine) {
	// Prometheus metrics 端点
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}