package resp

// RegisterResp 注册响应
type RegisterResp struct {
	UID      string `json:"uid"`      // 对外用户ID
	Username string `json:"username"` // 用户名
}

// LoginResp 登录响应
type LoginResp struct {
	UID          string `json:"uid"`           // 对外用户ID
	Username     string `json:"username"`      // 用户名
	AccessToken  string `json:"access_token"`  // Access Token
	RefreshToken string `json:"refresh_token"` // Refresh Token
	ExpiresIn    int64  `json:"expires_in"`    // Access Token 过期时间（秒）
}

// TokenPair 双 Token 响应
// @Description 双 Token 响应（Access Token + Refresh Token）
type TokenPair struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Access Token
	RefreshToken string `json:"refresh_token" example:"a1b2c3d4e5f6..."`                        // Refresh Token
	ExpiresIn    int64  `json:"expires_in" example:"1800"`                                      // Access Token 过期时间（秒）
}
