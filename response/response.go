package response

import (
	"time"
)

type ApiResponse struct {
	Code    int
	Message string
	Data    interface{}
	Date    string
}

func Result(apiCode *ApiCode, data interface{}) *ApiResponse {
	return ResultMsg(apiCode, data)
}

func ResultMsg(apiCode *ApiCode, data interface{}) *ApiResponse {
	return &ApiResponse{
		Code:    apiCode.Code,
		Message: apiCode.Message,
		Data:    data,
		Date:    time.Now().Format("2006-01-02 15:04:05"),
	}
}

func Ok(data interface{}) *ApiResponse {
	return Result(OK, data)
}

func Fail(apiCode *ApiCode) *ApiResponse {
	return Result(apiCode, nil)
}

func FailMessage(apiCode *ApiCode, message string) *ApiResponse {
	apiCode.Message = message
	return ResultMsg(apiCode, nil)
}

func FailData(apiCode *ApiCode, data interface{}) *ApiResponse {
	return ResultMsg(apiCode, data)
}
