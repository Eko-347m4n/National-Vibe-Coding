package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
	"swatantra-node/src/utils"
)

func init() {
	AddCommand(&Command{
		Name:        "publishmodel",
		Description: "Publikasikan model AI baru ke registry",
		Run:         publishModel,
		NeedsChain:  true,
	})
}

func publishModel(bc *core.Blockchain, args []string) {

publishModelCmd := flag.NewFlagSet("publishmodel", flag.ExitOnError)

publishModelFrom := publishModelCmd.String("from", "", "Alamat yang mendanai publikasi model")
publishModelContract := publishModelCmd.String("contract", "", "Alamat kontrak model registry")
publishModelName := publishModelCmd.String("name", "", "Nama unik untuk model")
publishModelFile := publishModelCmd.String("file", "", "Path ke file model")



	err := publishModelCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *publishModelFrom == "" || *publishModelContract == "" || *publishModelName == "" || *publishModelFile == "" {
		publishModelCmd.Usage()
		return
	}

	fmt.Printf("Menghitung hash untuk file model '%s'..\n", *publishModelFile)

	modelHash, err := utils.FileKeccak256(*publishModelFile)
	if err != nil {
		log.Panicf("Gagal menghitung hash file: %v", err)
	}
	fmt.Printf("Hash model: %s\n", modelHash)

	location := *publishModelFile // In a real system, this would be an IPFS CID or URL
	callArgs := fmt.Sprintf("%s|%s|%s", *publishModelName, modelHash, location)

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewContractCallTransaction(*publishModelFrom, *publishModelContract, "publish_model", callArgs, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Pemanggilan fungsi '%s' pada kontrak %s berhasil dikirim!\n", "publish_model", *publishModelContract)
}