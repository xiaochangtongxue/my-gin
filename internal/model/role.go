// Package model 数据模型
package model

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Code        string         `json:"code" gorm:"type:varchar(50);not null;uniqueIndex:uk_code"`
	Name        string         `json:"name" gorm:"type:varchar(50);not null"`
	Description string         `json:"description" gorm:"type:varchar(255);default:''"`
	Status      int8           `json:"status" gorm:"default:1;index:idx_status"`
	BuiltIn     int8           `json:"built_in" gorm:"default:0"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:datetime(3)"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:datetime(3)"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index:idx_deleted_at"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}