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

const (
	C_TX        = "tx"
	C_INV       = "inv"
	C_PING      = "ping"
	C_PONG      = "pong"
	C_ADDR      = "addr"
	C_BLOCK     = "block"
	C_ERROR     = "error"
	C_VERSION   = "version"
	C_GETDATA   = "getdata"
	C_GETBLOCKS = "getblocks"
	C_MESSAGE   = "msg"
	C_SYNCED    = "synced"
)

const (
	PROTOCOL       = "tcp"
	NODE_VERSION   = 1
	COMMAND_LENGTH = 12
)
