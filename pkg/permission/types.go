// Package permission 权限检查接口定义
// 支持多种权限模型：RBAC、ABAC、ACL
package permission

import "context"

// PermissionChecker 权限检查器接口
type PermissionChecker interface {
	Check(ctx context.Context, req *CheckRequest) (bool, error)
	BatchCheck(ctx context.Context, requests []*CheckRequest) ([]bool, error)
	Name() string
}

// PolicyManager 策略管理接口（可选）
type PolicyManager interface {
	AddPolicy(ctx context.Context, policy *Policy) error
	RemovePolicy(ctx context.Context, policy *Policy) error
	RemoveFilteredPolicy(ctx context.Context, subject string) error
	HasPolicy(ctx context.Context, subject, resource, action string) bool
	GetPolicies(ctx context.Context, filter *PolicyFilter) ([]*Policy, error)
}

// Action 操作类型
type Action string

const (
	ActionRead   Action = "read"
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionAny    Action = "*"
)

// CheckRequest 权限检查请求
type CheckRequest struct {
	Subject  *Subject  // 主体
	Resource *Resource // 资源
	Action   Action    // 操作
}

// Subject 主体
type Subject struct {
	ID    string   // 用户 ID
	Type  string   // 主体类型：user, service, etc.
	Roles []string // 角色ID列表（RBAC）
	Attrs Attrs    // 自定义属性（ABAC）
}

// Attrs 自定义属性
type Attrs map[string]string

// Resource 资源
type Resource struct {
	Type  string // 资源类型：article, user, etc.
	ID    string // 资源ID
	Path  string // 资源路径
	Owner string // 资源所有者ID（用于所有权检查）
	Attrs Attrs  // 自定义属性（ABAC）
}

// Policy 策略定义
type Policy struct {
	ID       string // 策略ID
	Type     string // 策略类型：p=权限, g=角色继承
	Subject  string // 主体：role_id, user_id, etc.
	Resource string // 资源
	Action   string // 操作
	Effect   string // 效果：allow, deny
	Attrs    Attrs  // 扩展属性
}

// PolicyFilter 策略过滤条件
type PolicyFilter struct {
	Type    string // 策略类型
	Subject string // 主体
}
