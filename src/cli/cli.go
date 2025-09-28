package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"swatantra-node/src/core"
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
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
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

// Run memulai pemrosesan perintah CLI
func (cli *CLI) Run() {
	cli.validateArgs()

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "Alamat yang menerima hadiah genesis")
	getBalanceAddress := getBalanceCmd.String("address", "", "Alamat wallet")

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
}