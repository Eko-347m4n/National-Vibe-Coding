package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "registernode",
		Description: "Mendaftarkan alamat sebagai node inference",
		Run:         registerNode,
		NeedsChain:  true,
	})
}

func registerNode(bc *core.Blockchain, args []string) {
	registerNodeCmd := flag.NewFlagSet("registernode", flag.ExitOnError)
	registerNodeFrom := registerNodeCmd.String("from", "", "Alamat yang mendaftar sebagai node")
	registerNodeContract := registerNodeCmd.String("contract", "", "Alamat kontrak oracle registry")

	err := registerNodeCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *registerNodeFrom == "" || *registerNodeContract == "" {
		registerNodeCmd.Usage()
		return
	}

	// The argument to the 'register' function is the node's address itself.
	callArgs := *registerNodeFrom

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewContractCallTransaction(*registerNodeFrom, *registerNodeContract, "register", callArgs, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Pemanggilan fungsi '%s' pada kontrak %s berhasil dikirim!\n", "register", *registerNodeContract)
}
