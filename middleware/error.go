package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/constant/apicode"
	"github.com/xiaochangtongxue/my-gin/model/response"
	"github.com/xiaochangtongxue/my-gin/myerrors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			e := ctx.Errors[0]
			err := e.Err
			if myErr, ok := err.(*myerrors.BusinessError); ok {
				ctx.JSON(http.StatusOK, response.Fail(myErr.ErrorCode))
			} else {
				ctx.JSON(http.StatusOK, response.FailMessage(apicode.InternalServerError, err.Error()))
			}
		}

	}
}
