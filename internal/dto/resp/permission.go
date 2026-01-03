// Package resp 响应 DTO
package resp

// RoleResp 角色响应
type RoleResp struct {
	ID          uint64 `json:"id" example:"1"`                            // 角色ID
	Code        string `json:"code" example:"editor"`                     // 角色编码
	Name        string `json:"name" example:"编辑"`                         // 角色名称
	Description string `json:"description" example:"内容编辑人员"`              // 角色描述
	Status      int8   `json:"status" example:"1"`                        // 状态：0禁用 1启用
	BuiltIn     int8   `json:"built_in" example:"0"`                      // 是否内置角色
	SortOrder   int    `json:"sort_order" example:"10"`                   // 排序
	CreatedAt   string `json:"created_at" example:"2024-01-01T00:00:00Z"` // 创建时间
}

// RoleWithPermissionsResp 带权限的角色响应
type RoleWithPermissionsResp struct {
	RoleResp
	Permissions []*PermissionResp `json:"permissions"` // 权限列表
}

// PermissionResp 权限响应
type PermissionResp struct {
	ID       uint64 `json:"id" example:"1"`                      // 策略ID
	Resource string `json:"resource" example:"/api/v1/articles"` // 资源路径
	Action   string `json:"action" example:"read"`               // 操作
}

// UserRolesResp 用户角色响应
type UserRolesResp struct {
	UserID uint64      `json:"user_id" example:"10001"` // 用户ID
	Roles  []*RoleResp `json:"roles"`                   // 角色列表
}
