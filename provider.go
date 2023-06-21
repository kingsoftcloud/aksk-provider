package provider

import (
	"newgit.op.ksyun.com/kce/aksk-provider/types"
)

type AKSKProvider interface {
	GetAKSK() (*types.AKSK, error)
	ReloadAKSK() (*types.AKSK, error)
}
