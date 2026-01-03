// Package model 数据模型
package model

import (
	"time"

	"gorm.io/gorm"
)

// UserRole 用户角色关联模型
type UserRole struct {
	ID        uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint64         `json:"user_id" gorm:"not null;index:idx_user_id,priority:1"`
	RoleID    uint64         `json:"role_id" gorm:"not null;index:idx_role_id,priority:1"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:datetime(3)"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:datetime(3)"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index:idx_deleted_at"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}