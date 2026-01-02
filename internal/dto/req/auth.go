package req

// RefreshTokenReq 刷新 Token 请求
type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"` // Refresh Token
}
