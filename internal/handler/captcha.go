package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/internal/dto/resp"
	"github.com/xiaochangtongxue/my-gin/internal/middleware"
	"github.com/xiaochangtongxue/my-gin/pkg/captcha"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"github.com/xiaochangtongxue/my-gin/pkg/response"
	"go.uber.org/zap"
)

// CaptchaHandler 验证码处理器
type CaptchaHandler struct {
	captcha *captcha.Captcha
}

// NewCaptchaHandler 创建验证码处理器
func NewCaptchaHandler(captcha *captcha.Captcha) *CaptchaHandler {
	return &CaptchaHandler{
		captcha: captcha,
	}
}

// GetCaptcha 获取验证码
// @Summary 获取验证码
// @Description 生成图形验证码
// @Tags 验证码
// @Produce json
// @Success 200 {object} response.Response{data=resp.CaptchaResp}
// @Router /api/v1/captcha [get]
func (h *CaptchaHandler) GetCaptcha(c *gin.Context) {
	captchaID, imageBase64, err := h.captcha.Generate(c.Request.Context())
	if err != nil {
		logger.Warn("生成验证码失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Error(err),
		)
		response.ServerError(c, "生成验证码失败")
		return
	}

	response.Success(c, &resp.CaptchaResp{
		CaptchaID:    captchaID,
		CaptchaImage: imageBase64,
	})
}
