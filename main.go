package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/xiaochangtongxue/my-gin/core"
	"github.com/xiaochangtongxue/my-gin/global"
	"github.com/xiaochangtongxue/my-gin/middleware"
	"github.com/xiaochangtongxue/my-gin/model/result"
)

type person struct {
	Age    int8
	Name   string
	Gender string
}

func main() {
	gin.SetMode(gin.DebugMode)
	global.MGIN_VIP = core.Viper()
	global.MGIN_ZAP = core.Zap()
	r := gin.Default()
	r.Use(middleware.XssHandler(nil))
	r.Use(middleware.ErrorHandler())
	r.GET("/user/:name/:action", func(ctx *gin.Context) {
		name := ctx.Param("name")
		action := ctx.Param("action")
		ctx.String(http.StatusOK, name+"-"+action)
	})
	r.GET("/user", func(ctx *gin.Context) {
		name := ctx.Query("name")
		env := os.Getenv("GO_ENV")
		fmt.Println(env)
		ctx.String(http.StatusOK, name)

	})

	r.POST("/person", func(ctx *gin.Context) {
		// var p1 person
		var p2 person
		// ctx.ShouldBindJSON(&p1)
		err := ctx.ShouldBindJSON(&p2)
		if err != nil {
			ctx.Error(err)
		} else {
			ctx.JSON(http.StatusOK, result.Ok(p2))
		}

	})

	r.GET("/log", func(ctx *gin.Context) {

		for i := 0; i < 100; i++ {
			global.MGIN_ZAP.Info("success", zap.String("test", "test"))
		}

	})
	r.Run(":8000")
}
