package p2p

import (
	"fmt"
	"io"
	"log"
	"net"
)

func StartServer(nodeID, minerAddress string) {
	ln, err := net.Listen(protocol, ":3000")
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	fmt.Printf("Node %s is listening on :3000\n", nodeID)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	request, err := io.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	HandleMessage(request)

	conn.Close()
}
