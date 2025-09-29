package cmd

import (
	"fmt"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "reindexutxo",
		Description: "Bangun ulang UTXO set dari awal",
		Run:         reindexUTXO,
		NeedsChain:  true,
	})
}

func reindexUTXO(bc *core.Blockchain, args []string) {
	UTXOSet := core.UTXOSet{bc}
	UTXOSet.Reindex()

	fmt.Println("Selesai! UTXO set telah diindeks ulang.")
}
