package core

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	"swatantra-node/src/wallet"
)

const (
	dbPath = "./tmp/blocks"
)

// Blockchain merepresentasikan rantai blok berbasis database
type Blockchain struct {
	tip []byte
	db  *badger.DB
}

func (bc *Blockchain) Database() *badger.DB {
	return bc.db
}

// Iterator mengembalikan BlockchainIterator untuk traversal
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

// CreateBlockchain membuat blockchain baru dengan blok genesis
func CreateBlockchain(address string) *Blockchain {
	if DbExists() {
		fmt.Println("Blockchain sudah ada.")
		os.Exit(1)
	}
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte
	err = db.Update(func(txn *badger.Txn) error {
		pubKeyHash := wallet.DecodeAddress(address)
		cbtx := NewCoinbaseTX(pubKeyHash, "")

		// Buat target kesulitan awal untuk genesis block
		target := big.NewInt(1)
		target.Lsh(target, uint(256-InitialTargetBits))

		genesis := NewBlock([]*Transaction{cbtx}, []byte{}, 0, target.Bytes())

		fmt.Println("Blok Genesis dibuat")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = txn.Set([]byte("lh"), genesis.Hash)
		tip = genesis.Hash
		return err
	})
	if err != nil {
		log.Panic(err)
	}

	bc := &Blockchain{tip, db}
	// Reindex UTXO set
	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	return bc
}

// NewBlockchain membuka blockchain yang sudah ada
func NewBlockchain() *Blockchain {
	if !DbExists() {
		fmt.Println("Tidak ada blockchain ditemukan. Buat dulu dengan 'createblockchain'.")
		os.Exit(1)
	}
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			tip = val
			return nil
		})
		return err
	})
	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}

// DbExists memeriksa apakah database blockchain sudah ada
func DbExists() bool {
	if _, err := os.Stat(dbPath + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}
	return true
}

// FindUTXO menemukan semua output transaksi yang belum dibelanjakan
func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					inTxID := hex.EncodeToString(in.TxID)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.OutIndex)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return UTXO
}

// FindTransaction menemukan sebuah transaksi berdasarkan ID
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

// NewUTXOTransaction membuat transaksi UTXO baru
func (bc *Blockchain) NewUTXOTransaction(from, to string, amount int, UTXOSet *UTXOSet) (*Transaction, error) {
	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := wallet.NewWallets()
	if err != nil {
		return nil, err
	}
	w := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(w.PublicKey)

	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		return nil, errors.New("Error: Not enough funds")
	}

	// Buat input
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			return nil, err
		}

		for _, out := range outs {
			input := TXInput{txID, out, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Buat output
	toPubKeyHash := wallet.DecodeAddress(to)
	outputs = append(outputs, TXOutput{amount, toPubKeyHash})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc-amount, pubKeyHash}) // kembalian
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, w.PrivateKey)

	return &tx, nil
}

// SignTransaction menandatangani input dari sebuah Transaction
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ed25519.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.TxID)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte
	var lastHeight int
	var parentBlock *Block

	err := bc.db.View(func(txn *badger.Txn) error {
		// Ambil hash dari blok terakhir
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		if err != nil {
			return err
		}

		// Ambil blok terakhir itu sendiri
		item, err = txn.Get(lastHash)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			parentBlock = DeserializeBlock(val)
			lastHeight = parentBlock.Height
			return nil
		})
		return err
	})
	if err != nil {
		log.Panic(err)
	}

	// Hitung target kesulitan baru menggunakan EMA
	newTarget := calculateNewTarget(parentBlock)

	// Buat blok baru
	newBlock := NewBlock(transactions, lastHash, lastHeight+1, newTarget.Bytes())

	// Simpan blok baru ke DB
	err = bc.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}

		// Perbarui pointer 'lh' (last hash)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		bc.tip = newBlock.Hash
		return err
	})

	if err != nil {
		log.Panic(err)
	}

	// Perbarui UTXO set
	UTXOSet := UTXOSet{bc}
	UTXOSet.Update(newBlock)
}

// calculateNewTarget menghitung target kesulitan berikutnya berdasarkan blok sebelumnya
func calculateNewTarget(parentBlock *Block) *big.Int {
	parentTarget := new(big.Int).SetBytes(parentBlock.Target)

	// Waktu pembuatan blok aktual
	actualBlockTime := time.Now().Unix() - parentBlock.Timestamp
	// Batasi agar tidak terlalu cepat atau lambat untuk mencegah fluktuasi ekstrem
	if actualBlockTime < TargetBlockTime/4 {
		actualBlockTime = TargetBlockTime / 4
	}
	if actualBlockTime > TargetBlockTime*4 {
		actualBlockTime = TargetBlockTime * 4
	}

	// Rumus EMA sederhana untuk menyesuaikan target:
	// newTarget = ( (N-1)*parentTarget + parentTarget * actual_time / target_time ) / N
	window := big.NewInt(DifficultyAdjustmentWindow)
	targetTime := big.NewInt(TargetBlockTime)

	// Term 1: (N-1) * parentTarget
	term1 := new(big.Int).Mul(parentTarget, new(big.Int).Sub(window, big.NewInt(1)))

	// Term 2: parentTarget * actual_time / target_time
	term2_num := new(big.Int).Mul(parentTarget, big.NewInt(actualBlockTime))
	term2 := new(big.Int).Div(term2_num, targetTime)

	// (Term 1 + Term 2) / N
	newTargetNum := new(big.Int).Add(term1, term2)
	newTarget := new(big.Int).Div(newTargetNum, window)

	return newTarget
}

// --- BlockchainIterator ---

type BlockchainIterator struct {
	currentHash []byte
	db          *badger.DB
}

// Next mengembalikan blok berikutnya dalam iterasi (dari tip ke genesis)
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(i.currentHash)
		if err != nil {
			log.Panic(err)
		}
		var encodedBlock []byte
		err = item.Value(func(val []byte) error {
			encodedBlock = val
			return nil
		})
		block = DeserializeBlock(encodedBlock)
		return err
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}