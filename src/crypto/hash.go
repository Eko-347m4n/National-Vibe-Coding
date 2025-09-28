package crypto

import (
	"golang.org/x/crypto/sha3"
)

// hash.go akan mengimplementasikan fungsi hashing yang dibutuhkan (Keccak256).

// Keccak256 menghitung dan mengembalikan hash Keccak-256 dari data.
func Keccak256(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}
