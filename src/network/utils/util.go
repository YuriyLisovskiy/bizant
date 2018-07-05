package utils

import (
	"io"
	"net"
	"fmt"
	"log"
	"bytes"
	"encoding/gob"
)

type Ping struct {
	AddrFrom string
}

type Pong struct {
	AddrFrom string
}

func SendData(addr string, data []byte, knownNodes *[]string) {
	conn, err := net.Dial(Protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string
		for _, node := range *knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}
		knownNodes = &updatedNodes
		fmt.Printf("Peers %d\n", len(*knownNodes))
		return
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func extractCommand(request []byte) []byte {
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

func SendPing(addrFrom, addrTo string, knownNodes *[]string) {
	data := GobEncode(Ping{AddrFrom: addrFrom})
	request := append(CommandToBytes("ping"), data...)
	SendData(addrTo, request, knownNodes)
}
