package aggregator

import (
	"testing"

	reportDB "github.com/0xBow-io/base-eas-asp/pkg/reportDB"
	"github.com/stretchr/testify/require"
)

func Test_Mock_GenSig(t *testing.T) {
	// hash original payload
	hash := genPayloadHash(mock_base_payload)
	require.Equal(t, "18d5c235669ab98e9d29d535ea8a587eb8ce107fe9215c0b0818994757adbbb1", hash)

	// generate signature
	sig := genPayloadSig("2024-03-12 04:04:14.11330824 +0000 UTC m=+17208.257632042", hash)
	require.Equal(t, "tsPXXWc53kXWob6ZbrHCxDSUrljtmk40d5vGGtFXvbs=", sig)
}

func Test_Mock_Aggregator(t *testing.T) {
	notifChan := make(chan string)
	rDB := reportDB.NewReportDB()
	rDB.SubscribeToNotif(notifChan)

	agg := NewReportAggregator(NewMockReportFeed(), rDB)
	for i := 0; i < 10; i++ {
		publicID := <-notifChan
		r := agg.GetLatestReportForPublicID(publicID)
		out, _, err := r.Parse()
		require.NoError(t, err)
		require.NotEmpty(t, out)
		require.Equal(t, out[0].Account.String(), publicID)
	}
}
