package myerrors

import (
	"github.com/xiaochangtongxue/my-gin/response"
)

type BusinessError struct {
	ErrorCode int
	Message   string
	Data      interface{}
}

// 实现接口
func (e *BusinessError) Error() string {
	return e.Message
}

func NewBusinessError(apiCode *response.ApiCode) *BusinessError {
	return &BusinessError{
		ErrorCode: apiCode.Code,
		Message:   apiCode.Message,
	}
}

func GetBusinessError(code int, msg string, data interface{}) *BusinessError {
	return &BusinessError{
		ErrorCode: code,
		Message:   msg,
		Data:      data,
	}
}
