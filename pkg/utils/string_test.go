package utils

import "testing"

// TestMaskMobile 测试手机号脱敏
func TestMaskMobile(t *testing.T) {
	tests := []struct {
		name     string
		mobile   string
		expected string
	}{
		{"正常手机号", "13812345678", "138****5678"},
		{"11位但非手机号", "12345678901", "123****8901"},
		{"不足11位", "123456789", "123456789"},
		{"超过11位", "138123456789", "138123456789"},
		{"空字符串", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskMobile(tt.mobile)
			if result != tt.expected {
				t.Errorf("MaskMobile(%q) = %q, want %q", tt.mobile, result, tt.expected)
			}
		})
	}
}

// TestMaskEmail 测试邮箱脱敏
func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{"正常邮箱", "example@gmail.com", "ex****@gmail.com"},
		{"短用户名", "ab@test.com", "ab@test.com"},
		{"无@", "examplegmail.com", "examplegmail.com"},
		{"多个@", "exa@mple@gmail.com", "exa@mple@gmail.com"},
		{"空字符串", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskEmail(tt.email)
			if result != tt.expected {
				t.Errorf("MaskEmail(%q) = %q, want %q", tt.email, result, tt.expected)
			}
		})
	}
}

// TestMaskIDCard 测试身份证号脱敏
func TestMaskIDCard(t *testing.T) {
	tests := []struct {
		name     string
		idCard   string
		expected string
	}{
		{"18位身份证", "110101199001011234", "110101********1234"},
		{"15位身份证", "110101900101123", "110101*****1123"},
		{"短于10位", "123456789", "123456789"},
		{"空字符串", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskIDCard(tt.idCard)
			if result != tt.expected {
				t.Errorf("MaskIDCard(%q) = %q, want %q", tt.idCard, result, tt.expected)
			}
		})
	}
}

// TestMaskBankCard 测试银行卡号脱敏
func TestMaskBankCard(t *testing.T) {
	tests := []struct {
		name     string
		cardNo   string
		expected string
	}{
		{"16位银行卡", "6222021234567890", "6222********7890"},
		{"19位银行卡", "6222021234567890123", "6222***********0123"},
		{"短于8位", "1234567", "1234567"},
		{"空字符串", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskBankCard(tt.cardNo)
			if result != tt.expected {
				t.Errorf("MaskBankCard(%q) = %q, want %q", tt.cardNo, result, tt.expected)
			}
		})
	}
}

// TestMaskString 测试通用字符串脱敏
func TestMaskString(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		start    int
		end      int
		maskChar string
		expected string
	}{
		{"正常脱敏", "1234567890", 2, 3, "*", "12*****890"},
		{"自定义掩码", "1234567890", 2, 3, "#", "12#####890"},
		{"start+end超过长度", "123", 2, 3, "*", "123"},
		{"空字符串", "", 2, 3, "*", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.maskChar != "" {
				result = MaskString(tt.s, tt.start, tt.end, tt.maskChar)
			} else {
				result = MaskString(tt.s, tt.start, tt.end)
			}
			if result != tt.expected {
				t.Errorf("MaskString(%q, %d, %d) = %q, want %q", tt.s, tt.start, tt.end, result, tt.expected)
			}
		})
	}
}

// TestRandomString 测试随机字符串生成
func TestRandomString(t *testing.T) {
	lengths := []int{0, 1, 8, 16, 32}

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
			// 验证只包含字母和数字
			for _, r := range result {
				if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
					t.Errorf("RandomString contains invalid character: %c", r)
				}
			}
		})
	}
}

// TestRandomHex 测试随机十六进制字符串
func TestRandomHex(t *testing.T) {
	lengths := []int{0, 1, 8, 16}

	for _, length := range lengths {
		t.Run("length_"+string(rune('0'+length)), func(t *testing.T) {
			result, err := RandomHex(length)
			if err != nil {
				t.Errorf("RandomHex(%d) error = %v", length, err)
				return
			}
			expectedLength := length * 2
			if len(result) != expectedLength {
				t.Errorf("RandomHex(%d) length = %d, want %d", length, len(result), expectedLength)
			}
		})
	}
}

// TestRandomDigits 测试随机数字字符串
func TestRandomDigits(t *testing.T) {
	lengths := []int{0, 1, 4, 6, 8}

	for _, length := range lengths {
		t.Run("length_"+string(rune('0'+length)), func(t *testing.T) {
			result, err := RandomDigits(length)
			if err != nil {
				t.Errorf("RandomDigits(%d) error = %v", length, err)
				return
			}
			if len(result) != length {
				t.Errorf("RandomDigits(%d) length = %d, want %d", length, len(result), length)
			}
			// 验证只包含数字
			for _, r := range result {
				if r < '0' || r > '9' {
					t.Errorf("RandomDigits contains invalid character: %c", r)
				}
			}
		})
	}
}

// TestCamelCaseToSnakeCase 测试驼峰转下划线
func TestCamelCaseToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"标准驼峰", "UserName", "user_name"},
		{"连续大写", "HTTPServer", "http_server"},
		{"带数字", "userID", "user_id"},
		{"小写开头", "firstName", "first_name"},
		{"全小写", "username", "username"},
		{"空字符串", "", ""},
		{"单个大写", "A", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CamelCaseToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("CamelCaseToSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSnakeCaseToCamelCase 测试下划线转驼峰
func TestSnakeCaseToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"标准下划线", "user_name", "UserName"},
		{"单个单词", "username", "Username"},
		{"多下划线", "first_name_last", "FirstNameLast"},
		{"带数字", "user_id", "UserId"},
		{"空字符串", "", ""},
		{"下划线开头", "_name", "Name"},
		{"下划线结尾", "name_", "Name"},
		{"连续下划线", "user__name", "UserName"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SnakeCaseToCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("SnakeCaseToCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSnakeCaseToLowerCamelCase 测试下划线转小驼峰
func TestSnakeCaseToLowerCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"标准下划线", "user_name", "userName"},
		{"单个单词", "username", "username"},
		{"多下划线", "first_name_last", "firstNameLast"},
		{"带数字", "user_id", "userId"},
		{"空字符串", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SnakeCaseToLowerCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("SnakeCaseToLowerCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsEmpty 测试字符串是否为空
func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{"空字符串", "", true},
		{"只有空格", "   ", true},
		{"制表符", "\t", true},
		{"换行符", "\n", true},
		{"混合空白", "  \t\n ", true},
		{"非空", "hello", false},
		{"带空格的内容", " hello ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmpty(tt.s)
			if result != tt.expected {
				t.Errorf("IsEmpty(%q) = %v, want %v", tt.s, result, tt.expected)
			}
		})
	}
}

// TestContainsIgnoreCase 测试忽略大小写的包含
func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{"小写包含", "hello world", "hello", true},
		{"大写包含", "HELLO WORLD", "hello", true},
		{"混合包含", "Hello World", "HELLO", true},
		{"不包含", "hello world", "goodbye", false},
		{"空子串", "hello", "", true},
		{"空字符串", "", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsIgnoreCase(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("ContainsIgnoreCase(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

// TestReverse 测试字符串反转
func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"普通字符串", "hello", "olleh"},
		{"带空格", "hello world", "dlrow olleh"},
		{"单个字符", "a", "a"},
		{"空字符串", "", ""},
		{"中文", "你好", "好你"},
		{"混合", "a你b好c", "c好b你a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reverse(tt.input)
			if result != tt.expected {
				t.Errorf("Reverse(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestTruncate 测试字符串截断
func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"需要截断", "hello world", 5, "hello..."},
		{"不需要截断", "hello", 10, "hello"},
		{"刚好相等", "hello", 5, "hello"},
		{"空字符串", "", 5, ""},
		{"中文字符", "你好世界", 3, "你好世..."},
		{"maxLen为0", "hello", 0, "..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Truncate(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

// TestFirstToUpper 测试首字母大写
func TestFirstToUpper(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"小写开头", "hello", "Hello"},
		{"大写开头", "Hello", "Hello"},
		{"空字符串", "", ""},
		{"单个字符", "a", "A"},
		{"中文", "你好", "你好"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FirstToUpper(tt.input)
			if result != tt.expected {
				t.Errorf("FirstToUpper(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestFirstToLower 测试首字母小写
func TestFirstToLower(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"大写开头", "Hello", "hello"},
		{"小写开头", "hello", "hello"},
		{"空字符串", "", ""},
		{"单个字符", "A", "a"},
		{"中文", "你好", "你好"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FirstToLower(tt.input)
			if result != tt.expected {
				t.Errorf("FirstToLower(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsAlpha 测试是否全是字母
func TestIsAlpha(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"全小写", "hello", true},
		{"全大写", "HELLO", true},
		{"混合", "Hello", true},
		{"含数字", "hello123", false},
		{"含空格", "hello world", false},
		{"空字符串", "", false},
		{"中文", "你好", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAlpha(tt.input)
			if result != tt.expected {
				t.Errorf("IsAlpha(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsNumeric 测试是否全是数字
func TestIsNumeric(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"纯数字", "12345", true},
		{"含字母", "123abc", false},
		{"含空格", "123 456", false},
		{"空字符串", "", false},
		{"小数", "123.45", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNumeric(tt.input)
			if result != tt.expected {
				t.Errorf("IsNumeric(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsAlphanumeric 测试是否全是字母或数字
func TestIsAlphanumeric(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"纯字母", "hello", true},
		{"纯数字", "12345", true},
		{"字母数字", "hello123", true},
		{"含空格", "hello 123", false},
		{"含特殊字符", "hello@123", false},
		{"空字符串", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAlphanumeric(tt.input)
			if result != tt.expected {
				t.Errorf("IsAlphanumeric(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// BenchmarkMaskMobile 基准测试
func BenchmarkMaskMobile(b *testing.B) {
	mobile := "13812345678"
	for i := 0; i < b.N; i++ {
		MaskMobile(mobile)
	}
}

// BenchmarkCamelCaseToSnakeCase 基准测试
func BenchmarkCamelCaseToSnakeCase(b *testing.B) {
	s := "UserName"
	for i := 0; i < b.N; i++ {
		CamelCaseToSnakeCase(s)
	}
}
