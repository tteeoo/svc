package svb

import (
	"bytes"
	"encoding/binary"
)

// UintToBytes converts a uint16 to two bytes.
func UintToBytes(u uint16) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, u)
	return buf.Bytes()
}

// BytesToUint converts two bytes to a uint16.
func BytesToUint(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}
