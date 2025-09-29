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
		Name:        "getjob",
		Description: "Lihat detail pekerjaan inferensi",
		Run:         getJob,
		NeedsChain:  true,
	})
}

func getJob(bc *core.Blockchain, args []string) {
	getJobCmd := flag.NewFlagSet("getjob", flag.ExitOnError)
	getJobContract := getJobCmd.String("contract", "", "Alamat kontrak inference market")
	getJobJobID := getJobCmd.String("jobid", "", "ID pekerjaan inferensi")

	err := getJobCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *getJobContract == "" || *getJobJobID == "" {
		getJobCmd.Usage()
		return
	}

	wallets, _ := wallet.NewWallets()
	addresses := wallets.GetAddresses()
	if len(addresses) == 0 {
		log.Panic("Tidak ada wallet yang ditemukan untuk melakukan panggilan getjob.")
	}
	// A get call doesn't need a transaction from a specific user, so we use the first available wallet.
	from := addresses[0]

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewContractCallTransaction(from, *getJobContract, "get_job", *getJobJobID, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Pemanggilan fungsi '%s' pada kontrak %s berhasil dikirim!\n", "get_job", *getJobContract)
}
