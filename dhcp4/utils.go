package dhcp4

import (
	"bytes"
	"math/rand"
	"time"
)

func Uint16ToBytes(data uint16) []byte {
	var buf bytes.Buffer
	buf.Grow(2)
	buf.WriteByte(byte(data >> 8))
	buf.WriteByte(byte(data))
	return buf.Bytes()
}

func Uint32ToBytes(data uint32) []byte {
	var buf bytes.Buffer
	buf.Grow(4)
	buf.WriteByte(byte(data >> 24))
	buf.WriteByte(byte(data >> 16))
	buf.WriteByte(byte(data >> 8))
	buf.WriteByte(byte(data))
	return buf.Bytes()
}

func Uint8ToBytes(data uint8) []byte {
	return []byte{data}
}

func BytesToUint16(data []byte) uint16 {
	if len(data) != 2 {
		return uint16(data[0])
	}
	return uint16(data[0])<<8 | uint16(data[1])
}

func BytesToUint32(data []byte) uint32 {
	if len(data) < 4 {
		return 0
	}
	return uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
}

func RandomTransactionID() uint32 {
	rand.Seed(time.Now().Unix())
	return rand.Uint32()
}
