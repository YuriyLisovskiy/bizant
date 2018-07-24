// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package leveldb

import "fmt"

type ErrorCode int

const (
	ErrorNotFound ErrorCode = iota
	ErrorCorruption
	ErrorNotImplemented
	ErrorInvalidArgument
	ErrorIO
)

var errMap = map[ErrorCode]string{
	ErrorNotFound:        "NotFound",
	ErrorCorruption:      "Corruption",
	ErrorNotImplemented:  "Not implemented",
	ErrorInvalidArgument: "Invalid argument",
	ErrorIO:              "IO error",
}

type Error struct {
	Code ErrorCode
	String string
}

func NewNotFoundError(s string) *Error {
	return &Error{Code: ErrorNotFound, String: s}
}

func NewCorruptionError(s string) *Error {
	return &Error{Code: ErrorCorruption, String: s}
}

func NewNotImplementedError(s string) *Error {
	return &Error{Code: ErrorNotImplemented, String: s}
}

func NewInvalidArgumentError(s string) *Error {
	return &Error{Code: ErrorInvalidArgument, String: s}
}

func NewIOError(s string) *Error {
	return &Error{Code: ErrorIO, String: s}
}

func (e *Error) Error() string {
	t, ok := errMap[e.Code]
	if !ok {
		return fmt.Sprintf("Unknown error(%d): %s", e.Code, e.String)
	}
	return fmt.Sprintf("%s: %s", t, e.String)
}
