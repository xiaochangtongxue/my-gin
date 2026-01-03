package config

import (
	"time"
)

// Config 全局配置结构
type Config struct {
	Server     ServerConfig      `mapstructure:"server" validate:"required"`
	Database   DatabaseConfig    `mapstructure:"database" validate:"required"`
	Redis      RedisConfig       `mapstructure:"redis" validate:"required"`
	Logger     LoggerConfig      `mapstructure:"logger" validate:"required"`
	JWT        JWTConfig         `mapstructure:"jwt" validate:"required"`
	Password   PasswordConfig    `mapstructure:"password" validate:"required"`
	Captcha    CaptchaConfig     `mapstructure:"captcha" validate:"required"`
	Bruteforce BruteforceConfig  `mapstructure:"bruteforce" validate:"required"`
	Middleware MiddlewareConfig  `mapstructure:"middleware" validate:"required"`
	Permission PermissionConfig  `mapstructure:"permission" validate:"required"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Mode      string `mapstructure:"mode" validate:"oneof=debug release test"`
	Host      string `mapstructure:"host" validate:"required"`
	Port      int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host" validate:"required"`
	Port            int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Username        string `mapstructure:"username" validate:"required"`
	Password        string `mapstructure:"password" validate:"required"`
	Database        string `mapstructure:"database" validate:"required"`
	Charset         string `mapstructure:"charset"`
	ParseTime       bool   `mapstructure:"parse_time"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" validate:"min=0"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" validate:"min=0"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	SlowThreshold   time.Duration `mapstructure:"slow_threshold"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host        string `mapstructure:"host" validate:"required"`
	Port        int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Password    string `mapstructure:"password"`
	DB          int    `mapstructure:"db" validate:"min=0,max=15"`
	PoolSize    int    `mapstructure:"pool_size" validate:"min=0"`
	MinIdleConn int    `mapstructure:"min_idle_conn" validate:"min=0"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level" validate:"oneof=debug info warn error fatal"`
	FileName   string `mapstructure:"file_name"`
	MaxSize    int    `mapstructure:"max_size" validate:"min=1"`
	MaxBackups int    `mapstructure:"max_backups" validate:"min=0"`
	MaxAge     int    `mapstructure:"max_age" validate:"min=0"`
	Compress   bool   `mapstructure:"compress"`
	Console    bool   `mapstructure:"console"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret        string        `mapstructure:"secret" validate:"required,min=16"`
	ExpireTime   time.Duration `mapstructure:"expire_time" validate:"required"` // Access Token 过期时间
	RefreshExpireTime time.Duration `mapstructure:"refresh_expire_time" validate:"required"` // Refresh Token 过期时间
	Issuer       string        `mapstructure:"issuer" validate:"required"`
}

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	Recovery  RecoveryConfig   `mapstructure:"recovery"`
	RequestID RequestIDConfig  `mapstructure:"request_id"`
	Logger    LoggerMiddlewareConfig `mapstructure:"logger"`
	CORS      CORSConfig       `mapstructure:"cors"`
	CSRF      CSRFConfig       `mapstructure:"csrf"`
	Security  SecurityConfig   `mapstructure:"security"`
	Auth      AuthConfig       `mapstructure:"auth"`
	Metrics   MetricsConfig    `mapstructure:"metrics"`
	RateLimit RateLimitConfig  `mapstructure:"ratelimit"`
}

// RecoveryConfig Recovery中间件配置
type RecoveryConfig struct {
	// EnableStacktrace 是否启用堆栈信息记录（生产环境建议关闭）
	EnableStacktrace bool `mapstructure:"enable_stacktrace"`
}

// RequestIDConfig RequestID中间件配置
type RequestIDConfig struct {
	// Header 请求ID的HTTP头字段名
	Header string `mapstructure:"header"`
	// SkipPaths 跳过生成RequestID的路径（如健康检查）
	SkipPaths []string `mapstructure:"skip_paths"`
}

// LoggerMiddlewareConfig Logger中间件配置
type LoggerMiddlewareConfig struct {
	// SlowThreshold 慢请求阈值
	SlowThreshold time.Duration `mapstructure:"slow_threshold"`
	// SkipPaths 跳过日志记录的路径
	SkipPaths []string `mapstructure:"skip_paths"`
	// LogQuery 是否记录查询参数
	LogQuery bool `mapstructure:"log_query"`
}

// CORSConfig CORS中间件配置
type CORSConfig struct {
	// AllowOrigins 允许的源
	AllowOrigins []string `mapstructure:"allow_origins"`
	// AllowMethods 允许的HTTP方法
	AllowMethods []string `mapstructure:"allow_methods"`
	// AllowHeaders 允许的请求头
	AllowHeaders []string `mapstructure:"allow_headers"`
	// ExposeHeaders 暴露的响应头
	ExposeHeaders []string `mapstructure:"expose_headers"`
	// AllowCredentials 是否允许携带凭证
	AllowCredentials bool `mapstructure:"allow_credentials"`
	// MaxAge 预检请求缓存时间（秒）
	MaxAge int `mapstructure:"max_age"`
}

// CSRFConfig CSRF中间件配置
type CSRFConfig struct {
	// Enable 是否启用CSRF保护
	Enable bool `mapstructure:"enable"`
	// TokenLength Token长度
	TokenLength int `mapstructure:"token_length"`
	// TokenExpiry Token过期时间
	TokenExpiry time.Duration `mapstructure:"token_expiry"`
	// CookieName Cookie名称
	CookieName string `mapstructure:"cookie_name"`
	// HeaderName Header字段名
	HeaderName string `mapstructure:"header_name"`
	// TrustedOrigins 信任的源（跳过CSRF验证）
	TrustedOrigins []string `mapstructure:"trusted_origins"`
}

// SecurityConfig 安全响应头配置
type SecurityConfig struct {
	// Enable 是否启用安全响应头
	Enable bool `mapstructure:"enable"`
	// XFrameOptions X-Frame-Options 响应头
	XFrameOptions string `mapstructure:"x_frame_options"`
	// XContent_type_Options X-Content-Type-Options 响应头
	XContentTypeOptions string `mapstructure:"x_content_type_options"`
	// XSSProtection X-XSS-Protection 响应头
	XSSProtection string `mapstructure:"xss_protection"`
	// HSTS HTTP Strict Transport Security
	HSTS struct {
		Enable        bool          `mapstructure:"enable"`
		MaxAge        time.Duration `mapstructure:"max_age"`
		IncludeSubDomains bool       `mapstructure:"include_subdomains"`
	} `mapstructure:"hsts"`
	// ContentSecurityPolicy CSP 响应头
	ContentSecurityPolicy string `mapstructure:"content_security_policy"`
	// ReferrerPolicy Referrer-Policy 响应头
	ReferrerPolicy string `mapstructure:"referrer_policy"`
}

// AuthConfig Auth认证中间件配置
type AuthConfig struct {
	// SkipPaths 跳过认证的路径（白名单）
	SkipPaths []string `mapstructure:"skip_paths"`
}

// MetricsConfig Metrics中间件配置
type MetricsConfig struct {
	// Enable 是否启用 Metrics
	Enable bool `mapstructure:"enable"`
	// SkipPaths 跳过统计的路径
	SkipPaths []string `mapstructure:"skip_paths"`
}

// RateLimitConfig 限流中间件配置
type RateLimitConfig struct {
	// Enable 是否启用限流
	Enable bool `mapstructure:"enable"`
	// Global 全局限流配置
	Global RateLimitItem `mapstructure:"global"`
	// IP IP限流配置
	IP RateLimitItem `mapstructure:"ip"`
	// User 用户限流配置
	User RateLimitItem `mapstructure:"user"`
}

// RateLimitItem 限流项配置
type RateLimitItem struct {
	Rate  float64 `mapstructure:"rate"`  // 每秒请求数
	Burst int      `mapstructure:"burst"` // 突发流量
}

// PasswordConfig 密码配置
type PasswordConfig struct {
	MinLength     int `mapstructure:"min_length" validate:"min=6"`     // 最小长度
	RequireLetter bool `mapstructure:"require_letter"`                // 必须包含字母
	RequireDigit  bool `mapstructure:"require_digit"`                 // 必须包含数字
	RequireSpecial bool `mapstructure:"require_special"`              // 必须包含特殊字符
	BcryptCost    int `mapstructure:"bcrypt_cost" validate:"min=4,max=12"` // bcrypt强度
}

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Length int           `mapstructure:"length" validate:"min=4,max=8"` // 验证码长度
	Width  int           `mapstructure:"width" validate:"min=50"`       // 图片宽度
	Height int           `mapstructure:"height" validate:"min=20"`      // 图片高度
	Expire time.Duration `mapstructure:"expire" validate:"required"`   // 过期时间
}

// BruteforceConfig 防暴力破解配置
type BruteforceConfig struct {
	Enable            bool          `mapstructure:"enable"`             // 是否启用
	MaxAttempts       int           `mapstructure:"max_attempts"`       // 最大尝试次数
	LockDuration      time.Duration `mapstructure:"lock_duration"`      // 锁定时长
	Window            time.Duration `mapstructure:"window"`             // 时间窗口
	BlacklistDuration time.Duration `mapstructure:"blacklist_duration"` // 黑名单时长
}

// PermissionConfig 权限控制配置
type PermissionConfig struct {
	// Model 权限模型：rbac, abac, acl
	Model string `mapstructure:"model" validate:"oneof=rbac abac acl"`
	// ModelFile Casbin 模型文件路径
	ModelFile string `mapstructure:"model_file"`
}