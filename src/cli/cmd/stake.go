package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "stake",
		Description: "Stake koin untuk menjadi node (koin akan di-burn)",
		Run:         stake,
		NeedsChain:  true,
	})
}

func stake(bc *core.Blockchain, args []string) {
	stakeCmd := flag.NewFlagSet("stake", flag.ExitOnError)
	stakeFrom := stakeCmd.String("from", "", "Alamat yang melakukan stake")
	stakeAmount := stakeCmd.Int("amount", 0, "Jumlah yang di-stake")

	err := stakeCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *stakeFrom == "" || *stakeAmount <= 0 {
		stakeCmd.Usage()
		return
	}

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewStakeTransaction(*stakeFrom, *stakeAmount, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Berhasil stake (burn) %d koin dari %s\n", *stakeAmount, *stakeFrom)
}
