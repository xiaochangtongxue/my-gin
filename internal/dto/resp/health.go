package resp

// CheckResponse 健康检查响应
// @description 简化版健康检查响应
type CheckResponse struct {
	Status string `json:"status" example:"ok"`
	Time   string `json:"time" example:"2025-01-02T15:04:05Z07:00"`
}

// ReadyResponse 就绪检查响应
// @description 完整版健康检查响应（包含各组件状态）
type ReadyResponse struct {
	Status    string            `json:"status" example:"healthy"`
	Timestamp string            `json:"timestamp" example:"2025-01-02T15:04:05Z07:00"`
	Version   string            `json:"version" example:"1.0.0"`
	Checks    map[string]Check  `json:"checks"`
}

// Check 单个组件检查结果
// @description 组件健康状态
type Check struct {
	Status  string `json:"status" example:"up"`  // up, down, degraded
	Error   string `json:"error,omitempty"`     // 错误信息
	Latency int64  `json:"latency" example:"5"`  // 响应耗时（毫秒）
}