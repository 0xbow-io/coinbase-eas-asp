package statedb

import (
	"math/big"

	"github.com/0xbow-io/go-iden3-crypto/poseidon"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

/*
TxHash : 32 bytes
LogIndex : 1 byte
Token : 20 bytes
From : 32 bytes
To : 32 bytes
Amount : 32 bytes
Total: 149 bytes
*/

const recordDataSize = 149
const namespaceLen = 32
const hashSize = 32

var (
	headerSize = namespaceLen + hashSize
)

type Records []Record

type Record [namespaceLen + hashSize + recordDataSize]byte

func (r Record) TxHash() common.Hash {
	return common.BytesToHash(r[headerSize : headerSize+32])
}

func (r Record) LogIndex() uint8 {
	return r[headerSize+32]
}

func (r Record) Token() common.Address {
	return common.BytesToAddress(r[headerSize+33 : headerSize+53])
}

func (r Record) From() common.Hash {
	return common.BytesToHash(r[headerSize+53 : headerSize+85])
}

func (r Record) To() common.Hash {
	return common.BytesToHash(r[headerSize+85 : headerSize+117])
}

func (r Record) Amount() common.Hash {
	return common.BytesToHash(r[headerSize+117 : headerSize+149])
}

func (r Record) AsEvent() (Event, error) {
	return Event{
		TxHash:   r.TxHash(),
		LogIndex: r.LogIndex(),
		Token:    r.Token(),
		From:     r.From(),
		To:       r.To(),
		Amount:   r.Amount(),
	}, nil
}

func (r Record) Ns() (ns [namespaceLen]byte) {
	copy(ns[:], r[:namespaceLen])
	return ns
}

func (r Record) NsHex() string {
	return common.Bytes2Hex(r[:namespaceLen])
}

func (r Record) Hash() []byte {
	return r[namespaceLen : namespaceLen+hashSize]
}

func (r Record) HashHex() string {
	return common.Bytes2Hex(r[namespaceLen : namespaceLen+hashSize])
}

func (r Record) ComputeHash() ([]byte, error) {
	hash, err := poseidon.Hash([]*big.Int{
		big.NewInt(0).SetBytes(r[headerSize : headerSize+32]),      // txHash
		big.NewInt(int64(r[headerSize+32])),                        // logIndex
		big.NewInt(0).SetBytes(r[headerSize+33 : headerSize+53]),   // token
		big.NewInt(0).SetBytes(r[headerSize+53 : headerSize+85]),   // from
		big.NewInt(0).SetBytes(r[headerSize+85 : headerSize+117]),  // to
		big.NewInt(0).SetBytes(r[headerSize+117 : headerSize+149]), // amount
	})

	if err != nil {
		return nil, err
	}

	bHash := hash.Bytes()
	if len(bHash) < 32 {
		pad := make([]byte, 32-len(bHash))
		bHash = append(pad, bHash...)
	} else if len(bHash) > 32 {
		return nil, errors.New("hash is greater than 32 bytes")
	}
	return bHash, nil
}
