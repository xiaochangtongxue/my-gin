package limiter

import "context"

// Limiter 限流器接口
type Limiter interface {
	// Allow 检查是否允许一个请求通过
	Allow(ctx context.Context, key string) bool

	// AllowN 检查是否允许 n 个请求通过
	AllowN(ctx context.Context, key string, n int) bool
}

// RateLimit 限流结果
type RateLimit struct {
	Allowed bool    // 是否允许
	Remaining int    // 剩余配额
	ResetAt   int64   // 重置时间（Unix 时间戳，秒）
	RetryAfter int   // 重试等待秒数（不允许时）
}

// Config 限流器配置
type Config struct {
	Enable bool // 是否启用限流
}