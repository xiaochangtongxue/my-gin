package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/jwt"
)

// RouteModule is implemented by vertical business modules.
type RouteModule interface {
	RegisterRoutes(*gin.Engine)
}

// App is the runtime dependency container.
type App struct {
	Config  *config.Config
	DB      *gorm.DB
	Cache   cache.Cache
	JWTMgr  *jwt.Manager
	Engine  *gin.Engine
	Modules []RouteModule
}

// New creates an application container.
func New(
	cfg *config.Config,
	db *gorm.DB,
	cache cache.Cache,
	jwtMgr *jwt.Manager,
	engine *gin.Engine,
	modules []RouteModule,
) *App {
	return &App{
		Config:  cfg,
		DB:      db,
		Cache:   cache,
		JWTMgr:  jwtMgr,
		Engine:  engine,
		Modules: modules,
	}
}

// RegisterRoutes registers every module route in deterministic order.
func (a *App) RegisterRoutes() {
	for _, module := range a.Modules {
		module.RegisterRoutes(a.Engine)
	}
}
