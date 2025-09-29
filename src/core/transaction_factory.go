package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"swatantra-node/src/wallet"
)

// createTransactionBase adalah helper internal untuk menangani logika pembuatan transaksi yang berulang.
// Ini menangani: pengambilan wallet, pencarian UTXO, pembuatan input, dan pembuatan output kembalian.
func createTransactionBase(from string, amount int, UTXOSet *UTXOSet) ([]TXInput, []TXOutput, int, error) {
	wallets, err := wallet.NewWallets()
	if err != nil {
		return nil, nil, 0, err
	}
	w := wallets.GetWallet(from)
	if w.PrivateKey == nil {
		return nil, nil, 0, fmt.Errorf("tidak dapat menemukan wallet untuk alamat pengirim: %s", from)
	}
	pubKeyHash := wallet.HashPubKey(w.PublicKey)

	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		return nil, nil, 0, fmt.Errorf("saldo tidak cukup. Dibutuhkan: %d, Tersedia: %d", amount, acc)
	}

	var inputs []TXInput
	var outputs []TXOutput

	// Buat input dari UTXO yang valid
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			return nil, nil, 0, err
		}
		for _, out := range outs {
			input := TXInput{txID, out, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Buat output kembalian jika ada sisa
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, pubKeyHash})
	}

	return inputs, outputs, acc, nil
}

// NewUTXOTransaction membuat transaksi UTXO baru (transfer koin)
func (bc *Blockchain) NewUTXOTransaction(from, to string, amount int, UTXOSet *UTXOSet) (*Transaction, error) {
	inputs, outputs, _, err := createTransactionBase(from, amount, UTXOSet)
	if err != nil {
		return nil, err
	}

	// Buat output utama untuk penerima
	toPubKeyHash := wallet.DecodeAddress(to)
	outputs = append(outputs, TXOutput{amount, toPubKeyHash})

	tx := Transaction{nil, inputs, outputs, TxNormal, nil}
	tx.ID = tx.Hash()
	
	wallets, _ := wallet.NewWallets()
	w := wallets.GetWallet(from)
	bc.SignTransaction(&tx, w.PrivateKey)

	return &tx, nil
}

// NewContractCreationTransaction membuat transaksi untuk menerbitkan kontrak baru
func (bc *Blockchain) NewContractCreationTransaction(from string, code []byte, UTXOSet *UTXOSet) (*Transaction, error) {
	const fee = 1 // Biaya untuk menerbitkan kontrak

	inputs, outputs, _, err := createTransactionBase(from, fee, UTXOSet)
	if err != nil {
		return nil, err
	}

	tx := Transaction{nil, inputs, outputs, TxContractCreation, code}
	tx.ID = tx.Hash()

	wallets, _ := wallet.NewWallets()
	w := wallets.GetWallet(from)
	bc.SignTransaction(&tx, w.PrivateKey)

	return &tx, nil
}

// NewContractCallTransaction membuat transaksi untuk memanggil fungsi pada kontrak
func (bc *Blockchain) NewContractCallTransaction(from, contractAddress, function, args string, UTXOSet *UTXOSet) (*Transaction, error) {
	const fee = 1 // Biaya untuk memanggil kontrak

	inputs, outputs, _, err := createTransactionBase(from, fee, UTXOSet)
	if err != nil {
		return nil, err
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

	wallets, _ := wallet.NewWallets()
	w := wallets.GetWallet(from)
	bc.SignTransaction(&tx, w.PrivateKey)

	return &tx, nil
}

// NewStakeTransaction membuat transaksi untuk 'membakar' koin sebagai stake
func (bc *Blockchain) NewStakeTransaction(from string, amount int, UTXOSet *UTXOSet) (*Transaction, error) {
	inputs, outputs, _, err := createTransactionBase(from, amount, UTXOSet)
	if err != nil {
		return nil, err
	}

	// Tambahkan output yang di-stake (dibakar) ke alamat kosong
	outputs = append(outputs, TXOutput{amount, []byte{}})

	tx := Transaction{nil, inputs, outputs, TxStake, nil}
	tx.ID = tx.Hash()

	wallets, _ := wallet.NewWallets()
	w := wallets.GetWallet(from)
	bc.SignTransaction(&tx, w.PrivateKey)

	return &tx, nil
}

// NewInferenceRequestTransaction membuat transaksi untuk permintaan inferensi
func (bc *Blockchain) NewInferenceRequestTransaction(from, modelName, inputData string, reward int, UTXOSet *UTXOSet) (*Transaction, error) {
	const fee = 1 // Biaya untuk membuat permintaan inferensi

	inputs, outputs, _, err := createTransactionBase(from, fee, UTXOSet)
	if err != nil {
		return nil, err
	}

	// Buat payload JSON untuk permintaan inferensi
	payloadData := struct {
		ModelName string `json:"model_name"`
		InputData string `json:"input_data"`
		Reward    int    `json:"reward"`
	}{
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

	wallets, _ := wallet.NewWallets()
	w := wallets.GetWallet(from)
	bc.SignTransaction(&tx, w.PrivateKey)

	return &tx, nil
}
