package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"swatantra-node/src/core"
	"swatantra-node/src/p2p"
	"swatantra-node/src/wallet"
)

func init() {
	AddCommand(&Command{
		Name:        "startnode",
		Description: "Mulai node, secara opsional dengan alamat miner",
		Run:         startNode,
		NeedsChain:  false,
	})
}

func startNode(bc *core.Blockchain, args []string) {
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	startNodeMiner := startNodeCmd.String("miner", "", "Aktifkan penambangan dan kirim hadiah ke alamat ini")

	err := startNodeCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Println("Error: NODE_ID environment variable not set.")
		startNodeCmd.Usage()
		return
	}

	fmt.Printf("Starting node %s\n", nodeID)
	if *startNodeMiner != "" {
		if wallet.ValidateAddress(*startNodeMiner) {
			fmt.Println("Mining is on. Address to receive rewards: ", *startNodeMiner)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	p2p.StartServer(nodeID, *startNodeMiner)
}
