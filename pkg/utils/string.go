// Package utils 通用工具包
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"strings"
	"unicode"
)

// MaskMobile 手机号脱敏（保留前3后4位）
// 例: 13812345678 -> 138****5678
func MaskMobile(mobile string) string {
	if len(mobile) != 11 {
		return mobile
	}
	return mobile[:3] + "****" + mobile[7:]
}

// MaskEmail 邮箱脱敏（保留前2和@后域名）
// 例: example@gmail.com -> ex****@gmail.com
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) < 3 {
		return email
	}
	username := parts[0]
	maskedUsername := username[:2] + "****"
	return maskedUsername + "@" + parts[1]
}

// MaskIDCard 身份证号脱敏（保留前6后4位）
// 例: 110101199001011234 -> 110101********1234
func MaskIDCard(idCard string) string {
	length := len(idCard)
	if length < 10 {
		return idCard
	}
	keepStart := 6
	keepEnd := 4
	if length <= keepStart+keepEnd {
		keepStart = length / 2
		keepEnd = length - keepStart
	}
	return idCard[:keepStart] + strings.Repeat("*", length-keepStart-keepEnd) + idCard[length-keepEnd:]
}

// MaskBankCard 银行卡号脱敏（保留前4后4位）
// 例: 6222021234567890123 -> 6222************0123
func MaskBankCard(cardNo string) string {
	length := len(cardNo)
	if length < 8 {
		return cardNo
	}
	keepStart := 4
	keepEnd := 4
	return cardNo[:keepStart] + strings.Repeat("*", length-keepStart-keepEnd) + cardNo[length-keepEnd:]
}

// MaskString 通用字符串脱敏
// start: 保留开头字符数
// end: 保留结尾字符数
// maskChar: 掩码字符（默认 *）
func MaskString(s string, start, end int, maskChar ...string) string {
	length := len(s)
	if length <= start+end {
		return s
	}
	mask := "*"
	if len(maskChar) > 0 && maskChar[0] != "" {
		mask = maskChar[0]
	}
	return s[:start] + strings.Repeat(mask, length-start-end) + s[length-end:]
}

// RandomString 生成随机字符串（字母+数字）
// length: 字符串长度
func RandomString(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

// RandomStringFromCharset 从指定字符集生成随机字符串
func RandomStringFromCharset(length int, charset string) (string, error) {
	if length <= 0 {
		return "", nil
	}
	if charset == "" {
		return "", nil
	}
	b := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

// RandomHex 生成随机十六进制字符串
// length: 字节长度（结果字符串长度为 length * 2）
func RandomHex(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// RandomDigits 生成随机数字字符串
func RandomDigits(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}
	const charset = "0123456789"
	return RandomStringFromCharset(length, charset)
}

// CamelCaseToSnakeCase 驼峰转下划线
// 例: UserName -> user_name, userID -> user_id
func CamelCaseToSnakeCase(s string) string {
	runes := []rune(s)
	var result []rune
	for i, r := range runes {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := runes[i-1]
			var next rune
			if i+1 < len(runes) {
				next = runes[i+1]
			}
			if (prev >= 'a' && prev <= 'z') ||
				(prev >= '0' && prev <= '9') ||
				(prev >= 'A' && prev <= 'Z' && next >= 'a' && next <= 'z') {
				result = append(result, '_')
			}
		}
		if r >= 'A' && r <= 'Z' {
			result = append(result, r+32) // 转小写
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// SnakeCaseToCamelCase 下划线转驼峰
// 例: user_name -> UserName, user_id -> UserID
func SnakeCaseToCamelCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		if word == "" {
			continue
		}
		// 首字母大写
		runes := []rune(word)
		if runes[0] >= 'a' && runes[0] <= 'z' {
			runes[0] = runes[0] - 32
		}
		words[i] = string(runes)
	}
	return strings.Join(words, "")
}

// SnakeCaseToLowerCamelCase 下划线转小驼峰
// 例: user_name -> userName, user_id -> userId
func SnakeCaseToLowerCamelCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		if word == "" {
			continue
		}
		// 第一个单词保持小写，其余首字母大写
		if i == 0 {
			words[i] = word
		} else {
			runes := []rune(word)
			if runes[0] >= 'a' && runes[0] <= 'z' {
				runes[0] = runes[0] - 32
			}
			words[i] = string(runes)
		}
	}
	return strings.Join(words, "")
}

// IsEmpty 判断字符串是否为空
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// TrimSpace 去除首尾空格
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// Contains 判断字符串是否包含子串（不区分大小写）
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Reverse 反转字符串
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Truncate 截断字符串到指定长度
func Truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

// ToUpper 首字母大写
func FirstToUpper(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	if runes[0] >= 'a' && runes[0] <= 'z' {
		runes[0] = runes[0] - 32
	}
	return string(runes)
}

// ToLower 首字母小写
func FirstToLower(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	if runes[0] >= 'A' && runes[0] <= 'Z' {
		runes[0] = runes[0] + 32
	}
	return string(runes)
}

// IsAlpha 判断是否全是字母
func IsAlpha(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return false
		}
	}
	return len(s) > 0
}

// IsNumeric 判断是否全是数字
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}

// IsAlphanumeric 判断是否全是字母或数字
func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || unicode.IsDigit(r)) {
			return false
		}
	}
	return len(s) > 0
}
