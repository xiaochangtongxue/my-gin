package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/internal/handler"
)

// RegisterAuthRoutes 注册认证相关路由
// Auth 中间件已在 middleware.Apply 中统一应用，白名单路径从 config.yaml 读取
func RegisterAuthRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
	auth := r.Group("/api/v1")
	{
		// 刷新 Token（白名单路径，不需要认证）
		auth.POST("/refresh", authHandler.RefreshToken)

		// 登出（需要认证）
		auth.POST("/logout", authHandler.Logout)

		// 登出所有设备（需要认证）
		auth.POST("/logout/all", authHandler.LogoutAll)
	}
}