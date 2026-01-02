package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/jwt"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"github.com/xiaochangtongxue/my-gin/pkg/response"
	"go.uber.org/zap"
)

const (
	// AuthorizationHeader Authorization 请求头字段名
	AuthorizationHeader = "Authorization"
	// BearerPrefix Bearer Token 前缀
	BearerPrefix = "Bearer "
	// UserIDKey 用户 ID 在 gin.Context 中的 key
	UserIDKey = "user_id"
	// UsernameKey 用户名在 gin.Context 中的 key
	UsernameKey = "username"
	// TokenBlacklistKeyPrefix Token黑名单键前缀
	TokenBlacklistKeyPrefix = "jwt:blacklist:"
)

// AuthConfig 认证中间件配置
type AuthConfig struct {
	SkipPaths []string // 跳过认证的路径（白名单）
	JWTMgr    *jwt.Manager
	Cache     cache.Cache // 黑名单缓存
}

var (
	// defaultAuthConfig 默认认证配置（兼容旧代码）
	defaultAuthConfig *AuthConfig
)

// InitAuth 初始化认证中间件（兼容全局配置）
func InitAuth(skipPaths []string, jwtMgr *jwt.Manager, cache cache.Cache) {
	defaultAuthConfig = &AuthConfig{
		SkipPaths: skipPaths,
		JWTMgr:    jwtMgr,
		Cache:     cache,
	}
}

// Auth JWT 认证中间件
// 使用默认配置
func Auth() gin.HandlerFunc {
	if defaultAuthConfig == nil {
		// 兼容：未初始化时使用全局 jwt 包
		return authWithOptions([]string{}, nil, nil)
	}
	return authWithOptions(defaultAuthConfig.SkipPaths, defaultAuthConfig.JWTMgr, defaultAuthConfig.Cache)
}

// AuthWithOptions 带配置的 JWT 认证中间件
func authWithOptions(skipPaths []string, jwtMgr *jwt.Manager, blacklistCache cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否在白名单中（支持前缀匹配）
		for _, path := range skipPaths {
			if c.Request.URL.Path == path || strings.HasPrefix(c.Request.URL.Path, path+"/") {
				c.Next()
				return
			}
		}

		// 获取请求头中的 Token
		token := extractToken(c)

		if token == "" {
			response.Unauthorized(c, "")
			c.Abort()
			return
		}

		// 检查 Token 是否在黑名单中
		if isTokenBlacklisted(c, token, blacklistCache) {
			logger.Warn("Token 在黑名单中",
				zap.String("request_id", GetRequestID(c)),
				zap.String("client_ip", c.ClientIP()),
			)
			response.Fail(c, response.CodeTokenInvalid, "Token 已失效")
			c.Abort()
			return
		}

		// 解析 Token
		var claims *jwt.Claims
		var err error
		if jwtMgr != nil {
			claims, err = jwtMgr.ParseToken(token)
		} else {
			// 兼容：使用全局 jwt 包
			claims, err = jwt.ParseToken(token)
		}

		if err != nil {
			// 记录日志
			logger.Warn("Token 解析失败",
				zap.String("request_id", GetRequestID(c)),
				zap.String("token", maskToken(token)),
				zap.Error(err),
			)

			// 根据错误类型返回不同的响应
			if errors.Is(err, jwt.ErrTokenExpired) {
				response.Fail(c, response.CodeTokenExpired, "")
				c.Abort()
				return
			}
			response.Fail(c, response.CodeTokenInvalid, "")
			c.Abort()
			return
		}

		// 将用户信息存入 Context，供后续处理器使用
		c.Set(UserIDKey, claims.UserID)
		c.Set(UsernameKey, claims.Username)

		logger.Debug("Token 验证成功",
			zap.String("request_id", GetRequestID(c)),
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
		)

		c.Next()
	}
}

// AuthOptional 可选的 JWT 认证中间件
// 如果有 Token 则验证，没有则跳过
func AuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		var blacklistCache cache.Cache
		if defaultAuthConfig != nil {
			blacklistCache = defaultAuthConfig.Cache
		}

		// 检查 Token 是否在黑名单中
		if isTokenBlacklisted(c, token, blacklistCache) {
			c.Next()
			return
		}

		// 尝试解析 Token
		var claims *jwt.Claims
		var err error
		if defaultAuthConfig != nil && defaultAuthConfig.JWTMgr != nil {
			claims, err = defaultAuthConfig.JWTMgr.ParseToken(token)
		} else {
			claims, err = jwt.ParseToken(token)
		}

		if err != nil {
			// Token 无效，但不阻止请求
			logger.Debug("可选认证 Token 解析失败，跳过认证",
				zap.String("request_id", GetRequestID(c)),
				zap.Error(err),
			)
			c.Next()
			return
		}

		// Token 有效，将用户信息存入 Context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UsernameKey, claims.Username)

		c.Next()
	}
}

// AddTokenToBlacklist 添加 Token 到黑名单
// 使用 Request Context，请求取消时自动取消 Redis 操作
func AddTokenToBlacklist(c *gin.Context, token string, claims *jwt.Claims) error {
	var blacklistCache cache.Cache
	if defaultAuthConfig != nil {
		blacklistCache = defaultAuthConfig.Cache
	}

	if blacklistCache == nil {
		return nil // 黑名单未初始化，跳过
	}

	// 计算 Token 的剩余有效时间
	ttl := getTokenTTL(claims)
	if ttl <= 0 {
		return nil // Token 已过期，无需添加到黑名单
	}

	// 使用 Token 的哈希值作为 key
	tokenKey := getTokenKey(token)

	// 使用 Request Context
	return blacklistCache.Set(c.Request.Context(), tokenKey, "1", time.Duration(ttl)*time.Second)
}

// isTokenBlacklisted 检查 Token 是否在黑名单中
// 使用 Request Context
func isTokenBlacklisted(c *gin.Context, token string, blacklistCache cache.Cache) bool {
	if blacklistCache == nil {
		return false // 黑名单未初始化，跳过检查
	}

	tokenKey := getTokenKey(token)
	exists, err := blacklistCache.Exists(c.Request.Context(), tokenKey)
	if err != nil {
		logger.Warn("检查 Token 黑名单失败",
			zap.String("request_id", GetRequestID(c)),
			zap.Error(err),
		)
		return false
	}
	return exists
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

// extractToken 从请求中提取 Token
func extractToken(c *gin.Context) string {
	// 从 Authorization 请求头获取
	authHeader := c.GetHeader(AuthorizationHeader)

	// 检查是否为 Bearer Token 格式
	if strings.HasPrefix(authHeader, BearerPrefix) {
		return strings.TrimPrefix(authHeader, BearerPrefix)
	}

	// 从查询参数获取（备用方案）
	if token := c.Query("token"); token != "" {
		return token
	}

	return ""
}

// maskToken 遮盖 Token 用于日志记录（只显示前6位和后4位）
func maskToken(token string) string {
	if len(token) < 10 {
		return "***"
	}
	return token[:6] + "***" + token[len(token)-4:]
}

// GetUserID 从 gin.Context 获取用户 ID
// 调用前需要确保使用了 Auth 中间件
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get(UserIDKey); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}

// GetUsername 从 gin.Context 获取用户名
// 调用前需要确保使用了 Auth 中间件
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get(UsernameKey); exists {
		if name, ok := username.(string); ok {
			return name
		}
	}
	return ""
}

// IsAuthenticated 检查当前请求是否已认证
func IsAuthenticated(c *gin.Context) bool {
	return GetUserID(c) > 0
}

// ExtractAuthToken 从请求中提取 Token（供外部使用）
func ExtractAuthToken(c *gin.Context) string {
	return extractToken(c)
}

// GetTokenClaims 从请求中解析 Token 获取 Claims（供外部使用）
func GetTokenClaims(c *gin.Context) (*jwt.Claims, error) {
	token := extractToken(c)
	if token == "" {
		return nil, errors.New("token 为空")
	}

	if defaultAuthConfig != nil && defaultAuthConfig.JWTMgr != nil {
		return defaultAuthConfig.JWTMgr.ParseToken(token)
	}
	return jwt.ParseToken(token)
}
