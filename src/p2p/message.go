package p2p

import (
	"bytes"
	"encoding/gob"
	"log"
)

const (
	commandLength = 12
	protocol      = "tcp"
)

type Message struct {
	Command []byte
	Payload []byte
}

func NewMessage(command string, payload []byte) *Message {
	return &Message{[]byte(command), payload}
}

func (msg *Message) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(msg)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func DeserializeMessage(data []byte) *Message {
	var msg Message

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&msg)
	if err != nil {
		log.Panic(err)
	}

	return &msg
}
