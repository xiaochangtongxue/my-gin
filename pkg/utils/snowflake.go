// Package utils 通用工具包
package utils

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrInvalidMachineID 机器ID无效
	ErrInvalidMachineID = errors.New("机器ID必须在0-1023之间")
	// ErrClockMovedBackwards 时钟回拨错误
	ErrClockMovedBackwards = errors.New("时钟回拨")
)

const (
	// 时间戳位数
	timestampBits = 41
	// 机器ID位数
	machineIDBits = 10
	// 序列号位数
	sequenceBits = 12

	// 最大值
	maxMachineID = -1 ^ (-1 << machineIDBits) // 1023
	maxSequence  = -1 ^ (-1 << sequenceBits)  // 4095

	// 位移
	machineIDShift = sequenceBits
	timestampShift = sequenceBits + machineIDBits

	// 纪元起始时间 (2024-01-01 00:00:00 UTC)
	epoch int64 = 1704067200000
)

// Snowflake Snowflake ID 生成器
type Snowflake struct {
	mu          sync.Mutex
	machineID   int64  // 机器ID (0-1023)
	sequence    int64  // 序列号
	lastTime    int64  // 上次生成时间戳
	maxSequence int64  // 最大序列号（用于测试或特殊场景）
}

// NewSnowflake 创建 Snowflake ID 生成器
// machineID: 机器ID (0-1023)
func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, ErrInvalidMachineID
	}
	return &Snowflake{
		machineID:   machineID,
		sequence:    0,
		lastTime:    0,
		maxSequence: maxSequence,
	}, nil
}

// Generate 生成 Snowflake ID
func (s *Snowflake) Generate() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentTime := s.currentTimeMillis()

	// 时钟回拨检测
	if currentTime < s.lastTime {
		return 0, ErrClockMovedBackwards
	}

	// 同一毫秒内，序列号自增
	if currentTime == s.lastTime {
		s.sequence = (s.sequence + 1) & s.maxSequence
		// 序列号溢出，等待下一毫秒
		if s.sequence == 0 {
			currentTime = s.waitNextMillis(currentTime)
		}
	} else {
		// 不同毫秒，序列号重置
		s.sequence = 0
	}

	s.lastTime = currentTime

	// 组装ID
	id := ((currentTime - epoch) << timestampShift) |
		(s.machineID << machineIDShift) |
		s.sequence

	return id, nil
}

// GenerateString 生成字符串格式的 Snowflake ID
func (s *Snowflake) GenerateString() (string, error) {
	id, err := s.Generate()
	if err != nil {
		return "", err
	}
	return Int64ToString(id), nil
}

// ParseSnowflake 解析 Snowflake ID
// 返回: 时间戳, 机器ID, 序列号
func ParseSnowflake(id int64) (int64, int64, int64) {
	timestamp := (id >> timestampShift) + epoch
	machineID := (id >> machineIDShift) & maxMachineID
	sequence := id & maxSequence
	return timestamp, machineID, sequence
}

// GetTimestamp 从 Snowflake ID 获取时间戳
func GetTimestamp(id int64) int64 {
	timestamp, _, _ := ParseSnowflake(id)
	return timestamp
}

// GetMachineID 从 Snowflake ID 获取机器ID
func GetMachineID(id int64) int64 {
	_, machineID, _ := ParseSnowflake(id)
	return machineID
}

// GetSequence 从 Snowflake ID 获取序列号
func GetSequence(id int64) int64 {
	_, _, sequence := ParseSnowflake(id)
	return sequence
}

// TimestampToTime 将 Snowflake 时间戳转换为 time.Time
func TimestampToTime(timestamp int64) time.Time {
	return time.UnixMilli(timestamp)
}

// GetTime 从 Snowflake ID 获取生成时间
func GetTime(id int64) time.Time {
	return TimestampToTime(GetTimestamp(id))
}

// currentTimeMillis 获取当前时间戳（毫秒）
func (s *Snowflake) currentTimeMillis() int64 {
	return time.Now().UnixMilli()
}

// waitNextMillis 等待下一毫秒
func (s *Snowflake) waitNextMillis(currentTime int64) int64 {
	for {
		now := s.currentTimeMillis()
		if now > currentTime {
			return now
		}
		// 短暂休眠避免CPU占用过高
		time.Sleep(time.Microsecond * 100)
	}
}

// Int64ToString 快速将 int64 转换为字符串
func Int64ToString(n int64) string {
	var buf [20]byte
	i := 20
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte(n%10) + '0'
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// 全局 Snowflake 实例
var defaultSnowflake *Snowflake

// InitSnowflake 初始化全局 Snowflake 生成器
func InitSnowflake(machineID int64) error {
	sf, err := NewSnowflake(machineID)
	if err != nil {
		return err
	}
	defaultSnowflake = sf
	return nil
}

// GenerateSnowflakeID 生成 Snowflake ID（使用全局实例）
func GenerateSnowflakeID() (int64, error) {
	if defaultSnowflake == nil {
		// 使用默认机器ID 0
		if err := InitSnowflake(0); err != nil {
			return 0, err
		}
	}
	return defaultSnowflake.Generate()
}

// GenerateSnowflakeIDString 生成字符串格式 Snowflake ID（使用全局实例）
func GenerateSnowflakeIDString() (string, error) {
	id, err := GenerateSnowflakeID()
	if err != nil {
		return "", err
	}
	return Int64ToString(id), nil
}

// MustGenerateSnowflakeID 生成 Snowflake ID，出错时 panic
func MustGenerateSnowflakeID() int64 {
	id, err := GenerateSnowflakeID()
	if err != nil {
		panic(err)
	}
	return id
}