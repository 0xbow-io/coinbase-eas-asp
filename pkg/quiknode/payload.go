package quiknode

import (
	eas "github.com/0xBow-io/base-eas-asp/pkg/base_eas"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

/*
Quiknode's QuickAlert Receipt + matched transactions structure
*/

type TxJSON struct {
	Type hexutil.Uint64 `json:"type"`

	ChainID              *hexutil.Big      `json:"chainId,omitempty"`
	Nonce                *hexutil.Uint64   `json:"nonce"`
	To                   *common.Address   `json:"to"`
	Gas                  *hexutil.Uint64   `json:"gas"`
	GasPrice             *hexutil.Big      `json:"gasPrice"`
	MaxPriorityFeePerGas *hexutil.Big      `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big      `json:"maxFeePerGas"`
	Value                *hexutil.Big      `json:"value"`
	Input                *hexutil.Bytes    `json:"input"`
	AccessList           *types.AccessList `json:"accessList,omitempty"`
	V                    *hexutil.Big      `json:"v"`
	R                    *hexutil.Big      `json:"r"`
	S                    *hexutil.Big      `json:"s"`

	// Only used for encoding:
	Hash common.Hash `json:"hash"`
}

type Receipt struct {
	TxHash common.Hash `json:"transactionHash" gencodec:"required"`

	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      string      `json:"blockNumber,omitempty"`
	TransactionIndex string      `json:"transactionIndex"`

	Status string       `json:"status"`
	Logs   []*types.Log `json:"logs"              gencodec:"required"`
}

type Payload struct {
	MatchedReceipts     []Receipt `json:"matchedReceipts"`
	MatchedTransactions []TxJSON  `json:"matchedTransactions"`
}

func ParsePayload(p *Payload) ([]eas.EAS, error) {
	// iterate through the payload
	// work out the attested addresses based on the logs
	// return the attested addresses
	var output []eas.EAS
	var easType eas.EAS_TYPE
	for _, receipt := range p.MatchedReceipts {
		for _, log := range receipt.Logs {
			if len(log.Topics) == 4 {

				switch log.Topics[0] {
				case common.HexToHash(eas.COINBASE_EAS_ATTEST_TOPIC):
					easType = eas.EAS_ATTEST
				case common.HexToHash(eas.COINBASE_EAS_REVOKE_TOPIC):
					easType = eas.EAS_REVOKE
				default:
					easType = eas.EAS_UNKNOWN
				}

				if easType != eas.EAS_UNKNOWN &&
					log.Topics[2] == common.HexToHash(eas.COINBASE_EAS_HASH) &&
					log.Topics[3] == common.HexToHash(eas.COINBASE_EAS_SCHEMA_ID) {
					output = append(output, eas.EAS{
						UUID:    common.Hash(log.Data),
						Account: log.Topics[1],
						Type:    easType,
					})
				}
			}
		}
	}
	return output, nil
}
