package provider

import (
	"newgit.op.ksyun.com/kce/aksk-provider/configmap"
	"newgit.op.ksyun.com/kce/aksk-provider/env"
	"newgit.op.ksyun.com/kce/aksk-provider/secret"
)

func SecretAKSKProvider(filePath, cipherKey string) *secret.SecretAKSKProvider {
	return secret.NewSecretAKSKProvider(filePath, cipherKey)
}

func CMAKSKProvider(filePath string) *configmap.CMAKSKProvider {
	return configmap.NewCMAKSKProvider(filePath)
}

func EnvAKSKProvider(encrypt bool, cipherKey string) *env.EnvAKSKProvider {
	return env.NewEnvAKSKProvider(encrypt, cipherKey)
}
