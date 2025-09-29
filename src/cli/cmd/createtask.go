package cmd

import (
	"flag"
	"fmt"
	"log"

	"swatantra-node/src/core"
)

func init() {
	AddCommand(&Command{
		Name:        "createtask",
		Description: "Buat tugas Federated Learning baru",
		Run:         createTask,
		NeedsChain:  true,
	})
}

func createTask(bc *core.Blockchain, args []string) {
	createTaskCmd := flag.NewFlagSet("createtask", flag.ExitOnError)
	createTaskFrom := createTaskCmd.String("from", "", "Alamat yang membuat tugas FL")
	createTaskContract := createTaskCmd.String("contract", "", "Alamat kontrak FL market")
	createTaskModel := createTaskCmd.String("model", "", "Nama model awal untuk dilatih")
	createTaskMinP := createTaskCmd.String("min", "", "Jumlah minimal peserta")
	createTaskReward := createTaskCmd.String("reward", "0", "Total hadiah untuk tugas FL")

	err := createTaskCmd.Parse(args)
	if err != nil {
		log.Panic(err)
	}

	if *createTaskFrom == "" || *createTaskContract == "" || *createTaskModel == "" || *createTaskMinP == "" || *createTaskReward == "" {
		createTaskCmd.Usage()
		return
	}

	callArgs := fmt.Sprintf("%s|%s|%s", *createTaskModel, *createTaskMinP, *createTaskReward)
	
	UTXOSet := core.UTXOSet{bc}

	tx, err := bc.NewContractCallTransaction(*createTaskFrom, *createTaskContract, "create_task", callArgs, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Pemanggilan fungsi '%s' pada kontrak %s berhasil dikirim!\n", "create_task", *createTaskContract)
}
