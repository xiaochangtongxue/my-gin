//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/xiaochangtongxue/my-gin/internal/handler"
	"github.com/xiaochangtongxue/my-gin/internal/middleware"
	"github.com/xiaochangtongxue/my-gin/internal/repository"
	"github.com/xiaochangtongxue/my-gin/internal/service"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/database"
	"github.com/xiaochangtongxue/my-gin/pkg/jwt"
)

// ==================== 基础设施 Provider ====================

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
	// 同时初始化全局实例（兼容便捷方法）
	jwt.Init(jwt.Config{
		Secret:     cfg.JWT.Secret,
		ExpireTime: cfg.JWT.ExpireTime,
		Issuer:     cfg.JWT.Issuer,
	})
	return mgr
}

// ==================== Repository Provider ====================

// ProvideRefreshTokenRepository 提供 Refresh Token 仓储
func ProvideRefreshTokenRepository(db *gorm.DB) repository.RefreshTokenRepository {
	return repository.NewRefreshTokenRepository(db)
}

// ==================== Service Provider ====================

// ProvideAuthService 提供认证服务
func ProvideAuthService(
	rtRepo repository.RefreshTokenRepository,
	redisCache cache.Cache,
	jwtMgr *jwt.Manager,
	cfg *config.Config,
) *service.AuthService {
	return service.NewAuthService(rtRepo, redisCache, jwtMgr, service.Config{
		AccessTokenExpire:  cfg.JWT.ExpireTime,
		RefreshTokenExpire: cfg.JWT.RefreshExpireTime,
	})
}

// ==================== Handler Provider ====================

// ProvideAuthHandler 提供认证处理器
func ProvideAuthHandler(authSvc *service.AuthService) *handler.AuthHandler {
	return handler.NewAuthHandler(authSvc)
}

// ProvideHealthHandler 提供健康检查处理器
func ProvideHealthHandler(db *gorm.DB, cache cache.Cache) *handler.HealthHandler {
	return handler.NewHealthHandler(db, cache)
}

// ==================== Gin Engine Provider ====================

// ProvideGinEngine 提供 Gin 引擎
func ProvideGinEngine(cfg *config.Config, jwtMgr *jwt.Manager, redisCache cache.Cache) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()

	// 应用中间件（依赖注入配置）
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
		AuthSkipPaths:   cfg.Middleware.Auth.SkipPaths,
		AuthJWTMgr:      jwtMgr,
		AuthCache:       redisCache,
	})

	return engine
}

// ==================== Wire ProviderSets ====================

// infraSet 基础设施 ProviderSet
var infraSet = wire.NewSet(
	ProvideConfig,
	ProvideDB,
	ProvideRedisCache,
	ProvideJWTManager,
)

// repoSet Repository ProviderSet
var repoSet = wire.NewSet(
	ProvideRefreshTokenRepository,
)

// serviceSet Service ProviderSet
var serviceSet = wire.NewSet(
	ProvideAuthService,
)

// handlerSet Handler ProviderSet
var handlerSet = wire.NewSet(
	ProvideAuthHandler,
	ProvideHealthHandler,
)

// engineSet Gin Engine ProviderSet
var engineSet = wire.NewSet(
	ProvideGinEngine,
)

// ==================== Application Container ====================

// App 应用程序依赖容器
type App struct {
	Config       *config.Config
	DB           *gorm.DB
	Cache        cache.Cache
	JWTMgr       *jwt.Manager
	Engine       *gin.Engine
	AuthHandler  *handler.AuthHandler
	HealthHandler *handler.HealthHandler
}

// InitializeApp 初始化应用程序（Wire 生成）
//go:generate wire
func InitializeApp(configFile string) (*App, error) {
	wire.Build(
		infraSet,
		repoSet,
		serviceSet,
		handlerSet,
		engineSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}