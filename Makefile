BINARY = blockchain
COVERAGE = coverage
FLAGS = main.go

all:
	@echo Building the binary for target platform...
	@go build -o bin/${BINARY} ${FLAGS}
	@echo Done.

dependencies:
	@bash dependencies.sh

CRYPTO = ./src/crypto
CORE = ./src/core
SHA3 = $(CRYPTO)/sha3
ACCOUNTS = ./src/accounts

PKG_SHA3_UTILS = $(SHA3)/utils/nist
PKG_SHA3 = $(PKG_SHA3_UTILS) $(SHA3)/blake512 $(SHA3)/bmw512 $(SHA3)/cubehash512 $(SHA3)/echo512 $(SHA3)/groestl512 $(SHA3)/jh512 $(SHA3)/keccak512 $(SHA3)/luffa512 $(SHA3)/shavite512 $(SHA3)/simd512 $(SHA3)/skein512
PKG_CRYPTO = $(CRYPTO)/secp256k1 $(PKG_SHA3) $(CRYPTO)/x11
PKG_CORE = $(CORE) $(CORE)/types $(CORE)/types/tx_io
PKG_ACCOUNTS = $(ACCOUNTS)/wallet $(ACCOUNTS)/auth/jwt

PACKAGES =  $(PKG_CORE) $(PKG_CRYPTO) $(PKG_ACCOUNTS) ./src/network/protocol ./src/utils ./src/encoding/base58 ./src/config ./src/db

test:
	@echo Running tests...
	@go test -v -timeout 3h -covermode=count -coverprofile=$(COVERAGE).out $(PACKAGES)
	@echo Generating coverage report...
	@go tool cover -html $(COVERAGE).out -o $(COVERAGE).html
	@echo Done

clean:
	@-rm -rf bin/ $(COVERAGE).out $(COVERAGE).html

