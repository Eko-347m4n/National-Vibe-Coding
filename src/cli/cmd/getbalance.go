package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
	"swatantra-node/src/wallet"
)

func init() {
	AddCommand(&Command{
		Name:        "getbalance",
		Description: "Dapatkan saldo sebuah alamat",
		Run:         getBalance,
		NeedsChain:  true,
	})
}

func getBalance(bc *core.Blockchain, args []string) {
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "Alamat wallet")

	err := getBalanceCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *getBalanceAddress == "" {
		getBalanceCmd.Usage()
		return
	}

	UTXOSet := core.UTXOSet{bc}

	balance := 0
	pubKeyHash := wallet.DecodeAddress(*getBalanceAddress)
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Saldo '%s': %d\n", *getBalanceAddress, balance)
}
