package utils

import (
	"io"
	"net"
	"fmt"
	"log"
	"bytes"
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

func SendPing(addrFrom, addrTo string, knownNodes *[]string) {
	data := GobEncode(Ping{AddrFrom: addrFrom})
	request := append(CommandToBytes("ping"), data...)
	SendData(addrTo, request, knownNodes)
}
