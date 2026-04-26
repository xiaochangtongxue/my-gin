package utils

import (
	"testing"
	"time"
)

// TestFormatDateTime 测试格式化日期时间
func TestFormatDateTime(t *testing.T) {
	tm := time.Date(2024, 1, 5, 12, 30, 45, 0, time.UTC)
	result := FormatDateTime(tm)
	expected := "2024-01-05 12:30:45"
	if result != expected {
		t.Errorf("FormatDateTime() = %q, want %q", result, expected)
	}
}

// TestFormatDateTimeString 测试格式化时间戳
func TestFormatDateTimeString(t *testing.T) {
	timestamp := int64(1704450645) // 2024-01-05 10:30:45 UTC
	result := FormatDateTimeString(timestamp)
	expected := "2024-01-05 10:30:45"
	if result != expected {
		t.Errorf("FormatDateTimeString(%d) = %q, want %q", timestamp, result, expected)
	}
}

// TestFormatDate 测试格式化日期
func TestFormatDate(t *testing.T) {
	tm := time.Date(2024, 1, 5, 12, 30, 45, 0, time.UTC)
	result := FormatDate(tm)
	expected := "2024-01-05"
	if result != expected {
		t.Errorf("FormatDate() = %q, want %q", result, expected)
	}
}

// TestFormatTime 测试格式化时间
func TestFormatTime(t *testing.T) {
	tm := time.Date(2024, 1, 5, 12, 30, 45, 0, time.UTC)
	result := FormatTime(tm)
	expected := "12:30:45"
	if result != expected {
		t.Errorf("FormatTime() = %q, want %q", result, expected)
	}
}

// TestParseDateTime 测试解析日期时间
func TestParseDateTime(t *testing.T) {
	input := "2024-01-05 12:30:45"
	result, err := ParseDateTime(input)
	if err != nil {
		t.Fatalf("ParseDateTime(%q) error = %v", input, err)
	}
	expected := time.Date(2024, 1, 5, 12, 30, 45, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("ParseDateTime(%q) = %v, want %v", input, result, expected)
	}
}

// TestParseDate 测试解析日期
func TestParseDate(t *testing.T) {
	input := "2024-01-05"
	result, err := ParseDate(input)
	if err != nil {
		t.Fatalf("ParseDate(%q) error = %v", input, err)
	}
	expected := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("ParseDate(%q) = %v, want %v", input, result, expected)
	}
}

// TestNowUnix 测试获取当前时间戳
func TestNowUnix(t *testing.T) {
	before := time.Now().Unix()
	result := NowUnix()
	after := time.Now().Unix()

	if result < before || result > after {
		t.Errorf("NowUnix() = %d, want between %d and %d", result, before, after)
	}
}

// TestNowUnixMilli 测试获取当前毫秒时间戳
func TestNowUnixMilli(t *testing.T) {
	result := NowUnixMilli()
	if result <= 0 {
		t.Errorf("NowUnixMilli() = %d, want > 0", result)
	}
}

// TestUnixToTime 测试时间戳转时间
func TestUnixToTime(t *testing.T) {
	timestamp := int64(1704450645)
	result := UnixToTime(timestamp)
	expected := time.Date(2024, 1, 5, 10, 30, 45, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("UnixToTime(%d) = %v, want %v", timestamp, result, expected)
	}
}

// TestTimeDiff 测试计算时间差
func TestTimeDiff(t *testing.T) {
	t1 := time.Date(2024, 1, 5, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 5, 12, 30, 45, 0, time.UTC)
	result := TimeDiff(t1, t2)
	expected := int64(9045) // 2小时30分45秒 = 9045秒
	if result != expected {
		t.Errorf("TimeDiff(%v, %v) = %d, want %d", t1, t2, result, expected)
	}
}

// TestTimeDiffAbs 测试计算时间差绝对值
func TestTimeDiffAbs(t *testing.T) {
	t1 := time.Date(2024, 1, 5, 12, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 5, 10, 0, 0, 0, time.UTC)
	result := TimeDiffAbs(t1, t2)
	expected := int64(7200) // 2小时
	if result != expected {
		t.Errorf("TimeDiffAbs(%v, %v) = %d, want %d", t1, t2, result, expected)
	}
}

// TestRemainingSeconds 测试计算剩余秒数
func TestRemainingSeconds(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		t        time.Time
		expected int64
	}{
		{"未来时间", now.Add(1 * time.Hour), 3600},
		{"过去时间", now.Add(-1 * time.Hour), 0},
		{"当前时间", now, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemainingSeconds(tt.t)
			// 允许1秒误差
			if tt.expected > 0 && (result < tt.expected-1 || result > tt.expected+1) {
				t.Errorf("RemainingSeconds(%v) = %d, want %d", tt.t, result, tt.expected)
			}
			if tt.expected == 0 && result != 0 {
				t.Errorf("RemainingSeconds(%v) = %d, want %d", tt.t, result, tt.expected)
			}
		})
	}
}

// TestIsExpired 测试判断时间是否过期
func TestIsExpired(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		t        time.Time
		expected bool
	}{
		{"过去时间", now.Add(-1 * time.Hour), true},
		{"未来时间", now.Add(1 * time.Hour), false},
		{"当前时间", now, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsExpired(tt.t)
			if result != tt.expected {
				t.Errorf("IsExpired(%v) = %v, want %v", tt.t, result, tt.expected)
			}
		})
	}
}

// TestIsExpiredUnix 测试判断时间戳是否过期
func TestIsExpiredUnix(t *testing.T) {
	now := time.Now().Unix()
	tests := []struct {
		name      string
		timestamp int64
		expected  bool
	}{
		{"过去时间", now - 3600, true},
		{"未来时间", now + 3600, false},
		{"当前时间", now, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsExpiredUnix(tt.timestamp)
			if result != tt.expected {
				t.Errorf("IsExpiredUnix(%d) = %v, want %v", tt.timestamp, result, tt.expected)
			}
		})
	}
}

// TestStartOfDay 测试获取一天开始时间
func TestStartOfDay(t *testing.T) {
	tm := time.Date(2024, 1, 5, 12, 30, 45, 123456789, time.UTC)
	result := StartOfDay(tm)
	expected := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("StartOfDay(%v) = %v, want %v", tm, result, expected)
	}
}

// TestEndOfDay 测试获取一天结束时间
func TestEndOfDay(t *testing.T) {
	tm := time.Date(2024, 1, 5, 12, 30, 45, 123456789, time.UTC)
	result := EndOfDay(tm)
	expected := time.Date(2024, 1, 5, 23, 59, 59, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("EndOfDay(%v) = %v, want %v", tm, result, expected)
	}
}

// TestStartOfWeek 测试获取一周开始时间
func TestStartOfWeek(t *testing.T) {
	// 2024-01-05 是星期五
	tm := time.Date(2024, 1, 5, 12, 30, 45, 0, time.UTC)
	result := StartOfWeek(tm)
	// 应该是周一 2024-01-01
	expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("StartOfWeek(%v) = %v (weekday: %v), want %v", tm, result, tm.Weekday(), expected)
	}
}

// TestEndOfWeek 测试获取一周结束时间
func TestEndOfWeek(t *testing.T) {
	// 2024-01-05 是星期五
	tm := time.Date(2024, 1, 5, 12, 30, 45, 0, time.UTC)
	result := EndOfWeek(tm)
	// 应该是周日 2024-01-07 23:59:59
	expected := time.Date(2024, 1, 7, 23, 59, 59, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("EndOfWeek(%v) = %v, want %v", tm, result, expected)
	}
}

// TestStartOfMonth 测试获取一月开始时间
func TestStartOfMonth(t *testing.T) {
	tm := time.Date(2024, 1, 15, 12, 30, 45, 0, time.UTC)
	result := StartOfMonth(tm)
	expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("StartOfMonth(%v) = %v, want %v", tm, result, expected)
	}
}

// TestEndOfMonth 测试获取一月结束时间
func TestEndOfMonth(t *testing.T) {
	tm := time.Date(2024, 1, 15, 12, 30, 45, 0, time.UTC)
	result := EndOfMonth(tm)
	expected := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("EndOfMonth(%v) = %v, want %v", tm, result, expected)
	}
}

// TestStartOfYear 测试获取一年开始时间
func TestStartOfYear(t *testing.T) {
	tm := time.Date(2024, 6, 15, 12, 30, 45, 0, time.UTC)
	result := StartOfYear(tm)
	expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("StartOfYear(%v) = %v, want %v", tm, result, expected)
	}
}

// TestEndOfYear 测试获取一年结束时间
func TestEndOfYear(t *testing.T) {
	tm := time.Date(2024, 6, 15, 12, 30, 45, 0, time.UTC)
	result := EndOfYear(tm)
	expected := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("EndOfYear(%v) = %v, want %v", tm, result, expected)
	}
}

// TestAddDays 测试添加天数
func TestAddDays(t *testing.T) {
	tm := time.Date(2024, 1, 5, 12, 30, 45, 0, time.UTC)
	result := AddDays(tm, 7)
	expected := time.Date(2024, 1, 12, 12, 30, 45, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("AddDays(%v, 7) = %v, want %v", tm, result, expected)
	}
}

// TestAddHours 测试添加小时数
func TestAddHours(t *testing.T) {
	tm := time.Date(2024, 1, 5, 12, 30, 45, 0, time.UTC)
	result := AddHours(tm, 2)
	expected := time.Date(2024, 1, 5, 14, 30, 45, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("AddHours(%v, 2) = %v, want %v", tm, result, expected)
	}
}

// TestAddMinutes 测试添加分钟数
func TestAddMinutes(t *testing.T) {
	tm := time.Date(2024, 1, 5, 12, 30, 45, 0, time.UTC)
	result := AddMinutes(tm, 30)
	expected := time.Date(2024, 1, 5, 13, 0, 45, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("AddMinutes(%v, 30) = %v, want %v", tm, result, expected)
	}
}

// TestAddSeconds 测试添加秒数
func TestAddSeconds(t *testing.T) {
	tm := time.Date(2024, 1, 5, 12, 30, 45, 0, time.UTC)
	result := AddSeconds(tm, 15)
	expected := time.Date(2024, 1, 5, 12, 31, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("AddSeconds(%v, 15) = %v, want %v", tm, result, expected)
	}
}

// TestBetween 测试判断时间是否在两个时间之间
func TestBetween(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
	middle := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	if !Between(middle, start, end) {
		t.Errorf("Between(%v, %v, %v) = false, want true", middle, start, end)
	}
	if Between(start, start, end) {
		t.Errorf("Between(%v, %v, %v) = true, want false", start, start, end)
	}
	if Between(end, start, end) {
		t.Errorf("Between(%v, %v, %v) = true, want false", end, start, end)
	}
}

// TestBetweenInclusive 测试判断时间是否在两个时间之间（包含边界）
func TestBetweenInclusive(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
	middle := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	if !BetweenInclusive(middle, start, end) {
		t.Errorf("BetweenInclusive(%v, %v, %v) = false, want true", middle, start, end)
	}
	if !BetweenInclusive(start, start, end) {
		t.Errorf("BetweenInclusive(%v, %v, %v) = false, want true", start, start, end)
	}
	if !BetweenInclusive(end, start, end) {
		t.Errorf("BetweenInclusive(%v, %v, %v) = false, want true", end, start, end)
	}
}

// BenchmarkFormatDateTime 基准测试
func BenchmarkFormatDateTime(b *testing.B) {
	tm := time.Now()
	for i := 0; i < b.N; i++ {
		FormatDateTime(tm)
	}
}
