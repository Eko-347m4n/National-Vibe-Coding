package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"os"

	"swatantra-node/src/crypto"
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

// FileKeccak256 menghitung hash Keccak256 dari sebuah file.
func FileKeccak256(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hashBytes := crypto.Keccak256(data)
	return hex.EncodeToString(hashBytes), nil
}
