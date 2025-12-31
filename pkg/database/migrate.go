package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/xiaochangtongxue/my-gin/db/migrations"
	"github.com/xiaochangtongxue/my-gin/pkg/config"
)

// NewMigrate 创建迁移实例
// 使用嵌入的迁移文件系统，避免容器内文件路径问题
func NewMigrate() (*migrate.Migrate, error) {
	// 从嵌入的文件系统创建迁移源
	d, err := iofs.New(migrations.MigrationFS, ".")
	if err != nil {
		return nil, fmt.Errorf("无法初始化内嵌迁移文件: %w", err)
	}

	// 创建迁移实例，使用数据库连接字符串
	return migrate.NewWithSourceInstance("iofs", d, "mysql://"+getDSN())
}

// Up 执行所有未执行的迁移
func Up() error {
	m, err := NewMigrate()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// UpN 执行 N 个迁移
func UpN(n int) error {
	m, err := NewMigrate()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(n); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Down 回滚最后一个迁移
func Down() error {
	m, err := NewMigrate()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// DownN 回滚 N 个迁移
func DownN(n int) error {
	m, err := NewMigrate()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(-n); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Version 获取当前版本
func Version() (uint, bool, error) {
	m, err := NewMigrate()
	if err != nil {
		return 0, false, err
	}
	defer m.Close()

	return m.Version()
}

// Status 查看迁移状态
func Status() error {
	m, err := NewMigrate()
	if err != nil {
		return err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil {
		return err
	}

	fmt.Printf("Schema version: %d, Dirty: %v\n", version, dirty)

	return nil
}

// Force 设置版本（慎用）
func Force(version int) error {
	m, err := NewMigrate()
	if err != nil {
		return err
	}
	defer m.Close()

	return m.Force(version)
}

// getDSN 获取数据库连接字符串
func getDSN() string {
	cfg := config.Get().Database
	charset := cfg.Charset
	if charset == "" {
		charset = "utf8mb4"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		charset,
		cfg.ParseTime,
	)
}
