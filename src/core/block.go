package core

import (
	"time"
)

// block.go akan berisi definisi struktur Block dan logika terkait.
// Seperti header, transaksi, dan metode untuk validasi.

// Block merepresentasikan satu unit dalam blockchain
type Block struct {
	Timestamp     int64         // Waktu pembuatan blok
	PrevBlockHash []byte        // Hash dari blok sebelumnya
	Hash          []byte        // Hash dari blok saat ini (dihitung dari header)
	Transactions  []*Transaction // Daftar transaksi dalam blok
	Nonce         int           // Angka yang digunakan dalam proses mining (Proof-of-Work)
	Height        int           // Ketinggian blok dalam chain
}

// NewBlock membuat dan mengembalikan block baru.
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Transactions:  transactions,
		Height:        height,
	}
	// TODO: Implementasikan proses Proof-of-Work untuk mendapatkan Hash dan Nonce
	// pow := NewProofOfWork(block)
	// nonce, hash := pow.Run()
	// block.Hash = hash[:]
	// block.Nonce = nonce

	return block
}
