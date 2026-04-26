package permission

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/xiaochangtongxue/my-gin/internal/handler"
	"github.com/xiaochangtongxue/my-gin/internal/middleware"
	rbac "github.com/xiaochangtongxue/my-gin/internal/permission"
	"github.com/xiaochangtongxue/my-gin/internal/repository"
	"github.com/xiaochangtongxue/my-gin/internal/service"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
)

// Module owns permission management routes.
type Module struct {
	handler      *handler.PermissionHandler
	checker      rbac.PermissionChecker
	userRoleRepo repository.UserRoleRepository
}

// ProviderSet contains permission module dependencies.
var ProviderSet = wire.NewSet(
	ProvideRoleRepository,
	ProvideUserRoleRepository,
	ProvidePermissionChecker,
	ProvidePermissionService,
	ProvidePermissionHandler,
	NewModule,
)

// NewModule creates a permission module.
func NewModule(
	handler *handler.PermissionHandler,
	checker rbac.PermissionChecker,
	userRoleRepo repository.UserRoleRepository,
) *Module {
	return &Module{handler: handler, checker: checker, userRoleRepo: userRoleRepo}
}

// RegisterRoutes registers RBAC management routes.
func (m *Module) RegisterRoutes(r *gin.Engine) {
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.LoadUserRoles(m.userRoleRepo))

	admin.POST("/roles",
		middleware.PermissionRequired(m.checker, "/api/v1/admin/roles", rbac.ActionCreate),
		m.handler.CreateRole,
	)
	admin.GET("/roles",
		middleware.PermissionRequired(m.checker, "/api/v1/admin/roles", rbac.ActionRead),
		m.handler.ListRoles,
	)
	admin.GET("/roles/:id",
		middleware.PermissionRequired(m.checker, "/api/v1/admin/roles", rbac.ActionRead),
		m.handler.GetRole,
	)
	admin.PUT("/roles/:id",
		middleware.PermissionRequired(m.checker, "/api/v1/admin/roles", rbac.ActionUpdate),
		m.handler.UpdateRole,
	)
	admin.DELETE("/roles/:id",
		middleware.PermissionRequired(m.checker, "/api/v1/admin/roles", rbac.ActionDelete),
		m.handler.DeleteRole,
	)
	admin.GET("/roles/:id/permissions",
		middleware.PermissionRequired(m.checker, "/api/v1/admin/roles", rbac.ActionRead),
		m.handler.GetRolePermissions,
	)
	admin.POST("/roles/:id/permissions",
		middleware.PermissionRequired(m.checker, "/api/v1/admin/roles", rbac.ActionCreate),
		m.handler.AddPermission,
	)
	admin.DELETE("/roles/:id/permissions",
		middleware.PermissionRequired(m.checker, "/api/v1/admin/roles", rbac.ActionDelete),
		m.handler.RemovePermission,
	)
	admin.GET("/users/:id/roles",
		middleware.PermissionRequired(m.checker, "/api/v1/admin/users", rbac.ActionRead),
		m.handler.GetUserRoles,
	)
	admin.POST("/users/:id/roles",
		middleware.PermissionRequired(m.checker, "/api/v1/admin/users", rbac.ActionUpdate),
		m.handler.AssignRole,
	)
	admin.DELETE("/users/:id/roles/:roleId",
		middleware.PermissionRequired(m.checker, "/api/v1/admin/users", rbac.ActionUpdate),
		m.handler.RemoveRole,
	)
}

// ProvideRoleRepository provides role storage.
func ProvideRoleRepository(db *gorm.DB) repository.RoleRepository {
	return repository.NewRoleRepository(db)
}

// ProvideUserRoleRepository provides user-role storage.
func ProvideUserRoleRepository(db *gorm.DB) repository.UserRoleRepository {
	return repository.NewUserRoleRepository(db)
}

// ProvidePermissionChecker provides the configured permission checker.
func ProvidePermissionChecker(cfg *config.Config, db *gorm.DB) (rbac.PermissionChecker, error) {
	return rbac.NewChecker(cfg, db)
}

// ProvidePermissionService provides the permission service.
func ProvidePermissionService(
	db *gorm.DB,
	roleRepo repository.RoleRepository,
	userRoleRepo repository.UserRoleRepository,
	checker rbac.PermissionChecker,
) service.PermissionService {
	policyMgr := checker.(rbac.PolicyManager)
	return service.NewPermissionService(db, roleRepo, userRoleRepo, policyMgr)
}

// ProvidePermissionHandler provides the permission HTTP handler.
func ProvidePermissionHandler(permissionSvc service.PermissionService) *handler.PermissionHandler {
	return handler.NewPermissionHandler(permissionSvc)
}
