package cache

import (
	"context"
	"time"
)

// Cache 缓存接口
type Cache interface {
	// Get 获取缓存（字符串）
	Get(ctx context.Context, key string) (string, error)
	// Set 设置缓存（字符串）
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Del 删除缓存
	Del(ctx context.Context, keys ...string) error
	// Exists 检查key是否存在
	Exists(ctx context.Context, key string) (bool, error)
	// Expire 设置过期时间
	Expire(ctx context.Context, key string, expiration time.Duration) error
	// TTL 获取剩余过期时间
	TTL(ctx context.Context, key string) (time.Duration, error)
	// Flush 清空所有缓存
	Flush(ctx context.Context) error

	// GetBytes 获取缓存（字节数组）
	GetBytes(ctx context.Context, key string) ([]byte, error)
	// SetBytes 设置缓存（字节数组）
	SetBytes(ctx context.Context, key string, value []byte, expiration time.Duration) error

	// GetObject 获取缓存（对象，自动反序列化）
	GetObject(ctx context.Context, key string, dest interface{}) error
	// SetObject 设置缓存（对象，自动序列化）
	SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Incr 自增
	Incr(ctx context.Context, key string) (int64, error)
	// IncrBy 自增指定值
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	// Decr 自减
	Decr(ctx context.Context, key string) (int64, error)
	// DecrBy 自减指定值
	DecrBy(ctx context.Context, key string, value int64) (int64, error)

	// HGet 哈希表获取字段
	HGet(ctx context.Context, key, field string) (string, error)
	// HSet 哈希表设置字段
	HSet(ctx context.Context, key, field string, value interface{}) error
	// HDel 哈希表删除字段
	HDel(ctx context.Context, key string, fields ...string) error
	// HGetAll 哈希表获取所有字段
	HGetAll(ctx context.Context, key string) (map[string]string, error)

	// LPush 列表左侧推入
	LPush(ctx context.Context, key string, values ...interface{}) error
	// RPush 列表右侧推入
	RPush(ctx context.Context, key string, values ...interface{}) error
	// LPop 列表左侧弹出
	LPop(ctx context.Context, key string) (string, error)
	// RPop 列表右侧弹出
	RPop(ctx context.Context, key string) (string, error)
	// LLen 列表长度
	LLen(ctx context.Context, key string) (int64, error)

	// SAdd 集合添加成员
	SAdd(ctx context.Context, key string, members ...interface{}) error
	// SRem 集合移除成员
	SRem(ctx context.Context, key string, members ...interface{}) error
	// SMembers 集合获取所有成员
	SMembers(ctx context.Context, key string) ([]string, error)
	// SIsMember 集合检查成员是否存在
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)

	// ZAdd 有序集合添加成员
	ZAdd(ctx context.Context, key string, members ...*ZMember) error
	// ZRem 有序集合移除成员
	ZRem(ctx context.Context, key string, members ...interface{}) error
	// ZRange 有序集合按分数范围获取成员
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	// ZRangeByScore 有序集合按分数范围获取成员
	ZRangeByScore(ctx context.Context, key string, min, max float64) ([]string, error)
	// ZIncrBy 有序集合成员分数增加
	ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error)

	// Close 关闭连接
	Close() error
	// Ping 检查连接健康状态
	Ping(ctx context.Context) error
}

// ZMember 有序集合成员
type ZMember struct {
	Score  float64
	Member interface{}
}