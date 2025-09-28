package core

import (
	"bytes"
	"math"
	"math/big"
	"swatantra-node/src/crypto"
	"swatantra-node/src/utils"
)

// ProofOfWork merepresentasikan proses pembuktian kerja

// EMA (Exponential Moving Average) Difficulty Adjustment constants
const (
	TargetBlockTime          = 15 // Detik
	DifficultyAdjustmentWindow = 10 // Jumlah blok untuk EMA
	InitialTargetBits        = 20 // Tingkat kesulitan awal (lebih mudah)
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork membuat instance PoW baru untuk sebuah blok
func NewProofOfWork(b *Block) *ProofOfWork {
	target := new(big.Int)
	target.SetBytes(b.Target)
	pow := &ProofOfWork{b, target}
	return pow
}

// prepareData menyiapkan data yang akan di-hash untuk proses mining
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(), // Menggunakan hash transaksi
			utils.IntToHex(pow.block.Timestamp),
			pow.block.Target, // Memasukkan target kesulitan saat ini
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
