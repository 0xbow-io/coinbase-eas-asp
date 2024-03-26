package quiknode

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_ParsePayload(t *testing.T) {
	// load testPayload.json
	file, err := os.Open("testPayload.json")
	require.NoError(t, err)

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	require.NoError(t, err)

	payload := new(Payload)
	err = json.Unmarshal(byteValue, &payload)
	require.NoError(t, err)

	output, err := ParsePayload(payload)
	require.NoError(t, err)

	require.Equal(t, 1, len(output))
	require.Equal(t, output[0].UUID, common.HexToHash("0x8133f214f7bdaf516f03655db7406ba3c7945e5e4849238e43fdea6ef7a25cd4"))
	require.Equal(t, output[0].Account, common.HexToHash("0x000000000000000000000000ff9418c67d18c8e067141bd77be43e32c4c3abe7"))
}
