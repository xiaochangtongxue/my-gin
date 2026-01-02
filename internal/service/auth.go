package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/xiaochangtongxue/my-gin/internal/model"
	"github.com/xiaochangtongxue/my-gin/internal/repository"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/jwt"
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

// Config 认证服务配置
type Config struct {
	AccessTokenExpire  time.Duration // Access Token 过期时间
	RefreshTokenExpire time.Duration // Refresh Token 过期时间
}

// TokenPair 双 Token 结构
// @description 双 Token 响应（Access Token + Refresh Token）
type TokenPair struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`  // Access Token
	RefreshToken string `json:"refresh_token" example:"a1b2c3d4e5f6..."`                         // Refresh Token
	ExpiresIn    int64  `json:"expires_in" example:"1800"`                                       // Access Token 过期时间（秒）
}

// AuthService 认证服务
type AuthService struct {
	rtRepo    repository.RefreshTokenRepository
	cache     cache.Cache
	jwtMgr    *jwt.Manager
	refreshTTL time.Duration
}

// NewAuthService 创建认证服务
func NewAuthService(
	rtRepo repository.RefreshTokenRepository,
	cache cache.Cache,
	jwtMgr *jwt.Manager,
	cfg Config,
) *AuthService {
	return &AuthService{
		rtRepo:    rtRepo,
		cache:     cache,
		jwtMgr:    jwtMgr,
		refreshTTL: cfg.RefreshTokenExpire,
	}
}

// Login 用户登录（生成双 Token）
// userID: 用户 ID
// username: 用户名
// 返回 TokenPair（AT + RT）
func (s *AuthService) Login(ctx context.Context, userID uint, username string) (*TokenPair, error) {
	// 1. 生成 Access Token
	accessToken, err := s.jwtMgr.GenerateToken(userID, username)
	if err != nil {
		return nil, err
	}

	// 2. 生成 Refresh Token（随机字符串）
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 3. 存储 Refresh Token 到数据库
	expiresAt := time.Now().Add(s.refreshTTL)

	rt := &model.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}

	if err := s.rtRepo.Create(ctx, rt); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtMgr.GetExpireTime().Seconds()),
	}, nil
}

// RefreshToken 刷新 Access Token
// refreshToken: 客户端携带的 Refresh Token
// 返回新的 TokenPair
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
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
	newRefreshToken, err := generateRefreshToken()
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

// Logout 登出（删除 Refresh Token）
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.rtRepo.DeleteByToken(ctx, refreshToken)
}

// LogoutByUser 登出用户的所有设备
func (s *AuthService) LogoutByUser(ctx context.Context, userID uint) error {
	return s.rtRepo.DeleteByUser(ctx, userID)
}

// generateRefreshToken 生成随机 Refresh Token
func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// LogoutWithAccessToken 登出（删除 RT 并将 AT 加入黑名单）
func (s *AuthService) LogoutWithAccessToken(ctx context.Context, accessToken, refreshToken string) error {
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

// LogoutAllWithAccessToken 登出所有设备（删除用户所有 RT 并将当前 AT 加入黑名单）
func (s *AuthService) LogoutAllWithAccessToken(ctx context.Context, accessToken string, userID uint) error {
	// 删除用户的所有 Refresh Token
	_ = s.rtRepo.DeleteByUser(ctx, userID)

	// 将当前 Access Token 加入黑名单
	if accessToken != "" && s.cache != nil {
		claims, err := s.jwtMgr.ParseToken(accessToken)
		if err == nil {
			_ = s.addTokenToBlacklist(ctx, accessToken, claims)
		}
	}
	return nil
}

// addTokenToBlacklist 将 Token 加入黑名单
func (s *AuthService) addTokenToBlacklist(ctx context.Context, token string, claims *jwt.Claims) error {
	if s.cache == nil {
		return nil
	}

	// 计算 Token 的剩余有效时间
	ttl := getTokenTTL(claims)
	if ttl <= 0 {
		return nil // Token 已过期，无需添加到黑名单
	}

	// 使用 Token 的哈希值作为 key
	tokenKey := getTokenKey(token)

	return s.cache.Set(ctx, tokenKey, "1", time.Duration(ttl)*time.Second)
}

// getTokenKey 获取 Token 的唯一标识（SHA256 哈希）
func getTokenKey(token string) string {
	hash := sha256.Sum256([]byte(token))
	return TokenBlacklistKeyPrefix + hex.EncodeToString(hash[:])
}

// getTokenTTL 计算 Token 的剩余有效时间（秒）
func getTokenTTL(claims *jwt.Claims) int64 {
	if claims.ExpiresAt == nil {
		return 0
	}
	expiresAt := claims.ExpiresAt.Time
	remaining := expiresAt.Unix() - time.Now().Unix()
	if remaining <= 0 {
		return 0
	}
	return remaining
}