package nmt

import (
	"bytes"
	"encoding/hex"
	"math/big"
)

const ElementSize = 32 // internal hash value of a node is 32 bytes

type Element [ElementSize]byte

func (e Element) Hex() string {
	return hex.EncodeToString(e[:])
}

func (e Element) BigInt() *big.Int {
	return big.NewInt(0).SetBytes(e[:])
}

func (e Element) Eq(x Element) bool {
	return bytes.Equal(e[:], x[:])
}

func ToElement(in []byte) (e Element) {
	size := len(in)
	// add buffer to make sure we store 32 bytes
	pos := 0
	if buffer := ElementSize - size; buffer > 0 {
		for n := 0; n < buffer; n++ {
			e[pos] = 0
			pos++
		}
	}
	// then insert element
	for j := 0; j < size; j++ {
		e[pos] = in[j]
		pos++
	}
	return e
}
