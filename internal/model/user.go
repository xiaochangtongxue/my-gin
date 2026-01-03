package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"-"`           // 内部ID，不对外暴露
	UID       uint64         `gorm:"not null;uniqueIndex;index:idx_uid" json:"uid"` // 对外展示的用户ID（14位纯数字）
	Username  string         `gorm:"type:varchar(20);not null;uniqueIndex" json:"username"` // 用户名
	Mobile    string         `gorm:"type:char(11);not null;uniqueIndex;index:idx_mobile" json:"mobile"` // 手机号
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`           // bcrypt哈希，不对外暴露
	CreatedAt time.Time      `gorm:"type:datetime(3)" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime(3)" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}