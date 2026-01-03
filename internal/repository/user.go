package repository

import (
	"context"

	"gorm.io/gorm"
	"github.com/xiaochangtongxue/my-gin/internal/model"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *model.User) error

	// FindByID 根据内部ID查找用户
	FindByID(ctx context.Context, id uint64) (*model.User, error)

	// FindByUID 根据对外UID查找用户
	FindByUID(ctx context.Context, uid uint64) (*model.User, error)

	// FindByMobile 根据手机号查找用户
	FindByMobile(ctx context.Context, mobile string) (*model.User, error)

	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*model.User, error)

	// UpdateUID 更新用户的UID
	UpdateUID(ctx context.Context, id uint64, uid uint64) error

	// ExistsByMobile 检查手机号是否已存在
	ExistsByMobile(ctx context.Context, mobile string) (bool, error)

	// ExistsByUsername 检查用户名是否已存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID 根据内部ID查找用户
func (r *userRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUID 根据对外UID查找用户
func (r *userRepository) FindByUID(ctx context.Context, uid uint64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByMobile 根据手机号查找用户
func (r *userRepository) FindByMobile(ctx context.Context, mobile string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("mobile = ?", mobile).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUID 更新用户的UID
func (r *userRepository) UpdateUID(ctx context.Context, id uint64, uid uint64) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("uid", uid).Error
}

// ExistsByMobile 检查手机号是否已存在
func (r *userRepository) ExistsByMobile(ctx context.Context, mobile string) (bool, error) {
	var exists int
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Select("1").
		Where("mobile = ?", mobile).
		First(&exists).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return err == nil, err
}

// ExistsByUsername 检查用户名是否已存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists int
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Select("1").
		Where("username = ?", username).
		First(&exists).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return err == nil, err
}