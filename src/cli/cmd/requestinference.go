package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "requestinference",
		Description: "Buat permintaan inferensi pada inference market",
		Run:         requestInference,
		NeedsChain:  true,
	})
}

func requestInference(bc *core.Blockchain, args []string) {
	requestInferenceCmd := flag.NewFlagSet("requestinference", flag.ExitOnError)
	requestInferenceFrom := requestInferenceCmd.String("from", "", "Alamat yang membuat permintaan")
	requestInferenceContract := requestInferenceCmd.String("contract", "", "Alamat kontrak inference market")
	requestInferenceModel := requestInferenceCmd.String("model", "", "Nama model yang akan digunakan")
	requestInferenceInput := requestInferenceCmd.String("input", "", "Data input untuk model")

	err := requestInferenceCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *requestInferenceFrom == "" || *requestInferenceContract == "" || *requestInferenceModel == "" || *requestInferenceInput == "" {
		requestInferenceCmd.Usage()
		return
	}

	reward := "10" // Default reward for now
	callArgs := fmt.Sprintf("%s|%s|%s", *requestInferenceModel, *requestInferenceInput, reward)

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewContractCallTransaction(*requestInferenceFrom, *requestInferenceContract, "request_inference", callArgs, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Pemanggilan fungsi '%s' pada kontrak %s berhasil dikirim!\n", "request_inference", *requestInferenceContract)
}
