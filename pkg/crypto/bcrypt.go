// Package crypto 加密解密工具包
package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost 默认 bcrypt 加密强度
	// 2^Cost 次迭代，值越大计算时间越长，推荐 10-12
	DefaultCost = bcrypt.DefaultCost
	// MinCost 最小加密强度
	MinCost = 4
	// MaxCost 最大加密强度
	MaxCost = 31
)

// HashPassword 使用 bcrypt 对密码进行哈希
// cost: 加密强度（4-31），推荐 10-12，传 0 使用默认值
func HashPassword(password string, cost ...int) (string, error) {
	c := DefaultCost
	if len(cost) > 0 && cost[0] > 0 {
		c = cost[0]
		if c < MinCost {
			c = MinCost
		}
		if c > MaxCost {
			c = MaxCost
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), c)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword 验证密码是否匹配哈希值
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// ValidateHash 验证哈希值是否有效（格式正确）
func ValidateHash(hash string) error {
	_, err := bcrypt.Cost([]byte(hash))
	return err
}

// GetHashCost 获取哈希值的 cost 参数
func GetHashCost(hash string) (int, error) {
	return bcrypt.Cost([]byte(hash))
}
