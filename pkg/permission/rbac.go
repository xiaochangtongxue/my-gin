// Package permission RBAC 权限模型实现
package permission

import (
	"context"
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// RBACChecker RBAC 权限检查器
type RBACChecker struct {
	enforcer *casbin.Enforcer
}

// 确保 RBACChecker 实现了接口
var _ PermissionChecker = (*RBACChecker)(nil)
var _ PolicyManager = (*RBACChecker)(nil)

// NewRBACChecker 创建 RBAC 检查器
func NewRBACChecker(db *gorm.DB, modelFile string) (*RBACChecker, error) {
	if modelFile == "" {
		modelFile = "pkg/permission/model.conf"
	}

	// 禁用 Casbin 自动建表，表结构由迁移文件管理（包含时间戳字段）
	gormadapter.TurnOffAutoMigrate(db)

	// 使用 GORM 适配器，指定表名为 casbin_rule
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rule")
	if err != nil {
		return nil, fmt.Errorf("创建 Casbin 适配器失败: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(modelFile, adapter)
	if err != nil {
		return nil, fmt.Errorf("创建 Casbin Enforcer 失败: %w", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("加载策略失败: %w", err)
	}

	return &RBACChecker{enforcer: enforcer}, nil
}

// Check 检查权限
func (c *RBACChecker) Check(ctx context.Context, req *CheckRequest) (bool, error) {
	// 构建资源路径
	resourcePath := c.buildResourcePath(req.Resource)

	// 遍历用户的所有角色
	for _, roleID := range req.Subject.Roles {
		// Casbin 检查：role_id, resource, action
		allowed, err := c.enforcer.Enforce(roleID, resourcePath, string(req.Action))
		if err != nil {
			return false, fmt.Errorf("权限检查失败: %w", err)
		}
		if allowed {
			return true, nil
		}
	}

	return false, nil
}

// BatchCheck 批量检查权限
func (c *RBACChecker) BatchCheck(ctx context.Context, requests []*CheckRequest) ([]bool, error) {
	results := make([]bool, len(requests))
	for i, req := range requests {
		allowed, err := c.Check(ctx, req)
		if err != nil {
			return nil, err
		}
		results[i] = allowed
	}
	return results, nil
}

// Name 返回权限模型名称
func (c *RBACChecker) Name() string {
	return "rbac"
}

// buildResourcePath 构建资源路径
func (c *RBACChecker) buildResourcePath(res *Resource) string {
	if res.ID != "" && strings.Contains(res.Path, ":id") {
		// 替换路径中的 :id 占位符
		return strings.Replace(res.Path, ":id", res.ID, 1)
	}
	return res.Path
}

// GetEnforcer 获取 Casbin Enforcer（用于策略管理）
func (c *RBACChecker) GetEnforcer() *casbin.Enforcer {
	return c.enforcer
}

// AddPolicy 添加策略（实现 PolicyManager 接口）
func (c *RBACChecker) AddPolicy(ctx context.Context, policy *Policy) error {
	_, err := c.enforcer.AddPolicy(policy.Subject, policy.Resource, policy.Action)
	return err
}

// RemovePolicy 删除策略（实现 PolicyManager 接口）
func (c *RBACChecker) RemovePolicy(ctx context.Context, policy *Policy) error {
	_, err := c.enforcer.RemovePolicy(policy.Subject, policy.Resource, policy.Action)
	return err
}

// GetPolicies 获取策略列表（实现 PolicyManager 接口）
func (c *RBACChecker) GetPolicies(ctx context.Context, filter *PolicyFilter) ([]*Policy, error) {
	var policies [][]string
	var err error
	if filter != nil && filter.Subject != "" {
		policies, err = c.enforcer.GetFilteredPolicy(0, filter.Subject)
	} else {
		policies, err = c.enforcer.GetPolicy()
	}
	if err != nil {
		return nil, err
	}

	result := make([]*Policy, 0, len(policies))
	for _, p := range policies {
		if len(p) >= 3 {
			result = append(result, &Policy{
				Type:     "p",
				Subject:  p[0],
				Resource: p[1],
				Action:   p[2],
				Effect:   "allow",
			})
		}
	}
	return result, nil
}

// RemoveFilteredPolicy 按主体删除策略（实现 PolicyManager 接口）
func (c *RBACChecker) RemoveFilteredPolicy(ctx context.Context, subject string) error {
	_, err := c.enforcer.RemoveFilteredPolicy(0, subject)
	return err
}

// HasPolicy 检查策略是否存在（实现 PolicyManager 接口）
func (c *RBACChecker) HasPolicy(ctx context.Context, subject, resource, action string) bool {
	has, _ := c.enforcer.HasPolicy(subject, resource, action)
	return has
}

// ReloadPolicy 重新从数据库加载策略（外部修改数据库后需要调用）
func (c *RBACChecker) ReloadPolicy() error {
	return c.enforcer.LoadPolicy()
}
