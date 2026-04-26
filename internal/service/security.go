package service

import (
	"context"
	"time"

	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	apperrors "github.com/xiaochangtongxue/my-gin/pkg/errors"
	"github.com/xiaochangtongxue/my-gin/pkg/response"
)

// SecurityService 安全服务接口
type SecurityService interface {
	// IsAccountLocked 检查账号是否被锁定
	IsAccountLocked(ctx context.Context, uid string) (bool, error)

	// NeedCaptcha 检查是否需要验证码
	NeedCaptcha(ctx context.Context, ip, uid string) (bool, error)

	// RecordFailure 记录登录失败
	RecordFailure(ctx context.Context, ip, uid string) error

	// ClearFailures 清除失败记录
	ClearFailures(ctx context.Context, ip, uid string) error

	// GetLockTimeRemaining 获取剩余锁定时间（分钟）
	GetLockTimeRemaining(ctx context.Context, uid string) (int, error)

	// UnlockAccount 解锁账号（管理员操作）
	UnlockAccount(ctx context.Context, uid string) error

	// IsIPBlacklisted 检查IP是否在黑名单中
	IsIPBlacklisted(ctx context.Context, ip string) (bool, error)
}

// Config 安全服务配置
type SecurityConfig struct {
	MaxAttempts       int           // 最大尝试次数
	LockDuration      time.Duration // 锁定时长
	Window            time.Duration // 时间窗口
	BlacklistDuration time.Duration // 黑名单时长
}

// securityService 安全服务实现
type securityService struct {
	cache cache.Cache
	cfg   SecurityConfig
}

// NewSecurityService 创建安全服务
func NewSecurityService(cache cache.Cache, cfg SecurityConfig) SecurityService {
	return &securityService{
		cache: cache,
		cfg:   cfg,
	}
}

// IsAccountLocked 检查账号是否被锁定
func (s *securityService) IsAccountLocked(ctx context.Context, uid string) (bool, error) {
	return s.cache.Exists(ctx, s.lockedKey(uid))
}

// NeedCaptcha 检查是否需要验证码
func (s *securityService) NeedCaptcha(ctx context.Context, ip, uid string) (bool, error) {
	return s.cache.Exists(ctx, s.needCaptchaKey(ip, uid))
}

// RecordFailure 记录登录失败
func (s *securityService) RecordFailure(ctx context.Context, ip, uid string) error {
	// 1. 设置需要验证码标记
	needCaptchaKey := s.needCaptchaKey(ip, uid)
	if err := s.cache.Set(ctx, needCaptchaKey, "1", s.cfg.Window); err != nil {
		return apperrors.Wrap(err, response.CodeRedisError, "设置验证码标记失败")
	}

	// 2. 增加失败次数
	failKey := s.failKey(ip, uid)
	failCount, err := s.cache.Incr(ctx, failKey)
	if err != nil {
		return apperrors.Wrap(err, response.CodeRedisError, "记录失败次数失败")
	}
	if failCount == 1 {
		if err := s.cache.Expire(ctx, failKey, s.cfg.Window); err != nil {
			return apperrors.Wrap(err, response.CodeRedisError, "设置失败次数过期时间失败")
		}
	}

	// 3. 检查是否需要锁定
	if failCount >= int64(s.cfg.MaxAttempts) {
		// 锁定账号
		if err := s.cache.Set(ctx, s.lockedKey(uid), "1", s.cfg.LockDuration); err != nil {
			return apperrors.Wrap(err, response.CodeRedisError, "锁定账号失败")
		}

		// 检查是否需要加入IP黑名单
		blacklistCountKey := s.blacklistCountKey(ip)
		blacklistCount, err := s.cache.Incr(ctx, blacklistCountKey)
		if err != nil {
			return apperrors.Wrap(err, response.CodeRedisError, "更新黑名单计数失败")
		}
		if blacklistCount == 1 {
			if err := s.cache.Expire(ctx, blacklistCountKey, 24*time.Hour); err != nil {
				return apperrors.Wrap(err, response.CodeRedisError, "设置黑名单计数过期时间失败")
			}
		}

		// 累计失败10次加入黑名单
		if blacklistCount >= 10 {
			if err := s.cache.Set(ctx, s.blacklistKey(ip), "1", s.cfg.BlacklistDuration); err != nil {
				return apperrors.Wrap(err, response.CodeRedisError, "加入黑名单失败")
			}
		}
	}

	return nil
}

// ClearFailures 清除失败记录
func (s *securityService) ClearFailures(ctx context.Context, ip, uid string) error {
	// 清除需要验证码标记
	_ = s.cache.Del(ctx, s.needCaptchaKey(ip, uid))

	// 清除失败次数
	_ = s.cache.Del(ctx, s.failKey(ip, uid))

	return nil
}

// GetLockTimeRemaining 获取剩余锁定时间（分钟）
func (s *securityService) GetLockTimeRemaining(ctx context.Context, uid string) (int, error) {
	key := s.lockedKey(uid)
	// 获取 TTL
	ttl, err := s.cache.TTL(ctx, key)
	if err != nil || ttl <= 0 {
		return 0, nil
	}
	// 转换为分钟
	return int(ttl.Minutes()), nil
}

// UnlockAccount 解锁账号（管理员操作）
func (s *securityService) UnlockAccount(ctx context.Context, uid string) error {
	// 删除锁定标记
	if err := s.cache.Del(ctx, s.lockedKey(uid)); err != nil {
		return apperrors.Wrap(err, response.CodeRedisError, "解锁账号失败")
	}
	return nil
}

// IsIPBlacklisted 检查IP是否在黑名单中
func (s *securityService) IsIPBlacklisted(ctx context.Context, ip string) (bool, error) {
	return s.cache.Exists(ctx, s.blacklistKey(ip))
}

// Redis key 生成方法
func (s *securityService) needCaptchaKey(ip, uid string) string {
	return "login:need_captcha:" + ip + ":" + uid
}

func (s *securityService) failKey(ip, uid string) string {
	return "login:fail:" + ip + ":" + uid
}

func (s *securityService) lockedKey(uid string) string {
	return "login:locked:" + uid
}

func (s *securityService) blacklistKey(ip string) string {
	return "login:blacklist:" + ip
}

func (s *securityService) blacklistCountKey(ip string) string {
	return "login:blacklist_count:" + ip
}
