package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "callcontract",
		Description: "Panggil sebuah fungsi pada smart contract yang sudah ada",
		Run:         callContract,
		NeedsChain:  true,
	})
}

func callContract(bc *core.Blockchain, args []string) {
	callContractCmd := flag.NewFlagSet("callcontract", flag.ExitOnError)
	callContractFrom := callContractCmd.String("from", "", "Alamat yang memanggil kontrak")
	callContractAddress := callContractCmd.String("contract", "", "Alamat smart contract yang akan dipanggil")
	callContractFunction := callContractCmd.String("function", "", "Fungsi di dalam kontrak yang akan dipanggil")
	callContractArgs := callContractCmd.String("args", "", "Argumen untuk fungsi kontrak (dipisahkan |)")

	err := callContractCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *callContractFrom == "" || *callContractAddress == "" || *callContractFunction == "" {
		callContractCmd.Usage()
		return
	}

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewContractCallTransaction(*callContractFrom, *callContractAddress, *callContractFunction, *callContractArgs, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Pemanggilan fungsi '%s' pada kontrak %s berhasil dikirim!\n", *callContractFunction, *callContractAddress)
}
