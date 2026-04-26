package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Module owns Prometheus metrics routes.
type Module struct{}

// ProviderSet contains metrics module dependencies.
var ProviderSet = wire.NewSet(NewModule)

// NewModule creates a metrics module.
func NewModule() *Module {
	return &Module{}
}

// RegisterRoutes registers metrics routes.
func (m *Module) RegisterRoutes(r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
