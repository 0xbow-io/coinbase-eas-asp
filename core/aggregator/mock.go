package aggregator

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	reportDB "github.com/0xBow-io/base-eas-asp/pkg/reportDB"

	"github.com/ethereum/go-ethereum/common"
)

const (
	mock_secret            = "qnsec_dFzHeJ5iQbefXDH1akAKow=="
	mock_url_path          = "/webhook/2057b8e5-de11-4c65-8e39-92507d20de80"
	mock_nonce             = "632e2d63-d253-4a06-ab77-d565e806e5e1"
	mock_base_payload      = `{"matchedReceipts":[{"blockHash":"0x6644e74c3ff45ac8d514cdfc7393bcd193067c47c138842427a46a49d528573a","blockNumber":"0xb2bbad","contractAddress":"","cumulativeGasUsed":"0x70ef7","effectiveGasPrice":"0x187d3","from":"0x8844591d47f17bca6f5df8f6b64f4a739f1c0080","gasUsed":"0x450b5","logs":[{"address":"0x4200000000000000000000000000000000000021","blockHash":"0x6644e74c3ff45ac8d514cdfc7393bcd193067c47c138842427a46a49d528573a","blockNumber":"0xb2bbad","data":"0x8133f214f7bdaf516f03655db7406ba3c7945e5e4849238e43fdea6ef7a25cd4","logIndex":"0x1","removed":false,"topics":["0x8bf46bf4cfd674fa735a3d63ec1c9ad4153f033c290341f3a588b75685141b35","0x000000000000000000000000ff9418c67d18c8e067141bd77be43e32c4c3abe7","0x000000000000000000000000357458739f90461b99789350868cd7cf330dd7ee","0xf8b05c79f090979bf4a80270aba232dff11a10d9ca55c4f88de95317970f0de9"],"transactionHash":"0x63a3aef220d84947b12b0d771fa0df510eb753cab0c218efa1861b0d7c3d4567","transactionIndex":"0x4"},{"address":"0x2c7ee1e5f416dff40054c27a62f7b357c4e8619c","blockHash":"0x6644e74c3ff45ac8d514cdfc7393bcd193067c47c138842427a46a49d528573a","blockNumber":"0xb2bbad","data":"0x000000000000000000000000d867cbed445c37b0f95cc956fe6b539bdef7f32f","logIndex":"0x2","removed":false,"topics":["0x7fd54fcc14543b4db08cef4cd9fb23a6670c072d8a44cb0f1817d35b474176ca","0x000000000000000000000000ff9418c67d18c8e067141bd77be43e32c4c3abe7","0xf8b05c79f090979bf4a80270aba232dff11a10d9ca55c4f88de95317970f0de9","0x8133f214f7bdaf516f03655db7406ba3c7945e5e4849238e43fdea6ef7a25cd4"],"transactionHash":"0x63a3aef220d84947b12b0d771fa0df510eb753cab0c218efa1861b0d7c3d4567","transactionIndex":"0x4"}],"logsBloom":"0x00000000000000000000000040000000100000000000000000000000000000000001000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000008000000000000000000020040000000000000000000000000000000000000002000000000800000000000000000000000000000000400001000000000000000000000000010000000000000000000008010800000000020000000000010000000000000002000000000000000800000000000000000000000000000000004000000000000000000000000000800010000200000000000000000000000000000000000000000080000","status":"0x1","to":"0x357458739f90461b99789350868cd7cf330dd7ee","transactionHash":"0x63a3aef220d84947b12b0d771fa0df510eb753cab0c218efa1861b0d7c3d4567","transactionIndex":"0x4","type":"0x2"}],"matchedTransactions":[{"accessList":[],"blockHash":"0x6644e74c3ff45ac8d514cdfc7393bcd193067c47c138842427a46a49d528573a","blockNumber":"0xb2bbad","chainId":"0x2105","from":"0x8844591d47f17bca6f5df8f6b64f4a739f1c0080","gas":"0x927c0","gasPrice":"0x187d3","hash":"0x63a3aef220d84947b12b0d771fa0df510eb753cab0c218efa1861b0d7c3d4567","input":"0x56feed5e000000000000000000000000ff9418c67d18c8e067141bd77be43e32c4c3abe7","maxFeePerGas":"0x31214","maxPriorityFeePerGas":"0x186a0","nonce":"0x11211","r":"0xf015d09c9d8b274b703fad93020d38300aba7a9a0cbc786f05dac65dd0e0795b","s":"0x2d77121dc4ed30832b960e025b2b4b961e5103c8ad7e8f629a1eb381977c556c","to":"0x357458739f90461b99789350868cd7cf330dd7ee","transactionIndex":"0x4","type":"0x2","v":"0x0","value":"0x0"}]}`
	mock_base_topics       = `"topics":["0x8bf46bf4cfd674fa735a3d63ec1c9ad4153f033c290341f3a588b75685141b35","0x000000000000000000000000ff9418c67d18c8e067141bd77be43e32c4c3abe7","0x000000000000000000000000357458739f90461b99789350868cd7cf330dd7ee","0xf8b05c79f090979bf4a80270aba232dff11a10d9ca55c4f88de95317970f0de9"]`
	mock_base_publicID     = `0x000000000000000000000000ff9418c67d18c8e067141bd77be43e32c4c3abe7`
	mock_base_commitmentID = `0x8133f214f7bdaf516f03655db7406ba3c7945e5e4849238e43fdea6ef7a25cd4`
)

type MockReportFeed struct {
	reportPeriod time.Duration `default:"5s"`
}

func NewMockReportFeed() *MockReportFeed {
	return &MockReportFeed{}
}

func genPayloadHash(payload string) string {
	// payloadHash is the hash of the SHA256 hash of url_path + payload
	hash := sha256.Sum256([]byte(mock_url_path + payload))
	return fmt.Sprintf("%x", hash)
}

func genPayloadSig(timestamp string, payloadHash string) string {
	// Create a new HMAC hasher with SHA256 and the secret key
	h := hmac.New(sha256.New, []byte(mock_secret))

	// Generate hash of nonce + bodyHash + timestamp
	h.Write([]byte(mock_nonce + payloadHash + timestamp))

	// Compute the HMAC
	result := h.Sum(nil)

	// Encode the result to Base64
	return base64.StdEncoding.EncodeToString(result)
}

// generate random payload from mock_base_payload
// just changing the content of topics & log data
func genRandomPayload() (newPayload string, newPayloadHash string) {
	randomHashBytes := make([]byte, 32)
	rand.Read(randomHashBytes)

	newCommitmentID := common.BytesToHash(randomHashBytes).String()
	newPayload = strings.Replace(mock_base_payload, mock_base_commitmentID, newCommitmentID, -1) // Replace all occurrences of the old commitment ID with the new one

	randomAddrBytes := make([]byte, 20)
	rand.Read(randomAddrBytes)
	newPublicID := common.BytesToHash(randomAddrBytes).String()

	newPayload = strings.Replace(newPayload, mock_base_publicID, newPublicID, -1) // Replace all occurrences of the old attested Address with the new one

	return newPayload, genPayloadHash(newPayload)
}

// generate random report for testing
func (m *MockReportFeed) genRandReport(id string) (r reportDB.Report) {
	r.Header = reportDB.ReportHeader{
		NotificationID: id,
		Nonce:          mock_nonce,
		Timestamp:      time.Now().Format(reportDB.Header_Time_Layout),
	}

	// generate random payload from basePayload
	r.Body, r.Header.ContentHash = genRandomPayload()
	r.Header.Signature = genPayloadSig(r.Header.Timestamp, r.Header.ContentHash)

	return r
}

func (m *MockReportFeed) SubscribeTo(id string, feedChan chan<- reportDB.Report) error {
	go func() {
		for {
			fmt.Printf("Generating random report for id: %s\n", id)
			feedChan <- m.genRandReport(id)
			time.Sleep(m.reportPeriod)
		}
	}()
	return nil
}
