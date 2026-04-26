package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"
	"strings"
)

// BusinessError 业务错误
type BusinessError struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误消息
	Err     error  `json:"-"`       // 原始错误
	File    string `json:"-"`       // 文件名
	Line    int    `json:"-"`       // 行号
}

// Error 实现error接口
func (e *BusinessError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持错误包装
func (e *BusinessError) Unwrap() error {
	return e.Err
}

// New 创建业务错误
func New(code int, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// AsBusinessError extracts a BusinessError from an error chain.
func AsBusinessError(err error) (*BusinessError, bool) {
	var businessErr *BusinessError
	if stderrors.As(err, &businessErr) {
		return businessErr, true
	}
	return nil, false
}

// Wrap 包装错误
func Wrap(err error, code int, message string) *BusinessError {
	if err == nil {
		return nil
	}
	return &BusinessError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Wrapf 包装错误（格式化消息）
func Wrapf(err error, code int, format string, args ...interface{}) *BusinessError {
	if err == nil {
		return nil
	}
	return &BusinessError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Err:     err,
	}
}

// WithCaller 添加调用者信息
func (e *BusinessError) WithCaller() *BusinessError {
	if pc, _, line, ok := runtime.Caller(1); ok {
		// 获取函数名
		if fn := runtime.FuncForPC(pc); fn != nil {
			parts := strings.Split(fn.Name(), ".")
			e.File = parts[len(parts)-1] + "()"
		}
		e.Line = line
	}
	return e
}

// Stack 获取堆栈信息
func (e *BusinessError) Stack() string {
	if e.File == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", e.File, e.Line)
}

// 预定义常用错误
var (
	ErrInvalidParam = &BusinessError{Code: 10001, Message: "参数错误"}
	ErrUnauthorized = &BusinessError{Code: 10002, Message: "未授权"}
	ErrForbidden    = &BusinessError{Code: 10003, Message: "禁止访问"}
	ErrNotFound     = &BusinessError{Code: 10004, Message: "资源不存在"}
	ErrServerError  = &BusinessError{Code: 10005, Message: "服务器内部错误"}
	ErrDBError      = &BusinessError{Code: 10006, Message: "数据库操作失败"}
	ErrTokenExpired = &BusinessError{Code: 10008, Message: "登录已过期"}
	ErrTokenInvalid = &BusinessError{Code: 10009, Message: "登录状态无效"}
)
