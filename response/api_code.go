package response

type ApiCode struct {
	Code    int    `json:"code" yaml:"code"`
	Message string `json:"msg" yaml:"msg"`
}

func apiCode(code int, msg string) *ApiCode {
	return &ApiCode{
		Code:    code,
		Message: msg,
	}
}

var (
	OK                  = apiCode(200, "操作成功")
	NoContent           = apiCode(204, "请求成功，没有资源")
	BadRequest          = apiCode(400, "无法解析客户端请求")
	UnAuthorized        = apiCode(401, "非法访问")
	Forbidden           = apiCode(403, "不允许访问")
	NoFound             = apiCode(404, "请求资源不存在")
	InternalServerError = apiCode(500, "请求失败")
	MyCode              = apiCode(1000, "自定义的错误")
)
