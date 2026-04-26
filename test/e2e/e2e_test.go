// Package e2e 端到端测试
// 测试完整的 HTTP 请求流程
//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/xiaochangtongxue/my-gin/cmd/server"
	"github.com/xiaochangtongxue/my-gin/internal/router"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/database"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
)

// E2ETestSuite E2E 测试套件
type E2ETestSuite struct {
	Server   *httptest.Server
	BaseURL  string
	Config   *config.Config
	App      *server.App
}

// SetupE2ESetup 初始化 E2E 测试
func SetupE2ESetup(t *testing.T) *E2ETestSuite {
	// 设置测试环境
	os.Setenv("APP_MODE", "test")
	os.Setenv("APP_JWT_SECRET", "test-jwt-secret-key-for-e2e-testing")
	os.Setenv("APP_DATABASE_HOST", "127.0.0.1")
	os.Setenv("APP_DATABASE_PORT", "3306")
	os.Setenv("APP_DATABASE_DATABASE", "my_gin_e2e")
	os.Setenv("APP_DATABASE_USERNAME", "root")
	os.Setenv("APP_DATABASE_PASSWORD", "123456")
	os.Setenv("APP_REDIS_HOST", "127.0.0.1:6379")
	os.Setenv("APP_REDIS_PASSWORD", "")
	os.Setenv("APP_REDIS_DB", "2")

	// 加载配置
	cfg, err := config.Init("configs/config.yaml")
	require.NoError(t, err)

	// 初始化日志
	if err := logger.Init(&cfg.Logger); err != nil {
		require.NoError(t, err)
	}

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		require.NoError(t, err)
	}

	// 初始化缓存
	if err := cache.InitRedis(&cfg.Redis); err != nil {
		require.NoError(t, err)
	}

	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 初始化应用
	app, err := server.InitializeApp("configs/config.yaml")
	require.NoError(t, err)

	// 注册路由
	engine := app.Engine
	router.RegisterSwaggerRoutes(engine)
	router.RegisterHealthRoutes(engine, app.HealthHandler)
	router.RegisterMetricsRoutes(engine)
	router.RegisterAuthRoutes(engine, app.AuthHandler, app.CaptchaHandler)

	// 创建测试服务器
	testServer := httptest.NewServer(engine)

	return &E2ETestSuite{
		Server:  testServer,
		BaseURL: testServer.URL,
		Config:  cfg,
		App:      app,
	}
}

// TeardownE2E 清理 E2E 测试
func (s *E2ETestSuite) TeardownE2E() {
	if s.Server != nil {
		s.Server.Close()
	}
	database.Close()
	cache.CloseRedis()
}

// Request 发送 HTTP 请求
func (s *E2ETestSuite) Request(method, path string, body interface{}, headers map[string]string) (*http.Response, []byte) {
	var reqBody io.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, _ := http.NewRequest(method, s.BaseURL+path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return resp, respBody
}

// Get 发送 GET 请求
func (s *E2ETestSuite) Get(path string, headers map[string]string) (*http.Response, []byte) {
	return s.Request(http.MethodGet, path, nil, headers)
}

// Post 发送 POST 请求
func (s *E2ETestSuite) Post(path string, body interface{}, headers map[string]string) (*http.Response, []byte) {
	return s.Request(http.MethodPost, path, body, headers)
}

// Delete 发送 DELETE 请求
func (s *E2ETestSuite) Delete(path string, headers map[string]string) (*http.Response, []byte) {
	return s.Request(http.MethodDelete, path, nil, headers)
}

// APIResponse API 响应结构
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// TestHealthCheck 测试健康检查接口
func TestHealthCheck(t *testing.T) {
	suite := SetupE2ESetup(t)
	defer suite.TeardownE2E()

	resp, body := suite.Get("/health", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

// TestPing 测试 Ping 接口
func TestPing(t *testing.T) {
	suite := SetupE2ESetup(t)
	defer suite.TeardownE2E()

	resp, body := suite.Get("/api/v1/ping", nil)

	var apiResp APIResponse
	err := json.Unmarshal(body, &apiResp)
	assert.NoError(t, err)
	assert.Equal(t, 0, apiResp.Code)
	assert.Equal(t, "pong", apiResp.Data)
}

// TestRegister 测试用户注册
func TestRegister(t *testing.T) {
	suite := SetupE2ESetup(t)
	defer suite.TeardownE2E()

	timestamp := time.Now().UnixNano()
	mobile := fmt.Sprintf("138%08d", timestamp%100000000)

	t.Run("ValidRegistration", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"username": "testuser",
			"mobile":   mobile,
			"password": "Password123",
		}
		resp, body := suite.Post("/api/v1/register", reqBody, nil)

		var apiResp APIResponse
		err := json.Unmarshal(body, &apiResp)
		assert.NoError(t, err)
		assert.Equal(t, 0, apiResp.Code)
	})

	t.Run("DuplicateMobile", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"username": "testuser2",
			"mobile":   mobile,
			"password": "Password123",
		}
		resp, body := suite.Post("/api/v1/register", reqBody, nil)

		var apiResp APIResponse
		err := json.Unmarshal(body, &apiResp)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, apiResp.Code)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"username": "testuser3",
			"mobile":   fmt.Sprintf("139%08d", time.Now().UnixNano()%100000000),
			"password": "123",
		}
		resp, body := suite.Post("/api/v1/register", reqBody, nil)

		var apiResp APIResponse
		err := json.Unmarshal(body, &apiResp)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, apiResp.Code)
	})
}

// TestLogin 测试用户登录
func TestLogin(t *testing.T) {
	suite := SetupE2ESetup(t)
	defer suite.TeardownE2E()

	timestamp := time.Now().UnixNano()
	mobile := fmt.Sprintf("138%08d", timestamp%100000000)

	// 先注册用户
	registerReq := map[string]interface{}{
		"username": "logintest",
		"mobile":   mobile,
		"password": "Password123",
	}
	suite.Post("/api/v1/register", registerReq, nil)

	t.Run("ValidLogin", func(t *testing.T) {
		loginReq := map[string]interface{}{
			"mobile":   mobile,
			"password": "Password123",
		}
		resp, body := suite.Post("/api/v1/login", loginReq, nil)

		var apiResp struct {
			Code    int `json:"code"`
			Message string `json:"message"`
			Data    struct {
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token"`
				ExpiresIn    int64  `json:"expires_in"`
			} `json:"data"`
		}
		err := json.Unmarshal(body, &apiResp)
		assert.NoError(t, err)
		assert.Equal(t, 0, apiResp.Code)
		assert.NotEmpty(t, apiResp.Data.AccessToken)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		loginReq := map[string]interface{}{
			"mobile":   mobile,
			"password": "WrongPassword",
		}
		resp, body := suite.Post("/api/v1/login", loginReq, nil)

		var apiResp APIResponse
		err := json.Unmarshal(body, &apiResp)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, apiResp.Code)
	})
}

// TestProtectedRoute 测试需要认证的路由
func TestProtectedRoute(t *testing.T) {
	suite := SetupE2ESetup(t)
	defer suite.TeardownE2E()

	// 注册并登录获取 token
	timestamp := time.Now().UnixNano()
	mobile := fmt.Sprintf("138%08d", timestamp%100000000)

	registerReq := map[string]interface{}{
		"username": "authtest",
		"mobile":   mobile,
		"password": "Password123",
	}
	suite.Post("/api/v1/register", registerReq, nil)

	loginReq := map[string]interface{}{
		"mobile":   mobile,
		"password": "Password123",
	}
	_, loginBody := suite.Post("/api/v1/login", loginReq, nil)

	var loginResp struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	json.Unmarshal(loginBody, &loginResp)

	t.Run("WithoutToken", func(t *testing.T) {
		resp, _ := suite.Get("/api/v1/user/info", nil)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("WithValidToken", func(t *testing.T) {
		headers := map[string]string{
			"Authorization": "Bearer " + loginResp.Data.AccessToken,
		}
		resp, body := suite.Get("/api/v1/user/info", headers)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err := json.Unmarshal(body, &apiResp)
		assert.NoError(t, err)
		assert.Equal(t, 0, apiResp.Code)
	})
}

// TestMain 测试主入口
func TestMain(m *testing.M) {
	// 如果不是在运行 E2E 测试，直接跳过
	if os.Getenv("E2E_TEST") != "1" {
		fmt.Println("Skipping E2E tests. Set E2E_TEST=1 to run.")
		os.Exit(0)
	}
	os.Exit(m.Run())
}