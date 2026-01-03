// Package repository 用户角色仓储
package repository

import (
	"context"

	"gorm.io/gorm"
	"github.com/xiaochangtongxue/my-gin/internal/model"
)

// UserRoleRepository 用户角色仓储接口
type UserRoleRepository interface {
	// Create 创建用户角色关联
	Create(ctx context.Context, userRole *model.UserRole) error

	// Delete 删除用户角色关联（软删除）
	Delete(ctx context.Context, userID, roleID uint64) error

	// DeleteByUserID 删除用户的所有角色
	DeleteByUserID(ctx context.Context, userID uint64) error

	// DeleteByRoleID 删除角色的所有用户关联
	DeleteByRoleID(ctx context.Context, roleID uint64) error

	// FindByUserID 查找用户的所有角色
	FindByUserID(ctx context.Context, userID uint64) ([]*model.Role, error)

	// FindRoleIDsByUserID 查找用户的角色ID列表
	FindRoleIDsByUserID(ctx context.Context, userID uint64) ([]uint64, error)

	// FindUserIDsByRoleID 查找角色下的所有用户ID
	FindUserIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error)

	// Exists 检查用户角色关联是否存在
	Exists(ctx context.Context, userID, roleID uint64) (bool, error)

	// Assign 分配角色给用户
	Assign(ctx context.Context, userID, roleID uint64) error
}

// userRoleRepository 用户角色仓储实现
type userRoleRepository struct {
	db *gorm.DB
}

// NewUserRoleRepository 创建用户角色仓储
func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}

// Create 创建用户角色关联
func (r *userRoleRepository) Create(ctx context.Context, userRole *model.UserRole) error {
	return r.db.WithContext(ctx).Create(userRole).Error
}

// Delete 删除用户角色关联（软删除）
func (r *userRoleRepository) Delete(ctx context.Context, userID, roleID uint64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&model.UserRole{}).Error
}

// DeleteByUserID 删除用户的所有角色
func (r *userRoleRepository) DeleteByUserID(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.UserRole{}).Error
}

// DeleteByRoleID 删除角色的所有用户关联
func (r *userRoleRepository) DeleteByRoleID(ctx context.Context, roleID uint64) error {
	return r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Delete(&model.UserRole{}).Error
}

// FindByUserID 查找用户的所有角色
func (r *userRoleRepository) FindByUserID(ctx context.Context, userID uint64) ([]*model.Role, error) {
	var roles []*model.Role
	err := r.db.WithContext(ctx).
		Table("roles").
		Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.status = 1", userID).
		Order("roles.sort_order ASC").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// FindRoleIDsByUserID 查找用户的角色ID列表
func (r *userRoleRepository) FindRoleIDsByUserID(ctx context.Context, userID uint64) ([]uint64, error) {
	var roleIDs []uint64
	err := r.db.WithContext(ctx).
		Model(&model.UserRole{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

// FindUserIDsByRoleID 查找角色下的所有用户ID
func (r *userRoleRepository) FindUserIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error) {
	var userIDs []uint64
	err := r.db.WithContext(ctx).
		Model(&model.UserRole{}).
		Where("role_id = ?", roleID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}

// Exists 检查用户角色关联是否存在
func (r *userRoleRepository) Exists(ctx context.Context, userID, roleID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error
	return count > 0, err
}

// Assign 分配角色给用户
func (r *userRoleRepository) Assign(ctx context.Context, userID, roleID uint64) error {
	userRole := &model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.WithContext(ctx).Create(userRole).Error
}