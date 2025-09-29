package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "sendinference",
		Description: "(Legacy) Buat permintaan inferensi sebagai transaksi, bukan pemanggilan kontrak",
		Run:         sendInference,
		NeedsChain:  true,
	})
}

func sendInference(bc *core.Blockchain, args []string) {
	sendInferenceCmd := flag.NewFlagSet("sendinference", flag.ExitOnError)
	sendInferenceFrom := sendInferenceCmd.String("from", "", "Alamat yang membuat permintaan")
	sendInferenceModel := sendInferenceCmd.String("model", "", "Nama model yang akan digunakan")
	sendInferenceInput := sendInferenceCmd.String("input", "", "Data input untuk model")
	sendInferenceReward := sendInferenceCmd.Int("reward", 0, "Jumlah hadiah untuk inferensi")

	err := sendInferenceCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *sendInferenceFrom == "" || *sendInferenceModel == "" || *sendInferenceInput == "" || *sendInferenceReward <= 0 {
		sendInferenceCmd.Usage()
		return
	}

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewInferenceRequestTransaction(*sendInferenceFrom, *sendInferenceModel, *sendInferenceInput, *sendInferenceReward, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Berhasil membuat permintaan inferensi untuk model '%s' dengan hadiah %d\n", *sendInferenceModel, *sendInferenceReward)
}
