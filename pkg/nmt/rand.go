package nmt

import (
	"testing"

	mock "github.com/0xBow-io/base-eas-asp/pkg/mock"
	sDB "github.com/0xBow-io/base-eas-asp/pkg/statedb"
	"github.com/stretchr/testify/require"
)

func gen_ngs(t *testing.T, groupSize int, recordSize int, withSort bool) (group *NsGroups, allRecords sDB.Records) {
	// test group with namespace length 32 bytes
	group = newNsGroups(32)

	for _, ns := range mock.GenRandomPublicIds(groupSize) {
		records, err := mock.GenRandomRecords(ns, recordSize)
		require.NoError(t, err)

		for _, r := range records {
			nsStr, _, err := group.Add(data(r[:]))
			require.NoError(t, err)
			require.Equal(t, ns.String(), "0x"+nsStr)
		}
		allRecords = append(allRecords, records...)
	}
	require.Equal(t, groupSize, group.Len())
	require.Equal(t, len(allRecords), group.Len()*recordSize)

	if withSort {
		group.Sort()
		require.NoError(t, group.ValidateOrder())
	}

	return group, allRecords
}
