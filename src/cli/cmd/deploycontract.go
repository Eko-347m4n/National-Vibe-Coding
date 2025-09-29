package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "deploycontract",
		Description: "Terbitkan smart contract dari sebuah file .lua",
		Run:         deployContract,
		NeedsChain:  true,
	})
}

func deployContract(bc *core.Blockchain, args []string) {
	deployContractCmd := flag.NewFlagSet("deploycontract", flag.ExitOnError)
	deployContractFrom := deployContractCmd.String("from", "", "Alamat yang mendanai penerbitan kontrak")
	deployContractFile := deployContractCmd.String("file", "", "Path ke file .lua smart contract")

	err := deployContractCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *deployContractFrom == "" || *deployContractFile == "" {
		deployContractCmd.Usage()
		return
	}

	code, err := os.ReadFile(*deployContractFile)
	if err != nil {
		log.Panicf("Gagal membaca file kontrak '%s': %v", *deployContractFile, err)
	}

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewContractCreationTransaction(*deployContractFrom, code, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Kontrak dari file '%s' berhasil diterbitkan!\n", *deployContractFile)
}
