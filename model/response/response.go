package response

import (
	"time"

	"github.com/xiaochangtongxue/my-gin/constant/apicode"
)

type apiResponse struct {
	Code    int
	Message string
	Data    interface{}
	Date    string
}

func Result(code int, data interface{}) *apiResponse {
	return ResultMsg(code, "", data)
}

func ResultMsg(code int, message string, data interface{}) *apiResponse {
	apiMessage := apicode.GetMsg(code)
	if len(message) == 0 && len(apiMessage) != 0 {
		message = apiMessage
	}
	return &apiResponse{
		Code:    code,
		Message: message,
		Data:    data,
		Date:    time.Now().Format("2006-01-02 15:04:05"),
	}
}

func Ok(data interface{}) *apiResponse {
	return Result(apicode.OK, data)
}

func Fail(code int) *apiResponse {
	return Result(code, nil)
}

func FailMessage(code int, message string) *apiResponse {
	return ResultMsg(code, message, nil)
}

func FailData(code int, data interface{}) *apiResponse {
	return ResultMsg(code, "", data)
}
