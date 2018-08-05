// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

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
