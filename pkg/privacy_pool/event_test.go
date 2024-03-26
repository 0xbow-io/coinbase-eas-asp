package privacypool

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_Event_Record_Element(t *testing.T) {
	e := Event{
		TxHash:   common.HexToHash("0x0d95bebae9f1b39ccc72830e42411cf6cbb29c184cc8e67ecac5a678fb256045"),
		LogIndex: 109,
		Token:    common.HexToAddress("0x9fb9b8c43232fbe999b5b66c2050ea6c70353c96"),
		From:     common.HexToHash("0xeF4fB24aD0916217251F553c0596F8Edc630EB66"),
		To:       common.HexToHash("0x6D7A3177f3500BEA64914642a49D0B5C0a7Dae6D"),
		Amount:   common.HexToHash("0x000000000000000000000000000000000000000000000000000000000bbbc803"),
	}

	hash, err := e.Hash()
	require.NoError(t, err)

	// test event to record
	rec, err := e.Serialize()
	require.NoError(t, err)
	require.Equal(t, e.TxHash, rec.TxHash())
	require.Equal(t, e.LogIndex, rec.LogIndex())
	require.Equal(t, e.Token, rec.Token())
	require.Equal(t, e.From, rec.From())
	require.Equal(t, e.To, rec.To())
	require.Equal(t, e.Amount, rec.Amount())

	_, err = rec.AsEvent()
	require.NoError(t, err)

	recHash, err := rec.ComputeHash()
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString(hash), hex.EncodeToString(recHash))
	require.Equal(t, hex.EncodeToString(hash), hex.EncodeToString(rec.Hash()))

	require.Equal(t, e.From.Hex()[2:], rec.NsHex())

}
