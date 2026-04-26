package security

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/xiaochangtongxue/my-gin/internal/handler"
	"github.com/xiaochangtongxue/my-gin/internal/middleware"
	rbac "github.com/xiaochangtongxue/my-gin/internal/permission"
	"github.com/xiaochangtongxue/my-gin/internal/repository"
	"github.com/xiaochangtongxue/my-gin/internal/service"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
)

// Module owns security management routes.
type Module struct {
	handler      *handler.SecurityHandler
	checker      rbac.PermissionChecker
	userRoleRepo repository.UserRoleRepository
}

// ProviderSet contains security module dependencies.
var ProviderSet = wire.NewSet(
	ProvideSecurityService,
	ProvideSecurityHandler,
	NewModule,
)

// NewModule creates a security module.
func NewModule(
	handler *handler.SecurityHandler,
	checker rbac.PermissionChecker,
	userRoleRepo repository.UserRoleRepository,
) *Module {
	return &Module{handler: handler, checker: checker, userRoleRepo: userRoleRepo}
}

// RegisterRoutes registers security administration routes.
func (m *Module) RegisterRoutes(r *gin.Engine) {
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.LoadUserRoles(m.userRoleRepo))
	admin.POST(
		"/unlock-account",
		middleware.PermissionRequired(m.checker, "/api/v1/admin/unlock-account", rbac.ActionUpdate),
		m.handler.UnlockAccount,
	)
}

// ProvideSecurityService provides the security service.
func ProvideSecurityService(cache cache.Cache, cfg *config.Config) service.SecurityService {
	return service.NewSecurityService(cache, service.SecurityConfig{
		MaxAttempts:       cfg.Bruteforce.MaxAttempts,
		LockDuration:      cfg.Bruteforce.LockDuration,
		Window:            cfg.Bruteforce.Window,
		BlacklistDuration: cfg.Bruteforce.BlacklistDuration,
	})
}

// ProvideSecurityHandler provides the security HTTP handler.
func ProvideSecurityHandler(securitySvc service.SecurityService) *handler.SecurityHandler {
	return handler.NewSecurityHandler(securitySvc)
}
