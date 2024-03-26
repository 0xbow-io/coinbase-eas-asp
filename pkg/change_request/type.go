package changeRequest

import (
	poa "github.com/0xBow-io/base-eas-asp/pkg/proof-of-audit"
	stateDB "github.com/0xBow-io/base-eas-asp/pkg/statedb"
)

type ChangeRequest struct {
	Ns         []byte                  `json:"nameSpace"`
	Membership stateDB.MEMBERSHIP_TYPE `json:"membership"`
	Proof      poa.SP1Proof            `json:"sp1Proof"`
}
