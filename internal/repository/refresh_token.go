package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/xiaochangtongxue/my-gin/internal/model"
)

// RefreshTokenRepository Refresh Token 存储接口
type RefreshTokenRepository interface {
	// Create 创建 Refresh Token
	Create(ctx context.Context, rt *model.RefreshToken) error
	// FindByToken 根据 Token 查找
	FindByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	// IsValidToken 检查 Refresh Token 是否有效（存在且未过期）
	IsValidToken(ctx context.Context, token string) (bool, error)
	// Delete 删除 Refresh Token
	Delete(ctx context.Context, id uint) error
	// DeleteByToken 根据 Token 删除
	DeleteByToken(ctx context.Context, token string) error
	// DeleteByUser 删除用户的所有 Refresh Token
	DeleteByUser(ctx context.Context, userID uint) error
	// DeleteOldAndCreate 删除旧的 Token 并创建新的（原子操作，返回新 RT）
	DeleteOldAndCreate(ctx context.Context, oldToken string, newRT *model.RefreshToken) (*model.RefreshToken, error)
	// CleanupExpired 清理过期的 Token
	CleanupExpired(ctx context.Context) error
}

// refreshTokenRepository Refresh Token 存储实现
type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository 创建 RefreshToken Repository/ 用户名可以从 Claims 中获取，但这里简化处理
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create 创建 Refresh Token
func (r *refreshTokenRepository) Create(ctx context.Context, rt *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(rt).Error
}

// FindByToken 根据 Token 查找
func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// IsValidToken 检查 Refresh Token 是否有效（存在且未过期）
func (r *refreshTokenRepository) IsValidToken(ctx context.Context, token string) (bool, error) {
	var valid bool
	err := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Select("1").
		Where("token = ? AND expires_at > ?", token, time.Now()).
		Limit(1).
		Scan(&valid).Error
	if err != nil {
		return false, err
	}
	return valid, nil
}

// Delete 删除 Refresh Token（软删除）
func (r *refreshTokenRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.RefreshToken{}, id).Error
}

// DeleteByToken 根据 Token 删除（软删除）
func (r *refreshTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&model.RefreshToken{}).Error
}

// DeleteByUser 删除用户的所有 Refresh Token
func (r *refreshTokenRepository) DeleteByUser(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.RefreshToken{}).Error
}

// DeleteOldAndCreate 删除旧的 Token 并创建新的（事务原子操作）
// 返回新创建的 RefreshToken
func (r *refreshTokenRepository) DeleteOldAndCreate(ctx context.Context, oldToken string, newRT *model.RefreshToken) (*model.RefreshToken, error) {
	var result *model.RefreshToken
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除旧的 Token
		if err := tx.Where("token = ?", oldToken).Delete(&model.RefreshToken{}).Error; err != nil {
			return err
		}
		// 创建新的 Token
		if err := tx.Create(newRT).Error; err != nil {
			return err
		}
		result = newRT
		return nil
	})
	return result, err
}

// CleanupExpired 清理过期的 Refresh Token
func (r *refreshTokenRepository) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&model.RefreshToken{}).Error
}
