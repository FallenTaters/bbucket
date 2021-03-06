package bbucket

import "encoding/binary"

func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func Btoi(v []byte) int {
	return int(binary.BigEndian.Uint64(v))
}
