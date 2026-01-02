package handler

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwaggerHandler Swagger 文档处理器
type SwaggerHandler struct{}

// NewSwaggerHandler 创建 Swagger 处理器
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// Index Swagger 文档首页
// @Summary Swagger API 文档
// @Description 访问 Swagger API 文档界面
// @Tags 文档
// @Accept json
// @Produce json
// @Router /swagger/index.html [get]
func (h *SwaggerHandler) Index(c *gin.Context) {
	ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
}

// Register 注册 Swagger 路由的便捷方法
func (h *SwaggerHandler) Register(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}