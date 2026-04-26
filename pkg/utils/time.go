// Package utils 通用工具包
package utils

import (
	"time"
)

const (
	// DateTimeFormat 标准日期时间格式
	DateTimeFormat = "2006-01-02 15:04:05"
	// DateFormat 日期格式
	DateFormat = "2006-01-02"
	// TimeFormat 时间格式
	TimeFormat = "15:04:05"
	// TimestampFormat 时间戳格式（秒）
	TimestampFormat = "2006-01-02 15:04:05"
)

// FormatDateTime 格式化时间为标准格式 YYYY-MM-DD HH:mm:ss
func FormatDateTime(t time.Time) string {
	return t.Format(DateTimeFormat)
}

// FormatDateTimeString 格式化时间戳为标准格式
func FormatDateTimeString(timestamp int64) string {
	return UnixToTime(timestamp).Format(DateTimeFormat)
}

// FormatDate 格式化日期 YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// FormatTime 格式化时间 HH:mm:ss
func FormatTime(t time.Time) string {
	return t.Format(TimeFormat)
}

// ParseDateTime 解析标准日期时间格式字符串
func ParseDateTime(s string) (time.Time, error) {
	return time.Parse(DateTimeFormat, s)
}

// ParseDate 解析日期格式字符串
func ParseDate(s string) (time.Time, error) {
	return time.Parse(DateFormat, s)
}

// NowUnix 获取当前时间戳（秒）
func NowUnix() int64 {
	return time.Now().Unix()
}

// NowUnixMilli 获取当前时间戳（毫秒）
func NowUnixMilli() int64 {
	return time.Now().UnixMilli()
}

// UnixToTime 时间戳转时间
func UnixToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).UTC()
}

// TimeDiff 计算两个时间之间的差值（秒）
func TimeDiff(t1, t2 time.Time) int64 {
	return int64(t2.Sub(t1) / time.Second)
}

// TimeDiffAbs 计算两个时间之间的差值绝对值（秒）
func TimeDiffAbs(t1, t2 time.Time) int64 {
	diff := t2.Sub(t1)
	if diff < 0 {
		diff = -diff
	}
	return int64(diff / time.Second)
}

// RemainingSeconds 计算从 t 到当前时间的剩余秒数
// 如果 t 在当前时间之前，返回 0
func RemainingSeconds(t time.Time) int64 {
	remaining := t.Unix() - time.Now().Unix()
	if remaining <= 0 {
		return 0
	}
	return remaining
}

// IsExpired 判断时间 t 是否已过期
func IsExpired(t time.Time) bool {
	return t.Before(time.Now())
}

// IsExpiredUnix 判断时间戳是否已过期
func IsExpiredUnix(timestamp int64) bool {
	return timestamp < time.Now().Unix()
}

// StartOfDay 获取一天的开始时间（00:00:00）
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay 获取一天的结束时间（23:59:59）
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}

// StartOfWeek 获取一周的开始时间（周一 00:00:00）
func StartOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	return StartOfDay(t.AddDate(0, 0, -int(weekday)+1))
}

// EndOfWeek 获取一周的结束时间（周日 23:59:59）
func EndOfWeek(t time.Time) time.Time {
	return StartOfWeek(t).AddDate(0, 0, 7).Add(-time.Second)
}

// StartOfMonth 获取一月的开始时间
func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 获取一月的结束时间
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Second)
}

// StartOfYear 获取一年的开始时间
func StartOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 获取一年的结束时间
func EndOfYear(t time.Time) time.Time {
	return StartOfYear(t).AddDate(1, 0, 0).Add(-time.Second)
}

// AddDays 添加天数
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddHours 添加小时数
func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

// AddMinutes 添加分钟数
func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// AddSeconds 添加秒数
func AddSeconds(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

// BeginningOfWeek 获取本周周一（兼容别名）
func BeginningOfWeek(t time.Time) time.Time {
	return StartOfWeek(t)
}

// BeginningOfMonth 获取本月第一天（兼容别名）
func BeginningOfMonth(t time.Time) time.Time {
	return StartOfMonth(t)
}

// BeginningOfDay 获取今天零点（兼容别名）
func BeginningOfDay(t time.Time) time.Time {
	return StartOfDay(t)
}

// Between 判断时间是否在两个时间之间
func Between(t, start, end time.Time) bool {
	return t.After(start) && t.Before(end)
}

// BetweenInclusive 判断时间是否在两个时间之间（包含边界）
func BetweenInclusive(t, start, end time.Time) bool {
	return (t.Equal(start) || t.After(start)) && (t.Equal(end) || t.Before(end))
}
