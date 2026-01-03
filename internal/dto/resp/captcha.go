package resp

// CaptchaResp 验证码响应
type CaptchaResp struct {
	CaptchaID    string `json:"captcha_id"`    // 验证码ID
	CaptchaImage string `json:"captcha_image"` // base64图片
}
