package utils

import (
	"bytes"
	"encoding/binary"
)

func UintToByte(num uint64) []byte {
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		return []byte{}
	}
	return buffer.Bytes()
}
