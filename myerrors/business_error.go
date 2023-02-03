package myerrors

import "github.com/xiaochangtongxue/my-gin/constant/apicode"

type BusinessError struct {
	ErrorCode int32
	Message   string
	Data      interface{}
}

// 实现接口
func (e *BusinessError) Error() string {
	return e.Message
}

func NewBusinessError(code int32) *BusinessError {
	return &BusinessError{
		ErrorCode: code,
		Message:   apicode.GetMsg(code),
	}
}

func GetBusinessError(code int32, msg string, data interface{}) *BusinessError {
	return &BusinessError{
		ErrorCode: code,
		Message:   msg,
		Data:      data,
	}
}
