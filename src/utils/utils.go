package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

// IntToHex mengubah int64 menjadi representasi byte slice
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
