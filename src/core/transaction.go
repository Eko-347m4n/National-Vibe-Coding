package core

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// TransactionType mendefinisikan jenis-jenis transaksi
type TransactionType int

const (
	TxNormal           TransactionType = iota // Transaksi transfer UTXO biasa
	TxContractCreation                        // Transaksi untuk membuat kontrak baru
	TxContractCall                            // Transaksi untuk memanggil fungsi di kontrak
	TxInferenceRequest                        // Transaksi untuk membuat permintaan inferensi
	TxStake                                   // Transaksi untuk staking (burn)
)

// Transaction merepresentasikan sebuah transaksi
type Transaction struct {
	ID      []byte
	Inputs  []TXInput
	Outputs []TXOutput
	Type    TransactionType // Jenis transaksi
	Payload []byte          // Data tambahan (misal: kode kontrak atau panggilan fungsi)
}

// TXOutputs adalah koleksi dari TXOutput, digunakan untuk serialisasi
type TXOutputs struct {
	Outputs []TXOutput
}

// TXInput merepresentasikan sebuah input transaksi
type TXInput struct {
	TxID      []byte
	OutIndex  int
	Signature []byte
	PubKey    []byte
}

// TXOutput merepresentasikan sebuah output transaksi
type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

// IsLockedWithKey memeriksa apakah output ini terkunci oleh sebuah public key hash
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// IsCoinbase memeriksa apakah transaksi ini adalah transaksi coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TxID) == 0 && tx.Inputs[0].OutIndex == -1
}

// Hash menghitung hash dari transaksi untuk ditandatangani
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// Serialize serializes the transaction
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// NewCoinbaseTX membuat transaksi coinbase baru (hadiah untuk miner)
func NewCoinbaseTX(pubKeyHash []byte, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%x'", pubKeyHash)
	}

	txInput := TXInput{[]byte{}, -1, nil, []byte(data)}
	txOutput := TXOutput{100, pubKeyHash}
	tx := Transaction{nil, []TXInput{txInput}, []TXOutput{txOutput}, TxNormal, nil}
	tx.ID = tx.Hash() // Menggunakan Hash() bukan SetID() lagi

	return &tx
}

// TrimmedCopy membuat salinan transaksi yang "dipangkas" untuk penandatanganan/verifikasi
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TXInput{in.TxID, in.OutIndex, nil, nil})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TXOutput{out.Value, out.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs, tx.Type, tx.Payload}

	return txCopy
}

// Sign menandatangani setiap input dari sebuah transaksi
func (tx *Transaction) Sign(privKey ed25519.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.TxID)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, in := range txCopy.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.TxID)]
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PubKey = prevTx.Outputs[in.OutIndex].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inID].PubKey = nil

		signature := ed25519.Sign(privKey, txCopy.ID)
		tx.Inputs[inID].Signature = signature
	}
}

// Verify memverifikasi tanda tangan dari setiap input transaksi
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.TxID)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, in := range tx.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.TxID)]
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PubKey = prevTx.Outputs[in.OutIndex].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inID].PubKey = nil

		if !ed25519.Verify(in.PubKey, txCopy.ID, in.Signature) {
			return false
		}
	}

	return true
}

// Serialize serializes TXOutputs
func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeOutputs deserializes TXOutputs
func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
