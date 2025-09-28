package core

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// block.go akan berisi definisi struktur Block dan logika terkait.
// Seperti header, transaksi, dan metode untuk validasi.

// Block merepresentasikan satu unit dalam blockchain
type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Transactions  []*Transaction
	Nonce         int
	Height        int
}

// NewBlock membuat dan mengembalikan block baru.
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Transactions:  transactions,
		Height:        height,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

// Serialize mengubah Block menjadi slice of bytes
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock mengubah slice of bytes kembali menjadi Block
func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}