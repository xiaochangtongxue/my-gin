package docs

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/xiaochangtongxue/my-gin/internal/handler"
)

// Module owns API documentation routes.
type Module struct{}

// ProviderSet contains docs module dependencies.
var ProviderSet = wire.NewSet(NewModule)

// NewModule creates a docs module.
func NewModule() *Module {
	return &Module{}
}

// RegisterRoutes registers Swagger routes.
func (m *Module) RegisterRoutes(r *gin.Engine) {
	handler.NewSwaggerHandler().Register(r)
}
