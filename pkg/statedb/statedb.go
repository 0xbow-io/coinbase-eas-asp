package statedb

import (
	nmt "github.com/0xBow-io/base-eas-asp/pkg/nmt"
)

type StateDB struct {
	// Contains event data stored as namespace records
	namespaceGroups nmt.NamespaceGroups
}

func NewStateDB() (*StateDB, error) {
	return nil, nil
}
