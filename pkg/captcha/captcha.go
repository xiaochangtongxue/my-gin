package captcha

import (
	"context"
	"fmt"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/xiaochangtongxue/my-gin/pkg/cache"
)

// Captcha 验证码管理器
type Captcha struct {
	driver base64Captcha.Driver
	store  base64Captcha.Store
	cache  cache.Cache
	expire time.Duration
}

// Config 验证码配置
type Config struct {
	Length int           // 验证码长度
	Width  int           // 图片宽度
	Height int           // 图片高度
	Expire time.Duration // 过期时间
}

// DefaultConfig 默认配置
var DefaultConfig = Config{
	Length: 4,
	Width:  120,
	Height: 40,
	Expire: 5 * time.Minute,
}

// New 创建验证码管理器
func New(cfg Config, c cache.Cache) *Captcha {
	if cfg.Length == 0 {
		cfg.Length = DefaultConfig.Length
	}
	if cfg.Width == 0 {
		cfg.Width = DefaultConfig.Width
	}
	if cfg.Height == 0 {
		cfg.Height = DefaultConfig.Height
	}
	if cfg.Expire == 0 {
		cfg.Expire = DefaultConfig.Expire
	}

	// 创建数字验证码驱动
	driver := base64Captcha.NewDriverDigit(
		cfg.Height,
		cfg.Width,
		cfg.Length,
		0.7, // 噪点干扰
		4,   // 配置项
	)

	return &Captcha{
		driver: driver,
		store:  base64Captcha.NewMemoryStore(1024, cfg.Expire),
		cache:  c,
		expire: cfg.Expire,
	}
}

// Generate 生成验证码
func (c *Captcha) Generate(ctx context.Context) (captchaID, imageBase64 string, err error) {
	// 创建验证码实例
	captcha := base64Captcha.NewCaptcha(c.driver, c.store)

	// 生成验证码，返回 (id, content, answer, err)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return "", "", fmt.Errorf("生成验证码失败: %w", err)
	}

	// 将答案存入 Redis（用于验证）
	key := c.captchaKey(id)
	if err := c.cache.Set(ctx, key, answer, c.expire); err != nil {
		return "", "", fmt.Errorf("存储验证码失败: %w", err)
	}

	return id, b64s, nil
}

// Verify 验证验证码
func (c *Captcha) Verify(ctx context.Context, captchaID, answer string) bool {
	if captchaID == "" || answer == "" {
		return false
	}

	key := c.captchaKey(captchaID)
	storedAnswer, err := c.cache.Get(ctx, key)
	if err != nil || storedAnswer == "" {
		return false
	}

	// 验证成功后删除验证码（一次性使用）
	_ = c.cache.Del(ctx, key)

	return storedAnswer == answer
}

// captchaKey 生成 Redis key
func (c *Captcha) captchaKey(captchaID string) string {
	return "captcha:" + captchaID
}
