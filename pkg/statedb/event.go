package statedb

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

func (e *Event) AsRecord() (Record, error) {
	var s Record

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
