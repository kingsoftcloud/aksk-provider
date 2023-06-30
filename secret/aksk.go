package secret

import (
	"fmt"
	"sync"

	prvd "newgit.op.ksyun.com/kce/aksk-provider"
	"newgit.op.ksyun.com/kce/aksk-provider/types"
	"newgit.op.ksyun.com/kce/aksk-provider/utils"
)

const (
	defaultAkskFilePath = "/var/lib/aksk"
)

var _ prvd.AKSKProvider = &SecretAKSKProvider{}

type SecretAKSKProvider struct {
	FilePath  string
	CipherKey string
	AkskMap   sync.Map
}

func NewSecretAKSKProvider(filePath, cipherKey string) prvd.AKSKProvider {
	if filePath == "" {
		filePath = defaultAkskFilePath
	}

	return &SecretAKSKProvider{
		FilePath:  filePath,
		CipherKey: cipherKey,
		AkskMap:   sync.Map{},
	}
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
	aksk, err := utils.ParseAkskDirectory(pvd.FilePath)
	if err != nil {
		return nil, err
	}

	if aksk.Cipher == "none" {
		return aksk, nil
	}

	aksk.SK, err = utils.AesDecrypt(aksk.SK, pvd.CipherKey)
	if err != nil {
		return nil, err
	}
	aksk.SecurityToken, err = utils.AesDecrypt(aksk.SecurityToken, pvd.CipherKey)
	if err != nil {
		return nil, err
	}

	pvd.AkskMap.Delete("aksk")
	pvd.AkskMap.Store("aksk", aksk)

	return aksk, nil
}
