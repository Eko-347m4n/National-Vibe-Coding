
package core

import (
	"fmt"
	"os"
	"testing"

	"swatantra-node/src/wallet"
)

// setupTest an teardownTest digunakan untuk memastikan lingkungan tes bersih.
func setupTest() {
	// Hapus database blockchain lama jika ada
	os.RemoveAll("./tmp/blocks")
	os.RemoveAll("./tmp/utxo")
}

func teardownTest() {
	// Hapus database setelah tes selesai
	os.RemoveAll("./tmp/blocks")
	os.RemoveAll("./tmp/utxo")
}

// TestBlockchainLifecycle adalah tes integrasi untuk alur kerja dasar blockchain.
func TestBlockchainLifecycle(t *testing.T) {
	setupTest()
	defer teardownTest()

	// 1. Buat wallet dan alamat
	wallets, err := wallet.NewWallets()
	if err != nil {
		t.Fatalf("Gagal membuat wallets: %v", err)
	}
	address := wallets.CreateWallet()
	wallets.SaveToFile() // Simpan wallet untuk digunakan nanti

	// 2. Buat blockchain baru, hadiah genesis dikirim ke alamat yang baru dibuat
	bc := CreateBlockchain(address)
	defer bc.Database().Close()

	// Verifikasi bahwa genesis block ada
	if bc.tip == nil {
		t.Fatal("Blockchain gagal dibuat, tip is nil")
	}
	t.Logf("Blockchain berhasil dibuat, genesis dikirim ke: %s", address)

	// 3. Periksa saldo awal
	UTXOSet := UTXOSet{bc}
	balance := 0
	pubKeyHash := wallet.DecodeAddress(address)
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)
	for _, out := range UTXOs {
		balance += out.Value
	}

	// Hadiah genesis default adalah 100 (dari NewCoinbaseTX)
	expectedBalance := 100
	if balance != expectedBalance {
		t.Fatalf("Saldo awal tidak sesuai. Harusnya %d, tapi dapat %d", expectedBalance, balance)
	}
	t.Logf("Saldo awal setelah genesis block: %d", balance)

	// 4. Buat transaksi baru
	wallets, _ = wallet.NewWallets() // Muat ulang wallets
	w := wallets.GetWallet(address)
	if w.PublicKey == nil {
		t.Fatalf("Gagal mendapatkan wallet untuk alamat: %s", address)
	}

	// Buat alamat penerima baru
	recipientAddress := wallets.CreateWallet()
	wallets.SaveToFile()

	amountToSend := 10
	tx, err := bc.NewUTXOTransaction(address, recipientAddress, amountToSend, &UTXOSet)
	if err != nil {
		t.Fatalf("Gagal membuat transaksi: %v", err)
	}
	t.Log("Transaksi baru berhasil dibuat.")

	// 5. Tambahkan blok baru dengan transaksi tersebut
	bc.AddBlock([]*Transaction{tx})
	t.Log("Blok baru berhasil ditambahkan.")

	// 6. Verifikasi saldo setelah transaksi
	// Saldo pengirim harus berkurang
	balanceFrom := 0
	pubKeyHashFrom := wallet.DecodeAddress(address)
	UTXOsFrom := UTXOSet.FindUTXO(pubKeyHashFrom)
	for _, out := range UTXOsFrom {
		balanceFrom += out.Value
	}
	if balanceFrom != expectedBalance-amountToSend {
		t.Fatalf("Saldo pengirim tidak sesuai. Harusnya %d, tapi dapat %d", expectedBalance-amountToSend, balanceFrom)
	}
	t.Logf("Saldo pengirim setelah transaksi: %d", balanceFrom)

	// Saldo penerima harus bertambah
	balanceTo := 0
	pubKeyHashTo := wallet.DecodeAddress(recipientAddress)
	UTXOsTo := UTXOSet.FindUTXO(pubKeyHashTo)
	for _, out := range UTXOsTo {
		balanceTo += out.Value
	}
	if balanceTo != amountToSend {
		t.Fatalf("Saldo penerima tidak sesuai. Harusnya %d, tapi dapat %d", amountToSend, balanceTo)
	}
	t.Logf("Saldo penerima setelah transaksi: %d", balanceTo)

	// 7. Verifikasi PoW dari blok baru
	bci := bc.Iterator()
	newBlock := bci.Next() // Blok terbaru
	pow := NewProofOfWork(newBlock)
	if !pow.Validate() {
		t.Fatal("Validasi Proof of Work untuk blok baru gagal.")
	}
	t.Log("Validasi Proof of Work untuk blok baru berhasil.")

	// Verifikasi kesulitan dinamis (minimal, periksa bahwa target tidak nol)
	if len(newBlock.Target) == 0 {
		t.Fatal("Target kesulitan di blok baru kosong.")
	}
	t.Logf("Target kesulitan blok baru (hex): %x", newBlock.Target)

	fmt.Println("Tes siklus hidup blockchain berhasil diselesaikan.")
}
