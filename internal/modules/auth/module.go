package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/xiaochangtongxue/my-gin/internal/handler"
	"github.com/xiaochangtongxue/my-gin/internal/repository"
	"github.com/xiaochangtongxue/my-gin/internal/service"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/captcha"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/jwt"
	"github.com/xiaochangtongxue/my-gin/pkg/notify"
	"gorm.io/gorm"
)

// Module owns authentication routes.
type Module struct {
	authHandler    *handler.AuthHandler
	captchaHandler *handler.CaptchaHandler
}

// ProviderSet contains auth module dependencies.
var ProviderSet = wire.NewSet(
	ProvideRefreshTokenRepository,
	ProvideUserRepository,
	ProvideCaptcha,
	ProvideNotifier,
	ProvideAuthService,
	ProvideAuthHandler,
	ProvideCaptchaHandler,
	NewModule,
)

// NewModule creates an auth module.
func NewModule(authHandler *handler.AuthHandler, captchaHandler *handler.CaptchaHandler) *Module {
	return &Module{authHandler: authHandler, captchaHandler: captchaHandler}
}

// RegisterRoutes registers /api/v1/auth routes.
func (m *Module) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/api/v1/auth")
	group.GET("/captcha", m.captchaHandler.GetCaptcha)
	group.POST("/register", m.authHandler.Register)
	group.POST("/login", m.authHandler.Login)
	group.POST("/refresh", m.authHandler.RefreshToken)
	group.POST("/logout", m.authHandler.Logout)
}

// ProvideRefreshTokenRepository provides refresh token storage.
func ProvideRefreshTokenRepository(db *gorm.DB) repository.RefreshTokenRepository {
	return repository.NewRefreshTokenRepository(db)
}

// ProvideUserRepository provides user storage.
func ProvideUserRepository(db *gorm.DB) repository.UserRepository {
	return repository.NewUserRepository(db)
}

// ProvideNotifier provides a notification sender.
func ProvideNotifier() notify.Notifier {
	return notify.NewNoopNotifier()
}

// ProvideCaptcha provides captcha generation and validation.
func ProvideCaptcha(cfg *config.Config, cache cache.Cache) *captcha.Captcha {
	return captcha.New(captcha.Config{
		Length: cfg.Captcha.Length,
		Width:  cfg.Captcha.Width,
		Height: cfg.Captcha.Height,
		Expire: cfg.Captcha.Expire,
	}, cache)
}

// ProvideAuthService provides the auth service.
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

// ProvideAuthHandler provides the auth HTTP handler.
func ProvideAuthHandler(authSvc service.AuthService) *handler.AuthHandler {
	return handler.NewAuthHandler(authSvc)
}

// ProvideCaptchaHandler provides the captcha HTTP handler.
func ProvideCaptchaHandler(captcha *captcha.Captcha) *handler.CaptchaHandler {
	return handler.NewCaptchaHandler(captcha)
}
