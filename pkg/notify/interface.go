package notify

import "context"

// Notifier 通知接口，用于发送各类通知（短信、邮件等）
type Notifier interface {
	// SendRegisterNotice 发送注册成功通知
	SendRegisterNotice(ctx context.Context, to, username string) error

	// SendAccountLockedNotice 发送账号锁定通知
	SendAccountLockedNotice(ctx context.Context, to, username string, lockDuration int) error

	// SendLoginAlertNotice 发送异地登录提醒通知
	SendLoginAlertNotice(ctx context.Context, to, location string) error

	// SendPasswordResetNotice 发送密码重置通知
	SendPasswordResetNotice(ctx context.Context, to string) error
}