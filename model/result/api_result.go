package result

import (
	"time"

	"github.com/xiaochangtongxue/my-gin/constant/apicode"
)

type apiResult struct {
	Code    int32
	Message string
	Data    interface{}
	Date    string
}

func Result(code int32, data interface{}) *apiResult {
	return ResultMsg(code, "", data)
}

func ResultMsg(code int32, message string, data interface{}) *apiResult {
	apiMessage := apicode.GetMsg(code)
	if len(message) == 0 && len(apiMessage) != 0 {
		message = apiMessage
	}
	return &apiResult{
		Code:    code,
		Message: message,
		Data:    data,
		Date:    time.Now().Format("2006-01-02 15:04:05"),
	}
}

func Ok(data interface{}) *apiResult {
	return Result(apicode.OK, data)
}

func Fail(code int32) *apiResult {
	return Result(code, nil)
}

func FailMessage(code int32, message string) *apiResult {
	return ResultMsg(code, message, nil)
}

func FailData(code int32, data interface{}) *apiResult {
	return ResultMsg(code, "", data)
}
