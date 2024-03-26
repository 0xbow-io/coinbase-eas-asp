package nmt

import (
	"crypto/sha256"
	"math/big"

	"github.com/0xbow-io/go-iden3-crypto/mimc7"
	"github.com/0xbow-io/go-iden3-crypto/poseidon"
)

type HashFunction func(left Element, right Element) Element

func SHA256Hash(left Element, right Element) Element {
	hash := sha256.New()
	sRight := right.Hex()
	sLeft := left.Hex()
	if sRight[0] == '0' && len(sRight) == 2 {
		sRight = sRight[1:]
	}
	if sLeft[0] == '0' && len(sLeft) == 2 {
		sLeft = sLeft[1:]
	}
	hash.Write([]byte(sLeft))
	hash.Write([]byte(sRight))
	return ToElement(hash.Sum(nil))
}

func Poseidon(left Element, right Element) Element {
	result, err := poseidon.Hash([]*big.Int{
		left.BigInt(),
		right.BigInt(),
	})

	if err != nil {
		panic(err.Error())
	}
	return ToElement(result.Bytes())
}

func Poseidon2(left Element, right Element) Element {
	result, err := poseidon.Poseidon2([]*big.Int{
		left.BigInt(),
		right.BigInt(),
	})

	if err != nil {
		panic(err.Error())
	}
	return ToElement(result.Bytes())
}

func MIMC7(left Element, right Element) Element {
	result := mimc7.MIMC7Hash(
		left.BigInt(),
		right.BigInt(),
	)
	return ToElement(result.Bytes())
}
