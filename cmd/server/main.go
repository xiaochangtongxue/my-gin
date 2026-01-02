package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	_ "github.com/xiaochangtongxue/my-gin/docs" // Swagger 文档
	"github.com/xiaochangtongxue/my-gin/internal/router"
	"github.com/xiaochangtongxue/my-gin/pkg/database"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"github.com/xiaochangtongxue/my-gin/pkg/validator"
)

// @title           My-Gin API
// @version         1.0
// @description     生产级 Gin 框架脚手架 API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

var (
	configFile = flag.String("c", "configs/config.yaml", "配置文件路径")
)

func main() {
	flag.Parse()

	// 使用 Wire 初始化应用程序（依赖注入）
	app, err := InitializeApp(*configFile)
	if err != nil {
		panic(fmt.Sprintf("应用程序初始化失败: %v", err))
	}

	// 初始化日志（需要在配置初始化后）
	if err := logger.Init(&logger.Config{
		Level:      app.Config.Logger.Level,
		FileName:   app.Config.Logger.FileName,
		MaxSize:    app.Config.Logger.MaxSize,
		MaxBackups: app.Config.Logger.MaxBackups,
		MaxAge:     app.Config.Logger.MaxAge,
		Compress:   app.Config.Logger.Compress,
		Console:    app.Config.Logger.Console,
	}); err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}
	defer logger.Sync()

	// 初始化验证器
	validator.Init()

	logger.Info("服务启动中...",
		zap.String("mode", app.Config.Server.Mode),
		zap.String("host", app.Config.Server.Host),
		zap.Int("port", app.Config.Server.Port),
	)

	// 注册路由
	setupRoutes(app)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", app.Config.Server.Host, app.Config.Server.Port),
		Handler:      app.Engine,
		ReadTimeout:  app.Config.Server.ReadTimeout,
		WriteTimeout: app.Config.Server.WriteTimeout,
	}

	// 数据库迁移
	database.Up()

	// 启动服务器
	go func() {
		logger.Info("服务启动成功")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务启动失败", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("服务正在关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务关闭失败", zap.Error(err))
	}

	// 关闭数据库连接
	if err := database.Close(); err != nil {
		logger.Error("数据库关闭失败", zap.Error(err))
	}

	logger.Info("服务已关闭")
}

// setupRoutes 配置路由
func setupRoutes(app *App) {
	engine := app.Engine

	// 注册各模块路由
	router.RegisterSwaggerRoutes(engine)            // Swagger 文档
	router.RegisterHealthRoutes(engine, app.HealthHandler) // 健康检查
	router.RegisterMetricsRoutes(engine)            // Prometheus Metrics
	router.RegisterAuthRoutes(engine, app.AuthHandler)      // 认证路由
}