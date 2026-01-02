package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/internal/handler"
)

// RegisterSwaggerRoutes 注册 Swagger 文档路由
func RegisterSwaggerRoutes(r *gin.Engine) {
	h := handler.NewSwaggerHandler()
	h.Register(r)
}