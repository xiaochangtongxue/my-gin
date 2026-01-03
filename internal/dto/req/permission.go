// Package req 请求 DTO
package req

// CreateRoleReq 创建角色请求
type CreateRoleReq struct {
	Code        string `json:"code" binding:"required,min=1,max=50" label:"角色编码" example:"editor"` // 角色编码
	Name        string `json:"name" binding:"required,min=1,max=50" label:"角色名称" example:"编辑"`     // 角色名称
	Description string `json:"description" binding:"max=255" label:"角色描述" example:"内容编辑人员"`        // 角色描述
	SortOrder   int    `json:"sort_order" binding:"min=0" label:"排序" example:"10"`                 // 排序
}

// UpdateRoleReq 更新角色请求
type UpdateRoleReq struct {
	Code        string `json:"code" binding:"required,min=1,max=50" label:"角色编码" example:"editor"` // 角色编码
	Name        string `json:"name" binding:"required,min=1,max=50" label:"角色名称" example:"编辑"`     // 角色名称
	Description string `json:"description" binding:"max=255" label:"角色描述" example:"内容编辑人员"`        // 角色描述
	Status      *int8  `json:"status" binding:"omitempty,oneof=0 1" label:"状态" example:"1"`        // 状态：0禁用 1启用
	SortOrder   int    `json:"sort_order" binding:"min=0" label:"排序" example:"10"`                 // 排序
}

// AddPermissionReq 添加权限请求
type AddPermissionReq struct {
	Resource string `json:"resource" binding:"required" label:"资源路径" example:"/api/v1/articles"`                   // 资源路径
	Action   string `json:"action" binding:"required,oneof=read create update delete *" label:"操作" example:"read"` // 操作
}

// AssignRoleReq 分配角色请求
type AssignRoleReq struct {
	RoleID uint64 `json:"role_id" binding:"required" label:"角色ID" example:"3"` // 角色ID
}

// ListRolesReq 查询角色列表请求
type ListRolesReq struct {
	Status *int8  `form:"status" label:"状态" example:"1"`      // 状态筛选（可选）
	Code   string `form:"code" label:"角色编码" example:"editor"` // 编码筛选（可选）
}
