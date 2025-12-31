package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// ctxKey 上下文key类型
type ctxKey string

const (
	TxKey ctxKey = "tx"
)

// WithTx 在context中设置事务
func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

// GetTx 从context中获取事务
func GetTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(TxKey).(*gorm.DB); ok {
		return tx
	}
	return nil
}

// HasTx 检查context中是否有事务
func HasTx(ctx context.Context) bool {
	return GetTx(ctx) != nil
}

// DB 获取数据库连接（优先使用事务）
func DB(ctx context.Context) *gorm.DB {
	if tx := GetTx(ctx); tx != nil {
		return tx
	}
	return Get()
}

// Transaction 执行事务
func Transaction(ctx context.Context, fn func(context.Context) error) error {
	return Get().Transaction(func(tx *gorm.DB) error {
		ctxWithTx := WithTx(ctx, tx)
		return fn(ctxWithTx)
	})
}

// TransactionWithResult 执行事务并返回结果
func TransactionWithResult[T any](ctx context.Context, fn func(context.Context) (T, error)) (T, error) {
	var result T
	var err error

	txErr := Get().Transaction(func(tx *gorm.DB) error {
		ctxWithTx := WithTx(ctx, tx)
		result, err = fn(ctxWithTx)
		return err
	})

	if txErr != nil {
		return result, txErr
	}

	return result, err
}

// Tx 事务接口
type Tx struct {
	tx *gorm.DB
}

// Begin 开始事务
func Begin() *Tx {
	return &Tx{
		tx: Get().Begin(),
	}
}

// Commit 提交事务
func (t *Tx) Commit() error {
	if t.tx == nil {
		return fmt.Errorf("事务未开始")
	}
	return t.tx.Commit().Error
}

// Rollback 回滚事务
func (t *Tx) Rollback() error {
	if t.tx == nil {
		return fmt.Errorf("事务未开始")
	}
	return t.tx.Rollback().Error
}

// DB 获取事务的数据库连接
func (t *Tx) DB() *gorm.DB {
	return t.tx
}

// Savepoint 设置保存点
func Savepoint(tx *gorm.DB, name string) error {
	return tx.SavePoint(name).Error
}

// RollbackTo 回滚到保存点
func RollbackTo(tx *gorm.DB, name string) error {
	return tx.RollbackTo(name).Error
}
