// Package service 权限服务
package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/xiaochangtongxue/my-gin/internal/dto/req"
	"github.com/xiaochangtongxue/my-gin/internal/dto/resp"
	"github.com/xiaochangtongxue/my-gin/internal/model"
	"github.com/xiaochangtongxue/my-gin/internal/repository"
	"github.com/xiaochangtongxue/my-gin/pkg/permission"
)

// PermissionService 权限服务接口
type PermissionService interface {
	// 角色管理
	CreateRole(ctx context.Context, r *req.CreateRoleReq) (*resp.RoleResp, error)
	UpdateRole(ctx context.Context, id uint64, r *req.UpdateRoleReq) error
	DeleteRole(ctx context.Context, id uint64) error
	GetRole(ctx context.Context, id uint64) (*resp.RoleResp, error)
	ListRoles(ctx context.Context, r *req.ListRolesReq) ([]*resp.RoleResp, error)

	// 角色权限管理
	AddPermission(ctx context.Context, roleID uint64, r *req.AddPermissionReq) error
	RemovePermission(ctx context.Context, roleID uint64, resource, action string) error
	GetRolePermissions(ctx context.Context, roleID uint64) ([]*resp.PermissionResp, error)

	// 用户角色管理
	AssignRole(ctx context.Context, userID, roleID uint64) error
	RemoveRole(ctx context.Context, userID, roleID uint64) error
	GetUserRoles(ctx context.Context, userID uint64) ([]*resp.RoleResp, error)
}

// permissionService 权限服务实现
type permissionService struct {
	db           *gorm.DB
	roleRepo     repository.RoleRepository
	userRoleRepo repository.UserRoleRepository
	policyMgr    permission.PolicyManager
}

// NewPermissionService 创建权限服务
func NewPermissionService(
	db *gorm.DB,
	roleRepo repository.RoleRepository,
	userRoleRepo repository.UserRoleRepository,
	policyMgr permission.PolicyManager,
) PermissionService {
	return &permissionService{
		db:           db,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		policyMgr:    policyMgr,
	}
}

// CreateRole 创建角色
func (s *permissionService) CreateRole(ctx context.Context, r *req.CreateRoleReq) (*resp.RoleResp, error) {
	// 检查编码是否已存在
	exist, err := s.roleRepo.ExistsByCode(ctx, r.Code)
	if err != nil {
		return nil, fmt.Errorf("检查角色编码失败: %w", err)
	}
	if exist {
		return nil, errors.New("角色编码已存在")
	}

	role := &model.Role{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		Status:      1,
		BuiltIn:     0,
		SortOrder:   r.SortOrder,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, fmt.Errorf("创建角色失败: %w", err)
	}

	return s.toRoleResp(role), nil
}

// UpdateRole 更新角色
func (s *permissionService) UpdateRole(ctx context.Context, id uint64, r *req.UpdateRoleReq) error {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}

	// 内置角色不允许修改编码
	if role.BuiltIn == 1 && role.Code != r.Code {
		return errors.New("内置角色不允许修改编码")
	}

	// 检查新编码是否被其他角色使用
	if r.Code != role.Code {
		exist, err := s.roleRepo.ExistsByCode(ctx, r.Code)
		if err == nil && exist {
			return errors.New("角色编码已存在")
		}
	}

	// 更新字段
	role.Code = r.Code
	role.Name = r.Name
	role.Description = r.Description
	if r.Status != nil {
		role.Status = *r.Status
	}
	role.SortOrder = r.SortOrder

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return fmt.Errorf("更新角色失败: %w", err)
	}

	return nil
}

// DeleteRole 删除角色
func (s *permissionService) DeleteRole(ctx context.Context, id uint64) error {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}

	// 内置角色不允许删除
	if role.BuiltIn == 1 {
		return errors.New("内置角色不允许删除")
	}

	// 使用事务处理
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除 Casbin 策略
		roleIDStr := strconv.FormatUint(role.ID, 10)
		if err := s.policyMgr.RemoveFilteredPolicy(ctx, roleIDStr); err != nil {
			return fmt.Errorf("删除策略失败: %w", err)
		}

		// 2. 删除用户角色关联
		if err := s.userRoleRepo.DeleteByRoleID(ctx, role.ID); err != nil {
			return fmt.Errorf("删除用户角色关联失败: %w", err)
		}

		// 3. 删除角色
		if err := s.roleRepo.Delete(ctx, role.ID); err != nil {
			return fmt.Errorf("删除角色失败: %w", err)
		}

		return nil
	})
}

// GetRole 获取角色详情
func (s *permissionService) GetRole(ctx context.Context, id uint64) (*resp.RoleResp, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	return s.toRoleResp(role), nil
}

// ListRoles 获取角色列表
func (s *permissionService) ListRoles(ctx context.Context, r *req.ListRolesReq) ([]*resp.RoleResp, error) {
	roles, err := s.roleRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询角色列表失败: %w", err)
	}

	result := make([]*resp.RoleResp, 0, len(roles))
	for _, role := range roles {
		// 过滤条件
		if r.Status != nil && role.Status != *r.Status {
			continue
		}
		if r.Code != "" && role.Code != r.Code {
			continue
		}
		result = append(result, s.toRoleResp(role))
	}

	return result, nil
}

// AddPermission 添加权限
func (s *permissionService) AddPermission(ctx context.Context, roleID uint64, r *req.AddPermissionReq) error {
	// 检查角色是否存在
	if _, err := s.roleRepo.FindByID(ctx, roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}

	roleIDStr := strconv.FormatUint(roleID, 10)

	// 检查权限是否已存在
	if s.policyMgr.HasPolicy(ctx, roleIDStr, r.Resource, r.Action) {
		return errors.New("权限已存在")
	}

	// 添加策略
	policy := &permission.Policy{
		Type:     "p",
		Subject:  roleIDStr,
		Resource: r.Resource,
		Action:   r.Action,
		Effect:   "allow",
	}
	if err := s.policyMgr.AddPolicy(ctx, policy); err != nil {
		return fmt.Errorf("添加策略失败: %w", err)
	}

	return nil
}

// RemovePermission 删除权限
func (s *permissionService) RemovePermission(ctx context.Context, roleID uint64, resource, action string) error {
	roleIDStr := strconv.FormatUint(roleID, 10)
	policy := &permission.Policy{
		Type:     "p",
		Subject:  roleIDStr,
		Resource: resource,
		Action:   action,
	}
	if err := s.policyMgr.RemovePolicy(ctx, policy); err != nil {
		return fmt.Errorf("删除策略失败: %w", err)
	}
	return nil
}

// GetRolePermissions 获取角色权限列表
func (s *permissionService) GetRolePermissions(ctx context.Context, roleID uint64) ([]*resp.PermissionResp, error) {
	// 检查角色是否存在
	if _, err := s.roleRepo.FindByID(ctx, roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	roleIDStr := strconv.FormatUint(roleID, 10)
	policies, err := s.policyMgr.GetPolicies(ctx, &permission.PolicyFilter{Subject: roleIDStr})
	if err != nil {
		return nil, fmt.Errorf("获取策略失败: %w", err)
	}

	result := make([]*resp.PermissionResp, 0, len(policies))
	for _, policy := range policies {
		result = append(result, &resp.PermissionResp{
			Resource: policy.Resource,
			Action:   policy.Action,
		})
	}

	return result, nil
}

// AssignRole 分配角色给用户
func (s *permissionService) AssignRole(ctx context.Context, userID, roleID uint64) error {
	// 检查角色是否存在
	if _, err := s.roleRepo.FindByID(ctx, roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}

	// 检查是否已分配
	exist, err := s.userRoleRepo.Exists(ctx, userID, roleID)
	if err != nil {
		return fmt.Errorf("检查用户角色失败: %w", err)
	}
	if exist {
		return errors.New("用户已拥有该角色")
	}

	if err := s.userRoleRepo.Assign(ctx, userID, roleID); err != nil {
		return fmt.Errorf("分配角色失败: %w", err)
	}

	return nil
}

// RemoveRole 移除用户角色
func (s *permissionService) RemoveRole(ctx context.Context, userID, roleID uint64) error {
	if err := s.userRoleRepo.Delete(ctx, userID, roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户角色关联不存在")
		}
		return fmt.Errorf("移除角色失败: %w", err)
	}
	return nil
}

// GetUserRoles 获取用户角色列表
func (s *permissionService) GetUserRoles(ctx context.Context, userID uint64) ([]*resp.RoleResp, error) {
	roles, err := s.userRoleRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户角色失败: %w", err)
	}

	result := make([]*resp.RoleResp, len(roles))
	for i, role := range roles {
		result[i] = s.toRoleResp(role)
	}

	return result, nil
}

// toRoleResp 转换为角色响应
func (s *permissionService) toRoleResp(role *model.Role) *resp.RoleResp {
	return &resp.RoleResp{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		Status:      role.Status,
		BuiltIn:     role.BuiltIn,
		SortOrder:   role.SortOrder,
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
	}
}
