// Package repository 角色仓储
package repository

import (
	"context"

	"gorm.io/gorm"
	"github.com/xiaochangtongxue/my-gin/internal/model"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	// Create 创建角色
	Create(ctx context.Context, role *model.Role) error

	// Update 更新角色
	Update(ctx context.Context, role *model.Role) error

	// Delete 删除角色（软删除）
	Delete(ctx context.Context, id uint64) error

	// FindByID 根据ID查找角色
	FindByID(ctx context.Context, id uint64) (*model.Role, error)

	// FindByCode 根据编码查找角色
	FindByCode(ctx context.Context, code string) (*model.Role, error)

	// List 获取角色列表
	List(ctx context.Context) ([]*model.Role, error)

	// ExistsByCode 检查编码是否存在
	ExistsByCode(ctx context.Context, code string) (bool, error)
}

// roleRepository 角色仓储实现
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓储
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// Create 创建角色
func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// Update 更新角色
func (r *roleRepository) Update(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Model(&model.Role{}).Where("id = ?", role.ID).Updates(role).Error
}

// Delete 删除角色（软删除）
func (r *roleRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

// FindByID 根据ID查找角色
func (r *roleRepository) FindByID(ctx context.Context, id uint64) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByCode 根据编码查找角色
func (r *roleRepository) FindByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List 获取角色列表
func (r *roleRepository) List(ctx context.Context) ([]*model.Role, error) {
	var roles []*model.Role
	err := r.db.WithContext(ctx).
		Order("sort_order ASC").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// ExistsByCode 检查编码是否存在
func (r *roleRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists int
	err := r.db.WithContext(ctx).Model(&model.Role{}).
		Select("1").
		Where("code = ?", code).
		First(&exists).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return err == nil, err
}