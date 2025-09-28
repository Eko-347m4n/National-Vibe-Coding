package main

import (
	"swatantra-node/src/cli"
)

func main() {
	// Serahkan eksekusi ke handler CLI
	cli := cli.CLI{}
	cli.Run()
}
