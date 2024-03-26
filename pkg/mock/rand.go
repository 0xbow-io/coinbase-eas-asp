package mock

import (
	"math/rand"
	mrand "math/rand"
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
