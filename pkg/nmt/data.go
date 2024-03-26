package nmt

import (
	"encoding/hex"
	"fmt"
)

type data []byte

func (d data) String(namespaceLen IDSize) string {
	return fmt.Sprintf("NID: %s Hash: %s", d.NID(namespaceLen).String(), d.Hash(namespaceLen).Hex())
}

func (d data) NID(namespaceLen IDSize) ID {
	return ID(d[:namespaceLen])
}

func (d data) Hash(namespaceLen IDSize) Element {
	return Element(d[namespaceLen : namespaceLen+ElementSize])
}

func (d data) Data(namespaceLen IDSize) []byte {
	return d[namespaceLen+ElementSize:]
}

func (d data) Hex() string {
	return hex.EncodeToString(d)
}
