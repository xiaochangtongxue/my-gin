package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/xiaochangtongxue/my-gin/middleware"
	"github.com/xiaochangtongxue/my-gin/model/result"
)

type person struct {
	Age    int8
	Name   string
	Gender string
}

func main() {
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
		logrus.WithFields(logrus.Fields{
			"name": "xiaotang",
			"age":  22,
		}).Info("日志信息")
	})

	r.Run(":8000")
}
