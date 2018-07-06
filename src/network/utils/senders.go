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

func SendData(addr string, data []byte, knownNodes *map[string]bool) bool {
	conn, err := net.Dial(PROTOCOL, addr)
	success := true
	if err != nil {
		success = false
		delete(*knownNodes, addr)
		fmt.Printf("Peers %d\n", len(*knownNodes))
		return success
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
	return success
}

func SendPing(addrFrom, addrTo string, knownNodes *map[string]bool) bool {
	data := GobEncode(Ping{AddrFrom: addrFrom})
	request := append(CommandToBytes("ping"), data...)
	return SendData(addrTo, request, knownNodes)
}
