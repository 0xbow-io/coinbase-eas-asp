package nmt

import (
	"bytes"
	"encoding/hex"
	"math"
	"math/big"
)

// IDSize is the number of bytes a namespace uses.
// Valid values are in [0,255].
type IDSize uint8

func (s IDSize) Size() int {
	return int(s)
}

// IDMaxSize defines the max. allowed namespace ID size in bytes.
const IDMaxSize = math.MaxUint8

type ID []byte

func (id ID) BigInt() *big.Int {
	return big.NewInt(0).SetBytes(id)
}

// Less returns true if nNS < other, otherwise, false.
func (id ID) Less(other ID) bool {
	return bytes.Compare(id, other) < 0
}

// Equal returns true if nNS == other, otherwise, false.
func (id ID) Equal(other ID) bool {
	return bytes.Equal(id, other)
}

// LessOrEqual returns true if nNS <= other, otherwise, false.
func (id ID) LessOrEqual(other ID) bool {
	return bytes.Compare(id, other) <= 0
}

// Size returns the byte size of the nNS.
func (id ID) Size() int {
	return len(id)
}

// String returns the hexadecimal encoding of the nNS. The output of
// nNS.String() is not equivalent to string(nNS).
func (id ID) String() string {
	return hex.EncodeToString(id)
}
