package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrTokenInvalid Token 无效
	ErrTokenInvalid = errors.New("token 无效")
	// ErrTokenExpired Token 已过期
	ErrTokenExpired = errors.New("token 已过期")
	// ErrTokenMalformed Token 格式错误
	ErrTokenMalformed = errors.New("token 格式错误")
)

// Claims JWT 声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`  // 用户 ID
	Username string `json:"username"` // 用户名
	jwt.RegisteredClaims
}

// Config JWT 配置
type Config struct {
	Secret     string        // 密钥
	ExpireTime time.Duration // Access Token 过期时间
	Issuer     string        // 签发者
}

// Manager JWT 管理器
type Manager struct {
	secret     []byte
	expireTime time.Duration
	issuer     string
}

// NewManager 创建 JWT 管理器
func NewManager(cfg Config) *Manager {
	return &Manager{
		secret:     []byte(cfg.Secret),
		expireTime: cfg.ExpireTime,
		issuer:     cfg.Issuer,
	}
}

// GenerateToken 生成 Token
func (m *Manager) GenerateToken(userID uint, username string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expireTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken 解析 Token
func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// ValidateToken 验证 Token 是否有效
func (m *Manager) ValidateToken(tokenString string) bool {
	_, err := m.ParseToken(tokenString)
	return err == nil
}

// GetExpireTime 获取 Token 过期时间
func (m *Manager) GetExpireTime() time.Duration {
	return m.expireTime
}

// 全局 JWT 管理器实例（兼容便捷方法）
var defaultManager *Manager

// Init 初始化 JWT 管理器
func Init(cfg Config) {
	defaultManager = NewManager(cfg)
}

// GetManager 获取 JWT 管理器
func GetManager() *Manager {
	return defaultManager
}

// GenerateToken 生成 Token（便捷方法）
func GenerateToken(userID uint, username string) (string, error) {
	if defaultManager == nil {
		return "", ErrTokenInvalid
	}
	return defaultManager.GenerateToken(userID, username)
}

// ParseToken 解析 Token（便捷方法）
func ParseToken(tokenString string) (*Claims, error) {
	if defaultManager == nil {
		return nil, ErrTokenInvalid
	}
	return defaultManager.ParseToken(tokenString)
}

// ValidateToken 验证 Token（便捷方法）
func ValidateToken(tokenString string) bool {
	if defaultManager == nil {
		return false
	}
	return defaultManager.ValidateToken(tokenString)
}

// GetExpireTime 获取过期时间（便捷方法）
func GetExpireTime() time.Duration {
	if defaultManager == nil {
		return 0
	}
	return defaultManager.GetExpireTime()
}