// Package handler 权限处理器
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/internal/dto/req"
	"github.com/xiaochangtongxue/my-gin/internal/middleware"
	"github.com/xiaochangtongxue/my-gin/internal/service"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"github.com/xiaochangtongxue/my-gin/pkg/response"
	appvalidator "github.com/xiaochangtongxue/my-gin/pkg/validator"
	"go.uber.org/zap"
)

// PermissionHandler 权限处理器
type PermissionHandler struct {
	permissionService service.PermissionService
}

// NewPermissionHandler 创建权限处理器
func NewPermissionHandler(permissionService service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新角色
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body req.CreateRoleReq true "创建角色请求"
// @Success 200 {object} response.Response{data=resp.RoleResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/roles [post]
func (h *PermissionHandler) CreateRole(c *gin.Context) {
	var r req.CreateRoleReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.ParamError(c, appvalidator.TranslateError(err))
		return
	}

	result, err := h.permissionService.CreateRole(c.Request.Context(), &r)
	if err != nil {
		logger.Warn("创建角色失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.String("code", r.Code),
			zap.Error(err),
		)
		response.Error(c, err)
		return
	}

	logger.Info("创建角色成功",
		zap.String("request_id", middleware.GetRequestID(c)),
		zap.String("code", r.Code),
		zap.String("name", r.Name),
	)

	response.Success(c, result)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param request body req.UpdateRoleReq true "更新角色请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/roles/{id} [put]
func (h *PermissionHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "无效的角色ID")
		return
	}

	var r req.UpdateRoleReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.ParamError(c, appvalidator.TranslateError(err))
		return
	}

	if err := h.permissionService.UpdateRole(c.Request.Context(), id, &r); err != nil {
		logger.Warn("更新角色失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Uint64("id", id),
			zap.Error(err),
		)
		response.Error(c, err)
		return
	}

	logger.Info("更新角色成功",
		zap.String("request_id", middleware.GetRequestID(c)),
		zap.Uint64("id", id),
	)

	response.Success(c, nil)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除角色（内置角色不可删除）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/roles/{id} [delete]
func (h *PermissionHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "无效的角色ID")
		return
	}

	if err := h.permissionService.DeleteRole(c.Request.Context(), id); err != nil {
		logger.Warn("删除角色失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Uint64("id", id),
			zap.Error(err),
		)
		response.Error(c, err)
		return
	}

	logger.Info("删除角色成功",
		zap.String("request_id", middleware.GetRequestID(c)),
		zap.Uint64("id", id),
	)

	response.Success(c, nil)
}

// GetRole 获取角色详情
// @Summary 获取角色详情
// @Description 获取角色详细信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=resp.RoleResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/roles/{id} [get]
func (h *PermissionHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "无效的角色ID")
		return
	}

	role, err := h.permissionService.GetRole(c.Request.Context(), id)
	if err != nil {
		logger.Warn("获取角色失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Uint64("id", id),
			zap.Error(err),
		)
		response.Error(c, err)
		return
	}

	response.Success(c, role)
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Description 获取所有角色列表
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query int false "状态筛选"
// @Param code query string false "编码筛选"
// @Success 200 {object} response.Response{data=[]resp.RoleResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/roles [get]
func (h *PermissionHandler) ListRoles(c *gin.Context) {
	var r req.ListRolesReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.ParamError(c, appvalidator.TranslateError(err))
		return
	}

	roles, err := h.permissionService.ListRoles(c.Request.Context(), &r)
	if err != nil {
		logger.Warn("获取角色列表失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Error(err),
		)
		response.Error(c, err)
		return
	}

	response.Success(c, roles)
}

// AddPermission 添加权限
// @Summary 为角色添加权限
// @Description 为指定角色添加权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param request body req.AddPermissionReq true "添加权限请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/roles/{id}/permissions [post]
func (h *PermissionHandler) AddPermission(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "无效的角色ID")
		return
	}

	var r req.AddPermissionReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.ParamError(c, appvalidator.TranslateError(err))
		return
	}

	if err := h.permissionService.AddPermission(c.Request.Context(), roleID, &r); err != nil {
		logger.Warn("添加权限失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Uint64("role_id", roleID),
			zap.String("resource", r.Resource),
			zap.String("action", r.Action),
			zap.Error(err),
		)
		response.Error(c, err)
		return
	}

	logger.Info("添加权限成功",
		zap.String("request_id", middleware.GetRequestID(c)),
		zap.Uint64("role_id", roleID),
		zap.String("resource", r.Resource),
		zap.String("action", r.Action),
	)

	response.Success(c, nil)
}

// RemovePermission 删除权限
// @Summary 删除角色权限
// @Description 删除指定角色的权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param resource query string true "资源路径"
// @Param action query string true "操作"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/roles/{id}/permissions [delete]
func (h *PermissionHandler) RemovePermission(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "无效的角色ID")
		return
	}

	resource := c.Query("resource")
	action := c.Query("action")

	if resource == "" || action == "" {
		response.ParamError(c, "resource 和 action 参数必填")
		return
	}

	if err := h.permissionService.RemovePermission(c.Request.Context(), roleID, resource, action); err != nil {
		logger.Warn("删除权限失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Uint64("role_id", roleID),
			zap.String("resource", resource),
			zap.String("action", action),
			zap.Error(err),
		)
		response.Error(c, err)
		return
	}

	logger.Info("删除权限成功",
		zap.String("request_id", middleware.GetRequestID(c)),
		zap.Uint64("role_id", roleID),
		zap.String("resource", resource),
		zap.String("action", action),
	)

	response.Success(c, nil)
}

// GetRolePermissions 获取角色权限列表
// @Summary 获取角色权限列表
// @Description 获取指定角色的所有权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=[]resp.PermissionResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/roles/{id}/permissions [get]
func (h *PermissionHandler) GetRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "无效的角色ID")
		return
	}

	permissions, err := h.permissionService.GetRolePermissions(c.Request.Context(), roleID)
	if err != nil {
		logger.Warn("获取角色权限失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Uint64("role_id", roleID),
			zap.Error(err),
		)
		response.Error(c, err)
		return
	}

	response.Success(c, permissions)
}

// AssignRole 分配角色给用户
// @Summary 分配角色给用户
// @Description 为指定用户分配角色
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body req.AssignRoleReq true "分配角色请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/users/{id}/roles [post]
func (h *PermissionHandler) AssignRole(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "无效的用户ID")
		return
	}

	var r req.AssignRoleReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.ParamError(c, appvalidator.TranslateError(err))
		return
	}

	if err := h.permissionService.AssignRole(c.Request.Context(), userID, r.RoleID); err != nil {
		logger.Warn("分配角色失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Uint64("user_id", userID),
			zap.Uint64("role_id", r.RoleID),
			zap.Error(err),
		)
		response.Error(c, err)
		return
	}

	logger.Info("分配角色成功",
		zap.String("request_id", middleware.GetRequestID(c)),
		zap.Uint64("user_id", userID),
		zap.Uint64("role_id", r.RoleID),
	)

	response.Success(c, nil)
}

// RemoveRole 移除用户角色
// @Summary 移除用户角色
// @Description 移除用户的指定角色
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param role_id path int true "角色ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/users/{id}/roles/{roleId} [delete]
func (h *PermissionHandler) RemoveRole(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "无效的用户ID")
		return
	}

	roleIDStr := c.Param("roleId")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
	if err != nil {
		response.ParamError(c, "无效的角色ID")
		return
	}

	if err := h.permissionService.RemoveRole(c.Request.Context(), userID, roleID); err != nil {
		logger.Warn("移除角色失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Uint64("user_id", userID),
			zap.Uint64("role_id", roleID),
			zap.Error(err),
		)
		response.Error(c, err)
		return
	}

	logger.Info("移除角色成功",
		zap.String("request_id", middleware.GetRequestID(c)),
		zap.Uint64("user_id", userID),
		zap.Uint64("role_id", roleID),
	)

	response.Success(c, nil)
}

// GetUserRoles 获取用户角色列表
// @Summary 获取用户角色列表
// @Description 获取指定用户的所有角色
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=[]resp.RoleResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/users/{id}/roles [get]
func (h *PermissionHandler) GetUserRoles(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "无效的用户ID")
		return
	}

	roles, err := h.permissionService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		logger.Warn("获取用户角色失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Uint64("user_id", userID),
			zap.Error(err),
		)
		response.Error(c, err)
		return
	}

	response.Success(c, roles)
}
