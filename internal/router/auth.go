package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/internal/handler"
	"github.com/xiaochangtongxue/my-gin/internal/middleware"
	"github.com/xiaochangtongxue/my-gin/internal/repository"
	"github.com/xiaochangtongxue/my-gin/pkg/permission"
)

// RegisterAuthRoutes 注册认证相关路由
// Auth 中间件已在 middleware.Apply 中统一应用，白名单路径从 config.yaml 读取
func RegisterAuthRoutes(r *gin.Engine, authHandler *handler.AuthHandler, captchaHandler *handler.CaptchaHandler) {
	auth := r.Group("/api/v1")
	{
		// 验证码（白名单路径，不需要认证）
		auth.GET("/captcha", captchaHandler.GetCaptcha)

		// 注册（白名单路径，不需要认证）
		auth.POST("/register", authHandler.Register)

		// 登录（白名单路径，不需要认证）
		auth.POST("/login", authHandler.Login)

		// 刷新 Token（白名单路径，不需要认证）
		auth.POST("/refresh", authHandler.RefreshToken)

		// 登出（需要认证）
		auth.POST("/logout", authHandler.Logout)

	}
}

// RegisterSecurityRoutes 注册安全管理相关路由
func RegisterSecurityRoutes(r *gin.Engine, securityHandler *handler.SecurityHandler, checker permission.PermissionChecker, userRoleRepo repository.UserRoleRepository) {
	admin := r.Group("/api/v1/admin")
	{
		// 加载用户角色中间件（在权限检查前执行）
		admin.Use(middleware.LoadUserRoles(userRoleRepo))
		admin.POST("/unlock-account", middleware.PermissionRequired(checker, "/api/v1/admin/roles", permission.ActionCreate), securityHandler.UnlockAccount)
	}
}
