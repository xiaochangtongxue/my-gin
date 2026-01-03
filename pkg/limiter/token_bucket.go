package limiter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/xiaochangtongxue/my-gin/pkg/cache"
)

// TokenBucketLimiter 令牌桶限流器（通过 Cache 接口实现）
type TokenBucketLimiter struct {
	cache    cache.Cache
	rate     float64 // 每秒产生的令牌数
	capacity int64   // 桶容量（最大令牌数）
}

// NewTokenBucketLimiter 创建令牌桶限流器
// cache: 缓存实例
// rate: 每秒产生的令牌数
// capacity: 桶容量（最大突发令牌数）
func NewTokenBucketLimiter(cache cache.Cache, rate float64, capacity int64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		cache:    cache,
		rate:     rate,
		capacity: capacity,
	}
}

// Allow 检查是否允许一个请求通过
func (t *TokenBucketLimiter) Allow(ctx context.Context, key string) bool {
	return t.AllowN(ctx, key, 1)
}

// AllowN 检查是否允许 n 个请求通过
// 使用 Lua 脚本保证原子性
func (t *TokenBucketLimiter) AllowN(ctx context.Context, key string, n int) bool {
	// Lua 脚本：令牌桶算法
	script := `
		local key = KEYS[1]
		local rate = tonumber(ARGV[1])
		local capacity = tonumber(ARGV[2])
		local now = tonumber(ARGV[3])
		local requested = tonumber(ARGV[4])
		local expire = tonumber(ARGV[5])

		-- 获取当前桶状态：[当前令牌数, 上次更新时间]
		local info = redis.call("HMGET", key, "tokens", "last_time")
		local tokens = tonumber(info[1])
		local last_time = tonumber(info[2])

		-- 首次访问，初始化为满桶
		if tokens == nil then
			tokens = capacity
			last_time = now
		else
			-- 计算自上次更新以来生成的令牌数
			local delta = math.max(0, now - last_time)
			local generated = delta * rate
			tokens = math.min(capacity, tokens + generated)
		end
		local result = 0
		-- 检查令牌是否足够
		if tokens >= requested then
			tokens = tokens - requested
            last_time = now -- 只有成功扣减才更新时间戳，保证下一次计算准确
            result = 1
            -- 只有成功时才写入
            redis.call("HMSET", key, "tokens", tokens, "last_time", last_time)
            redis.call("EXPIRE", key, expire)
		end
		return result
	`

	// 2. 毫秒级时间戳
	now := time.Now().UnixNano() / int64(time.Second)

	// 3. 过期时间：计算桶填满需要多少秒，再加上缓冲
	// 如果 rate 是 10/s, capacity 是 100, 填满需 10s
	expireSeconds := int(float64(t.capacity)/t.rate) + 2

	result, err := t.cache.Eval(ctx, script, []string{t.getKey(key)},
		t.rate, t.capacity, now, n, expireSeconds)
	if err != nil {
		return true // 缓存出错时放行（降级策略）
	}

	if r, ok := result.(int64); ok {
		return r == 1
	}
	return false
}

// getKey 获取缓存 key
func (t *TokenBucketLimiter) getKey(key string) string {
	return fmt.Sprintf("token_bucket:%s", key)
}

// Reset 重置指定 key 的限流器
func (t *TokenBucketLimiter) Reset(ctx context.Context, key string) error {
	return t.cache.Del(ctx, t.getKey(key))
}

// GetTokens 获取当前可用令牌数（用于调试或监控）
func (t *TokenBucketLimiter) GetTokens(ctx context.Context, key string) (int64, error) {
	result, err := t.cache.HMGet(ctx, t.getKey(key), "tokens")
	if err != nil {
		return 0, err
	}
	if len(result) == 0 || result[0] == nil {
		return t.capacity, nil // 初始状态为满桶
	}
	return parseInt64(result[0].(string))
}

// parseInt64 字符串转 int64
func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
