// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package network

import (
	"fmt"
	"log"
	"bytes"
	"encoding/gob"
)

func ExtractCommand(request []byte) []byte {
	return request[:COMMAND_LENGTH]
}

func CommandToBytes(command string) []byte {
	var b [COMMAND_LENGTH]byte
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

func makeRequest(data interface{}, cmd string) []byte {
	return append(CommandToBytes(cmd), GobEncode(data)...)
}
