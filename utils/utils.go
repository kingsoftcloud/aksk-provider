package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/kingsoftcloud/aksk-provider/types"
	utils "github.com/kingsoftcloud/aksk-provider/utils/encryption"
)

const (
	TimeLayoutStr    = "2006-01-02T15:04:05"
	DefaultExpiredAt = 100 * 365 * 24 * time.Hour
)

func createEncryptorConfig(key string, cipher string) (*utils.EncryptorConfig, error) {
	cipher = strings.ToUpper(cipher)
	switch cipher {
	case "RSA":
		privateKey, err := utils.LoadPrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to load privateKey: %v", err)
		}
		return &utils.EncryptorConfig{
			Cipher: cipher,
			RSAKeys: &utils.RSAKeysConfig{
				PrivateKey: privateKey,
			},
		}, nil
	default:
		// default cipher
		return &utils.EncryptorConfig{
			Cipher: "AES256",
			AESKey: key,
		}, nil
	}
}

func DecryptData(cryted string, key string, cipher string) (string, error) {
	encryptorConfig, err := createEncryptorConfig(key, cipher)
	if err != nil {
		return "", err
	}

	decryptor, err := utils.NewEncryptorWithOptions(encryptorConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create decryptor: %v", err)
	}
	return decryptor.Decrypt(cryted)
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
		ts, err := time.Parse(TimeLayoutStr, strings.TrimSpace(string(tsFile)))
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

func GetAkskSecret(akskName, akskNameSpace string, clientset *kubernetes.Clientset) (*types.AKSK, error) {
	secret, err := clientset.CoreV1().Secrets(akskNameSpace).Get(context.TODO(), akskName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get aksk secret %s error: %v", akskName, err)
	}
	ts, err := time.Parse(TimeLayoutStr, strings.TrimSpace(string(secret.Data["expired_at"])))
	if err != nil {
		ts = time.Now().Add(DefaultExpiredAt)
	}
	aksk := &types.AKSK{
		AK:            string(secret.Data["ak"]),
		SK:            string(secret.Data["sk"]),
		Cipher:        string(secret.Data["cipher"]),
		ExpiredAt:     ts,
		SecurityToken: string(secret.Data["securityToken"]),
	}
	return aksk, nil
}

func GetAkskConfigMap(akskName, akskNameSpace string, clientset *kubernetes.Clientset) (*types.AKSK, error) {
	configMap, err := clientset.CoreV1().ConfigMaps(akskNameSpace).Get(context.TODO(), akskName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get aksk configMap %s error: %v", akskName, err)
	}
	ts, err := time.Parse(TimeLayoutStr, strings.TrimSpace(configMap.Data["expired_at"]))
	if err != nil {
		ts = time.Now().Add(DefaultExpiredAt)
	}
	aksk := &types.AKSK{
		AK:            configMap.Data["ak"],
		SK:            configMap.Data["sk"],
		Cipher:        configMap.Data["cipher"],
		ExpiredAt:     ts,
		SecurityToken: configMap.Data["securityToken"],
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
		ts, err := time.Parse(TimeLayoutStr, akskMap["expired_at"])
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
