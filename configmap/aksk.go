package configmap

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

type CMAKSKProvider struct {
	FilePath string
	AkskMap  sync.Map
}

func NewCMAKSKProvider(filePath string) (*CMAKSKProvider, error) {
	if filePath == "" {
		filePath = defaultAkskFilePath
	}
	provider := &CMAKSKProvider{
		FilePath: filePath,
		AkskMap:  sync.Map{},
	}

	return provider, nil
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
	content, err := os.ReadFile(pvd.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", pvd.FilePath, err)
	}

	var aksk types.AKSK
	if err = json.Unmarshal([]byte(content), &aksk); err != nil {
		return nil, fmt.Errorf("json unmarshal file %s error: %v", pvd.FilePath, err)
	}

	pvd.AkskMap.Delete("aksk")
	pvd.AkskMap.Store("aksk", &aksk)

	return &aksk, nil
}
