package config

import (
	"time"
)

// Config 全局配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server" validate:"required"`
	Database DatabaseConfig `mapstructure:"database" validate:"required"`
	Redis    RedisConfig    `mapstructure:"redis" validate:"required"`
	Logger   LoggerConfig   `mapstructure:"logger" validate:"required"`
	JWT      JWTConfig      `mapstructure:"jwt" validate:"required"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Mode      string `mapstructure:"mode" validate:"oneof=debug release test"`
	Host      string `mapstructure:"host" validate:"required"`
	Port      int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host" validate:"required"`
	Port            int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Username        string `mapstructure:"username" validate:"required"`
	Password        string `mapstructure:"password" validate:"required"`
	Database        string `mapstructure:"database" validate:"required"`
	Charset         string `mapstructure:"charset"`
	ParseTime       bool   `mapstructure:"parse_time"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" validate:"min=0"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" validate:"min=0"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	SlowThreshold   time.Duration `mapstructure:"slow_threshold"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host        string `mapstructure:"host" validate:"required"`
	Port        int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Password    string `mapstructure:"password"`
	DB          int    `mapstructure:"db" validate:"min=0,max=15"`
	PoolSize    int    `mapstructure:"pool_size" validate:"min=0"`
	MinIdleConn int    `mapstructure:"min_idle_conn" validate:"min=0"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level" validate:"oneof=debug info warn error fatal"`
	FileName   string `mapstructure:"file_name"`
	MaxSize    int    `mapstructure:"max_size" validate:"min=1"`
	MaxBackups int    `mapstructure:"max_backups" validate:"min=0"`
	MaxAge     int    `mapstructure:"max_age" validate:"min=0"`
	Compress   bool   `mapstructure:"compress"`
	Console    bool   `mapstructure:"console"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `mapstructure:"secret" validate:"required,min=16"`
	ExpireTime time.Duration `mapstructure:"expire_time" validate:"required"`
	Issuer     string        `mapstructure:"issuer" validate:"required"`
}