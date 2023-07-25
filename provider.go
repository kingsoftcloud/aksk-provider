package provider

import (
	"ezone.ksyun.com/code/kce/aksk-provider/types"
)

type AKSKProvider interface {
	GetAKSK() (*types.AKSK, error)
	ReloadAKSK() (*types.AKSK, error)
}
