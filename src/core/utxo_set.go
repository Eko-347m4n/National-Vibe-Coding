package core

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/dgraph-io/badger/v4"
)

const utxoPrefix = "utxo-"

// UTXOSet merepresentasikan UTXO set.
type UTXOSet struct {
	Blockchain *Blockchain
}

// FindSpendableOutputs menemukan UTXO yang dapat dibelanjakan untuk sebuah alamat
func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.db

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(utxoPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			var v []byte
			err := item.Value(func(val []byte) error {
				v = val
				return nil
			})
			if err != nil {
				log.Panic(err)
			}

			k = bytes.TrimPrefix(k, prefix)
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)

			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

// FindUTXO menemukan semua UTXO untuk sebuah public key hash
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.db

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(utxoPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			var v []byte
			err := item.Value(func(val []byte) error {
				v = val
				return nil
			})
			if err != nil {
				log.Panic(err)
			}

			outs := DeserializeOutputs(v)

			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// Reindex membangun kembali UTXO set dari awal.
func (u UTXOSet) Reindex() {
	db := u.Blockchain.db
	prefix := []byte(utxoPrefix)

	// Hapus UTXO set lama
	err := db.DropPrefix(prefix)
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Blockchain.FindUTXO()

	err = db.Update(func(txn *badger.Txn) error {
		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}
			key = append(prefix, key...)
			err = txn.Set(key, outs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// Update memperbarui UTXO set dengan transaksi dari blok baru.
func (u UTXOSet) Update(block *Block) {
	db := u.Blockchain.db

	err := db.Update(func(txn *badger.Txn) error {
		for _, tx := range block.Transactions {
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					updatedOuts := TXOutputs{}
					inID := append([]byte(utxoPrefix), in.TxID...)
					item, err := txn.Get(inID)
					if err != nil {
						log.Panic(err)
					}
					var v []byte
					err = item.Value(func(val []byte) error {
						v = val
						return nil
					})
					if err != nil {
						log.Panic(err)
					}

					outs := DeserializeOutputs(v)

					for outIdx, out := range outs.Outputs {
						if outIdx != in.OutIndex {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}

					if len(updatedOuts.Outputs) == 0 {
						if err := txn.Delete(inID); err != nil {
							log.Panic(err)
						}
					} else {
						if err := txn.Set(inID, updatedOuts.Serialize()); err != nil {
							log.Panic(err)
						}
					}
				}
			}

			newOutputs := TXOutputs{}
			for _, out := range tx.Outputs {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}

			txID := append([]byte(utxoPrefix), tx.ID...)
			if err := txn.Set(txID, newOutputs.Serialize()); err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}