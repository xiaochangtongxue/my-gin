package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/xiaochangtongxue/my-gin/pkg/config"
)

var (
	db         *gorm.DB
	dbConfig   *config.DatabaseConfig // 存储配置供后续使用
)

// Init 初始化数据库连接
func Init(cfg *config.DatabaseConfig) error {
	// 存储配置供后续使用
	dbConfig = cfg

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
		cfg.ParseTime,
	)

	// GORM 配置
	gormConfig := &gorm.Config{
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 跳过默认事务
		SkipDefaultTransaction: true,
	}

	// 设置日志级别（根据环境）
	if cfg.SlowThreshold > 0 {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	// 连接数据库
	var err error
	db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return nil
}

// Get 获取数据库实例
func Get() *gorm.DB {
	if db == nil {
		panic("数据库未初始化")
	}
	return db
}

// Close 关闭数据库连接
func Close() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// Ping 检查数据库连接健康状态
func Ping() error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
