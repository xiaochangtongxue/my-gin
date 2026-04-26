package health

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/xiaochangtongxue/my-gin/internal/handler"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
)

// Module owns health check routes.
type Module struct {
	handler *handler.HealthHandler
}

// ProviderSet contains health module dependencies.
var ProviderSet = wire.NewSet(ProvideHealthHandler, NewModule)

// NewModule creates a health module.
func NewModule(handler *handler.HealthHandler) *Module {
	return &Module{handler: handler}
}

// RegisterRoutes registers health check routes.
func (m *Module) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", m.handler.Check)
	r.GET("/health/live", m.handler.Live)
	r.GET("/health/ready", m.handler.Ready)
}

// ProvideHealthHandler provides the health HTTP handler.
func ProvideHealthHandler(db *gorm.DB, cache cache.Cache) *handler.HealthHandler {
	return handler.NewHealthHandler(db, cache)
}
