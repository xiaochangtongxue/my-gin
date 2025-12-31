package response

// 业务状态码定义
const (
	CodeSuccess      = 0     // 成功
	CodeInvalidParam = 10001 // 无效参数
	CodeUnauthorized = 10002 // 未授权
	CodeForbidden    = 10003 // 禁止访问
	CodeNotFound     = 10004 // 资源不存在
	CodeServerError  = 10005 // 服务器错误
	CodeDBError      = 10006 // 数据库错误
	CodeRedisError   = 10007 // Redis错误
	CodeTokenExpired = 10008 // Token过期
	CodeTokenInvalid = 10009 // Token无效
)

// 错误码映射
var codeMessages = map[int]string{
	CodeSuccess:      "success",
	CodeInvalidParam: "参数错误",
	CodeUnauthorized: "未授权，请先登录",
	CodeForbidden:    "无权访问",
	CodeNotFound:     "资源不存在",
	CodeServerError:  "服务器内部错误",
	CodeDBError:      "数据库操作失败",
	CodeRedisError:   "缓存操作失败",
	CodeTokenExpired: "登录已过期，请重新登录",
	CodeTokenInvalid: "登录状态无效",
}

// GetMessage 获取错误消息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// SetMessage 自定义错误消息
func SetMessage(code int, message string) {
	codeMessages[code] = message
}
