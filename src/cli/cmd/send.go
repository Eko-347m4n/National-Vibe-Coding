package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "send",
		Description: "Kirim koin dari satu alamat ke alamat lain",
		Run:         send,
		NeedsChain:  true,
	})
}

func send(bc *core.Blockchain, args []string) {
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := sendCmd.String("from", "", "Alamat pengirim")
	sendTo := sendCmd.String("to", "", "Alamat penerima")
	sendAmount := sendCmd.Int("amount", 0, "Jumlah yang dikirim")

	err := sendCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
		sendCmd.Usage()
		return
	}

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewUTXOTransaction(*sendFrom, *sendTo, *sendAmount, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}
	bc.AddBlock([]*core.Transaction{tx})
	fmt.Println("Transaksi berhasil!")
}
