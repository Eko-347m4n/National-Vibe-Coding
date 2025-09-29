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
		Name:        "submithash",
		Description: "Kirim hash model terenkripsi untuk tugas FL",
		Run:         submitHash,
		NeedsChain:  true,
	})
}

func submitHash(bc *core.Blockchain, args []string) {
	submitHashCmd := flag.NewFlagSet("submithash", flag.ExitOnError)
	submitHashFrom := submitHashCmd.String("from", "", "Alamat peserta yang mengirim")
	submitHashContract := submitHashCmd.String("contract", "", "Alamat kontrak FL market")
	submitHashTaskID := submitHashCmd.String("taskid", "", "ID Tugas FL")
	submitHashHash := submitHashCmd.String("hash", "", "Hash dari model yang diperbarui (plaintext, akan dienkripsi)")
	submitHashParticipant := submitHashCmd.String("participant", "", "(Opsional) Alamat peserta jika berbeda dari pengirim")

	err := submitHashCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *submitHashFrom == "" || *submitHashContract == "" || *submitHashTaskID == "" || *submitHashHash == "" {
		submitHashCmd.Usage()
		return
	}

	participant := *submitHashParticipant
	if participant == "" {
		participant = *submitHashFrom
	}

	fmt.Println("Mengenkripsi hash model...")
	encryptedHash, err := utils.Encrypt([]byte(*submitHashHash))
	if err != nil {
		log.Panicf("Gagal mengenkripsi hash: %v", err)
	}
	fmt.Printf("Hash terenkripsi (Base64): %s\n", encryptedHash)

	callArgs := fmt.Sprintf("%s|%s|%s", *submitHashTaskID, encryptedHash, participant)

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewContractCallTransaction(*submitHashFrom, *submitHashContract, "submit_hash", callArgs, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Pemanggilan fungsi '%s' pada kontrak %s berhasil dikirim!\n", "submit_hash", *submitHashContract)
}
