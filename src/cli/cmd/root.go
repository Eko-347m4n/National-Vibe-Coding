package cmd

import (
	"fmt"
	"os"

	"swatantra-node/src/core"
)

// Command mendefinisikan struktur untuk sebuah perintah CLI
type Command struct {
	Name        string
	Description string
	Run         func(chain *core.Blockchain, args []string)
	NeedsChain  bool // Flag untuk menandakan apakah perintah ini butuh instance blockchain
}

// commands menyimpan semua perintah yang terdaftar
var commands = make(map[string]*Command)

// AddCommand mendaftarkan sebuah perintah baru
func AddCommand(cmd *Command) {
	commands[cmd.Name] = cmd
}

// Execute adalah titik masuk utama untuk menjalankan perintah
func Execute() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmd, ok := commands[cmdName]
	if !ok {
		fmt.Printf("Error: Perintah tidak dikenal '%s'\n", cmdName)
		printUsage()
		os.Exit(1)
	}

	// Hanya buat instance blockchain jika perintah membutuhkannya
	if cmd.NeedsChain {
		bc := core.NewBlockchain()
		defer bc.Database().Close()
		cmd.Run(bc, os.Args[2:])
	} else {
		// Untuk perintah seperti 'createblockchain', kita tidak perlu instance yang sudah ada
		cmd.Run(nil, os.Args[2:])
	}
}

// printUsage mencetak daftar semua perintah yang tersedia
func printUsage() {
	fmt.Println("Penggunaan: swatantra-node <perintah> [argumen]")
	fmt.Println("\nPerintah yang tersedia:")
	for _, cmd := range commands {
		fmt.Printf("  %-20s %s\n", cmd.Name, cmd.Description)
	}
}
