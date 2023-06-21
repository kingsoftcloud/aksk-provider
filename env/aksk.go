package env

import (
	"fmt"
	"os"
	"sync"
	"time"

	"newgit.op.ksyun.com/kce/aksk-provider/types"
	"newgit.op.ksyun.com/kce/aksk-provider/utils"
)

const (
	defaultAkEnv            = "AK"
	defaultSkEnv            = "SK"
	defaultSecurityTokenEnv = "SECURITY_TOKEN"
	defaultExpiredAt        = 100 * 365 * 24 * time.Hour
)

type EnvAKSKProvider struct {
	Encrypt   bool
	CipherKey string
	AkskMap   sync.Map
}

func NewEnvAKSKProvider(encrypt bool, cipherKey string) *EnvAKSKProvider {
	provider := &EnvAKSKProvider{
		Encrypt:   encrypt,
		CipherKey: cipherKey,
		AkskMap:   sync.Map{},
	}

	return provider
}

func (pvd *EnvAKSKProvider) GetAKSK() (*types.AKSK, error) {
	if v, ok := pvd.AkskMap.Load("aksk"); ok && !utils.IsExpired(v.(*types.AKSK).ExpiredAt) {
		return v.(*types.AKSK), nil
	}

	aksk, err := pvd.ReloadAKSK()
	if err != nil {
		return nil, fmt.Errorf("reload aksk from env error: %v", err)
	}

	return aksk, nil
}

func (pvd *EnvAKSKProvider) ReloadAKSK() (*types.AKSK, error) {
	if os.Getenv(defaultAkEnv) == "" {
		return nil, fmt.Errorf("get ak from env %s failed: nil", defaultAkEnv)
	}

	if os.Getenv(defaultSkEnv) == "" {
		return nil, fmt.Errorf("get sk from env %s failed: nil", defaultSkEnv)
	}

	aksk := &types.AKSK{}
	aksk.AK = os.Getenv(defaultAkEnv)
	aksk.SK = os.Getenv(defaultSkEnv)
	aksk.SecurityToken = os.Getenv(defaultSecurityTokenEnv)
	if pvd.Encrypt {
		aksk.SK = utils.AesDecrypt(aksk.SK, pvd.CipherKey)
		if aksk.SecurityToken != "" {
			aksk.SecurityToken = utils.AesDecrypt(aksk.SecurityToken, pvd.CipherKey)
		}
	}

	aksk.ExpiredAt = time.Now().Add(defaultExpiredAt)

	pvd.AkskMap.Delete("aksk")
	pvd.AkskMap.Store("aksk", &aksk)

	return aksk, nil
}
