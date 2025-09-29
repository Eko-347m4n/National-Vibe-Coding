package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "jointask",
		Description: "Bergabung dengan tugas Federated Learning yang ada",
		Run:         joinTask,
		NeedsChain:  true,
	})
}

func joinTask(bc *core.Blockchain, args []string) {
	joinTaskCmd := flag.NewFlagSet("jointask", flag.ExitOnError)
	joinTaskFrom := joinTaskCmd.String("from", "", "Alamat peserta yang bergabung")
	joinTaskContract := joinTaskCmd.String("contract", "", "Alamat kontrak FL market")
	joinTaskTaskID := joinTaskCmd.String("taskid", "", "ID Tugas FL")
	joinTaskParticipant := joinTaskCmd.String("participant", "", "(Opsional) Alamat peserta jika berbeda dari pengirim")
	joinTaskStake := joinTaskCmd.String("stake", "0", "(Opsional) Jumlah stake yang dimiliki peserta")

	err := joinTaskCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *joinTaskFrom == "" || *joinTaskContract == "" || *joinTaskTaskID == "" {
		joinTaskCmd.Usage()
		return
	}

	participant := *joinTaskParticipant
	if participant == "" {
		participant = *joinTaskFrom
	}

	callArgs := fmt.Sprintf("%s|%s|%s", *joinTaskTaskID, participant, *joinTaskStake)

	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewContractCallTransaction(*joinTaskFrom, *joinTaskContract, "join_task", callArgs, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Pemanggilan fungsi '%s' pada kontrak %s berhasil dikirim!\n", "join_task", *joinTaskContract)
}
