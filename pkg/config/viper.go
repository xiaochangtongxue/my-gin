package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// globalCfg 全局配置实例
var globalCfg *Config

// Init 初始化配置
func Init(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件
	if configPath == "" {
		// 默认从环境变量获取配置文件路径
		configPath = os.Getenv("CONFIG_PATH")
		if configPath == "" {
			configPath = "configs/config.yaml"
		}
	}

	// 获取运行环境
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = os.Getenv("GO_ENV")
		if env == "" {
			env = "dev"
		}
	}

	// 根据环境选择配置文件
	if env != "dev" {
		configPath = fmt.Sprintf("configs/config.%s.yaml", env)
	}

	// 设置配置文件路径
	v.SetConfigFile(configPath)

	// 设置配置文件类型
	v.SetConfigType("yaml")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 环境变量前缀
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 解析配置
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 配置校验
	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("配置校验失败: %w", err)
	}

	globalCfg = cfg
	return cfg, nil
}

// Get 获取全局配置
func Get() *Config {
	if globalCfg == nil {
		panic("配置未初始化，请先调用 Init()")
	}
	return globalCfg
}

// MustGet 获取全局配置，如果未初始化则panic
func MustGet() *Config {
	return Get()
}

// validate 配置校验
func validate(cfg *Config) error {
	// 基础校验
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("server.port 无效: %d", cfg.Server.Port)
	}

	if cfg.Database.Host == "" {
		return fmt.Errorf("database.host 不能为空")
	}

	if cfg.Database.Database == "" {
		return fmt.Errorf("database.database 不能为空")
	}

	if cfg.Redis.Host == "" {
		return fmt.Errorf("redis.host 不能为空")
	}

	if cfg.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret 不能为空")
	}

	if len(cfg.JWT.Secret) < 16 {
		return fmt.Errorf("jwt.secret 长度不能少于16位")
	}

	return nil
}