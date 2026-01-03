// Package permission 权限模块工厂方法
package permission

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/xiaochangtongxue/my-gin/pkg/config"
)

// NewChecker 根据配置创建权限检查器
// 支持多种权限模型：rbac, abac, acl
func NewChecker(cfg *config.Config, db *gorm.DB) (PermissionChecker, error) {
	switch cfg.Permission.Model {
	case "rbac":
		return NewRBACChecker(db, cfg.Permission.ModelFile)
	// 未来可扩展其他模型
	// case "abac":
	// 	return NewABACChecker(db, cfg.Permission.ModelFile)
	// case "acl":
	// 	return NewACLChecker(db, cfg.Permission.ModelFile)
	default:
		// 默认使用 RBAC
		return NewRBACChecker(db, "pkg/permission/model.conf")
	}
}

// MustNewChecker 创建检查器，panic on error
func MustNewChecker(cfg *config.Config, db *gorm.DB) PermissionChecker {
	checker, err := NewChecker(cfg, db)
	if err != nil {
		panic(fmt.Sprintf("创建权限检查器失败: %v", err))
	}
	return checker
}
