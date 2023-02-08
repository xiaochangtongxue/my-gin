package apicode

var codeMsg = make(map[int]string)

const (
	OK                     = 200
	OKMsg                  = "操作成功"
	NoContent              = 204
	NoContentMsg           = "请求成功，没有资源"
	BadRequest             = 400
	BadRequestMsg          = "无法解析客户端请求"
	UnAuthorized           = 401
	UnAuthorizedMsg        = "非法访问"
	Forbidden              = 403
	ForbiddenMsg           = "不允许访问"
	NoFound                = 404
	NoFoundMsg             = "请求的资源不存在"
	InternalServerError    = 500
	InternalServerErrorMsg = "操作失败"
)

const (
	MyCode    = 1000
	MyCodeMsg = "自定义的错误"
)

func GetMsg(code int) string {
	return codeMsg[code]
}

func init() {
	codeMsg[OK] = OKMsg
	codeMsg[NoContent] = NoContentMsg
	codeMsg[BadRequest] = BadRequestMsg
	codeMsg[UnAuthorized] = UnAuthorizedMsg
	codeMsg[Forbidden] = ForbiddenMsg
	codeMsg[NoFound] = NoFoundMsg
	codeMsg[InternalServerError] = InternalServerErrorMsg
	codeMsg[MyCode] = MyCodeMsg
}
