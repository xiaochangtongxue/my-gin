package notify

import (
	"context"

	"go.uber.org/zap"
	"github.com/xiaochangtongxue/my-gin/pkg/logger"
)

// NoopNotifier 空实现的通知器，用于开发环境或暂时不需要通知的场景
type NoopNotifier struct{}

// NewNoopNotifier 创建空实现的通知器
func NewNoopNotifier() *NoopNotifier {
	return &NoopNotifier{}
}

// SendRegisterNotice 发送注册成功通知（仅记录日志）
func (n *NoopNotifier) SendRegisterNotice(ctx context.Context, to, username string) error {
	logger.Info("注册成功通知",
		zap.String("to", to),
		zap.String("username", username))
	return nil
}

// SendAccountLockedNotice 发送账号锁定通知（仅记录日志）
func (n *NoopNotifier) SendAccountLockedNotice(ctx context.Context, to, username string, lockDuration int) error {
	logger.Warn("账号锁定通知",
		zap.String("to", to),
		zap.String("username", username),
		zap.Int("lock_duration", lockDuration))
	return nil
}

// SendLoginAlertNotice 发送异地登录提醒通知（仅记录日志）
func (n *NoopNotifier) SendLoginAlertNotice(ctx context.Context, to, location string) error {
	logger.Info("异地登录提醒",
		zap.String("to", to),
		zap.String("location", location))
	return nil
}

// SendPasswordResetNotice 发送密码重置通知（仅记录日志）
func (n *NoopNotifier) SendPasswordResetNotice(ctx context.Context, to string) error {
	logger.Info("密码重置通知",
		zap.String("to", to))
	return nil
}