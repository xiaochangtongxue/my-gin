package utils

import (
	"math/bits"
)

const (
	uidPrefix = 10000000000000 // 前缀，确保14位
)

// feistelEncode 使用 Feistel 网络对数字进行可逆混淆
func feistelEncode(id uint64) uint64 {
	const rounds = 3
	const mask = 0xFFFFFFFFFFF // 44位

	l := id >> 22
	r := id & 0x3FFFFF

	for i := 0; i < rounds; i++ {
		next := l ^ roundFunction(r, i)
		l, r = r, next
	}

	return (l << 22) | r
}

// roundFunction Feistel 轮函数
func roundFunction(val uint64, round int) uint64 {
	key := uint64(374601907 + round*31) // 质数
	return bits.RotateLeft64(val*key+key, 11) & 0x3FFFFF
}

// EncodeUID 将数据库自增ID混淆为14位UID
func EncodeUID(id uint64) uint64 {
	return uidPrefix + feistelEncode(id)
}

// DecodeUID 将14位UID解码为原始自增ID（仅供调试使用）
func DecodeUID(uid uint64) uint64 {
	x := uid - uidPrefix
	const rounds = 3
	l := x >> 22
	r := x & 0x3FFFFF

	// Feistel 解密（反向执行）
	for i := rounds - 1; i >= 0; i-- {
		next := r ^ roundFunction(l, i)
		r, l = l, next
	}

	return (l << 22) | r
}