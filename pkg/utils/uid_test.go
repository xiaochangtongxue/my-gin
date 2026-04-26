package utils

import (
	"testing"
)

// TestEncodeUID 测试UID编码
func TestEncodeUID(t *testing.T) {
	tests := []struct {
		name     string
		id       uint64
		want14digits bool
	}{
		{"ID为1", 1, true},
		{"ID为100", 100, true},
		{"ID为10000", 10000, true},
		{"ID为0", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid := EncodeUID(tt.id)
			// 验证结果是14位数字
			if uid < 10000000000000 || uid > 99999999999999 {
				t.Errorf("EncodeUID(%d) = %d, not 14 digits", tt.id, uid)
			}
		})
	}
}

// TestDecodeUID 测试UID解码
func TestDecodeUID(t *testing.T) {
	tests := []struct {
		name string
		id   uint64
	}{
		{"ID为1", 1},
		{"ID为100", 100},
		{"ID为10000", 10000},
		{"ID为0", 0},
		{"ID为1000000", 1000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid := EncodeUID(tt.id)
			decoded := DecodeUID(uid)
			if decoded != tt.id {
				t.Errorf("EncodeUID(%d) -> DecodeUID() = %d, want %d", tt.id, decoded, tt.id)
			}
		})
	}
}

// TestUIDUniqueness 测试UID唯一性
func TestUIDUniqueness(t *testing.T) {
	// 测试不同的ID产生不同的UID
	ids := make(map[uint64]bool)
	for i := uint64(1); i <= 10000; i++ {
		uid := EncodeUID(i)
		if ids[uid] {
			t.Errorf("UID collision detected: ID %d produced UID that already exists", i)
		}
		ids[uid] = true
	}
}

// TestUIDSequentialPattern 测试UID不会暴露原始ID的顺序
func TestUIDSequentialPattern(t *testing.T) {
	// 测试连续的ID产生的UID不会看起来连续
	uids := make([]uint64, 10)
	for i := 0; i < 10; i++ {
		uids[i] = EncodeUID(uint64(i + 1))
	}

	// 检查相邻UID的差值不是常数1（即不是简单的连续）
	// Feistel混淆应该使相邻差值变大
	for i := 1; i < 10; i++ {
		diff := uids[i] - uids[i-1]
		if diff == 1 {
			t.Errorf("UIDs are sequential: %d and %d differ by 1", uids[i-1], uids[i])
		}
	}
}

// TestUIDEncodeDecode 测试大量编码解码
func TestUIDEncodeDecode(t *testing.T) {
	testIDs := []uint64{
		1, 2, 3, 100, 1000, 10000, 100000, 1000000, 10000000,
		123456789, 987654321, 0, 9999999999,
	}

	for _, id := range testIDs {
		t.Run("", func(t *testing.T) {
			uid := EncodeUID(id)
			decoded := DecodeUID(uid)
			if decoded != id {
				t.Errorf("Encode/Decode round trip failed: %d -> %d -> %d", id, uid, decoded)
			}
		})
	}
}

// TestUIDRange 测试UID范围
func TestUIDRange(t *testing.T) {
	// 测试边界情况
	testCases := []uint64{0, 1, 100, 1000000}

	for _, id := range testCases {
		uid := EncodeUID(id)
		// UID应该是14位数字
		if uid < 10000000000000 {
			t.Errorf("UID %d is less than 14 digits", uid)
		}
		if uid > 99999999999999 {
			t.Errorf("UID %d exceeds 14 digits", uid)
		}
	}
}

// BenchmarkEncodeUID 基准测试
func BenchmarkEncodeUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EncodeUID(12345)
	}
}

// BenchmarkDecodeUID 基准测试
func BenchmarkDecodeUID(b *testing.B) {
	uid := EncodeUID(12345)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DecodeUID(uid)
	}
}
