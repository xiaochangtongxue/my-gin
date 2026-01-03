package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/xiaochangtongxue/my-gin/pkg/cache"
)

// SlidingWindowLimiter 滑动窗口限流器（通过 Cache 接口实现）
type SlidingWindowLimiter struct {
	cache       cache.Cache
	maxRequests int64         // 窗口内允许的最大请求数
	window      time.Duration // 时间窗口大小
}

// NewSlidingWindowLimiter 创建滑动窗口限流器
// maxRequests: 时间窗口内允许的最大请求数
// window: 时间窗口大小
func NewSlidingWindowLimiter(cache cache.Cache, maxRequests int64, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		cache:       cache,
		maxRequests: maxRequests,
		window:      window,
	}
}

// Allow 检查是否允许一个请求通过
func (s *SlidingWindowLimiter) Allow(ctx context.Context, key string) bool {
	return s.AllowN(ctx, key, 1)
}

// AllowN 检查是否允许 n 个请求通过
func (s *SlidingWindowLimiter) AllowN(ctx context.Context, key string, n int) bool {
	now := float64(time.Now().UnixMicro()) / 1e6 // 转换为秒（带小数）
	windowStart := now - s.window.Seconds()

	cacheKey := s.getKey(key)

	// 使用 Lua 脚本保证原子性
	script := `
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local window_start = tonumber(ARGV[2])
		local requested = tonumber(ARGV[3])
		local expire = tonumber(ARGV[4])

		-- 移除窗口外的记录
		redis.call("ZREMRANGEBYSCORE", key, 0, window_start)

		-- 获取当前窗口内的请求数
		local current = redis.call("ZCARD", key)

		-- 检查是否超过限制
		if current + requested <= tonumber(ARGV[5]) then
			-- 添加当前请求
			redis.call("ZADD", key, now, now .. ":" .. math.random(1, 1000000))
			-- 设置过期时间
			redis.call("EXPIRE", key, expire)
			return 1
		else
			return 0
		end
	`

	result, err := s.cache.Eval(ctx, script, []string{cacheKey},
		now, windowStart, n, int(s.window.Seconds())+1, s.maxRequests)
	if err != nil {
		return true // 缓存出错时放行（降级策略）
	}

	// 结果可能是 int64 或其他类型
	if r, ok := result.(int64); ok {
		return r == 1
	}
	return false
}

// getKey 获取缓存 key
func (s *SlidingWindowLimiter) getKey(key string) string {
	return fmt.Sprintf("sliding_window:%s", key)
}

// Reset 重置指定 key 的限流器
func (s *SlidingWindowLimiter) Reset(ctx context.Context, key string) error {
	return s.cache.Del(ctx, s.getKey(key))
}