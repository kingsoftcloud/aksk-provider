package file

import (
	"fmt"
	"sync"

	prvd "ezone.ksyun.com/ezone/kce/aksk-provider"
	"ezone.ksyun.com/ezone/kce/aksk-provider/types"
	"ezone.ksyun.com/ezone/kce/aksk-provider/utils"
)

const (
	defaultAkskFilePath = "/var/lib/aksk"
)

var _ prvd.AKSKProvider = &FileAKSKProvider{}

type FileAKSKProvider struct {
	FilePath  string
	CipherKey string
	AkskMap   sync.Map
}

func NewFileAKSKProvider(filePath, cipherKey string) prvd.AKSKProvider {
	if filePath == "" {
		filePath = defaultAkskFilePath
	}

	return &FileAKSKProvider{
		FilePath:  filePath,
		CipherKey: cipherKey,
		AkskMap:   sync.Map{},
	}
}

func (pvd *FileAKSKProvider) GetAKSK() (*types.AKSK, error) {
	if v, ok := pvd.AkskMap.Load("aksk"); ok && !utils.IsExpired(v.(*types.AKSK).ExpiredAt) {
		return v.(*types.AKSK), nil
	}

	aksk, err := pvd.ReloadAKSK()
	if err != nil {
		return nil, fmt.Errorf("reload aksk from file %s error: %v", pvd.FilePath, err)
	}

	return aksk, nil
}

func (pvd *FileAKSKProvider) ReloadAKSK() (*types.AKSK, error) {
	aksk, err := utils.ParseAkskDirectory(pvd.FilePath)
	if err != nil {
		return nil, err
	}

	if aksk.Cipher == "none" || aksk.Cipher == "" {
		return aksk, nil
	}

	aksk.SK, err = utils.AesDecrypt(aksk.SK, pvd.CipherKey)
	if err != nil {
		return nil, err
	}

	if aksk.SecurityToken != "" {
		aksk.SecurityToken, err = utils.AesDecrypt(aksk.SecurityToken, pvd.CipherKey)
		if err != nil {
			return nil, err
		}
	}

	pvd.AkskMap.Delete("aksk")
	pvd.AkskMap.Store("aksk", aksk)

	return aksk, nil
}
