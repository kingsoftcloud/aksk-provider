package utils

import (
	"crypto/rsa"
	"fmt"
)

// Encryptor 定义了加解密接口
type Encryptor interface {
	Encrypt(plainText string) (string, error)
	Decrypt(cipherText string) (string, error)
}

// RSAEncryptor 使用 RSA 算法加解密
type RSAEncryptor struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// NewRSAEncryptor 创建一个支持 RSA 加解密的实例
func NewRSAEncryptor(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) *RSAEncryptor {
	return &RSAEncryptor{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

// Encrypt 实现 RSA 加密
func (r *RSAEncryptor) Encrypt(plainData string) (string, error) {
	// 使用公钥加密
	return RsaEncrypt(plainData, r.publicKey)
}

// Decrypt 实现 RSA 解密
func (r *RSAEncryptor) Decrypt(cipherText string) (string, error) {
	// 使用私钥解密
	return RsaDecrypt(cipherText, r.privateKey)
}

// AESEncryptor 使用 AES 算法加解密
type AESEncryptor struct {
	key string
}

// NewAESEncryptor 创建一个支持 AES 加解密的实例
func NewAESEncryptor(key string) *AESEncryptor {
	return &AESEncryptor{
		key: key,
	}
}

// Encrypt 实现 AES 加密
func (a *AESEncryptor) Encrypt(plainText string) (string, error) {
	return AesEncrypt(plainText, a.key)
}

// Decrypt 实现 AES 解密
func (a *AESEncryptor) Decrypt(cipherText string) (string, error) {
	return AesDecrypt(cipherText, a.key)
}

// EncryptorConfig 用于传递加密算法的配置
type EncryptorConfig struct {
	Cipher  string
	RSAKeys *RSAKeysConfig
	AESKey  string
}

// RSAKeysConfig 用于传递 RSA 密钥对的配置
type RSAKeysConfig struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

// Option 是用于配置 EncryptorConfig 的函数类型
type Option func(*EncryptorConfig)

// WithRSAKeys 配置 RSA 密钥对
func WithRSAKeys(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) Option {
	return func(config *EncryptorConfig) {
		config.RSAKeys = &RSAKeysConfig{
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		}
	}
}

// WithAESKey 配置 AES 密钥
func WithAESKey(key string) Option {
	return func(config *EncryptorConfig) {
		config.AESKey = key
	}
}

// NewEncryptorWithOptions 根据传入的配置和选项创建相应的加解密算法提供者
func NewEncryptorWithOptions(config *EncryptorConfig, options ...Option) (Encryptor, error) {
	// 应用所有选项配置
	for _, option := range options {
		option(config)
	}

	// 根据算法类型初始化不同的加解密器
	switch config.Cipher {
	case "RSA":
		if config.RSAKeys == nil {
			return nil, fmt.Errorf("invalid RSA keys configuration")
		}
		if config.RSAKeys.PublicKey != nil && config.RSAKeys.PrivateKey == nil {
			return NewRSAEncryptor(config.RSAKeys.PublicKey, nil), nil
		} else if config.RSAKeys.PrivateKey != nil && config.RSAKeys.PublicKey == nil {
			return NewRSAEncryptor(nil, config.RSAKeys.PrivateKey), nil
		} else if config.RSAKeys.PublicKey != nil && config.RSAKeys.PrivateKey != nil {
			return NewRSAEncryptor(config.RSAKeys.PublicKey, config.RSAKeys.PrivateKey), nil
		} else {
			return nil, fmt.Errorf("invalid RSA keys configuration")
		}
	case "AES256":
		if len(config.AESKey) == 0 {
			return nil, fmt.Errorf("invalid AES key")
		}
		return NewAESEncryptor(config.AESKey), nil

	default:
		return nil, fmt.Errorf("unsupported encryption cipher")
	}
}
