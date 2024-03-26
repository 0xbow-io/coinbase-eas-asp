package privacypool

import (
	"math/big"

	"github.com/0xbow-io/go-iden3-crypto/poseidon"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type Event struct {
	TxHash   common.Hash `json:"txHash"`
	LogIndex uint8       `json:"logIndex"`

	Token common.Address `json:"token"`

	From   common.Hash `json:"from"`
	To     common.Hash `json:"to"`
	Amount common.Hash `json:"amount"`
}

func (e *Event) Hash() ([]byte, error) {
	hash, err := poseidon.Hash([]*big.Int{
		e.TxHash.Big(),
		big.NewInt(int64(e.LogIndex)),
		big.NewInt(0).SetBytes(e.Token.Bytes()),
		e.From.Big(),
		e.To.Big(),
		e.Amount.Big(),
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

func (e *Event) Serialize() (SerialEvent, error) {
	var s SerialEvent

	hash, err := e.Hash()
	if err != nil {
		return s, err
	}

	// pouplate header
	copy(s[:namespaceLen], e.From[:])         // namespace
	copy(s[namespaceLen:headerSize], hash[:]) // hash

	// populate body
	copy(s[headerSize:headerSize+32], e.TxHash[:])
	s[headerSize+32] = e.LogIndex
	copy(s[headerSize+33:headerSize+53], e.Token[:])
	copy(s[headerSize+53:headerSize+85], e.From[:])
	copy(s[headerSize+85:headerSize+117], e.To[:])
	copy(s[headerSize+117:headerSize+149], e.Amount[:])
	return s, nil
}

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

type SerialEvents []SerialEvent

type SerialEvent [namespaceLen + hashSize + recordDataSize]byte

func (se SerialEvent) TxHash() common.Hash {
	return common.BytesToHash(se[headerSize : headerSize+32])
}

func (se SerialEvent) LogIndex() uint8 {
	return se[headerSize+32]
}

func (se SerialEvent) Token() common.Address {
	return common.BytesToAddress(se[headerSize+33 : headerSize+53])
}

func (se SerialEvent) From() common.Hash {
	return common.BytesToHash(se[headerSize+53 : headerSize+85])
}

func (se SerialEvent) To() common.Hash {
	return common.BytesToHash(se[headerSize+85 : headerSize+117])
}

func (se SerialEvent) Amount() common.Hash {
	return common.BytesToHash(se[headerSize+117 : headerSize+149])
}

func (se SerialEvent) AsEvent() (Event, error) {
	return Event{
		TxHash:   se.TxHash(),
		LogIndex: se.LogIndex(),
		Token:    se.Token(),
		From:     se.From(),
		To:       se.To(),
		Amount:   se.Amount(),
	}, nil
}

func (se SerialEvent) Ns() (ns [namespaceLen]byte) {
	copy(ns[:], se[:namespaceLen])
	return ns
}

func (se SerialEvent) NsHex() string {
	return common.Bytes2Hex(se[:namespaceLen])
}

func (se SerialEvent) Hash() []byte {
	return se[namespaceLen : namespaceLen+hashSize]
}

func (se SerialEvent) HashHex() string {
	return common.Bytes2Hex(se[namespaceLen : namespaceLen+hashSize])
}

func (se SerialEvent) ComputeHash() ([]byte, error) {
	hash, err := poseidon.Hash([]*big.Int{
		big.NewInt(0).SetBytes(se[headerSize : headerSize+32]),      // txHash
		big.NewInt(int64(se[headerSize+32])),                        // logIndex
		big.NewInt(0).SetBytes(se[headerSize+33 : headerSize+53]),   // token
		big.NewInt(0).SetBytes(se[headerSize+53 : headerSize+85]),   // from
		big.NewInt(0).SetBytes(se[headerSize+85 : headerSize+117]),  // to
		big.NewInt(0).SetBytes(se[headerSize+117 : headerSize+149]), // amount
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
