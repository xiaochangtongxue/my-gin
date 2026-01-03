package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/internal/middleware"
	"github.com/xiaochangtongxue/my-gin/internal/service"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"github.com/xiaochangtongxue/my-gin/pkg/response"
	"go.uber.org/zap"
)

// SecurityHandler 安全处理器
type SecurityHandler struct {
	securityService service.SecurityService
}

// NewSecurityHandler 创建安全处理器
func NewSecurityHandler(securityService service.SecurityService) *SecurityHandler {
	return &SecurityHandler{
		securityService: securityService,
	}
}

// UnlockAccount 解锁账号
// @Summary 解锁用户账号
// @Description 管理员解锁被锁定的用户账号
// @Tags 安全管理
// @Param uid query string true "用户UID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/unlock-account [post]
func (h *SecurityHandler) UnlockAccount(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		response.ParamError(c, "用户UID不能为空")
		return
	}

	if err := h.securityService.UnlockAccount(c.Request.Context(), uid); err != nil {
		logger.Warn("解锁账号失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.String("uid", uid),
			zap.Error(err),
		)
		response.ServerError(c, "解锁账号失败")
		return
	}

	logger.Info("账号解锁成功",
		zap.String("request_id", middleware.GetRequestID(c)),
		zap.String("uid", uid),
	)

	response.Success(c, nil)
}
