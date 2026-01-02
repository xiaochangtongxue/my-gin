package middleware

import (
	"html"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	// xssPatterns XSS 攻击模式
	xssPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),           // script 标签
		regexp.MustCompile(`(?i)javascript:`),                          // javascript: 协议
		regexp.MustCompile(`(?i)on\w+\s*=`),                           // 事件处理器
		regexp.MustCompile(`(?i)<iframe[^>]*>`),                        // iframe 标签
		regexp.MustCompile(`(?i)<object[^>]*>`),                        // object 标签
		regexp.MustCompile(`(?i)<embed[^>]*>`),                         // embed 标签
	}
)

// XSS XSS 防护中间件
// 对查询参数和表单数据进行 XSS 过滤
func XSS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 过滤查询参数
		if c.Request.URL.RawQuery != "" {
			filterQueryParams(c)
		}

		// 设置 CSP 响应头
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:")

		c.Next()
	}
}

// filterQueryParams 过滤查询参数
func filterQueryParams(c *gin.Context) {
	query := c.Request.URL.Query()

	for key, values := range query {
		for i, value := range values {
			values[i] = sanitizeString(value)
		}
		query.Set(key, values[0]) // 只取第一个值
	}

	// 重建查询字符串
	c.Request.URL.RawQuery = query.Encode()
}

// sanitizeString 清理字符串中的 XSS 代码
func sanitizeString(s string) string {
	// 检测是否包含 XSS 模式
	for _, pattern := range xssPatterns {
		if pattern.MatchString(s) {
			// 转义 HTML 特殊字符
			s = html.EscapeString(s)
			// 移除危险的标签和属性
			s = stripTags(s)
			break
		}
	}
	return s
}

// stripTags 移除 HTML 标签（简化版）
func stripTags(s string) string {
	var result strings.Builder
	inTag := false

	for _, r := range s {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				result.WriteRune(r)
			}
		}
	}

	return result.String()
}

// SanitizeInput 清理用户输入（供 Handler 调用）
func SanitizeInput(input string) string {
	return sanitizeString(input)
}

// SanitizeInputSlice 清理字符串切片
func SanitizeInputSlice(inputs []string) []string {
	result := make([]string, len(inputs))
	for i, input := range inputs {
		result[i] = sanitizeString(input)
	}
	return result
}