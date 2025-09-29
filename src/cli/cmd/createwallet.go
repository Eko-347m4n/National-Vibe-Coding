package cmd

import (
	"fmt"

	"swatantra-node/src/core"
	"swatantra-node/src/wallet"
)

func init() {
	AddCommand(&Command{
		Name:        "createwallet",
		Description: "Buat wallet baru",
		Run:         createWallet,
		NeedsChain:  false, // Tidak butuh instance blockchain
	})
}

func createWallet(bc *core.Blockchain, args []string) {
	wallets, _ := wallet.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Alamat baru Anda: %s\n", address)
}
