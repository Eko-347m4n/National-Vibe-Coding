package core

import (
	"bytes"
	"crypto/ed25519"
	"encoding/gob"
	"encoding/json"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	lua "github.com/yuin/gopher-lua"
	"swatantra-node/src/vm"
	"swatantra-node/src/wallet"
)

// Alamat kontrak oracle registry yang di-hardcode
const OracleRegistryAddress = "0000000000000000000000000000000000000000000000000000000000000001"

// ContractCallPayload adalah struktur data untuk payload pemanggilan fungsi kontrak
type ContractCallPayload struct {
	ContractAddress string   `json:"contract"`
	FunctionName    string   `json:"function"`
	Args            []string `json:"args"`
}

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

	bc := &Blockchain{nil, db}

	err = db.Update(func(txn *badger.Txn) error {
		// 1. Buat blok genesis
		pubKeyHash := wallet.DecodeAddress(address)
		cbtx := NewCoinbaseTX(pubKeyHash, "")
		target := big.NewInt(1)
		target.Lsh(target, uint(256-InitialTargetBits))
		genesis := NewBlock([]*Transaction{cbtx}, []byte{}, 0, target.Bytes(), time.Now().Unix())
		fmt.Println("Blok Genesis dibuat")
		if err := txn.Set(genesis.Hash, genesis.Serialize()); err != nil {
			return err
		}
		if err := txn.Set([]byte("lh"), genesis.Hash); err != nil {
			return err
		}
		bc.tip = genesis.Hash

		// 2. Terbitkan kontrak oracle_registry secara otomatis
		fmt.Println("Menerbitkan kontrak Oracle Registry...")
		oracleCode, err := os.ReadFile("oracle_registry.lua")
		if err != nil {
			return fmt.Errorf("gagal membaca oracle_registry.lua: %w", err)
		}
		oracleTxID, _ := hex.DecodeString(OracleRegistryAddress)
		oracleTx := Transaction{ID: oracleTxID, Type: TxContractCreation, Payload: oracleCode}
		bc.ProcessTransactions([]*Transaction{&oracleTx}, txn)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	// Reindex UTXO set setelah genesis
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
	err = db.View(func(txn *badger.Txn) error {
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

// InferenceRequestPayload adalah struktur data untuk payload permintaan inferensi
type InferenceRequestPayload struct {
	ModelName string `json:"model_name"`
	InputData string `json:"input_data"`
	Reward    int    `json:"reward"`
}

// NewInferenceRequestTransaction membuat transaksi untuk permintaan inferensi
func (bc *Blockchain) NewInferenceRequestTransaction(from, modelName, inputData string, reward int, UTXOSet *UTXOSet) (*Transaction, error) {
	wallets, err := wallet.NewWallets()
	if err != nil {
		return nil, err
	}
	w := wallets.GetWallet(from)
	if w.PublicKey == nil {
		return nil, fmt.Errorf("Wallet untuk alamat %s tidak ditemukan", from)
	}
	pubKeyHash := wallet.HashPubKey(w.PublicKey)

	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, reward)
	if acc < reward {
		return nil, fmt.Errorf("Error: Saldo tidak cukup untuk hadiah. Dibutuhkan: %d, Tersedia: %d", reward, acc)
	}

	var inputs []TXInput
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

	// Buat output kembalian jika ada. Tidak ada output tujuan (ini adalah burn).
	var outputs []TXOutput
	if acc > reward {
		outputs = append(outputs, TXOutput{acc - reward, pubKeyHash})
	}

	// Buat payload JSON
	payloadData := InferenceRequestPayload{
		ModelName: modelName,
		InputData: inputData,
		Reward:    reward,
	}
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat payload permintaan inferensi: %w", err)
	}

	tx := Transaction{nil, inputs, outputs, TxInferenceRequest, payloadBytes}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, w.PrivateKey)

	return &tx, nil
}

// NewStakeTransaction membuat transaksi untuk staking (membakar koin)
func (bc *Blockchain) NewStakeTransaction(from string, amount int, UTXOSet *UTXOSet) (*Transaction, error) {
	wallets, err := wallet.NewWallets()
	if err != nil {
		return nil, err
	}
	w := wallets.GetWallet(from)
	if w.PublicKey == nil {
		return nil, fmt.Errorf("Wallet untuk alamat %s tidak ditemukan", from)
	}
	pubKeyHash := wallet.HashPubKey(w.PublicKey)

	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)
	if acc < amount {
		return nil, fmt.Errorf("Error: Saldo tidak cukup untuk stake. Dibutuhkan: %d, Tersedia: %d", amount, acc)
	}

	var inputs []TXInput
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

	// Buat output kembalian jika ada. Tidak ada output tujuan (ini adalah burn).
	var outputs []TXOutput
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, pubKeyHash})
	}

	// Payload berisi jumlah yang di-stake
	payload := []byte(fmt.Sprintf("%d", amount))

	tx := Transaction{nil, inputs, outputs, TxStake, payload}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, w.PrivateKey)

	return &tx, nil
}


// NewContractCallTransaction membuat transaksi untuk memanggil fungsi pada kontrak
func (bc *Blockchain) NewContractCallTransaction(from, contractAddress, function, args string, UTXOSet *UTXOSet) (*Transaction, error) {
	const fee = 1 // Biaya sementara untuk memanggil kontrak

	wallets, err := wallet.NewWallets()
	if err != nil {
		return nil, err
	}
	w := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(w.PublicKey)

	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, fee)
	if acc < fee {
		return nil, errors.New("Error: Saldo tidak cukup untuk membayar biaya pemanggilan kontrak")
	}

	var inputs []TXInput
	var outputs []TXOutput

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

	// Buat output (hanya kembalian jika ada)
	if acc > fee {
		outputs = append(outputs, TXOutput{acc - fee, pubKeyHash})
	}

	// Buat payload JSON
	argSlice := strings.Split(args, "|")
	payloadData := ContractCallPayload{
		ContractAddress: contractAddress,
		FunctionName:    function,
		Args:            argSlice,
	}
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat payload kontrak: %w", err)
	}

	tx := Transaction{nil, inputs, outputs, TxContractCall, payloadBytes}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, w.PrivateKey)

	return &tx, nil
}


// NewContractCreationTransaction membuat transaksi untuk menerbitkan kontrak baru
func (bc *Blockchain) NewContractCreationTransaction(from string, code []byte, UTXOSet *UTXOSet) (*Transaction, error) {
	var inputs []TXInput
	var outputs []TXOutput
	const fee = 1 // Biaya sementara untuk menerbitkan kontrak

	wallets, err := wallet.NewWallets()
	if err != nil {
		return nil, err
	}
	w := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(w.PublicKey)

	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, fee)

	if acc < fee {
		return nil, errors.New("Error: Saldo tidak cukup untuk membayar biaya penerbitan kontrak")
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

	// Buat output (hanya kembalian jika ada)
	if acc > fee {
		outputs = append(outputs, TXOutput{acc - fee, pubKeyHash}) // kembalian
	}

	tx := Transaction{nil, inputs, outputs, TxContractCreation, code}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, w.PrivateKey)

	return &tx, nil
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

	tx := Transaction{nil, inputs, outputs, TxNormal, nil}
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

	// Buat blok baru dengan timestamp saat ini
	newTimestamp := time.Now().Unix()

	// Validasi Timestamp
	if newTimestamp <= parentBlock.Timestamp {
		log.Panic("Error: Timestamp blok baru harus lebih besar dari timestamp blok sebelumnya.")
	}
	// Ini adalah validasi sederhana, idealnya node tidak akan membuat blok dengan timestamp masa depan

	newBlock := NewBlock(transactions, lastHash, lastHeight+1, newTarget.Bytes(), newTimestamp)

	// Simpan blok baru dan proses transaksi di dalam satu DB transaction
	err = bc.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}

		err = txn.Set([]byte("lh"), newBlock.Hash)
		if err != nil {
			return err
		}
		bc.tip = newBlock.Hash

		// Perbarui UTXO set dan proses transaksi kontrak
		UTXOSet := UTXOSet{bc}
		bc.ProcessTransactions(newBlock.Transactions, txn)
		UTXOSet.Update(newBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// serializeState mengubah map state menjadi byte slice menggunakan gob
func serializeState(state map[string]string) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(state)
	return buffer.Bytes(), err
}

// deserializeState mengubah byte slice kembali menjadi map state
func deserializeState(data []byte) (map[string]string, error) {
	var state map[string]string
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&state)
	return state, err
}

// ProcessTransactions memproses setiap transaksi dalam satu blok, termasuk eksekusi smart contract
func (bc *Blockchain) ProcessTransactions(txs []*Transaction, txn *badger.Txn) {
	for _, tx := range txs {
		switch tx.Type {
		case TxStake:
			fmt.Println("Mendeteksi transaksi stake...")
			// Dapatkan alamat pengirim dari public key input pertama
			wallets, _ := wallet.NewWallets()
			stakerAddress := wallets.GetAddress(tx.Inputs[0].PubKey)
			stakeAmount := string(tx.Payload)

			// Panggil fungsi register_stake di kontrak oracle registry
			payloadData := ContractCallPayload{
				ContractAddress: OracleRegistryAddress,
				FunctionName:    "register_stake",
				Args:            []string{stakerAddress, stakeAmount},
			}
			bc.processContractCall(payloadData, txn)

		case TxContractCreation:
			fmt.Println("Mendeteksi transaksi pembuatan kontrak, mengeksekusi VM...")
			const creationGasLimit = 1000000 // Batas gas untuk pembuatan kontrak
			vmInstance := vm.NewVM(creationGasLimit)
			// Muat kode untuk mendefinisikan fungsi
			err := vmInstance.Execute(string(tx.Payload))
			if err != nil {
				vmInstance.Close()
				log.Panicf("Gagal mengeksekusi kontrak: %v", err)
			}

			// Panggil fungsi init() jika ada (by convention)
			initFn := vmInstance.L.GetGlobal("init")
			if initFn.Type() == lua.LTFunction {
				vmInstance.L.Push(initFn)
				if err := vmInstance.L.PCall(0, 0, nil); err != nil {
					vmInstance.Close()
					log.Panicf("Gagal memanggil fungsi init() pada kontrak: %v", err)
				}
			}

			// Dapatkan state setelah init, lalu tutup VM
			finalState := vmInstance.State
			vmInstance.Close()

			// Jika eksekusi berhasil, simpan state awal dan kode kontrak ke DB
			contractAddress := hex.EncodeToString(tx.ID)
			stateKey := []byte("contract:" + contractAddress)
			codeKey := []byte("contract_code:" + contractAddress)

			serializedState, err := serializeState(finalState)
			if err != nil {
				log.Panicf("Gagal serialisasi state kontrak: %v", err)
			}

			// Simpan state
			err = txn.Set(stateKey, serializedState)
			if err != nil {
				log.Panicf("Gagal menyimpan state kontrak ke DB: %v", err)
			}

			// Simpan kode
			err = txn.Set(codeKey, tx.Payload)
			if err != nil {
				log.Panicf("Gagal menyimpan kode kontrak ke DB: %v", err)
			}

			fmt.Printf("State dan kode untuk kontrak baru %s berhasil disimpan.\n", contractAddress)

		case TxContractCall:
			var payload ContractCallPayload
			if err := json.Unmarshal(tx.Payload, &payload); err != nil {
				log.Panicf("Gagal unmarshal payload: %v", err)
			}
			bc.processContractCall(payload, txn)
		case TxInferenceRequest:
			fmt.Println("Mendeteksi transaksi permintaan inferensi...")
			// Logika untuk menangani permintaan inferensi akan ditambahkan di sini
			// Misalnya, memvalidasi model, memilih oracle, dll.
			var payload InferenceRequestPayload
			if err := json.Unmarshal(tx.Payload, &payload); err != nil {
				log.Printf("Gagal unmarshal payload permintaan inferensi: %v", err)
				// Lanjutkan ke transaksi berikutnya jika payload tidak valid
				continue
			}
			fmt.Printf("Permintaan inferensi diterima untuk model '%s' dengan hadiah %d\n", payload.ModelName, payload.Reward)
		}
	}
}

func (bc *Blockchain) processContractCall(payload ContractCallPayload, txn *badger.Txn) {
	fmt.Println("Mengeksekusi pemanggilan kontrak...")

	const callGasLimit = 500000 // Batas gas untuk pemanggilan fungsi

	codeKey := []byte("contract_code:" + payload.ContractAddress)
	stateKey := []byte("contract:" + payload.ContractAddress)

	codeItem, err := txn.Get(codeKey)
	if err != nil {
		log.Panicf("Gagal mendapatkan kode kontrak %s: %v", payload.ContractAddress, err)
	}
	contractCode, _ := codeItem.ValueCopy(nil)

	stateItem, err := txn.Get(stateKey)
	if err != nil {
		log.Panicf("Gagal mendapatkan state kontrak %s: %v", payload.ContractAddress, err)
	}
	serializedState, _ := stateItem.ValueCopy(nil)
	currentState, err := deserializeState(serializedState)
	if err != nil {
		log.Panicf("Gagal deserialisasi state: %v", err)
	}

	vmInstance := vm.NewVM(callGasLimit)
	vmInstance.State = currentState
	defer vmInstance.Close()

	if err := vmInstance.Execute(string(contractCode)); err != nil {
		log.Panicf("Gagal memuat kode kontrak ke VM: %v", err)
	}

	// Logika pemanggilan fungsi yang benar
	fnToCall := vmInstance.L.GetGlobal(payload.FunctionName)
	vmInstance.L.Push(fnToCall)

	nArgs := 0
	for _, arg := range payload.Args {
		if arg != "" {
			vmInstance.L.Push(lua.LString(arg))
			nArgs++
		}
	}

	if err := vmInstance.L.PCall(nArgs, 1, nil); err != nil {
		log.Panicf("Gagal memanggil fungsi '%s': %v", payload.FunctionName, err)
	}

	// Ambil nilai kembali (jika ada)
	ret := vmInstance.L.Get(-1)
	vmInstance.L.Pop(1)

	// Simpan state baru
	newSerializedState, err := serializeState(vmInstance.State)
	if err != nil {
		log.Panicf("Gagal serialisasi state baru: %v", err)
	}
	if err := txn.Set(stateKey, newSerializedState); err != nil {
		log.Panicf("Gagal menyimpan state baru: %v", err)
	}

	fmt.Printf("Fungsi '%s' pada kontrak %s berhasil dieksekusi.\n", payload.FunctionName, payload.ContractAddress)
	if ret != lua.LNil {
		fmt.Printf("Nilai kembali: %s\n", ret.String())
	}
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