// Package middleware 权限中间件
package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/internal/repository"
	"github.com/xiaochangtongxue/my-gin/internal/permission"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"github.com/xiaochangtongxue/my-gin/pkg/response"
	"go.uber.org/zap"
)

// PermissionRequired 权限检查中间件
func PermissionRequired(checker permission.PermissionChecker, resource string, action permission.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := getSubjectFromContext(c)
		if subject == nil {
			response.Fail(c, response.CodeUnauthorized, "未登录")
			c.Abort()
			return
		}

		req := &permission.CheckRequest{
			Subject: subject,
			Resource: &permission.Resource{
				Path: resource,
				ID:   c.Param("id"),
			},
			Action: action,
		}

		allowed, err := checker.Check(c.Request.Context(), req)
		if err != nil {
			logger.Warn("权限检查失败",
				zap.String("request_id", GetRequestID(c)),
				zap.String("user_id", subject.ID),
				zap.String("resource", resource),
				zap.Error(err),
			)
			response.Fail(c, response.CodeServerError, "权限检查失败")
			c.Abort()
			return
		}

		if !allowed {
			logger.Warn("权限不足",
				zap.String("request_id", GetRequestID(c)),
				zap.String("user_id", subject.ID),
				zap.String("resource", resource),
				zap.String("action", string(action)),
			)
			response.Fail(c, response.CodePermissionDenied, "权限不足，无法访问该资源")
			c.Abort()
			return
		}

		c.Next()
	}
}

// OwnerRequired 资源所有者权限检查
func OwnerRequired(checker permission.PermissionChecker, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := getSubjectFromContext(c)
		if subject == nil {
			response.Fail(c, response.CodeUnauthorized, "未登录")
			c.Abort()
			return
		}

		resourceID := c.Param("id")

		req := &permission.CheckRequest{
			Subject: subject,
			Resource: &permission.Resource{
				Type: resourceType,
				ID:   resourceID,
				Path: "/api/v1/" + resourceType + "/:id",
			},
			Action: permission.ActionAny,
		}

		allowed, err := checker.Check(c.Request.Context(), req)
		if err != nil {
			logger.Warn("权限检查失败",
				zap.String("request_id", GetRequestID(c)),
				zap.String("user_id", subject.ID),
				zap.String("resource_type", resourceType),
				zap.String("resource_id", resourceID),
				zap.Error(err),
			)
			response.Fail(c, response.CodeServerError, "权限检查失败")
			c.Abort()
			return
		}

		if !allowed {
			logger.Warn("资源所有者权限检查失败",
				zap.String("request_id", GetRequestID(c)),
				zap.String("user_id", subject.ID),
				zap.String("resource_type", resourceType),
				zap.String("resource_id", resourceID),
			)
			response.Fail(c, response.CodePermissionDenied, "无权访问此资源")
			c.Abort()
			return
		}

		// 将资源所有者信息存入上下文，供后续使用
		c.Set("resource_owner_id", subject.ID)
		c.Next()
	}
}

// getSubjectFromContext 从上下文获取用户主体信息
func getSubjectFromContext(c *gin.Context) *permission.Subject {
	// 从 JWT 中间件获取 user_id
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return nil
	}

	var userID string
	switch v := userIDVal.(type) {
	case uint64:
		userID = strconv.FormatUint(v, 10)
	case string:
		userID = v
	default:
		return nil
	}

	if userID == "" {
		return nil
	}

	// 获取用户角色
	var roleIDs []string
	if rolesVal, exists := c.Get("user_roles"); exists {
		switch v := rolesVal.(type) {
		case []string:
			roleIDs = v
		case []uint64:
			for _, id := range v {
				roleIDs = append(roleIDs, strconv.FormatUint(id, 10))
			}
		}
	}

	// 如果没有角色，使用默认访客角色
	if len(roleIDs) == 0 {
		roleIDs = []string{"4"} // guest
	}

	return &permission.Subject{
		ID:    userID,
		Type:  "user",
		Roles: roleIDs,
	}
}

// LoadUserRoles 加载用户角色中间件
// 用于在需要权限检查的接口前加载用户角色
func LoadUserRoles(userRoleRepo repository.UserRoleRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		var userID uint64
		switch v := userIDVal.(type) {
		case uint64:
			userID = v
		case string:
			id, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				c.Next()
				return
			}
			userID = id
		default:
			c.Next()
			return
		}

		// 查询用户角色
		roleIDs, err := userRoleRepo.FindRoleIDsByUserID(c.Request.Context(), userID)
		if err != nil {
			logger.Warn("查询用户角色失败",
				zap.String("request_id", GetRequestID(c)),
				zap.Uint64("user_id", userID),
				zap.Error(err),
			)
			// 查询失败不影响，使用默认访客角色
			roleIDs = []uint64{4}
		}

		// 如果没有角色，使用默认访客角色
		if len(roleIDs) == 0 {
			roleIDs = []uint64{4}
		}

		// 转换为字符串数组
		roleIDStrs := make([]string, len(roleIDs))
		for i, id := range roleIDs {
			roleIDStrs[i] = strconv.FormatUint(id, 10)
		}

		c.Set("user_roles", roleIDStrs)
		c.Next()
	}
}

// SetupUserRoles 从请求头获取用户角色（兼容性中间件）
// 如果没有使用 LoadUserRoles，可以尝试从 JWT Token 中解析角色
func SetupUserRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果已经设置了 user_roles，跳过
		if _, exists := c.Get("user_roles"); exists {
			c.Next()
			return
		}

		// 从 Authorization header 获取 token 并解析
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			// 这里可以扩展解析 JWT 中的 roles claim
			// 暂时使用默认访客角色
			c.Set("user_roles", []string{"4"})
		}

		c.Next()
	}
}
