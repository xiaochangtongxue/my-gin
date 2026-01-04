package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/xiaochangtongxue/my-gin/internal/model"
	"github.com/xiaochangtongxue/my-gin/internal/repository"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/captcha"
	"github.com/xiaochangtongxue/my-gin/pkg/crypto"
	"github.com/xiaochangtongxue/my-gin/pkg/jwt"
	"github.com/xiaochangtongxue/my-gin/pkg/notify"
	"github.com/xiaochangtongxue/my-gin/pkg/utils"
	"gorm.io/gorm"
)

const (
	// TokenBlacklistKeyPrefix Token黑名单键前缀
	TokenBlacklistKeyPrefix = "jwt:blacklist:"
)

var (
	// ErrRefreshTokenInvalid Refresh Token 无效
	ErrRefreshTokenInvalid = errors.New("refresh token 无效")
	// ErrRefreshTokenExpired Refresh Token 已过期
	ErrRefreshTokenExpired = errors.New("refresh token 已过期")
)

type AuthService interface {
	// Register 用户注册
	Register(ctx context.Context, mobile, username, password string) (*RegisterResult, error)

	// Login 用户登录
	Login(ctx context.Context, mobile, password, captchaID, captchaCode, ip string) (*LoginResult, error)

	//返回新的 TokenPair
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)

	// LogoutWithAccessToken 登出（删除 RT 并将 AT 加入黑名单）
	LogoutWithAccessToken(ctx context.Context, accessToken, refreshToken string) error
}

// Config 认证服务配置
type Config struct {
	AccessTokenExpire  time.Duration // Access Token 过期时间
	RefreshTokenExpire time.Duration // Refresh Token 过期时间
}

// TokenPair 双 Token 结构
// @description 双 Token 响应（Access Token + Refresh Token）
type TokenPair struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Access Token
	RefreshToken string `json:"refresh_token" example:"a1b2c3d4e5f6..."`                        // Refresh Token
	ExpiresIn    int64  `json:"expires_in" example:"1800"`                                      // Access Token 过期时间（秒）
}

// RegisterResult 注册结果
type RegisterResult struct {
	UID      uint64
	Username string
}

// LoginResult 登录结果
type LoginResult struct {
	UID          uint64
	Username     string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// AuthService 认证服务
type authService struct {
	rtRepo      repository.RefreshTokenRepository
	cache       cache.Cache
	jwtMgr      *jwt.Manager
	refreshTTL  time.Duration
	userRepo    repository.UserRepository
	securitySvc SecurityService
	captcha     *captcha.Captcha
	notify      notify.Notifier
	bcryptCost  int
}

// NewAuthService 创建认证服务
func NewAuthService(
	rtRepo repository.RefreshTokenRepository,
	cache cache.Cache,
	jwtMgr *jwt.Manager,
	cfg Config,
	userRepo repository.UserRepository,
	securitySvc SecurityService,
	captcha *captcha.Captcha,
	notify notify.Notifier,
	bcryptCost int,
) *authService {
	if bcryptCost == 0 {
		bcryptCost = crypto.DefaultCost
	}
	return &authService{
		rtRepo:      rtRepo,
		cache:       cache,
		jwtMgr:      jwtMgr,
		refreshTTL:  cfg.RefreshTokenExpire,
		userRepo:    userRepo,
		securitySvc: securitySvc,
		captcha:     captcha,
		notify:      notify,
		bcryptCost:  bcryptCost,
	}
}

func (s *authService) Register(ctx context.Context, mobile, username, password string) (*RegisterResult, error) {
	// 1. 检查手机号是否已注册
	exist, err := s.userRepo.ExistsByMobile(ctx, mobile)
	if err != nil {
		return nil, fmt.Errorf("检查手机号失败: %w", err)
	}
	if exist {
		return nil, errors.New("该手机号已注册")
	}

	// 2. 检查用户名是否已存在
	exist, err = s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if exist {
		return nil, errors.New("用户名已存在")
	}

	// 3. 密码加密
	hash, err := crypto.HashPassword(password, s.bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 4. 创建用户
	user := &model.User{
		Mobile:   mobile,
		Username: username,
		Password: string(hash),
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 5. 生成UID（Feistel混淆自增ID）
	user.UID = utils.EncodeUID(user.ID)
	if err := s.userRepo.UpdateUID(ctx, user.ID, user.UID); err != nil {
		return nil, fmt.Errorf("更新用户UID失败: %w", err)
	}

	// 6. 发送注册成功通知（异步，不阻塞）
	go func() {
		_ = s.notify.SendRegisterNotice(context.Background(), mobile, username)
	}()

	return &RegisterResult{
		UID:      user.UID,
		Username: user.Username,
	}, nil
}

// Login 用户登录
func (s *authService) Login(ctx context.Context, mobile, password, captchaID, captchaCode, ip string) (*LoginResult, error) {
	// 1. 根据手机号查用户
	user, err := s.userRepo.FindByMobile(ctx, mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	uid := strconv.FormatUint(user.UID, 10)

	// 2. 检查IP是否在黑名单中
	blacklisted, err := s.securitySvc.IsIPBlacklisted(ctx, ip)
	if err != nil {
		return nil, fmt.Errorf("检查黑名单失败: %w", err)
	}
	if blacklisted {
		return nil, errors.New("IP已被封禁，请联系管理员")
	}

	// 3. 检查账号是否被锁定
	locked, err := s.securitySvc.IsAccountLocked(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("检查锁定状态失败: %w", err)
	}
	if locked {
		ttl, _ := s.securitySvc.GetLockTimeRemaining(ctx, uid)
		return nil, fmt.Errorf("账号已锁定，请%d分钟后再试", ttl)
	}

	// 4. 检查是否需要验证码
	needCaptcha, err := s.securitySvc.NeedCaptcha(ctx, ip, uid)
	if err != nil {
		return nil, fmt.Errorf("检查验证码状态失败: %w", err)
	}
	if needCaptcha {
		if captchaID == "" || captchaCode == "" {
			return nil, errors.New("请输入验证码")
		}
		if !s.captcha.Verify(ctx, captchaID, captchaCode) {
			return nil, errors.New("验证码错误或已过期")
		}
	}

	// 5. 验证密码
	if !crypto.VerifyPassword(user.Password, password) {
		// 密码错误：记录失败
		_ = s.securitySvc.RecordFailure(ctx, ip, uid)
		return nil, errors.New("用户名或密码错误")
	}

	// 6. 密码正确：清除失败记录
	_ = s.securitySvc.ClearFailures(ctx, ip, uid)

	// 7. 生成Token
	accessToken, err := s.jwtMgr.GenerateToken(user.UID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %w", err)
	}

	// 8. 生成RefreshToken
	refreshToken, err := crypto.RandomString(64)
	if err != nil {
		return nil, fmt.Errorf("生成刷新Token失败: %w", err)
	}

	// 9. 把生成的RefreshToken存储到数据库
	newRT := &model.RefreshToken{
		UserID:    user.UID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	s.rtRepo.Create(ctx, newRT)

	return &LoginResult{
		UID:          user.UID,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtMgr.GetExpireTime().Seconds()),
	}, nil
}

// RefreshToken 刷新 Access Token
// refreshToken: 客户端携带的 Refresh Token
// 返回新的 TokenPair
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// 1. 验证 Refresh Token 是否有效（未过期且未删除）
	valid, err := s.rtRepo.IsValidToken(ctx, refreshToken)
	if err != nil || !valid {
		return nil, ErrRefreshTokenInvalid
	}

	// 2. 获取 Refresh Token 信息（需要 UserID）
	rt, err := s.rtRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	// 3. 生成新的 Access Token
	accessToken, err := s.jwtMgr.GenerateToken(rt.UserID, "")
	if err != nil {
		return nil, err
	}

	// 4. 生成新的 Refresh Token（滚动更新）
	newRefreshToken, err := crypto.RandomString(64)
	if err != nil {
		return nil, err
	}

	newRT := &model.RefreshToken{
		UserID:    rt.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}

	// 5. 原子操作：删除旧的 RT，创建新的 RT
	_, err = s.rtRepo.DeleteOldAndCreate(ctx, refreshToken, newRT)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.jwtMgr.GetExpireTime().Seconds()),
	}, nil
}

// LogoutWithAccessToken 登出（删除 RT 并将 AT 加入黑名单）
func (s *authService) LogoutWithAccessToken(ctx context.Context, accessToken, refreshToken string) error {
	// 删除 Refresh Token
	_ = s.rtRepo.DeleteByToken(ctx, refreshToken)

	// 将 Access Token 加入黑名单
	if accessToken != "" && s.cache != nil {
		claims, err := s.jwtMgr.ParseToken(accessToken)
		if err == nil {
			_ = s.addTokenToBlacklist(ctx, accessToken, claims)
		}
	}
	return nil
}

// addTokenToBlacklist 将 Token 加入黑名单
func (s *authService) addTokenToBlacklist(ctx context.Context, token string, claims *jwt.Claims) error {
	if s.cache == nil {
		return nil
	}

	// 计算 Token 的剩余有效时间
	ttl := utils.RemainingSeconds(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil // Token 已过期，无需添加到黑名单
	}

	// 使用 Token 的哈希值作为 key
	tokenKey := TokenBlacklistKeyPrefix + crypto.SHA256(token)

	return s.cache.Set(ctx, tokenKey, "1", time.Duration(ttl)*time.Second)
}
