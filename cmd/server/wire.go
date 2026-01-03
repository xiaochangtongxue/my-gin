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
	"github.com/xiaochangtongxue/my-gin/pkg/captcha"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/database"
	"github.com/xiaochangtongxue/my-gin/pkg/jwt"
	"github.com/xiaochangtongxue/my-gin/pkg/notify"
	"github.com/xiaochangtongxue/my-gin/pkg/permission"
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

// ProvideUserRepository 提供用户仓储
func ProvideUserRepository(db *gorm.DB) repository.UserRepository {
	return repository.NewUserRepository(db)
}

// ProvideRoleRepository 提供角色仓储
func ProvideRoleRepository(db *gorm.DB) repository.RoleRepository {
	return repository.NewRoleRepository(db)
}

// ProvideUserRoleRepository 提供用户角色仓储
func ProvideUserRoleRepository(db *gorm.DB) repository.UserRoleRepository {
	return repository.NewUserRoleRepository(db)
}

// ==================== Notification Provider ====================

// ProvideNotifier 提供通知器（使用 NoopNotifier）
func ProvideNotifier() notify.Notifier {
	return notify.NewNoopNotifier()
}

// ProvideCaptcha 提供验证码
func ProvideCaptcha(cfg *config.Config, cache cache.Cache) *captcha.Captcha {
	return captcha.New(captcha.Config{
		Length: cfg.Captcha.Length,
		Width:  cfg.Captcha.Width,
		Height: cfg.Captcha.Height,
		Expire: cfg.Captcha.Expire,
	}, cache)
}

// ProvidePermissionChecker 提供权限检查器
func ProvidePermissionChecker(cfg *config.Config, db *gorm.DB) (permission.PermissionChecker, error) {
	return permission.NewChecker(cfg, db)
}

// ==================== Service Provider ====================

// ProvideAuthService 提供认证服务
func ProvideAuthService(
	rtRepo repository.RefreshTokenRepository,
	redisCache cache.Cache,
	jwtMgr *jwt.Manager,
	cfg *config.Config,
	userRepo repository.UserRepository,
	securitySvc service.SecurityService,
	captcha *captcha.Captcha,
	notify notify.Notifier,
) service.AuthService {
	return service.NewAuthService(rtRepo, redisCache, jwtMgr, service.Config{
		AccessTokenExpire:  cfg.JWT.ExpireTime,
		RefreshTokenExpire: cfg.JWT.RefreshExpireTime,
	}, userRepo, securitySvc, captcha, notify, cfg.Password.BcryptCost)
}

// ProvideSecurityService 提供安全服务
func ProvideSecurityService(cache cache.Cache, cfg *config.Config) service.SecurityService {
	return service.NewSecurityService(cache, service.SecurityConfig{
		MaxAttempts:       5,
		LockDuration:      30 * 60, // 从配置读取
		Window:            15 * 60,
		BlacklistDuration: 60 * 60,
	})
}

// ProvidePermissionService 提供权限服务
func ProvidePermissionService(
	db *gorm.DB,
	roleRepo repository.RoleRepository,
	userRoleRepo repository.UserRoleRepository,
	checker permission.PermissionChecker,
) service.PermissionService {
	// PermissionChecker 同时实现了 PolicyManager 接口
	policyMgr := checker.(permission.PolicyManager)
	return service.NewPermissionService(db, roleRepo, userRoleRepo, policyMgr)
}

// ==================== Handler Provider ====================

// ProvideAuthHandler 提供认证处理器
func ProvideAuthHandler(authSvc service.AuthService) *handler.AuthHandler {
	return handler.NewAuthHandler(authSvc)
}

// ProvideCaptchaHandler 提供验证码处理器
func ProvideCaptchaHandler(captcha *captcha.Captcha) *handler.CaptchaHandler {
	return handler.NewCaptchaHandler(captcha)
}

// ProvideSecurityHandler 提供安全处理器
func ProvideSecurityHandler(securitySvc service.SecurityService) *handler.SecurityHandler {
	return handler.NewSecurityHandler(securitySvc)
}

// ProvideHealthHandler 提供健康检查处理器
func ProvideHealthHandler(db *gorm.DB, cache cache.Cache) *handler.HealthHandler {
	return handler.NewHealthHandler(db, cache)
}

// ProvidePermissionHandler 提供权限处理器
func ProvidePermissionHandler(permissionSvc service.PermissionService) *handler.PermissionHandler {
	return handler.NewPermissionHandler(permissionSvc)
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
		EnableRateLimit: cfg.Middleware.RateLimit.Enable,
		AuthSkipPaths:   cfg.Middleware.Auth.SkipPaths,
		AuthJWTMgr:      jwtMgr,
		AuthCache:       redisCache,
		RateLimitCache:  redisCache,
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
	ProvideNotifier,
	ProvideCaptcha,
	ProvidePermissionChecker,
)

// repoSet Repository ProviderSet
var repoSet = wire.NewSet(
	ProvideRefreshTokenRepository,
	ProvideUserRepository,
	ProvideRoleRepository,
	ProvideUserRoleRepository,
)

// serviceSet Service ProviderSet
var serviceSet = wire.NewSet(
	ProvideAuthService,
	ProvideSecurityService,
	ProvidePermissionService,
)

// handlerSet Handler ProviderSet
var handlerSet = wire.NewSet(
	ProvideAuthHandler,
	ProvideCaptchaHandler,
	ProvideSecurityHandler,
	ProvideHealthHandler,
	ProvidePermissionHandler,
)

// engineSet Gin Engine ProviderSet
var engineSet = wire.NewSet(
	ProvideGinEngine,
)

// ==================== Application Container ====================

// App 应用程序依赖容器
type App struct {
	Config            *config.Config
	DB                *gorm.DB
	Cache             cache.Cache
	JWTMgr            *jwt.Manager
	Engine            *gin.Engine
	AuthHandler       *handler.AuthHandler
	CaptchaHandler    *handler.CaptchaHandler
	SecurityHandler   *handler.SecurityHandler
	HealthHandler     *handler.HealthHandler
	PermissionHandler *handler.PermissionHandler
	PermissionChecker permission.PermissionChecker
	UserRoleRepo      repository.UserRoleRepository
}

// InitializeApp 初始化应用程序（Wire 生成）
//
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
