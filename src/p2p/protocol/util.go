// Copyright (c) 2018 Yuriy Lisovskiy
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package protocol

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
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

func MakeRequest(data interface{}, cmd string) []byte {
	return append(CommandToBytes(cmd), GobEncode(data)...)
}
