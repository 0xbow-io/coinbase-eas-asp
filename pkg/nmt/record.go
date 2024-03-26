package nmt

import (
	"encoding/hex"
	"fmt"
)

type Record []byte

func (r Record) String(namespaceLen IDSize) string {
	return fmt.Sprintf("NID: %s Hash: %s", r.NID(namespaceLen).String(), r.Hash(namespaceLen).Hex())
}

func (r Record) NID(namespaceLen IDSize) ID {
	return ID(r[:namespaceLen])
}

func (r Record) Hash(namespaceLen IDSize) Element {
	return Element(r[namespaceLen : namespaceLen+ElementSize])
}

func (r Record) Data(namespaceLen IDSize) []byte {
	return r[namespaceLen+ElementSize:]
}

func (r Record) Hex() string {
	return hex.EncodeToString(r)
}
