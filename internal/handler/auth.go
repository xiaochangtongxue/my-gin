package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/internal/dto/req"
	"github.com/xiaochangtongxue/my-gin/internal/middleware"
	"github.com/xiaochangtongxue/my-gin/internal/service"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"github.com/xiaochangtongxue/my-gin/pkg/response"
	"go.uber.org/zap"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RefreshToken 刷新 Token
// @Summary 刷新 Access Token
// @Description 使用 Refresh Token 获取新的 Access Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body req.RefreshTokenReq true "刷新请求"
// @Success 200 {object} response.Response{data=service.TokenPair}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req req.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 调用 Service 刷新 Token
	tokenPair, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		logger.Warn("刷新 Token 失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Error(err),
		)
		response.Fail(c, response.CodeTokenInvalid, "Refresh Token 无效或已过期")
		return
	}

	logger.Info("Token 刷新成功",
		zap.String("request_id", middleware.GetRequestID(c)),
	)

	response.Success(c, tokenPair)
}

// Logout 登出
// @Summary 用户登出
// @Description 删除 Refresh Token 并将 Access Token 加入黑名单
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body req.RefreshTokenReq true "登出请求"
// @Success 200 {object} response.Response
// @Router /api/v1/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req req.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 提取 Access Token
	accessToken := extractAccessToken(c)

	// 调用 Service 处理登出（删除 RT 并将 AT 加入黑名单）
	if err := h.authService.LogoutWithAccessToken(c.Request.Context(), accessToken, req.RefreshToken); err != nil {
		logger.Warn("登出失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Error(err),
		)
	}

	response.Success(c, nil)
}

// LogoutAll 登出所有设备
// @Summary 登出所有设备
// @Description 删除用户的所有 Refresh Token 并将当前 Access Token 加入黑名单
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/v1/logout/all [post]
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	// 提取 Access Token
	accessToken := extractAccessToken(c)

	// 调用 Service 处理登出（删除所有 RT 并将当前 AT 加入黑名单）
	if err := h.authService.LogoutAllWithAccessToken(c.Request.Context(), accessToken, userID); err != nil {
		logger.Warn("登出所有设备失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
	}

	response.Success(c, nil)
}

// extractAccessToken 从请求中提取 Access Token
func extractAccessToken(c *gin.Context) string {
	// 从 Authorization 请求头获取
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}