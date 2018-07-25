// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

var (
	InvalidError               = &Error{"Invalid database", nil}
	VersionMismatchError       = &Error{"version mismatch", nil}
	DatabaseNotOpenError       = &Error{"db is not open", nil}
	DatabaseAlreadyOpenedError = &Error{"db already open", nil}
	TransactionInProgressError = &Error{"writable transaction is already in progress", nil}
	InvalidTransactionError    = &Error{"txn is invalid", nil}
	BucketAlreadyExistsError   = &Error{"bucket already exists", nil}
)

type Error struct {
	message string
	cause   error
}

func (e *Error) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}
