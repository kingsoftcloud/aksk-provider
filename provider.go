package provider

import (
	"github.com/kingsoftcloud/aksk-provider/types"
)

type AKSKProvider interface {
	GetAKSK() (*types.AKSK, error)
	ReloadAKSK() (*types.AKSK, error)
}
