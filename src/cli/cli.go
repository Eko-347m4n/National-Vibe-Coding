package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"swatantra-node/src/core"
	"swatantra-node/src/p2p"
	"swatantra-node/src/utils"
	"swatantra-node/src/wallet"
)

// CLI bertanggung jawab untuk memproses perintah dari baris perintah
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Penggunaan:")
	fmt.Println("  createblockchain -address ADDRESS - Buat blockchain baru")
	fmt.Println("  getbalance -address ADDRESS     - Dapatkan saldo alamat")
	fmt.Println("  createwallet                  - Buat wallet baru")
	fmt.Println("  listaddresses               - Tampilkan semua alamat wallet")
	fmt.Println("  printchain                  - Cetak semua blok di blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Kirim koin")
	fmt.Println("  sendinference -from FROM -model MODEL_NAME -input INPUT_DATA -reward REWARD - Buat permintaan inferensi")
	fmt.Println("  stake -from FROM -amount AMOUNT - Stake koin untuk menjadi node (koin akan di-burn)")
	fmt.Println("  deploycontract -from FROM -file FILE_PATH - Terbitkan smart contract")
	fmt.Println("  callcontract -from FROM -contract ADDR -function FUNC -args ARGS - Panggil fungsi kontrak")
	fmt.Println("  publishmodel -from FROM -contract ADDR -name NAME -file FILE - Publikasikan model AI")
	fmt.Println("  requestinference -from FROM -contract ADDR -model NAME -input DATA - Buat permintaan inferensi")
	fmt.Println("  submitresponse -from FROM -contract ADDR -jobid ID -result DATA - Kirim hasil inferensi")
	fmt.Println("  getjob -contract ADDR -jobid ID - Lihat detail pekerjaan inferensi")
	fmt.Println("  createtask -from FROM -contract ADDR -model NAME -min P - Buat tugas Federated Learning")
	fmt.Println("  jointask -from FROM -contract ADDR -taskid ID -stake STAKE [-participant ADDR] - Bergabung dengan tugas FL")
	fmt.Println("  submithash -from FROM -contract ADDR -taskid ID -hash HASH [-participant ADDR] - Kirim hash model FL")
	fmt.Println("  registernode -from FROM -contract ADDR - Mendaftarkan alamat sebagai node inference")
	fmt.Println("  startnode -miner ADDRESS      - Mulai node")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) stake(from string, amount int) {
	bc := core.NewBlockchain()
	UTXOSet := core.UTXOSet{bc}
	defer bc.Database().Close()

	tx, err := bc.NewStakeTransaction(from, amount, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Berhasil stake (burn) %d koin dari %s\n", amount, from)
}

func (cli *CLI) sendInference(from, modelName, inputData string, reward int) {
	bc := core.NewBlockchain()
	UTXOSet := core.UTXOSet{bc}
	defer bc.Database().Close()

	tx, err := bc.NewInferenceRequestTransaction(from, modelName, inputData, reward, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Berhasil membuat permintaan inferensi untuk model '%s' dengan hadiah %d\n", modelName, reward)
}

func (cli *CLI) createTask(from, contract, model, minParticipants string) {
	args := fmt.Sprintf("%s|%s", model, minParticipants)
	cli.callContract(from, contract, "create_task", args)
}

func (cli *CLI) joinTask(from, contract, taskID, participant, stake string) {
	if participant == "" {
		participant = from
	}
	args := fmt.Sprintf("%s|%s|%s", taskID, participant, stake)
	cli.callContract(from, contract, "join_task", args)
}

func (cli *CLI) submitHash(from, contract, taskID, hash, participant string) {
	if participant == "" {
		participant = from
	}
	args := fmt.Sprintf("%s|%s|%s", taskID, hash, participant)
	cli.callContract(from, contract, "submit_hash", args)
}

func (cli *CLI) requestInference(from, contract, model, input string) {
	reward := "10"
	args := fmt.Sprintf("%s|%s|%s", model, input, reward)
	cli.callContract(from, contract, "request_inference", args)
}

func (cli *CLI) submitResponse(from, contract, jobID, result string) {
	args := fmt.Sprintf("%s|%s|%s", jobID, result, from)
	cli.callContract(from, contract, "submit_response", args)
}

func (cli *CLI) getJob(contract, jobID string) {
	wallets, _ := wallet.NewWallets()
	addresses := wallets.GetAddresses()
	if len(addresses) == 0 {
		log.Panic("Tidak ada wallet yang ditemukan untuk melakukan panggilan getjob.")
	}
	cli.callContract(addresses[0], contract, "get_job", jobID)
}

func (cli *CLI) registerNode(from, contractAddress string) {
    // Argumen untuk fungsi register(node_address)
    // Alamat yang mendaftar adalah alamat yang membayar fee
    args := from
    cli.callContract(from, contractAddress, "register", args)
}

func (cli *CLI) publishModel(from, contractAddress, modelName, filePath string) {
	fmt.Printf("Menghitung hash untuk file model '%s'...\n", filePath)
	modelHash, err := utils.FileKeccak256(filePath)
	if err != nil {
		log.Panicf("Gagal menghitung hash file: %v", err)
	}
	fmt.Printf("Hash model: %s\n", modelHash)

	location := filePath
	args := fmt.Sprintf("%s|%s|%s", modelName, modelHash, location)
	cli.callContract(from, contractAddress, "publish_model", args)
}

func (cli *CLI) callContract(from, contractAddress, function, args string) {
	bc := core.NewBlockchain()
	UTXOSet := core.UTXOSet{bc}
	defer bc.Database().Close()

	tx, err := bc.NewContractCallTransaction(from, contractAddress, function, args, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Pemanggilan fungsi '%s' pada kontrak %s berhasil dikirim!\n", function, contractAddress)
}

func (cli *CLI) deployContract(from, filePath string) {
	code, err := os.ReadFile(filePath)
	if err != nil {
		log.Panicf("Gagal membaca file kontrak '%s': %v", filePath, err)
	}

	bc := core.NewBlockchain()
	UTXOSet := core.UTXOSet{bc}
	defer bc.Database().Close()

	tx, err := bc.NewContractCreationTransaction(from, code, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Kontrak dari file '%s' berhasil diterbitkan!\n", filePath)
}

func (cli *CLI) createBlockchain(address string) {
	bc := core.CreateBlockchain(address)
	defer bc.Database().Close()
	fmt.Println("Blockchain baru berhasil dibuat!")
}

func (cli *CLI) getBalance(address string) {
	bc := core.NewBlockchain()
	UTXOSet := core.UTXOSet{bc}
	defer bc.Database().Close()

	balance := 0
	pubKeyHash := wallet.DecodeAddress(address)
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Saldo '%s': %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
	bc := core.NewBlockchain()
	UTXOSet := core.UTXOSet{bc}
	defer bc.Database().Close()

	tx, err := bc.NewUTXOTransaction(from, to, amount, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}
	bc.AddBlock([]*core.Transaction{tx})
	fmt.Println("Transaksi berhasil!")
}

func (cli *CLI) reindexUTXO() {
	bc := core.NewBlockchain()
	UTXOSet := core.UTXOSet{bc}
	UTXOSet.Reindex()

	fmt.Println("Selesai! UTXO set telah diindeks ulang.")
}

func (cli *CLI) createWallet() {
	wallets, _ := wallet.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Alamat baru Anda: %s\n", address)
}

func (cli *CLI) listAddresses() {
	wallets, _ := wallet.NewWallets()
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CLI) printChain() {
	bc := core.NewBlockchain()
	defer bc.Database().Close()
	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ Block %d ============\n", block.Height)
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		if len(block.Transactions) > 0 {
			fmt.Printf("Data Transaksi: %s\n", string(block.Transactions[0].Outputs[0].PubKeyHash))
		}
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := core.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) startNode(nodeID, minerAddress string) {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	p2p.StartServer(nodeID, minerAddress)
}

// Run memulai pemrosesan perintah CLI
func (cli *CLI) Run() {
	cli.validateArgs()

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	stakeCmd := flag.NewFlagSet("stake", flag.ExitOnError)
	deployContractCmd := flag.NewFlagSet("deploycontract", flag.ExitOnError)
	callContractCmd := flag.NewFlagSet("callcontract", flag.ExitOnError)
	publishModelCmd := flag.NewFlagSet("publishmodel", flag.ExitOnError)
	requestInferenceCmd := flag.NewFlagSet("requestinference", flag.ExitOnError)
	submitResponseCmd := flag.NewFlagSet("submitresponse", flag.ExitOnError)
	getJobCmd := flag.NewFlagSet("getjob", flag.ExitOnError)
	createTaskCmd := flag.NewFlagSet("createtask", flag.ExitOnError)
	joinTaskCmd := flag.NewFlagSet("jointask", flag.ExitOnError)
	submitHashCmd := flag.NewFlagSet("submithash", flag.ExitOnError)
	registerNodeCmd := flag.NewFlagSet("registernode", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	sendInferenceCmd := flag.NewFlagSet("sendinference", flag.ExitOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "Alamat yang menerima hadiah genesis")
	getBalanceAddress := getBalanceCmd.String("address", "", "Alamat wallet")
	sendFrom := sendCmd.String("from", "", "Alamat pengirim")
	sendTo := sendCmd.String("to", "", "Alamat penerima")
	sendAmount := sendCmd.Int("amount", 0, "Jumlah yang dikirim")
	stakeFrom := stakeCmd.String("from", "", "Alamat yang melakukan stake")
	stakeAmount := stakeCmd.Int("amount", 0, "Jumlah yang di-stake")
	deployContractFrom := deployContractCmd.String("from", "", "Alamat yang mendanai penerbitan kontrak")
	deployContractFile := deployContractCmd.String("file", "", "Path ke file .lua smart contract")
	callContractFrom := callContractCmd.String("from", "", "Alamat yang memanggil kontrak")
	callContractAddress := callContractCmd.String("contract", "", "Alamat smart contract yang akan dipanggil")
	callContractFunction := callContractCmd.String("function", "", "Fungsi di dalam kontrak yang akan dipanggil")
	callContractArgs := callContractCmd.String("args", "", "Argumen untuk fungsi kontrak (dipisahkan koma)")
	publishModelFrom := publishModelCmd.String("from", "", "Alamat yang mendanai publikasi model")
	publishModelContract := publishModelCmd.String("contract", "", "Alamat kontrak model registry")
	publishModelName := publishModelCmd.String("name", "", "Nama unik untuk model")
	publishModelFile := publishModelCmd.String("file", "", "Path ke file model")
	requestInferenceFrom := requestInferenceCmd.String("from", "", "Alamat yang membuat permintaan")
	requestInferenceContract := requestInferenceCmd.String("contract", "", "Alamat kontrak inference market")
	requestInferenceModel := requestInferenceCmd.String("model", "", "Nama model yang akan digunakan")
	requestInferenceInput := requestInferenceCmd.String("input", "", "Data input untuk model")
	sendInferenceFrom := sendInferenceCmd.String("from", "", "Alamat yang membuat permintaan")
	sendInferenceModel := sendInferenceCmd.String("model", "", "Nama model yang akan digunakan")
	sendInferenceInput := sendInferenceCmd.String("input", "", "Data input untuk model")
	sendInferenceReward := sendInferenceCmd.Int("reward", 0, "Jumlah hadiah untuk inferensi")
	submitResponseFrom := submitResponseCmd.String("from", "", "Alamat node yang mengirimkan hasil")
	submitResponseContract := submitResponseCmd.String("contract", "", "Alamat kontrak inference market")
	submitResponseJobID := submitResponseCmd.String("jobid", "", "ID pekerjaan inferensi")
	submitResponseResult := submitResponseCmd.String("result", "", "Hasil inferensi")
	getJobContract := getJobCmd.String("contract", "", "Alamat kontrak inference market")
	getJobJobID := getJobCmd.String("jobid", "", "ID pekerjaan inferensi")
	createTaskFrom := createTaskCmd.String("from", "", "Alamat yang membuat tugas FL")
	createTaskContract := createTaskCmd.String("contract", "", "Alamat kontrak FL market")
	createTaskModel := createTaskCmd.String("model", "", "Nama model awal untuk dilatih")
	createTaskMinP := createTaskCmd.String("min", "", "Jumlah minimal peserta")
	joinTaskFrom := joinTaskCmd.String("from", "", "Alamat peserta yang bergabung")
	joinTaskContract := joinTaskCmd.String("contract", "", "Alamat kontrak FL market")
	joinTaskTaskID := joinTaskCmd.String("taskid", "", "ID Tugas FL")
	joinTaskParticipant := joinTaskCmd.String("participant", "", "(Opsional) Alamat peserta jika berbeda dari pengirim")
	joinTaskStake := joinTaskCmd.String("stake", "0", "(Opsional) Jumlah stake yang dimiliki peserta")
	submitHashFrom := submitHashCmd.String("from", "", "Alamat peserta yang mengirim")
	submitHashContract := submitHashCmd.String("contract", "", "Alamat kontrak FL market")
	submitHashTaskID := submitHashCmd.String("taskid", "", "ID Tugas FL")
	submitHashHash := submitHashCmd.String("hash", "", "Hash dari model yang diperbarui")
	submitHashParticipant := submitHashCmd.String("participant", "", "(Opsional) Alamat peserta jika berbeda dari pengirim")
	registerNodeFrom := registerNodeCmd.String("from", "", "Alamat yang mendaftar sebagai node")
    registerNodeContract := registerNodeCmd.String("contract", "", "Alamat kontrak oracle registry")
	startNodeMiner := startNodeCmd.String("miner", "", "Aktifkan penambangan dan kirim hadiah ke alamat ini")

	switch os.Args[1] {
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "stake":
		err := stakeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "deploycontract":
		err := deployContractCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "callcontract":
		err := callContractCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "publishmodel":
		err := publishModelCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "requestinference":
		err := requestInferenceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "submitresponse":
		err := submitResponseCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getjob":
		err := getJobCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createtask":
		err := createTaskCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "jointask":
		err := joinTaskCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "submithash":
		err := submitHashCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "registernode":
		err := registerNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "sendinference":
		err := sendInferenceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO()
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if stakeCmd.Parsed() {
		if *stakeFrom == "" || *stakeAmount <= 0 {
			stakeCmd.Usage()
			os.Exit(1)
		}
		cli.stake(*stakeFrom, *stakeAmount)
	}

	if deployContractCmd.Parsed() {
		if *deployContractFrom == "" || *deployContractFile == "" {
			deployContractCmd.Usage()
			os.Exit(1)
		}
		cli.deployContract(*deployContractFrom, *deployContractFile)
	}

	if callContractCmd.Parsed() {
		if *callContractFrom == "" || *callContractAddress == "" || *callContractFunction == "" {
			callContractCmd.Usage()
			os.Exit(1)
		}
		cli.callContract(*callContractFrom, *callContractAddress, *callContractFunction, *callContractArgs)
	}

	if publishModelCmd.Parsed() {
		if *publishModelFrom == "" || *publishModelContract == "" || *publishModelName == "" || *publishModelFile == "" {
			publishModelCmd.Usage()
			os.Exit(1)
		}
		cli.publishModel(*publishModelFrom, *publishModelContract, *publishModelName, *publishModelFile)
	}

	if requestInferenceCmd.Parsed() {
		if *requestInferenceFrom == "" || *requestInferenceContract == "" || *requestInferenceModel == "" || *requestInferenceInput == "" {
			requestInferenceCmd.Usage()
			os.Exit(1)
		}
		cli.requestInference(*requestInferenceFrom, *requestInferenceContract, *requestInferenceModel, *requestInferenceInput)
	}

	if sendInferenceCmd.Parsed() {
		if *sendInferenceFrom == "" || *sendInferenceModel == "" || *sendInferenceInput == "" || *sendInferenceReward <= 0 {
			sendInferenceCmd.Usage()
			os.Exit(1)
		}
		cli.sendInference(*sendInferenceFrom, *sendInferenceModel, *sendInferenceInput, *sendInferenceReward)
	}

	if submitResponseCmd.Parsed() {
		if *submitResponseFrom == "" || *submitResponseContract == "" || *submitResponseJobID == "" || *submitResponseResult == "" {
			submitResponseCmd.Usage()
			os.Exit(1)
		}
		cli.submitResponse(*submitResponseFrom, *submitResponseContract, *submitResponseJobID, *submitResponseResult)
	}

	if getJobCmd.Parsed() {
		if *getJobContract == "" || *getJobJobID == "" {
			getJobCmd.Usage()
			os.Exit(1)
		}
		cli.getJob(*getJobContract, *getJobJobID)
	}

	if createTaskCmd.Parsed() {
		if *createTaskFrom == "" || *createTaskContract == "" || *createTaskModel == "" || *createTaskMinP == "" {
			createTaskCmd.Usage()
			os.Exit(1)
		}
		cli.createTask(*createTaskFrom, *createTaskContract, *createTaskModel, *createTaskMinP)
	}

	if joinTaskCmd.Parsed() {
		if *joinTaskFrom == "" || *joinTaskContract == "" || *joinTaskTaskID == "" {
			joinTaskCmd.Usage()
			os.Exit(1)
		}
		cli.joinTask(*joinTaskFrom, *joinTaskContract, *joinTaskTaskID, *joinTaskParticipant, *joinTaskStake)
	}

	if submitHashCmd.Parsed() {
		if *submitHashFrom == "" || *submitHashContract == "" || *submitHashTaskID == "" || *submitHashHash == "" {
			submitHashCmd.Usage()
			os.Exit(1)
		}
		cli.submitHash(*submitHashFrom, *submitHashContract, *submitHashTaskID, *submitHashHash, *submitHashParticipant)
	}

	if registerNodeCmd.Parsed() {
        if *registerNodeFrom == "" || *registerNodeContract == "" {
            registerNodeCmd.Usage()
            os.Exit(1)
        }
        cli.registerNode(*registerNodeFrom, *registerNodeContract)
    }

	if startNodeCmd.Parsed() {
		nodeID := os.Getenv("NODE_ID")
		if nodeID == "" {
			startNodeCmd.Usage()
			os.Exit(1)
		}
		cli.startNode(nodeID, *startNodeMiner)
	}
}