package utils

import (
	"fmt"
	"log"
	"bytes"
	"encoding/gob"
)

func ExtractCommand(request []byte) []byte {
	return request[:CommandLength]
}

func CommandToBytes(command string) []byte {
	var b [CommandLength]byte
	for i, c := range command {
		b[i] = byte(c)
	}
	return b[:]
}

func BytesToCommand(bytes []byte) string {
	var command []byte
	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func NodeIsKnown(addr string, knownNodes []string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}
	return false
}
