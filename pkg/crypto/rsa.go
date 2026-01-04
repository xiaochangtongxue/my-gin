// Package crypto 加密解密工具包
package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	// ErrInvalidPublicKey 无效的公钥
	ErrInvalidPublicKey = errors.New("无效的公钥")
	// ErrInvalidPrivateKey 无效的私钥
	ErrInvalidPrivateKey = errors.New("无效的私钥")
	// ErrDataTooLong 数据过长
	ErrDataTooLong = errors.New("数据过长，无法加密")
)

// GenerateRSAKey 生成 RSA 密钥对
// bits: 密钥位数（推荐 2048 或 4096）
func GenerateRSAKey(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("生成 RSA 密钥失败: %w", err)
	}
	return privateKey, &privateKey.PublicKey, nil
}

// PrivateKeyToPEM 将私钥转换为 PEM 格式字符串
func PrivateKeyToPEM(privateKey *rsa.PrivateKey) string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return string(privateKeyPEM)
}

// PublicKeyToPEM 将公钥转换为 PEM 格式字符串
func PublicKeyToPEM(publicKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("序列化公钥失败: %w", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM), nil
}

// ParsePrivateKeyFromPEM 从 PEM 字符串解析私钥
func ParsePrivateKeyFromPEM(pemData string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, ErrInvalidPrivateKey
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// 尝试 PKCS8 格式
		pkcs8Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, ErrInvalidPrivateKey
		}
		rsaKey, ok := pkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, ErrInvalidPrivateKey
		}
		return rsaKey, nil
	}
	return privateKey, nil
}

// ParsePublicKeyFromPEM 从 PEM 字符串解析公钥
func ParsePublicKeyFromPEM(pemData string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, ErrInvalidPublicKey
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, ErrInvalidPublicKey
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidPublicKey
	}

	return rsaPublicKey, nil
}

// RSAEncrypt RSA 公钥加密
func RSAEncrypt(publicKey *rsa.PublicKey, plaintext []byte) (string, error) {
	// 使用 OAEP 填充模式
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plaintext, nil)
	if err != nil {
		return "", fmt.Errorf("RSA 加密失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// RSAEncryptString RSA 公钥加密字符串
func RSAEncryptString(publicKey *rsa.PublicKey, plaintext string) (string, error) {
	return RSAEncrypt(publicKey, []byte(plaintext))
}

// RSAEncryptWithPEM 使用 PEM 格式公钥加密
func RSAEncryptWithPEM(publicKeyPEM, plaintext string) (string, error) {
	publicKey, err := ParsePublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return "", err
	}
	return RSAEncrypt(publicKey, []byte(plaintext))
}

// RSADecrypt RSA 私钥解密
func RSADecrypt(privateKey *rsa.PrivateKey, ciphertext string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("Base64 解码失败: %w", err)
	}

	// 使用 OAEP 填充模式
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("RSA 解密失败: %w", err)
	}

	return plaintext, nil
}

// RSADecryptString RSA 私钥解密为字符串
func RSADecryptString(privateKey *rsa.PrivateKey, ciphertext string) (string, error) {
	plaintext, err := RSADecrypt(privateKey, ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// RSADecryptWithPEM 使用 PEM 格式私钥解密
func RSADecryptWithPEM(privateKeyPEM, ciphertext string) (string, error) {
	privateKey, err := ParsePrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return "", err
	}
	return RSADecryptString(privateKey, ciphertext)
}

// RSASign RSA 私钥签名（使用 PSS）
func RSASign(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	hashed := sha256.Sum256(data)

	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hashed[:], &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
	})
	if err != nil {
		return "", fmt.Errorf("RSA 签名失败: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// RSASignString RSA 私钥签名字符串
func RSASignString(privateKey *rsa.PrivateKey, data string) (string, error) {
	return RSASign(privateKey, []byte(data))
}

// RSASignWithPEM 使用 PEM 格式私钥签名
func RSASignWithPEM(privateKeyPEM, data string) (string, error) {
	privateKey, err := ParsePrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return "", err
	}
	return RSASignString(privateKey, data)
}

// RSAVerify RSA 公钥验签（使用 PSS）
func RSAVerify(publicKey *rsa.PublicKey, data []byte, signature string) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("Base64 解码签名失败: %w", err)
	}

	hashed := sha256.Sum256(data)

	err = rsa.VerifyPSS(publicKey, crypto.SHA256, hashed[:], sig, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
	})
	if err != nil {
		return fmt.Errorf("RSA 验签失败: %w", err)
	}

	return nil
}

// RSAVerifyString RSA 公钥验签字符串
func RSAVerifyString(publicKey *rsa.PublicKey, data, signature string) error {
	return RSAVerify(publicKey, []byte(data), signature)
}

// RSAVerifyWithPEM 使用 PEM 格式公钥验签
func RSAVerifyWithPEM(publicKeyPEM, data, signature string) error {
	publicKey, err := ParsePublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return err
	}
	return RSAVerifyString(publicKey, data, signature)
}
