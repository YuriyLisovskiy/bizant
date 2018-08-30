# Copyright (c) 2018 Yuriy Lisovskiy
# Distributed under the GNU General Public License v3.0 software license,
# see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

BINARY = blockchain
COVERAGE = coverage
FLAGS = main.go

all: coverage
	@echo Building the binary for target platform...
	@go build -o bin/${BINARY} ${FLAGS}
	@echo Done.

dependencies:
	@bash dependencies.sh

PACKAGES = ./src/consensus ./src/core ./src/core/types ./src/core/types/tx_io ./src/network/protocol ./src/crypto/secp256k1 ./src/utils ./src/wallet ./src/db

coverage: test
	@echo Generating coverage report...
	@go tool cover -html $(COVERAGE).out -o $(COVERAGE).html
	@echo Done

test:
	@echo Running tests...
	@go test -v -timeout 1h -covermode=count -coverprofile=$(COVERAGE).out $(PACKAGES)
	@echo Done

clean:
	@-rm -rf bin/ $(COVERAGE).out $(COVERAGE).html
