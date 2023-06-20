package secret

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"newgit.op.ksyun.com/kce/aksk-provider/types"
	"newgit.op.ksyun.com/kce/aksk-provider/utils"
)

const (
	defaultAkskFilePath = "/var/lib/aksk"
)

type SecretAKSKProvider struct {
	FilePath  string
	CipherKey string
	AkskMap   sync.Map
}

func NewSecretAKSKProvider(filePath, cipherKey string) (*SecretAKSKProvider, error) {
	if filePath == "" {
		filePath = defaultAkskFilePath
	}
	provider := &SecretAKSKProvider{
		FilePath:  filePath,
		CipherKey: cipherKey,
		AkskMap:   sync.Map{},
	}

	return provider, nil
}

func (pvd *SecretAKSKProvider) GetAKSK() (*types.AKSK, error) {
	if v, ok := pvd.AkskMap.Load("aksk"); ok && !utils.IsExpired(v.(*types.AKSK).ExpiredAt) {
		return v.(*types.AKSK), nil
	}

	aksk, err := pvd.ReloadAKSK()
	if err != nil {
		return nil, fmt.Errorf("reload aksk from file %s error: %v", pvd.FilePath, err)
	}

	return aksk, nil
}

func (pvd *SecretAKSKProvider) ReloadAKSK() (*types.AKSK, error) {
	content, err := os.ReadFile(pvd.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", pvd.FilePath, err)
	}

	var aksk *types.AKSK
	if err = json.Unmarshal([]byte(content), &aksk); err != nil {
		return nil, fmt.Errorf("json unmarshal file %s error: %v", pvd.FilePath, err)
	}

	if aksk.Cipher == "none" {
		return aksk, nil
	}

	aksk.SK = utils.AesDecrypt(aksk.SK, pvd.CipherKey)
	aksk.SecurityToken = utils.AesDecrypt(aksk.SecurityToken, pvd.CipherKey)

	pvd.AkskMap.Delete("aksk")
	pvd.AkskMap.Store("aksk", &aksk)

	return aksk, nil
}
