package configmap

import (
	"fmt"
	"sync"

	prvd "newgit.op.ksyun.com/kce/aksk-provider"
	"newgit.op.ksyun.com/kce/aksk-provider/types"
	"newgit.op.ksyun.com/kce/aksk-provider/utils"
)

const (
	defaultAkskFilePath = "/etc/aksk"
)

var _ prvd.AKSKProvider = &CMAKSKProvider{}

type CMAKSKProvider struct {
	FilePath string
	AkskMap  sync.Map
}

func NewCMAKSKProvider(filePath string) *CMAKSKProvider {
	if filePath == "" {
		filePath = defaultAkskFilePath
	}
	provider := &CMAKSKProvider{
		FilePath: filePath,
		AkskMap:  sync.Map{},
	}

	return provider
}

func (pvd *CMAKSKProvider) GetAKSK() (*types.AKSK, error) {
	if v, ok := pvd.AkskMap.Load("aksk"); ok && !utils.IsExpired(v.(*types.AKSK).ExpiredAt) {
		return v.(*types.AKSK), nil
	}

	aksk, err := pvd.ReloadAKSK()
	if err != nil {
		return nil, fmt.Errorf("reload aksk from file %s error: %v", pvd.FilePath, err)
	}

	return aksk, nil
}

func (pvd *CMAKSKProvider) ReloadAKSK() (*types.AKSK, error) {
	aksk, err := utils.ParseAkskDirectory(pvd.FilePath)
	if err != nil {
		return nil, err
	}

	pvd.AkskMap.Delete("aksk")
	pvd.AkskMap.Store("aksk", aksk)

	return aksk, nil
}
