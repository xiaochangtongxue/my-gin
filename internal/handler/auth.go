package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaochangtongxue/my-gin/internal/dto/req"
	"github.com/xiaochangtongxue/my-gin/internal/dto/resp"
	"github.com/xiaochangtongxue/my-gin/internal/middleware"
	"github.com/xiaochangtongxue/my-gin/internal/service"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"github.com/xiaochangtongxue/my-gin/pkg/response"
	appvalidator "github.com/xiaochangtongxue/my-gin/pkg/validator"
	"go.uber.org/zap"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 使用手机号、用户名和密码注册新用户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body req.RegisterReq true "注册请求"
// @Success 200 {object} response.Response{data=resp.RegisterResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var r req.RegisterReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.ParamError(c, appvalidator.TranslateError(err))
		return
	}

	// 调用服务注册
	result, err := h.authService.Register(c.Request.Context(), r.Mobile, r.Username, r.Password)
	if err != nil {
		logger.Warn("注册失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.String("mobile", r.Mobile),
			zap.Error(err),
		)
		response.Fail(c, response.CodeInvalidParam, err.Error())
		return
	}

	logger.Info("用户注册成功",
		zap.String("request_id", middleware.GetRequestID(c)),
		zap.String("uid", strconv.FormatUint(result.UID, 10)),
		zap.String("username", result.Username),
	)

	response.Success(c, &resp.RegisterResp{
		UID:      strconv.FormatUint(result.UID, 10),
		Username: result.Username,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用手机号和密码登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body req.LoginReq true "登录请求"
// @Success 200 {object} response.Response{data=resp.LoginResp}
// @Failure 400 {object} response.Response
// @Router /api/v1/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var r req.LoginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.ParamError(c, appvalidator.TranslateError(err))
		return
	}

	// 获取客户端IP
	ip := c.ClientIP()

	// 调用服务登录
	result, err := h.authService.Login(c.Request.Context(), r.Mobile, r.Password, r.CaptchaID, r.CaptchaCode, ip)
	if err != nil {
		logger.Warn("登录失败",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.String("mobile", r.Mobile),
			zap.String("ip", ip),
			zap.Error(err),
		)
		// 根据错误信息返回不同的错误码
		msg := err.Error()
		switch msg {
		case "用户不存在":
			response.Fail(c, response.CodeNotFound, msg)
		case "IP已被封禁，请联系管理员":
			response.Fail(c, response.CodeForbidden, msg)
		case "请输入验证码":
			response.Fail(c, response.CodeCaptchaRequired, msg)
		case "验证码错误或已过期":
			response.Fail(c, response.CodeCaptchaError, msg)
		default:
			if len(msg) > 5 && msg[:5] == "账号已锁定" {
				response.Fail(c, response.CodeAccountLocked, msg)
			} else {
				response.Fail(c, response.CodePasswordError, msg)
			}
		}
		return
	}

	logger.Info("用户登录成功",
		zap.String("request_id", middleware.GetRequestID(c)),
		zap.String("uid", strconv.FormatUint(result.UID, 10)),
		zap.String("username", result.Username),
		zap.String("ip", ip),
	)

	response.Success(c, &resp.LoginResp{
		UID:          strconv.FormatUint(result.UID, 10),
		Username:     result.Username,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	})
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

// extractAccessToken 从请求中提取 Access Token
func extractAccessToken(c *gin.Context) string {
	// 从 Authorization 请求头获取
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
