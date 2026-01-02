package resp

// TokenPair 双 Token 响应
// @Description 双 Token 响应（Access Token + Refresh Token）
type TokenPair struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`  // Access Token
	RefreshToken string `json:"refresh_token" example:"a1b2c3d4e5f6..."`                         // Refresh Token
	ExpiresIn    int64  `json:"expires_in" example:"1800"`                                       // Access Token 过期时间（秒）
}