package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/xiaochangtongxue/my-gin/internal/dto/resp"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(db *gorm.DB, cache cache.Cache) *HealthHandler {
	return &HealthHandler{
		db:    db,
		cache: cache,
	}
}

// Check 简化健康检查
// @Summary 健康检查
// @Description 检查服务运行状态（仅返回基本状态）
// @Tags 健康
// @Accept json
// @Produce json
// @Success 200 {object} resp.CheckResponse
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, resp.CheckResponse{
		Status: "ok",
		Time:   time.Now().Format(time.RFC3339),
	})
}

// Live 存活检查
// @Summary 存活检查
// @Description 检查服务是否存活（轻量级，不检查依赖）
// @Tags 健康
// @Accept json
// @Produce json
// @Success 200 {object} resp.CheckResponse
// @Failure 503 {object} resp.CheckResponse
// @Router /health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, resp.CheckResponse{
		Status: "ok",
		Time:   time.Now().Format(time.RFC3339),
	})
}

// Ready 就绪检查
// @Summary 就绪检查
// @Description 检查服务及其依赖（数据库、Redis）的就绪状态
// @Tags 健康
// @Accept json
// @Produce json
// @Success 200 {object} resp.ReadyResponse
// @Failure 503 {object} resp.ReadyResponse
// @Router /health/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	checks := make(map[string]resp.Check)
	overallStatus := "healthy"

	// 检查数据库
	dbCheck := h.checkDatabase(ctx)
	checks["database"] = dbCheck
	if dbCheck.Status != "up" {
		overallStatus = "unhealthy"
	}

	// 检查 Redis
	cacheCheck := h.checkCache(ctx)
	checks["redis"] = cacheCheck
	if cacheCheck.Status != "up" && overallStatus != "unhealthy" {
		overallStatus = "degraded"
	}

	// 根据整体状态返回相应的 HTTP 状态码
	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, resp.ReadyResponse{
		Status:    overallStatus,
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
		Checks:    checks,
	})
}

// checkDatabase 检查数据库连接状态
func (h *HealthHandler) checkDatabase(ctx context.Context) resp.Check {
	if h.db == nil {
		return resp.Check{Status: "down", Error: "database not initialized"}
	}

	start := time.Now()

	// 使用 Ping 检查连接
	sqlDB, err := h.db.DB()
	if err != nil {
		return resp.Check{Status: "down", Error: err.Error(), Latency: time.Since(start).Milliseconds()}
	}

	err = sqlDB.PingContext(ctx)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return resp.Check{Status: "down", Error: err.Error(), Latency: latency}
	}

	return resp.Check{Status: "up", Latency: latency}
}

// checkCache 检查 Redis 连接状态
func (h *HealthHandler) checkCache(ctx context.Context) resp.Check {
	if h.cache == nil {
		return resp.Check{Status: "down", Error: "cache not initialized"}
	}

	start := time.Now()

	// 使用简单的 Set/Get 操作检查连接
	testKey := "health:check"
	testValue := "1"

	// 尝试设置一个测试键
	err := h.cache.Set(ctx, testKey, testValue, 5*time.Second)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return resp.Check{Status: "down", Error: fmt.Sprintf("cache set failed: %v", err), Latency: latency}
	}

	// 尝试获取测试键
	val, err := h.cache.Get(ctx, testKey)
	if err != nil {
		return resp.Check{Status: "down", Error: fmt.Sprintf("cache get failed: %v", err), Latency: latency}
	}

	if val != testValue {
		return resp.Check{Status: "degraded", Error: "cache value mismatch", Latency: latency}
	}

	return resp.Check{Status: "up", Latency: latency}
}