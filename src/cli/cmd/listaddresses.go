package cmd

import (
	"fmt"

	"swatantra-node/src/core"
	"swatantra-node/src/wallet"
)

func init() {
	AddCommand(&Command{
		Name:        "listaddresses",
		Description: "Tampilkan semua alamat wallet",
		Run:         listAddresses,
		NeedsChain:  false, // Tidak butuh instance blockchain
	})
}

func listAddresses(bc *core.Blockchain, args []string) {
	wallets, _ := wallet.NewWallets()
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
