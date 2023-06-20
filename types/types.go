package types

import (
	"time"
)

type AKSK struct {
	AK            string    `json:"ak"`
	SK            string    `json:"sk"`
	SecurityToken string    `json:"securityToken"`
	ExpiredAt     time.Time `json:"expired_at"`
	Cipher        string    `json:"cipher"`
}
