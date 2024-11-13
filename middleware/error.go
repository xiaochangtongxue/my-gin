package middleware

import (
	"github.com/xiaochangtongxue/my-gin/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/myerrors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			e := ctx.Errors[0]
			err := e.Err
			if _, ok := err.(*myerrors.BusinessError); ok {
				ctx.JSON(http.StatusOK, response.Fail(response.MyCode))
			} else {
				ctx.JSON(http.StatusOK, response.FailMessage(response.InternalServerError, err.Error()))
			}
		}

	}
}
