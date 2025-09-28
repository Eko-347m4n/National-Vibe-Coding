package p2p

import (
	"fmt"
)

func HandleMessage(request []byte) {
	msg := DeserializeMessage(request)

	switch string(msg.Command) {
	case "version":
		fmt.Println("Received version message")
	case "getblocks":
		fmt.Println("Received getblocks message")
	case "inv":
		fmt.Println("Received inv message")
	case "getdata":
		fmt.Println("Received getdata message")
	case "block":
		fmt.Println("Received block message")
	case "tx":
		fmt.Println("Received tx message")
	default:
		fmt.Println("Unknown message")
	}
}
