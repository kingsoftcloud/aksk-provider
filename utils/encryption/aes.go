package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func AesEncrypt(plainData string, key string) (string, error) {
	originData := []byte(plainData)
	k := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", fmt.Errorf("create new aes cipher failed: %v", err)
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	originData = PKCS7Padding(originData, blockSize)
	// 加密模式
	iv := k[:blockSize]
	blockMode := cipher.NewCBCEncrypter(block, iv)
	// 创建数组
	cryted := make([]byte, len(originData))
	// 加密
	blockMode.CryptBlocks(cryted, originData)
	return base64.StdEncoding.EncodeToString(cryted), nil
}

// 补码
// AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesDecrypt(cryted string, key string) (string, error) {
	// 转成字节数组
	crytedByte, err := base64.StdEncoding.DecodeString(cryted) //base解码
	if err != nil {
		return "", fmt.Errorf("decodeString failed: %v", err)
	}
	k := []byte(key)
	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", fmt.Errorf("create new aes cipher failed: %v", err)
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig), nil
}

// 去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
