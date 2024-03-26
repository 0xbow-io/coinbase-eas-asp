package mock

import (
	"math/rand"
	mrand "math/rand"

	stateDB "github.com/0xBow-io/base-eas-asp/pkg/statedb"
	"github.com/ethereum/go-ethereum/common"
)

func GenRandBytes(random *mrand.Rand, n int) []byte {
	out := make([]byte, n)
	_, err := random.Read(out)
	if err != nil {
		panic(err)
	}
	return out
}

func GenRandomHash(size int) (out []byte) {
	out = make([]byte, size)
	m := rand.Intn(size)
	for i := size - 1; i > m; i-- {
		out[i] = byte(rand.Uint32())
	}
	return out

}

func GenRandomPublicIds(n int) (ids []common.Hash) {
	for i := 0; i < n; i++ {
		ids = append(ids, common.BytesToHash(GenRandomHash(32)))
	}
	return ids
}

func GenRandomRecord(publicID common.Hash) (e stateDB.Record, err error) {
	event := stateDB.Event{
		TxHash:   common.BytesToHash(GenRandBytes(rand.New(rand.NewSource(0)), 32)),
		LogIndex: uint8(rand.Uint32()),
		Token:    common.BytesToAddress(GenRandBytes(rand.New(rand.NewSource(0)), 20)),
		From:     publicID,
		To:       common.BytesToHash(GenRandBytes(rand.New(rand.NewSource(0)), 32)),
		Amount:   common.BytesToHash(GenRandBytes(rand.New(rand.NewSource(0)), 32)),
	}
	return event.AsRecord()
}

func GenRandomRecords(publicID common.Hash, n int) (stateDB.Records, error) {
	var recs stateDB.Records
	for i := 0; i < n; i++ {
		rec, err := GenRandomRecord(publicID)
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}
