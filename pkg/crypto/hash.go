// Package crypto 加密解密工具包
package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash"
)

// RandomString 生成指定长度的随机十六进制字符串
// length: 期望的字符串长度（实际字节长度为 length/2）
func RandomString(length int) (string, error) {
	if length <= 0 {
		length = 64
	}
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// RandomBytes 生成指定字节数的随机数据
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MD5 计算字符串的 MD5 哈希值
// 注意：MD5 不安全，仅用于兼容性或非安全场景
func MD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// MD5Bytes 计算字节数组的 MD5 哈希值
func MD5Bytes(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 计算字符串的 SHA256 哈希值
func SHA256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256Bytes 计算字节数组的 SHA256 哈希值
func SHA256Bytes(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA256 计算 HMAC-SHA256 签名
// key: 密钥
// data: 待签名数据
func HMACSHA256(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA256Bytes 计算 HMAC-SHA256 签名（字节数组）
func HMACSHA256Bytes(key, data []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// HMAC 使用指定的哈希算法计算 HMAC
// hashFunc: 哈希函数构造器（如 sha256.New）
// key: 密钥
// data: 待签名数据
func HMAC(hashFunc func() hash.Hash, key, data []byte) []byte {
	h := hmac.New(hashFunc, key)
	h.Write(data)
	return h.Sum(nil)
}

// Base64Encode 进行 Base64 编码
func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Base64EncodeBytes 对字节数组进行 Base64 编码
func Base64EncodeBytes(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode 进行 Base64 解码
func Base64Decode(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// Base64DecodeBytes 对字节数组进行 Base64 解码
func Base64DecodeBytes(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

// Base64URLEncode 进行 URL Safe Base64 编码
func Base64URLEncode(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

// Base64URLDecode 进行 URL Safe Base64 解码
func Base64URLDecode(data string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
