package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"ezone.ksyun.com/ezone/kce/aksk-provider/types"
)

const (
	timeLayoutStr    = "2006-01-02T15:04:05"
	DefaultExpiredAt = 100 * 365 * 24 * time.Hour
)

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

//去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func IsExpired(expiredAt time.Time) bool {
	if expiredAt.IsZero() || expiredAt.Before(time.Now()) {
		return true
	}

	return false
}

func ParseAkskDirectory(directory string) (*types.AKSK, error) {
	aksk := &types.AKSK{}

	akFile, err := os.ReadFile(directory + "/ak")
	if err != nil {
		return nil, fmt.Errorf("read file %s error: %v", directory+"/ak", err)
	}
	aksk.AK = strings.TrimSpace(string(akFile))

	skFile, err := os.ReadFile(directory + "/sk")
	if err != nil {
		return nil, fmt.Errorf("read file %s error: %v", directory+"/ak", err)
	}
	aksk.SK = strings.TrimSpace(string(skFile))

	_, err = os.Stat(directory + "/securityToken")
	if err == nil {
		tokenFile, err := os.ReadFile(directory + "/securityToken")
		if err != nil {
			return nil, fmt.Errorf("read file %s error: %v", directory+"/securityToken", err)
		}
		aksk.SecurityToken = strings.TrimSpace(string(tokenFile))
	}

	_, err = os.Stat(directory + "/expired_at")
	if err == nil {
		tsFile, err := os.ReadFile(directory + "/expired_at")
		if err != nil {
			return nil, fmt.Errorf("read file %s error: %v", directory+"/", err)
		}
		ts, err := time.Parse(timeLayoutStr, strings.TrimSpace(string(tsFile)))
		if err != nil {
			return nil, err
		}
		aksk.ExpiredAt = ts
	} else {
		aksk.ExpiredAt = time.Now().Add(DefaultExpiredAt)
	}

	_, err = os.Stat(directory + "/cipher")
	if err == nil {
		cipherFile, err := os.ReadFile(directory + "/cipher")
		if err != nil {
			return nil, fmt.Errorf("read file %s error: %v", directory+"/cipher", err)
		}
		aksk.Cipher = strings.TrimSpace(string(cipherFile))
	}

	return aksk, nil
}

func ParseAkskFile(filePath string) (*types.AKSK, error) {
	fd, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	bytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	akskMap := make(map[string]string)
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		array := strings.SplitN(line, ":", 2)
		key := strings.Trim(array[0], "\" ")
		value := strings.Trim(array[1], "\" ")
		akskMap[key] = value
	}

	aksk := &types.AKSK{
		AK:     akskMap["ak"],
		SK:     akskMap["sk"],
		Cipher: akskMap["cipher"],
	}

	if _, ok := akskMap["expired_at"]; ok {
		ts, err := time.Parse(timeLayoutStr, akskMap["expired_at"])
		if err != nil {
			return nil, err
		}
		aksk.ExpiredAt = ts
	}

	if _, ok := akskMap["securityToken"]; ok {
		aksk.SecurityToken = akskMap["securityToken"]
	}

	return aksk, nil
}
