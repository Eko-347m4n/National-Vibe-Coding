package wallet

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

const walletFile = "./tmp/wallets.dat"

// Wallets menyimpan koleksi Wallet
type Wallets struct {
	Wallets map[string]*Wallet
}

// NewWallets membuat Wallets dan memuatnya dari file jika ada
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile()

	return &wallets, err
}

// CreateWallet membuat dan menambahkan Wallet baru ke koleksi
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

// GetAddresses mengembalikan semua alamat dari koleksi wallet
func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetWallet mengembalikan sebuah Wallet berdasarkan alamat
func (ws Wallets) GetWallet(address string) Wallet {
	wallet := ws.Wallets[address]
	if wallet == nil {
		return Wallet{} // Kembalikan wallet kosong jika tidak ditemukan
	}
	return *wallet
}

// GetAddress mengembalikan alamat untuk sebuah public key
func (ws Wallets) GetAddress(pubKey []byte) string {
	w := Wallet{PublicKey: pubKey}
	return string(w.Address())
}

// LoadFromFile memuat wallets dari file
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return nil
	}

	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}

// SaveToFile menyimpan wallets ke dalam file
func (ws *Wallets) SaveToFile() {
	if err := os.MkdirAll("./tmp", 0755); err != nil {
		log.Panic(err)
	}
	var content bytes.Buffer

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
