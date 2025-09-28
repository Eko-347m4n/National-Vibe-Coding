package core

import (
	"bytes"
	"math"
	"math/big"
	"swatantra-node/src/crypto"
	"swatantra-node/src/utils"
)

// TargetBits menentukan tingkat kesulitan. Angka yang lebih kecil berarti lebih sulit.
// Ini akan digantikan oleh algoritma kesulitan dinamis (EMA) nanti.
const TargetBits = 16

// ProofOfWork merepresentasikan proses pembuktian kerja
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork membuat instance PoW baru untuk sebuah blok
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// Geser bit ke kiri untuk menentukan target.
	// Lsh (Left Shift) sama dengan target * 2^(256 - TargetBits)
	target.Lsh(target, uint(256-TargetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

// prepareData menyiapkan data yang akan di-hash untuk proses mining
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			// TODO: Hash transaksi juga harus dimasukkan di sini
			utils.IntToHex(pow.block.Timestamp),
			utils.IntToHex(int64(TargetBits)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

// Run menjalankan loop mining untuk menemukan nonce yang valid
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash []byte
	nonce := 0

	// Loop tak terbatas sampai nonce yang valid ditemukan
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = crypto.Keccak256(data)

		hashInt.SetBytes(hash[:])

		// Cmp membandingkan hashInt dengan target.
		// Jika hashInt < target, maka nonce ditemukan.
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

// Validate memvalidasi apakah PoW sebuah blok valid
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := crypto.Keccak256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
