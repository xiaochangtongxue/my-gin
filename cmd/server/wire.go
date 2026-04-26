//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"

	apppkg "github.com/xiaochangtongxue/my-gin/internal/app"
	"github.com/xiaochangtongxue/my-gin/internal/middleware"
	authmodule "github.com/xiaochangtongxue/my-gin/internal/modules/auth"
	docsmodule "github.com/xiaochangtongxue/my-gin/internal/modules/docs"
	healthmodule "github.com/xiaochangtongxue/my-gin/internal/modules/health"
	metricsmodule "github.com/xiaochangtongxue/my-gin/internal/modules/metrics"
	permissionmodule "github.com/xiaochangtongxue/my-gin/internal/modules/permission"
	securitymodule "github.com/xiaochangtongxue/my-gin/internal/modules/security"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/database"
	"github.com/xiaochangtongxue/my-gin/pkg/jwt"
	"github.com/xiaochangtongxue/my-gin/pkg/metrics"
)

// ProvideConfig 提供配置
func ProvideConfig(configFile string) (*config.Config, error) {
	return config.Init(configFile)
}

// ProvideDB 提供数据库连接
func ProvideDB(cfg *config.Config) (*gorm.DB, error) {
	if err := database.Init(&cfg.Database); err != nil {
		return nil, err
	}
	return database.Get(), nil
}

// ProvideRedisCache 提供 Redis 缓存
func ProvideRedisCache(cfg *config.Config) (cache.Cache, error) {
	if err := cache.InitRedis(&cfg.Redis); err != nil {
		return nil, err
	}
	return cache.NewRedisCache(), nil
}

// ProvideJWTManager 提供 JWT 管理器
func ProvideJWTManager(cfg *config.Config) *jwt.Manager {
	mgr := jwt.NewManager(jwt.Config{
		Secret:     cfg.JWT.Secret,
		ExpireTime: cfg.JWT.ExpireTime,
		Issuer:     cfg.JWT.Issuer,
	})
	jwt.Init(jwt.Config{
		Secret:     cfg.JWT.Secret,
		ExpireTime: cfg.JWT.ExpireTime,
		Issuer:     cfg.JWT.Issuer,
	})
	return mgr
}

// ProvideGinEngine 提供 Gin 引擎
func ProvideGinEngine(cfg *config.Config, jwtMgr *jwt.Manager, redisCache cache.Cache) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()

	metrics.AppStartTime.SetToCurrentTime()
	middleware.ApplyWithConfig(engine, cfg, middleware.DefaultMiddleware{
		EnableRecovery:  true,
		EnableRequestID: true,
		EnableLogger:    true,
		EnableCORS:      true,
		EnableSecurity:  true,
		EnableXSS:       false,
		EnableCSRF:      cfg.Middleware.CSRF.Enable,
		EnableAuth:      true,
		EnableMetrics:   cfg.Middleware.Metrics.Enable,
		EnableRateLimit: cfg.Middleware.RateLimit.Enable,
		AuthSkipPaths:   cfg.Middleware.Auth.SkipPaths,
		AuthJWTMgr:      jwtMgr,
		AuthCache:       redisCache,
		RateLimitCache:  redisCache,
	})

	return engine
}

// ProvideRouteModules keeps route registration order explicit.
func ProvideRouteModules(
	docs *docsmodule.Module,
	health *healthmodule.Module,
	metrics *metricsmodule.Module,
	auth *authmodule.Module,
	security *securitymodule.Module,
	permission *permissionmodule.Module,
) []apppkg.RouteModule {
	return []apppkg.RouteModule{
		docs,
		health,
		metrics,
		auth,
		security,
		permission,
	}
}

var infraSet = wire.NewSet(
	ProvideConfig,
	ProvideDB,
	ProvideRedisCache,
	ProvideJWTManager,
	ProvideGinEngine,
)

var moduleSet = wire.NewSet(
	docsmodule.ProviderSet,
	healthmodule.ProviderSet,
	metricsmodule.ProviderSet,
	permissionmodule.ProviderSet,
	securitymodule.ProviderSet,
	authmodule.ProviderSet,
	ProvideRouteModules,
)

// InitializeApp 初始化应用程序（Wire 生成）
//
//go:generate wire
func InitializeApp(configFile string) (*apppkg.App, error) {
	wire.Build(
		infraSet,
		moduleSet,
		apppkg.New,
	)
	return nil, nil
}
