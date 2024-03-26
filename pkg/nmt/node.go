package nmt

import (
	"encoding/hex"
	"fmt"
)

// Min ID || Max ID || Hash
// ID is namespaceLen bytes
// Hash is ElementSize bytes
type Node []byte

func (n Node) String() string {
	return fmt.Sprintf("Min: %s Max: %s Hash: %s", n.MinNs(32).String(), n.MaxNs(32).String(), n.Hash(32).Hex())
}

func (n Node) Hex() string {
	return hex.EncodeToString(n)
}

func (n Node) MinNs(namespaceLen IDSize) ID {
	return ID(n[:namespaceLen])
}

func (n Node) MaxNs(namespaceLen IDSize) ID {
	return ID(n[namespaceLen : namespaceLen*2])
}

func (n Node) Hash(namespaceLen IDSize) Element {
	return Element(n[namespaceLen*2 : namespaceLen*2+ElementSize])
}

func (n Node) Equal(other Node) bool {
	return n.Hex() == other.Hex()
}

func NodeValueFromZero(namespaceLen IDSize, zero Element) (out Node) {
	out = make([]byte, namespaceLen*2+ElementSize)

	for i := 0; i < int(namespaceLen); i++ {
		out[i] = 0
		out[i+int(namespaceLen)] = 0
	}

	for i := 0; i < ElementSize; i++ {
		out[i+int(namespaceLen*2)] = zero[i]
	}
	return out
}

/*
Hash the data according to nmt specs:
  - hash = hash(data)
  - minNs = data.NID
  - maxNs = data.NID

Used to generate leaf nodes from set leaf data
*/
func DataToNode(namespaceLen IDSize, data data) (out Node) {
	out = make([]byte, namespaceLen*2+ElementSize)

	nsID := data.NID(namespaceLen)
	for i := 0; i < int(namespaceLen); i++ {
		out[i] = nsID[i]
		out[i+int(namespaceLen)] = nsID[i]
	}

	hash := data.Hash(namespaceLen)
	for i := 0; i < ElementSize; i++ {
		out[i+int(namespaceLen*2)] = hash[i]
	}
	return out
}

/*
BuildNode creates a new node from the given left and right nodes.
Unless one of the node is a zero node, the zeroSide parameter should be set to 0.
If the left node is a zero node, the zeroSide parameter should be set to 1.
If the right node is a zero node, the zeroSide parameter should be set to 2.
*/
func BuildNode(namespaceLen IDSize, left Node, right Node, zeroSide int, hashFn HashFunction) (n Node) {
	n = make([]byte, (namespaceLen*2)+ElementSize)

	nsOffset := int(namespaceLen)
	for i := 0; i < nsOffset; i++ {

		if zeroSide == 1 {
			// left node is a zero value
			n[i] = right[i] //minNs = rghtMinNs
		} else {
			n[i] = left[i] //minNs = leftMinNs
		}

		if zeroSide == 2 {
			// right node is a zero value
			n[i+nsOffset] = left[i+nsOffset] //maxNs = leftMaxNs
		} else {
			// otherwise
			n[i+nsOffset] = right[i+nsOffset] //maxNs = rightMaxNs
		}

	}
	hash := hashFn(left.Hash(namespaceLen), right.Hash(namespaceLen))

	for i := 0; i < ElementSize; i++ {
		n[nsOffset*2+i] = hash[i]
	}

	return n
}
