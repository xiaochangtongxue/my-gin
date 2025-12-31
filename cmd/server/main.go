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

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/xiaochangtongxue/my-gin/pkg/cache"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
	"github.com/xiaochangtongxue/my-gin/pkg/database"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
	"github.com/xiaochangtongxue/my-gin/pkg/response"
	"github.com/xiaochangtongxue/my-gin/pkg/validator"
)

var (
	configFile = flag.String("c", "configs/config.yaml", "配置文件路径")
)

func main() {
	flag.Parse()

	// 初始化配置
	cfg, err := config.Init(*configFile)
	if err != nil {
		panic(fmt.Sprintf("配置初始化失败: %v", err))
	}

	// 初始化日志
	if err := logger.Init(&logger.Config{
		Level:      cfg.Logger.Level,
		FileName:   cfg.Logger.FileName,
		MaxSize:    cfg.Logger.MaxSize,
		MaxBackups: cfg.Logger.MaxBackups,
		MaxAge:     cfg.Logger.MaxAge,
		Compress:   cfg.Logger.Compress,
		Console:    cfg.Logger.Console,
	}); err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}
	defer logger.Sync()

	// 初始化验证器
	validator.Init()

	// 初始化数据库（可选，失败不阻止启动）
	if err := database.Init(&cfg.Database); err != nil {
		logger.Warn("数据库初始化失败（可选）", zap.Error(err))
	} else {
		logger.Info("数据库连接成功")
		defer database.Close()
	}

	// 初始化Redis（可选，失败不阻止启动）
	if err := cache.InitRedis(&cfg.Redis); err != nil {
		logger.Warn("Redis初始化失败（可选）", zap.Error(err))
	} else {
		logger.Info("Redis连接成功")
	}

	logger.Info("服务启动中...",
		zap.String("mode", cfg.Server.Mode),
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port),
	)

	// 设置gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建引擎
	engine := gin.New()

	// 使用gin自带中间件（阶段3会替换为自定义中间件）
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	// 注册路由
	setupRouter(engine)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// sql数据迁移
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

	logger.Info("服务已关闭")
}

// setupRouter 配置路由
func setupRouter(engine *gin.Engine) {
	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API路由组
	api := engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// 示例接口
			v1.GET("/ping", func(c *gin.Context) {
				response.Success(c, gin.H{
					"message": "pong",
				})
			})

			// 示例：参数错误
			v1.GET("/error", func(c *gin.Context) {
				response.ParamError(c, "这是一个参数错误示例")
			})

			// 示例：业务错误
			v1.GET("/business", func(c *gin.Context) {
				response.Fail(c, 10001, "业务处理失败")
			})

			v1.POST("/register", func(c *gin.Context) {
				var req RegisterRequest

				// ShouldBindJSON 会自动触发我们在 validator.Init 中注册的验证器
				if err := c.ShouldBindJSON(&req); err != nil {
					// 调用封装好的翻译函数
					errMsg := validator.TranslateError(err)
					response.ParamError(c, errMsg)
					return
				}

				// 验证通过，执行业务逻辑（如保存到 MySQL）
				response.Success(c, gin.H{
					"name": req.Username,
				})
			})
		}
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3" label:"用户名"`
	Phone    string `json:"phone" binding:"required,phone" label:"手机号"`
	Password string `json:"password" binding:"required,password" label:"密码"`
	Email    string `json:"email" binding:"required,email" label:"电子邮箱"`
}
