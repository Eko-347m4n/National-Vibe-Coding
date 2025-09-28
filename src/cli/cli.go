package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"swatantra-node/src/core"
	"swatantra-node/src/p2p"
	"swatantra-node/src/wallet"
)

// CLI bertanggung jawab untuk memproses perintah dari baris perintah
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Penggunaan:")
	fmt.Println("  createblockchain -address ADDRESS - Buat blockchain baru dan kirim hadiah genesis ke alamat")
	fmt.Println("  getbalance -address ADDRESS     - Dapatkan saldo alamat")
	fmt.Println("  createwallet                  - Buat wallet baru")
	fmt.Println("  listaddresses               - Tampilkan semua alamat wallet")
	fmt.Println("  reindexutxo                 - Bangun ulang UTXO set")
	fmt.Println("  printchain                  - Cetak semua blok di blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Kirim sejumlah koin dari satu alamat ke alamat lain")
	fmt.Println("  deploycontract -from FROM -file FILE_PATH - Terbitkan smart contract dari file Lua")
	fmt.Println("  callcontract -from FROM -contract ADDR -function FUNC -args ARGS - Panggil fungsi pada smart contract")
	fmt.Println("  startnode -miner ADDRESS      - Mulai node dengan ID yang diatur dari variabel lingkungan NODE_ID dan mulai menambang")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) callContract(from, contractAddress, function, args string) {
	bc := core.NewBlockchain()
	UTXOSet := core.UTXOSet{bc}
	defer bc.Database().Close()

	// Buat transaksi pemanggilan kontrak
	tx, err := bc.NewContractCallTransaction(from, contractAddress, function, args, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	// Tambahkan transaksi ke blok baru
	bc.AddBlock([]*core.Transaction{tx})
	fmt.Printf("Pemanggilan fungsi '%s' pada kontrak %s berhasil dikirim!\n", function, contractAddress)
}

func (cli *CLI) deployContract(from, filePath string) {
	// Baca kode kontrak dari file
	code, err := os.ReadFile(filePath)
	if err != nil {
		log.Panicf("Gagal membaca file kontrak '%s': %v", filePath, err)
	}

	bc := core.NewBlockchain()
	UTXOSet := core.UTXOSet{bc}
	defer bc.Database().Close()

	// Buat transaksi pembuatan kontrak
	tx, err := bc.NewContractCreationTransaction(from, code, &UTXOSet)
	if err != nil {
		log.Panic(err)
	}

	// Tambahkan transaksi ke blok baru
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
	deployContractCmd := flag.NewFlagSet("deploycontract", flag.ExitOnError)
	callContractCmd := flag.NewFlagSet("callcontract", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "Alamat yang menerima hadiah genesis")
	getBalanceAddress := getBalanceCmd.String("address", "", "Alamat wallet")
	sendFrom := sendCmd.String("from", "", "Alamat pengirim")
	sendTo := sendCmd.String("to", "", "Alamat penerima")
	sendAmount := sendCmd.Int("amount", 0, "Jumlah yang dikirim")
	deployContractFrom := deployContractCmd.String("from", "", "Alamat yang mendanai penerbitan kontrak")
	deployContractFile := deployContractCmd.String("file", "", "Path ke file .lua smart contract")
	callContractFrom := callContractCmd.String("from", "", "Alamat yang memanggil kontrak")
	callContractAddress := callContractCmd.String("contract", "", "Alamat smart contract yang akan dipanggil")
	callContractFunction := callContractCmd.String("function", "", "Fungsi di dalam kontrak yang akan dipanggil")
	callContractArgs := callContractCmd.String("args", "", "Argumen untuk fungsi kontrak (dipisahkan koma)")
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
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
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

	if startNodeCmd.Parsed() {
		nodeID := os.Getenv("NODE_ID")
		if nodeID == "" {
			startNodeCmd.Usage()
			os.Exit(1)
		}
		cli.startNode(nodeID, *startNodeMiner)
	}
}