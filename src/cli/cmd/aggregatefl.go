package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
	"swatantra-node/src/utils"
)

func init() {
	AddCommand(&Command{
		Name:        "aggregatefl",
		Description: "(Simulasi Aggregator) Finalisasi tugas FL dan distribusikan hadiah",
		Run:         aggregateFL,
		NeedsChain:  true,
	})
}

func aggregateFL(bc *core.Blockchain, args []string) {
	aggregateFLCmd := flag.NewFlagSet("aggregatefl", flag.ExitOnError)
	aggregateFLFrom := aggregateFLCmd.String("from", "", "Alamat yang memicu agregasi")
	aggregateFLContract := aggregateFLCmd.String("contract", "", "Alamat kontrak FL market")
	aggregateFLTaskID := aggregateFLCmd.String("taskid", "", "ID Tugas FL yang akan diagregasi")

	err := aggregateFLCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *aggregateFLFrom == "" || *aggregateFLContract == "" || *aggregateFLTaskID == "" {
		aggregateFLCmd.Usage()
		return
	}

	log.Println("[AGGREGATOR DAEMON] Memulai proses agregasi untuk tugas", *aggregateFLTaskID)

	log.Println("[AGGREGATOR DAEMON] Mensimulasikan pengambilan hash terenkripsi dari state kontrak...")
	simulatedEncryptedHash := "CiCCM4boNmyg3T5+t+P2UdGvYyEwV+i5P9vj/ROl3a2b+XfA1A==" // Placeholder
	log.Printf("[AGGREGATOR DAEMON] Ditemukan hash terenkripsi (contoh): %s", simulatedEncryptedHash)

	log.Println("[AGGREGATOR DAEMON] Mensimulasikan dekripsi hash...")
	decryptedHash, err := utils.Decrypt(simulatedEncryptedHash)
	if err != nil {
		log.Printf("[AGGREGATOR DAEMON] Peringatan: Gagal mendekripsi contoh hash: %v. Melanjutkan dengan data simulasi.", err)
		decryptedHash = []byte("fallback_decrypted_hash")
	}
	log.Printf("[AGGREGATOR DAEMON] Hash berhasil didekripsi (contoh): %s", string(decryptedHash))

	aggregatedHash := "simulated_aggregated_hash_from_" + string(decryptedHash)
	log.Println("[AGGREGATOR DAEMON] Hasil agregasi (simulasi):", aggregatedHash)

	callArgs := fmt.Sprintf("%s|%s", *aggregateFLTaskID, aggregatedHash)

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewContractCallTransaction(*aggregateFLFrom, *aggregateFLContract, "finalize_task", callArgs, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Pemanggilan fungsi '%s' pada kontrak %s berhasil dikirim!\n", "finalize_task", *aggregateFLContract)
}
