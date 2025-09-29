package cmd

import (
	"fmt"
	"strconv"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "printchain",
		Description: "Cetak semua blok di blockchain",
		Run:         printChain,
		NeedsChain:  true,
	})
}

func printChain(bc *core.Blockchain, args []string) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ Block %d ============\n", block.Height)
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		if len(block.Transactions) > 0 {
			fmt.Printf("Data Transaksi: %s\n", string(block.Transactions[0].Outputs[0].PubKeyHash))
		}
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := core.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
