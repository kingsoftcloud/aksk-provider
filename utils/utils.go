package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"newgit.op.ksyun.com/kce/aksk-provider/types"
)

const (
	timeLayoutStr = "2006-01-02T15:04:05"
)

func AesDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted) //base解码
	k := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
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
	return string(orig)
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
