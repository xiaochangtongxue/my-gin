// Package router 权限路由注册
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/internal/handler"
	"github.com/xiaochangtongxue/my-gin/internal/middleware"
	"github.com/xiaochangtongxue/my-gin/internal/repository"
	"github.com/xiaochangtongxue/my-gin/internal/permission"
)

// RegisterPermissionRoutes 注册权限管理路由
func RegisterPermissionRoutes(r *gin.Engine, permissionHandler *handler.PermissionHandler, checker permission.PermissionChecker, userRoleRepo repository.UserRoleRepository) {
	admin := r.Group("/api/v1/admin")
	{
		// 加载用户角色中间件（在权限检查前执行）
		admin.Use(middleware.LoadUserRoles(userRoleRepo))

		// ====== 角色管理 ======

		// 创建角色
		admin.POST("/roles",
			middleware.PermissionRequired(checker, "/api/v1/admin/roles", permission.ActionCreate),
			permissionHandler.CreateRole,
		)

		// 获取角色列表
		admin.GET("/roles",
			middleware.PermissionRequired(checker, "/api/v1/admin/roles", permission.ActionRead),
			permissionHandler.ListRoles,
		)

		// 获取角色详情
		admin.GET("/roles/:id",
			middleware.PermissionRequired(checker, "/api/v1/admin/roles", permission.ActionRead),
			permissionHandler.GetRole,
		)

		// 更新角色
		admin.PUT("/roles/:id",
			middleware.PermissionRequired(checker, "/api/v1/admin/roles", permission.ActionUpdate),
			permissionHandler.UpdateRole,
		)

		// 删除角色
		admin.DELETE("/roles/:id",
			middleware.PermissionRequired(checker, "/api/v1/admin/roles", permission.ActionDelete),
			permissionHandler.DeleteRole,
		)

		// ====== 角色权限管理 ======

		// 获取角色权限列表
		admin.GET("/roles/:id/permissions",
			middleware.PermissionRequired(checker, "/api/v1/admin/roles", permission.ActionRead),
			permissionHandler.GetRolePermissions,
		)

		// 添加角色权限
		admin.POST("/roles/:id/permissions",
			middleware.PermissionRequired(checker, "/api/v1/admin/roles", permission.ActionCreate),
			permissionHandler.AddPermission,
		)

		// 删除角色权限
		admin.DELETE("/roles/:id/permissions",
			middleware.PermissionRequired(checker, "/api/v1/admin/roles", permission.ActionDelete),
			permissionHandler.RemovePermission,
		)

		// ====== 用户角色管理 ======

		// 获取用户角色列表
		admin.GET("/users/:id/roles",
			middleware.PermissionRequired(checker, "/api/v1/admin/users", permission.ActionRead),
			permissionHandler.GetUserRoles,
		)

		// 分配角色给用户
		admin.POST("/users/:id/roles",
			middleware.PermissionRequired(checker, "/api/v1/admin/users", permission.ActionUpdate),
			permissionHandler.AssignRole,
		)

		// 移除用户角色
		admin.DELETE("/users/:id/roles/:roleId",
			middleware.PermissionRequired(checker, "/api/v1/admin/users", permission.ActionUpdate),
			permissionHandler.RemoveRole,
		)
	}
}
