package req

// RefreshTokenReq 刷新 Token 请求
type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required" label:"刷新Token"` // Refresh Token
}

// RegisterReq 注册请求
type RegisterReq struct {
	Mobile   string `json:"mobile" binding:"required,mobile" label:"手机号"`     // 手机号
	Username string `json:"username" binding:"required,username" label:"用户名"` // 用户名
	Password string `json:"password" binding:"required,password" label:"密码"`  // 密码
}

// LoginReq 登录请求
type LoginReq struct {
	Mobile      string `json:"mobile" binding:"required,mobile" label:"手机号"`    // 手机号
	Password    string `json:"password" binding:"required,password" label:"密码"` // 密码
	CaptchaID   string `json:"captcha_id,omitempty" label:"验证码ID"`              // 验证码ID（可选）
	CaptchaCode string `json:"captcha_code,omitempty" label:"验证码"`              // 验证码（可选）
}
