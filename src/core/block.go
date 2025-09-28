package core

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"swatantra-node/src/crypto"
)

// Block merepresentasikan satu unit dalam blockchain
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Height        int
	Target        []byte // Target kesulitan untuk blok ini
}

// HashTransactions mengembalikan hash dari transaksi-transaksi di dalam blok
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.ID)
	}
	// Di masa depan, ini akan menjadi Merkle Tree
	hash := crypto.Keccak256(bytes.Join(transactions, []byte{}))

	return hash
}

// NewBlock membuat dan mengembalikan sebuah Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int, target []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0, height, target}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:];
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