package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// MemoryCache 内存缓存实现
type MemoryCache struct {
	mu    sync.RWMutex
	data  map[string]*item
	close chan struct{}
}

// item 缓存项
type item struct {
	value      interface{}
	expiration time.Time
}

// isExpired 检查是否过期
func (i *item) isExpired() bool {
	return !i.expiration.IsZero() && time.Now().After(i.expiration)
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{
		data:  make(map[string]*item),
		close: make(chan struct{}),
	}
	// 启动清理过期数据的协程
	go c.cleanup()
	return c
}

// cleanup 定期清理过期数据
func (m *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.deleteExpired()
		case <-m.close:
			return
		}
	}
}

// deleteExpired 删除过期数据
func (m *MemoryCache) deleteExpired() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for k, v := range m.data {
		if !v.expiration.IsZero() && now.After(v.expiration) {
			delete(m.data, k)
		}
	}
}

// Get 获取缓存（字符串）
func (m *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if item, ok := m.data[key]; ok {
		if item.isExpired() {
			delete(m.data, key)
			return "", nil
		}
		if str, ok := item.value.(string); ok {
			return str, nil
		}
		return "", nil
	}
	return "", nil
}

// Set 设置缓存（字符串）
func (m *MemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var exp time.Time
	if expiration > 0 {
		exp = time.Now().Add(expiration)
	}

	m.data[key] = &item{
		value:      value,
		expiration: exp,
	}
	return nil
}

// Del 删除缓存
func (m *MemoryCache) Del(ctx context.Context, keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, key := range keys {
		delete(m.data, key)
	}
	return nil
}

// Exists 检查key是否存在
func (m *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if item, ok := m.data[key]; ok {
		if item.isExpired() {
			delete(m.data, key)
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

// Expire 设置过期时间
func (m *MemoryCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if item, ok := m.data[key]; ok {
		if expiration > 0 {
			item.expiration = time.Now().Add(expiration)
		} else {
			item.expiration = time.Time{}
		}
		return nil
	}
	return nil
}

// TTL 获取剩余过期时间
func (m *MemoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if item, ok := m.data[key]; ok {
		if item.expiration.IsZero() {
			return -1, nil // 永不过期
		}
		ttl := time.Until(item.expiration)
		if ttl < 0 {
			return 0, nil // 已过期
		}
		return ttl, nil
	}
	return -2, nil // key不存在
}

// Flush 清空所有缓存
func (m *MemoryCache) Flush(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]*item)
	return nil
}

// GetBytes 获取缓存（字节数组）
func (m *MemoryCache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if item, ok := m.data[key]; ok {
		if item.isExpired() {
			delete(m.data, key)
			return nil, nil
		}
		if bytes, ok := item.value.([]byte); ok {
			return bytes, nil
		}
	}
	return nil, nil
}

// SetBytes 设置缓存（字节数组）
func (m *MemoryCache) SetBytes(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return m.Set(ctx, key, value, expiration)
}

// GetObject 获取缓存（对象，自动反序列化）
func (m *MemoryCache) GetObject(ctx context.Context, key string, dest interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if item, ok := m.data[key]; ok {
		if item.isExpired() {
			delete(m.data, key)
			return nil
		}
		if bytes, ok := item.value.([]byte); ok {
			return json.Unmarshal(bytes, dest)
		}
	}
	return nil
}

// SetObject 设置缓存（对象，自动序列化）
func (m *MemoryCache) SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.SetBytes(ctx, key, data, expiration)
}

// Incr 自增
func (m *MemoryCache) Incr(ctx context.Context, key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var val int64
	if item, ok := m.data[key]; ok && !item.isExpired() {
		switch v := item.value.(type) {
		case int:
			val = int64(v)
		case int64:
			val = v
		case float64:
			val = int64(v)
		}
	}
	val++

	var exp time.Time
	if item, ok := m.data[key]; ok {
		exp = item.expiration
	}

	m.data[key] = &item{
		value:      val,
		expiration: exp,
	}
	return val, nil
}

// IncrBy 自增指定值
func (m *MemoryCache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var val int64
	if item, ok := m.data[key]; ok && !item.isExpired() {
		switch v := item.value.(type) {
		case int:
			val = int64(v)
		case int64:
			val = v
		case float64:
			val = int64(v)
		}
	}
	val += value

	var exp time.Time
	if item, ok := m.data[key]; ok {
		exp = item.expiration
	}

	m.data[key] = &item{
		value:      val,
		expiration: exp,
	}
	return val, nil
}

// Decr 自减
func (m *MemoryCache) Decr(ctx context.Context, key string) (int64, error) {
	return m.IncrBy(ctx, key, -1)
}

// DecrBy 自减指定值
func (m *MemoryCache) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return m.IncrBy(ctx, key, -value)
}

// HGet 哈希表获取字段
func (m *MemoryCache) HGet(ctx context.Context, key, field string) (string, error) {
	// 内存缓存简化实现，使用嵌套map
	m.mu.RLock()
	defer m.mu.RUnlock()

	hashKey := "hash:" + key
	if item, ok := m.data[hashKey]; ok && !item.isExpired() {
		if hash, ok := item.value.(map[string]interface{}); ok {
			if val, ok := hash[field].(string); ok {
				return val, nil
			}
		}
	}
	return "", nil
}

// HSet 哈希表设置字段
func (m *MemoryCache) HSet(ctx context.Context, key, field string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	hashKey := "hash:" + key
	var exp time.Time
	var hash map[string]interface{}

	if item, ok := m.data[hashKey]; ok && !item.isExpired() {
		if h, ok := item.value.(map[string]interface{}); ok {
			hash = h
		}
		exp = item.expiration
	}

	if hash == nil {
		hash = make(map[string]interface{})
	}

	hash[field] = value
	m.data[hashKey] = &item{
		value:      hash,
		expiration: exp,
	}
	return nil
}

// HDel 哈希表删除字段
func (m *MemoryCache) HDel(ctx context.Context, key string, fields ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	hashKey := "hash:" + key
	if item, ok := m.data[hashKey]; ok && !item.isExpired() {
		if hash, ok := item.value.(map[string]interface{}); ok {
			for _, field := range fields {
				delete(hash, field)
			}
			item.value = hash
		}
	}
	return nil
}

// HGetAll 哈希表获取所有字段
func (m *MemoryCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hashKey := "hash:" + key
	if item, ok := m.data[hashKey]; ok && !item.isExpired() {
		if hash, ok := item.value.(map[string]interface{}); ok {
			result := make(map[string]string)
			for k, v := range hash {
				if str, ok := v.(string); ok {
					result[k] = str
				}
			}
			return result, nil
		}
	}
	return make(map[string]string), nil
}

// LPush 列表左侧推入
func (m *MemoryCache) LPush(ctx context.Context, key string, values ...interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	listKey := "list:" + key
	var exp time.Time
	var list []interface{}

	if item, ok := m.data[listKey]; ok && !item.isExpired() {
		if l, ok := item.value.([]interface{}); ok {
			list = l
		}
		exp = item.expiration
	}

	list = append(values, list...)
	m.data[listKey] = &item{
		value:      list,
		expiration: exp,
	}
	return nil
}

// RPush 列表右侧推入
func (m *MemoryCache) RPush(ctx context.Context, key string, values ...interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	listKey := "list:" + key
	var exp time.Time
	var list []interface{}

	if item, ok := m.data[listKey]; ok && !item.isExpired() {
		if l, ok := item.value.([]interface{}); ok {
			list = l
		}
		exp = item.expiration
	}

	list = append(list, values...)
	m.data[listKey] = &item{
		value:      list,
		expiration: exp,
	}
	return nil
}

// LPop 列表左侧弹出
func (m *MemoryCache) LPop(ctx context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	listKey := "list:" + key
	if item, ok := m.data[listKey]; ok && !item.isExpired() {
		if list, ok := item.value.([]interface{}); ok && len(list) > 0 {
			val := list[0]
			item.value = list[1:]
			if str, ok := val.(string); ok {
				return str, nil
			}
		}
	}
	return "", nil
}

// RPop 列表右侧弹出
func (m *MemoryCache) RPop(ctx context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	listKey := "list:" + key
	if item, ok := m.data[listKey]; ok && !item.isExpired() {
		if list, ok := item.value.([]interface{}); ok && len(list) > 0 {
			n := len(list) - 1
			val := list[n]
			item.value = list[:n]
			if str, ok := val.(string); ok {
				return str, nil
			}
		}
	}
	return "", nil
}

// LLen 列表长度
func (m *MemoryCache) LLen(ctx context.Context, key string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	listKey := "list:" + key
	if item, ok := m.data[listKey]; ok && !item.isExpired() {
		if list, ok := item.value.([]interface{}); ok {
			return int64(len(list)), nil
		}
	}
	return 0, nil
}

// SAdd 集合添加成员
func (m *MemoryCache) SAdd(ctx context.Context, key string, members ...interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	setKey := "set:" + key
	var exp time.Time
	set := make(map[interface{}]struct{})

	if item, ok := m.data[setKey]; ok && !item.isExpired() {
		if s, ok := item.value.(map[interface{}]struct{}); ok {
			set = s
		}
		exp = item.expiration
	}

	for _, member := range members {
		set[member] = struct{}{}
	}

	m.data[setKey] = &item{
		value:      set,
		expiration: exp,
	}
	return nil
}

// SRem 集合移除成员
func (m *MemoryCache) SRem(ctx context.Context, key string, members ...interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	setKey := "set:" + key
	if item, ok := m.data[setKey]; ok && !item.isExpired() {
		if set, ok := item.value.(map[interface{}]struct{}); ok {
			for _, member := range members {
				delete(set, member)
			}
		}
	}
	return nil
}

// SMembers 集合获取所有成员
func (m *MemoryCache) SMembers(ctx context.Context, key string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	setKey := "set:" + key
	if item, ok := m.data[setKey]; ok && !item.isExpired() {
		if set, ok := item.value.(map[interface{}]struct{}); ok {
			result := make([]string, 0, len(set))
			for member := range set {
				if str, ok := member.(string); ok {
					result = append(result, str)
				}
			}
			return result, nil
		}
	}
	return []string{}, nil
}

// SIsMember 集合检查成员是否存在
func (m *MemoryCache) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	setKey := "set:" + key
	if item, ok := m.data[setKey]; ok && !item.isExpired() {
		if set, ok := item.value.(map[interface{}]struct{}); ok {
			_, exists := set[member]
			return exists, nil
		}
	}
	return false, nil
}

// ZAdd 有序集合添加成员
func (m *MemoryCache) ZAdd(ctx context.Context, key string, members ...*ZMember) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	zKey := "zset:" + key
	var exp time.Time
	zset := make(map[string]float64)

	if item, ok := m.data[zKey]; ok && !item.isExpired() {
		if z, ok := item.value.(map[string]float64); ok {
			zset = z
		}
		exp = item.expiration
	}

	for _, member := range members {
		if str, ok := member.Member.(string); ok {
			zset[str] = member.Score
		}
	}

	m.data[zKey] = &item{
		value:      zset,
		expiration: exp,
	}
	return nil
}

// ZRem 有序集合移除成员
func (m *MemoryCache) ZRem(ctx context.Context, key string, members ...interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	zKey := "zset:" + key
	if item, ok := m.data[zKey]; ok && !item.isExpired() {
		if zset, ok := item.value.(map[string]float64); ok {
			for _, member := range members {
				if str, ok := member.(string); ok {
					delete(zset, str)
				}
			}
		}
	}
	return nil
}

// ZRange 有序集合按索引范围获取成员
func (m *MemoryCache) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	// 内存缓存简化实现，返回所有成员
	m.mu.RLock()
	defer m.mu.RUnlock()

	zKey := "zset:" + key
	if item, ok := m.data[zKey]; ok && !item.isExpired() {
		if zset, ok := item.value.(map[string]float64); ok {
			result := make([]string, 0, len(zset))
			for member := range zset {
				result = append(result, member)
			}
			return result, nil
		}
	}
	return []string{}, nil
}

// ZRangeByScore 有序集合按分数范围获取成员
func (m *MemoryCache) ZRangeByScore(ctx context.Context, key string, min, max float64) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	zKey := "zset:" + key
	if item, ok := m.data[zKey]; ok && !item.isExpired() {
		if zset, ok := item.value.(map[string]float64); ok {
			result := []string{}
			for member, score := range zset {
				if score >= min && score <= max {
					result = append(result, member)
				}
			}
			return result, nil
		}
	}
	return []string{}, nil
}

// ZIncrBy 有序集合成员分数增加
func (m *MemoryCache) ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	zKey := "zset:" + key
	var exp time.Time
	zset := make(map[string]float64)

	if item, ok := m.data[zKey]; ok && !item.isExpired() {
		if z, ok := item.value.(map[string]float64); ok {
			zset = z
		}
		exp = item.expiration
	}

	zset[member] += increment

	m.data[zKey] = &item{
		value:      zset,
		expiration: exp,
	}
	return zset[member], nil
}

// Close 关闭缓存
func (m *MemoryCache) Close() error {
	close(m.close)
	return nil
}

// Ping 检查健康状态
func (m *MemoryCache) Ping(ctx context.Context) error {
	return nil
}