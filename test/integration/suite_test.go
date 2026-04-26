// Package integration 集成测试
// 测试模块间的集成，需要完整的依赖环境（MySQL、Redis等）
//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/xiaochangtongxue/my-gin/internal/model"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/database"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// TestSuite 集成测试套件
type TestSuite struct {
	DB       *gorm.DB
	Cache    cache.Cache
	Config   *config.Config
	Cleanup  func()
}

// SetupSuite 初始化测试套件
func SetupSuite(t *testing.T) *TestSuite {
	// 设置测试环境
	os.Setenv("APP_MODE", "test")
	os.Setenv("APP_JWT_SECRET", "test-jwt-secret-key-for-integration-testing")
	os.Setenv("APP_DATABASE_HOST", getEnv("TEST_DATABASE_HOST", "127.0.0.1"))
	os.Setenv("APP_DATABASE_PORT", getEnv("TEST_DATABASE_PORT", "3306"))
	os.Setenv("APP_DATABASE_DATABASE", getEnv("TEST_DATABASE_NAME", "my_gin_test"))
	os.Setenv("APP_DATABASE_USERNAME", getEnv("TEST_DATABASE_USER", "root"))
	os.Setenv("APP_DATABASE_PASSWORD", getEnv("TEST_DATABASE_PASSWORD", "123456"))
	os.Setenv("APP_REDIS_HOST", getEnv("TEST_REDIS_HOST", "127.0.0.1:6379"))
	os.Setenv("APP_REDIS_PASSWORD", getEnv("TEST_REDIS_PASSWORD", ""))
	os.Setenv("APP_REDIS_DB", "1")

	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	log := logger.New(cfg.Logger)

	// 初始化数据库
	db, err := database.NewDB(cfg.Database, log)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// 初始化缓存
	cacheClient, err := cache.NewCache(cfg.Cache, nil)
	if err != nil {
		t.Fatalf("Failed to connect to cache: %v", err)
	}

	// 自动迁移测试表
	if err := db.AutoMigrate(
		&model.User{},
	); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建清理函数
	cleanup := func() {
		// 清理测试数据
		db.Exec("DELETE FROM users WHERE mobile LIKE '138%'")
	}

	return &TestSuite{
		DB:      db,
		Cache:   cacheClient,
		Config:  cfg,
		Cleanup: cleanup,
	}
}

// TeardownSuite 清理测试套件
func (s *TestSuite) TeardownSuite() {
	if s.Cleanup != nil {
		s.Cleanup()
	}
	if s.DB != nil {
		sqlDB, _ := s.DB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
	if s.Cache != nil {
		s.Cache.Close()
	}
}

// CreateTestUser 创建测试用户
func (s *TestSuite) CreateTestUser(t *testing.T) *model.User {
	user := &model.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Mobile:   fmt.Sprintf("138%08d", time.Now().UnixNano()%100000000),
		Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // bcrypt hash of "password123"
		Status:   1,
	}
	if err := s.DB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

// TruncateUsers 清空用户表
func (s *TestSuite) TruncateUsers(t *testing.T) {
	if err := s.DB.Exec("DELETE FROM users").Error; err != nil {
		t.Fatalf("Failed to truncate users table: %v", err)
	}
}

// WaitForDB 等待数据库就绪
func (s *TestSuite) WaitForDB(maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		sqlDB, _ := s.DB.DB()
		if sqlDB != nil {
			if err = sqlDB.Ping(); err == nil {
				return nil
			}
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("database not ready after %d retries: %w", maxRetries, err)
}

// WaitForRedis 等待 Redis 就绪
func (s *TestSuite) WaitForRedis(maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err := s.Cache.Ping(ctx)
		cancel()
		if err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("redis not ready after %d retries", maxRetries)
}

// getEnv 获取环境变量，带默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestMain 测试主入口
func TestMain(m *testing.M) {
	// 如果不是在运行集成测试，直接跳过
	if os.Getenv("INTEGRATION_TEST") != "1" {
		fmt.Println("Skipping integration tests. Set INTEGRATION_TEST=1 to run.")
		os.Exit(0)
	}
	os.Exit(m.Run())
}