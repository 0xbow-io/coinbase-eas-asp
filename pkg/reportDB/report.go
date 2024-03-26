package reportdb

import (
	"encoding/json"
	"errors"
	"time"

	eas "github.com/0xBow-io/base-eas-asp/pkg/base_eas"
	quiknode "github.com/0xBow-io/base-eas-asp/pkg/quiknode"
)

const Header_Time_Layout string = "2006-01-02 15:04:05.999999999 -0700 MST"

type ReportHeader struct {
	NotificationID string `json:"x-qn-notification-id"`
	ContentHash    string `json:"x-qn-content-hash"`
	Nonce          string `json:"x-qn-nonce"`
	Signature      string `json:"x-qn-signature"`
	Timestamp      string `json:"x-qn-timestamp"`
}

type Report struct {
	Header ReportHeader `json:"header"`
	Body   string       `json:"body"` // payload body
}

func (r *Report) GetTimeStamp() int64 {
	parsedTime, err := time.Parse(Header_Time_Layout, r.Header.Timestamp)
	if err != nil {
		return 0
	}
	return parsedTime.Unix()
}

// get public IDs & commitments (attested wallet address) from report

func (r *Report) Parse() ([]eas.EAS, int64, error) {
	// do basic validation
	if r.GetTimeStamp() == 0 || r.Header.NotificationID == "" || r.Header.ContentHash == "" || r.Header.Nonce == "" || r.Header.Signature == "" {
		return nil, 0, errors.New("incorrect header")
	}

	if r.Body == "" {
		return nil, 0, errors.New("empty body")
	}

	ts := r.GetTimeStamp()
	if ts == 0 {
		return nil, ts, errors.New("failed to parse timestamp")
	}

	payload := new(quiknode.Payload)
	err := json.Unmarshal([]byte(r.Body), &payload)
	if err != nil {
		return nil,
			r.GetTimeStamp(), err
	}

	out, err := quiknode.ParsePayload(payload)
	return out, ts, err
}
