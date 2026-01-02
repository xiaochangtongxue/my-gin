package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/internal/handler"
)

// RegisterHealthRoutes 注册健康检查路由
func RegisterHealthRoutes(r *gin.Engine, healthHandler *handler.HealthHandler) {
	// 简化版健康检查
	r.GET("/health", healthHandler.Check)

	// 存活检查
	r.GET("/health/live", healthHandler.Live)

	// 就绪检查（完整依赖检查）
	r.GET("/health/ready", healthHandler.Ready)
}