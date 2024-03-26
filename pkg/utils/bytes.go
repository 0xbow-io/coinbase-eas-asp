package utils

func Ui32tob(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[3-i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}

func Ui32BtoUi32(b []byte) uint32 {
	var v uint32
	for i := uint32(0); i < 4; i++ {
		v |= uint32(b[i]) << (8 * (3 - i))
	}
	return v
}

func Force32(e []byte) []byte {
	if len(e) < 32 {
		return append(make([]byte, 32-len(e)), e...)
	}
	return e
}
