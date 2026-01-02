package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 标准响应结构
// @description 统一 API 响应格式
type Response struct {
	Code      int         `json:"code" example:"0"`                        // 业务状态码
	Message   string      `json:"message" example:"success"`               // 响应消息
	Data      interface{} `json:"data,omitempty"`                          // 响应数据
	Timestamp int64       `json:"timestamp" example:"1704067200"`          // 时间戳
	RequestID string      `json:"request_id,omitempty" example:"abc123"`  // 请求ID
}

// PageData 分页数据结构
type PageData struct {
	Items     interface{} `json:"items"`      // 数据列表
	Total     int64       `json:"total"`      // 总数
	Page      int         `json:"page"`       // 当前页
	PageSize  int         `json:"page_size"`  // 每页数量
	TotalPage int         `json:"total_page"` // 总页数
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      CodeSuccess,
		Message:   GetMessage(CodeSuccess),
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, message string) {
	if message == "" {
		message = GetMessage(code)
	}
	statusCode := getStatusCode(code)
	c.JSON(statusCode, Response{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
	})
}

// FailWithData 失败响应（带数据）
func FailWithData(c *gin.Context, code int, message string, data interface{}) {
	if message == "" {
		message = GetMessage(code)
	}
	statusCode := getStatusCode(code)
	c.JSON(statusCode, Response{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	totalPage := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPage++
	}
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: GetMessage(CodeSuccess),
		Data: PageData{
			Items:     items,
			Total:     total,
			Page:      page,
			PageSize:  pageSize,
			TotalPage: totalPage,
		},
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
	})
}

// ParamError 参数错误响应
func ParamError(c *gin.Context, message string) {
	Fail(c, CodeInvalidParam, message)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = GetMessage(CodeUnauthorized)
	}
	Fail(c, CodeUnauthorized, message)
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = GetMessage(CodeForbidden)
	}
	Fail(c, CodeForbidden, message)
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = GetMessage(CodeNotFound)
	}
	Fail(c, CodeNotFound, message)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, message string) {
	if message == "" {
		message = GetMessage(CodeServerError)
	}
	Fail(c, CodeServerError, message)
}

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("RequestID"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// getStatusCode 根据业务状态码获取HTTP状态码
func getStatusCode(code int) int {
	switch code {
	case CodeSuccess:
		return http.StatusOK
	case CodeInvalidParam:
		return http.StatusBadRequest
	case CodeUnauthorized, CodeTokenExpired, CodeTokenInvalid:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
