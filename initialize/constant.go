package initialize

import "github.com/xiaochangtongxue/my-gin/constant/apicode"

//初始化apicode中的code-message组合
func InitApiCode() {
	apicode.CodeMsg[apicode.OK] = apicode.OKMsg
	apicode.CodeMsg[apicode.NoContent] = apicode.NoContentMsg
	apicode.CodeMsg[apicode.BadRequest] = apicode.BadRequestMsg
	apicode.CodeMsg[apicode.UnAuthorized] = apicode.UnAuthorizedMsg
	apicode.CodeMsg[apicode.Forbidden] = apicode.ForbiddenMsg
	apicode.CodeMsg[apicode.NoFound] = apicode.NoFoundMsg
	apicode.CodeMsg[apicode.InternalServerError] = apicode.InternalServerErrorMsg
	apicode.CodeMsg[apicode.MyCode] = apicode.MyCodeMsg
}
