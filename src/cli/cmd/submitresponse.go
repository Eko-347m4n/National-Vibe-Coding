package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "submitresponse",
		Description: "Kirim hasil inferensi untuk sebuah pekerjaan",
		Run:         submitResponse,
		NeedsChain:  true,
	})
}

func submitResponse(bc *core.Blockchain, args []string) {
	submitResponseCmd := flag.NewFlagSet("submitresponse", flag.ExitOnError)
	submitResponseFrom := submitResponseCmd.String("from", "", "Alamat node yang mengirimkan hasil")
	submitResponseContract := submitResponseCmd.String("contract", "", "Alamat kontrak inference market")
	submitResponseJobID := submitResponseCmd.String("jobid", "", "ID pekerjaan inferensi")
	submitResponseResult := submitResponseCmd.String("result", "", "Hasil inferensi")

	err := submitResponseCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *submitResponseFrom == "" || *submitResponseContract == "" || *submitResponseJobID == "" || *submitResponseResult == "" {
		submitResponseCmd.Usage()
		return
	}

	callArgs := fmt.Sprintf("%s|%s|%s", *submitResponseJobID, *submitResponseResult, *submitResponseFrom)

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewContractCallTransaction(*submitResponseFrom, *submitResponseContract, "submit_response", callArgs, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Pemanggilan fungsi '%s' pada kontrak %s berhasil dikirim!\n", "submit_response", *submitResponseContract)
}
