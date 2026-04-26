package crypto

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

// TestMD5 测试MD5哈希
func TestMD5(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "空字符串",
			input:    "",
			expected: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:     "hello",
			input:    "hello",
			expected: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name:     "hello world",
			input:    "hello world",
			expected: "5eb63bbbe01eeed093cb22bb8f5acdc3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MD5(tt.input)
			if result != tt.expected {
				t.Errorf("MD5(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSHA256 测试SHA256哈希
func TestSHA256(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "空字符串",
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "hello",
			input:    "hello",
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SHA256(tt.input)
			if result != tt.expected {
				t.Errorf("SHA256(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestHMACSHA256 测试HMAC-SHA256签名
func TestHMACSHA256(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		data     string
		expected string
	}{
		{
			name:     "基本签名",
			key:      "secret-key",
			data:     "message",
			expected: "287a3bd8a4fc7731a94c722079055323644d8798bd291bf9878abc9b8fd4b1d0",
		},
		{
			name:     "空数据",
			key:      "secret-key",
			data:     "",
			expected: "345fba21f06a4f75ed673fb93dc16cd47d8dc7a69f52e84e3016fcf69835fdb8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HMACSHA256(tt.key, tt.data)
			if result != tt.expected {
				t.Errorf("HMACSHA256(%q, %q) = %q, want %q", tt.key, tt.data, result, tt.expected)
			}
		})
	}
}

// TestHMACSHA256Consistency 测试HMAC签名一致性
func TestHMACSHA256Consistency(t *testing.T) {
	key := "test-key"
	data := "test-data"

	result1 := HMACSHA256(key, data)
	result2 := HMACSHA256(key, data)

	if result1 != result2 {
		t.Errorf("HMACSHA256 produced different results for same input: %s != %s", result1, result2)
	}
}

// TestHMACSHA256DifferentKeys 测试不同密钥产生不同签名
func TestHMACSHA256DifferentKeys(t *testing.T) {
	data := "test-data"

	result1 := HMACSHA256("key1", data)
	result2 := HMACSHA256("key2", data)

	if result1 == result2 {
		t.Errorf("HMACSHA256 produced same results for different keys")
	}
}

// TestBase64Encode 测试Base64编码
func TestBase64Encode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "基本编码",
			input:    "hello",
			expected: "aGVsbG8=",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "中文字符",
			input:    "你好",
			expected: "5L2g5aW9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base64Encode(tt.input)
			if result != tt.expected {
				t.Errorf("Base64Encode(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestBase64Decode 测试Base64解码
func TestBase64Decode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "基本解码",
			input:    "aGVsbG8=",
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "中文字符",
			input:    "5L2g5aW9",
			expected: "你好",
			wantErr:  false,
		},
		{
			name:     "无效Base64",
			input:    "!!!invalid!!!",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Base64Decode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Base64Decode(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("Base64Decode(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestBase64EncodeDecode 测试Base64编解码往返
func TestBase64EncodeDecode(t *testing.T) {
	tests := []string{
		"",
		"hello",
		"hello world",
		"特殊字符 !@#$%",
		"你好世界",
		"aGVsbG8gd29ybGQ=", // 已编码的字符串
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			encoded := Base64Encode(input)
			decoded, err := Base64Decode(encoded)
			if err != nil {
				t.Errorf("Base64Decode(%q) error = %v", encoded, err)
				return
			}
			if decoded != input {
				t.Errorf("Base64Encode/Decode round trip: %q -> %q -> %q", input, encoded, decoded)
			}
		})
	}
}

// TestBase64URLEncode 测试URL Safe Base64编码
func TestBase64URLEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "基本编码",
			input:    "hello?",
			expected: "aGVsbG8_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base64URLEncode(tt.input)
			// URL Safe 编码会把 + 变成 -，/ 变成 _
			if result != tt.expected {
				t.Errorf("Base64URLEncode(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestBase64URLDecode 测试URL Safe Base64解码
func TestBase64URLDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL Safe字符",
			input:    "aGVsbG8",
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Base64URLDecode(tt.input)
			if err != nil {
				t.Errorf("Base64URLDecode(%q) error = %v", tt.input, err)
				return
			}
			if result != tt.expected {
				t.Errorf("Base64URLDecode(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRandomString 测试随机字符串生成
func TestRandomString(t *testing.T) {
	// 测试生成的字符串长度
	lengths := []int{32, 64, 128}

	for _, length := range lengths {
		t.Run("length_"+string(rune('0'+length)), func(t *testing.T) {
			result, err := RandomString(length)
			if err != nil {
				t.Errorf("RandomString(%d) error = %v", length, err)
				return
			}
			if len(result) != length {
				t.Errorf("RandomString(%d) length = %d, want %d", length, len(result), length)
			}
			// 验证是有效的十六进制字符串
			_, err = hex.DecodeString(result)
			if err != nil {
				t.Errorf("RandomString(%d) produced invalid hex: %v", length, err)
			}
		})
	}
}

// TestRandomBytes 测试随机字节生成
func TestRandomBytes(t *testing.T) {
	sizes := []int{16, 32, 64}

	for _, size := range sizes {
		t.Run("size_"+string(rune('0'+size)), func(t *testing.T) {
			result, err := RandomBytes(size)
			if err != nil {
				t.Errorf("RandomBytes(%d) error = %v", size, err)
				return
			}
			if len(result) != size {
				t.Errorf("RandomBytes(%d) length = %d, want %d", size, len(result), size)
			}
		})
	}
}

// TestRandomBytesUniqueness 测试随机字节唯一性
func TestRandomBytesUniqueness(t *testing.T) {
	results := make(map[string]bool)
	for i := 0; i < 100; i++ {
		b, err := RandomBytes(16)
		if err != nil {
			t.Fatalf("RandomBytes error = %v", err)
		}
		s := base64.StdEncoding.EncodeToString(b)
		if results[s] {
			t.Errorf("Duplicate random value detected")
		}
		results[s] = true
	}
}

// BenchmarkMD5 基准测试
func BenchmarkMD5(b *testing.B) {
	data := "hello world"
	for i := 0; i < b.N; i++ {
		MD5(data)
	}
}

// BenchmarkSHA256 基准测试
func BenchmarkSHA256(b *testing.B) {
	data := "hello world"
	for i := 0; i < b.N; i++ {
		SHA256(data)
	}
}

// BenchmarkHMACSHA256 基准测试
func BenchmarkHMACSHA256(b *testing.B) {
	key := "secret-key"
	data := "message"
	for i := 0; i < b.N; i++ {
		HMACSHA256(key, data)
	}
}
