package nmt

import (
	"math/rand"
	"testing"

	mock "github.com/0xBow-io/base-eas-asp/pkg/mock"
	pp "github.com/0xBow-io/base-eas-asp/pkg/privacy_pool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func gen_ngs(t *testing.T, groupSize int, recordSize int, withSort bool) (group *NsGroups) {
	// test group with namespace length 32 bytes
	group = newNsGroups(32)

	for _, ns := range GenRandomPublicIds(groupSize) {
		records, err := GenRandomRecords(ns, recordSize)
		require.NoError(t, err)

		for _, r := range records {
			nsStr, _, err := group.Add(Record(r[:]))
			require.NoError(t, err)
			require.Equal(t, ns.String(), "0x"+nsStr)
		}
	}
	require.Equal(t, groupSize*recordSize, group.Size())

	if withSort {
		group.Sort()
		require.NoError(t, group.ValidateOrder())
	}

	return group
}

func GenRandomPublicIds(n int) (ids []common.Hash) {
	for i := 0; i < n; i++ {
		ids = append(ids, common.BytesToHash(mock.GenRandomHash(32)))
	}
	return ids
}

func GenRandomRecord(publicID common.Hash) (e Record, err error) {
	event := pp.Event{
		TxHash:   common.BytesToHash(mock.GenRandBytes(rand.New(rand.NewSource(0)), 32)),
		LogIndex: uint8(rand.Uint32()),
		Token:    common.BytesToAddress(mock.GenRandBytes(rand.New(rand.NewSource(0)), 20)),
		From:     publicID,
		To:       common.BytesToHash(mock.GenRandBytes(rand.New(rand.NewSource(0)), 32)),
		Amount:   common.BytesToHash(mock.GenRandBytes(rand.New(rand.NewSource(0)), 32)),
	}
	if se, err := event.Serialize(); err != nil {
		return nil, err
	} else {
		// wrap the serialized event in a record
		return Record(se[:]), nil
	}
}

func GenRandomRecords(publicID common.Hash, n int) ([]Record, error) {
	var recs []Record
	for i := 0; i < n; i++ {
		rec, err := GenRandomRecord(publicID)
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}
