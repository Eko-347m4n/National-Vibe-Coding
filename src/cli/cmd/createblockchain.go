package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "createblockchain",
		Description: "Buat blockchain baru dan kirim hadiah genesis ke alamat",
		Run:         createBlockchain,
		NeedsChain:  false, // Perintah ini tidak butuh instance blockchain yang sudah ada
	})
}

func createBlockchain(bc *core.Blockchain, args []string) {
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createBlockchainAddress := createBlockchainCmd.String("address", "", "Alamat yang menerima hadiah genesis")

	err := createBlockchainCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *createBlockchainAddress == "" {
		createBlockchainCmd.Usage()
		return
	}

	// Perintah ini khusus, ia membuat blockchain, bukan menggunakan yang sudah ada.
	newBC := core.CreateBlockchain(*createBlockchainAddress)
	defer newBC.Database().Close()
	fmt.Println("Blockchain baru berhasil dibuat!")
}
