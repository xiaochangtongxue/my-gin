// Package crypto 加密解密工具包
package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

var (
	// ErrInvalidBlockSize 数据块大小无效
	ErrInvalidBlockSize = errors.New("数据块大小无效")
	// ErrInvalidPKCS7Data PKCS7 填充数据无效
	ErrInvalidPKCS7Data = errors.New("PKCS7 填充数据无效")
)

// AESCBCEncrypt AES CBC 模式加密
// key: 密钥（16/24/32 字节对应 AES-128/192/256）
// iv: 初始化向量（16 字节）
// plaintext: 明文
// 返回 Base64 编码的密文
func AESCBCEncrypt(key, iv, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建 AES cipher 失败: %w", err)
	}

	// PKCS7 填充
	plaintext = pkcs7Padding(plaintext, block.BlockSize())

	// CBC 模式加密
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESCBCEncryptString AES CBC 模式加密字符串
// key: 密钥（16/24/32 字节）
// iv: 初始化向量（16 字节）
// plaintext: 明文字符串
func AESCBCEncryptString(key, iv, plaintext string) (string, error) {
	return AESCBCEncrypt([]byte(key), []byte(iv), []byte(plaintext))
}

// AESCBCDecrypt AES CBC 模式解密
// key: 密钥
// iv: 初始化向量
// ciphertext: Base64 编码的密文
// 返回明文
func AESCBCDecrypt(key, iv []byte, ciphertext string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("Base64 解码失败: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建 AES cipher 失败: %w", err)
	}

	if len(data)%block.BlockSize() != 0 {
		return nil, ErrInvalidBlockSize
	}

	// CBC 模式解密
	plaintext := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, data)

	// 去除 PKCS7 填充
	plaintext, err = pkcs7UnPadding(plaintext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// AESCBCDecryptString AES CBC 模式解密字符串
func AESCBCDecryptString(key, iv, ciphertext string) (string, error) {
	plaintext, err := AESCBCDecrypt([]byte(key), []byte(iv), ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// AECBCEncrypt AES ECB 模式加密（不推荐，仅用于兼容）
func AECBCEncrypt(key, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建 AES cipher 失败: %w", err)
	}

	plaintext = pkcs7Padding(plaintext, block.BlockSize())
	ciphertext := make([]byte, len(plaintext))

	// ECB 模式需要手动实现（每次加密一个块）
	bs := block.BlockSize()
	for i := 0; i < len(plaintext); i += bs {
		block.Encrypt(ciphertext[i:i+bs], plaintext[i:i+bs])
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AECBCEncryptString AES ECB 模式加密字符串
func AECBCEncryptString(key, plaintext string) (string, error) {
	return AECBCEncrypt([]byte(key), []byte(plaintext))
}

// AECBDecrypt AES ECB 模式解密
func AECBDecrypt(key []byte, ciphertext string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("Base64 解码失败: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建 AES cipher 失败: %w", err)
	}

	if len(data)%block.BlockSize() != 0 {
		return nil, ErrInvalidBlockSize
	}

	plaintext := make([]byte, len(data))
	bs := block.BlockSize()
	for i := 0; i < len(data); i += bs {
		block.Decrypt(plaintext[i:i+bs], data[i:i+bs])
	}

	plaintext, err = pkcs7UnPadding(plaintext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// AECBDecryptString AES ECB 模式解密字符串
func AECBDecryptString(key, ciphertext string) (string, error) {
	plaintext, err := AECBDecrypt([]byte(key), ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// pkcs7Padding PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7UnPadding 去除 PKCS7 填充
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	padding := int(data[length-1])
	if padding > length || padding > aes.BlockSize {
		return nil, ErrInvalidPKCS7Data
	}
	// 验证填充
	for i := length - padding; i < length; i++ {
		if data[i] != byte(padding) {
			return nil, ErrInvalidPKCS7Data
		}
	}
	return data[:length-padding], nil
}
