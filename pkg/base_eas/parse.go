package baseeas

import (
	sDB "github.com/0xBow-io/base-eas-asp/pkg/statedb"
	"github.com/ethereum/go-ethereum/common"
)

type EAS_TYPE string

const (
	EAS_REVOKE  = EAS_TYPE("revoke")
	EAS_ATTEST  = EAS_TYPE("attest")
	EAS_UNKNOWN = EAS_TYPE("unknown")
)

func EasTypeToMembership(t EAS_TYPE) sDB.MEMBERSHIP_TYPE {
	switch t {
	case EAS_REVOKE:
		return sDB.EXCLUSION
	case EAS_ATTEST:
		return sDB.INCLUSION
	default:
		return sDB.PARTIAL_INCLUSION
	}
}

type EAS struct {
	UUID    common.Hash `json:"uuid"`
	Account common.Hash `json:"address"`
	Type    EAS_TYPE    `json:"type"`
}
