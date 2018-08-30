BINARY = blockchain
COVERAGE = coverage
FLAGS = main.go

all: coverage
	@echo Building the binary for target platform...
	@go build -o bin/${BINARY} ${FLAGS}
	@echo Done.

dependencies:
	@bash dependencies.sh


PKG_SHA3_UTILS = ./src/crypto/sha3/utils/nist
PKG_SHA3 = $(PKG_SHA3_UTILS) ./src/crypto/sha3/blake ./src/crypto/sha3/bmw ./src/crypto/sha3/cubehash ./src/crypto/sha3/echo ./src/crypto/sha3/groestl ./src/crypto/sha3/jh ./src/crypto/sha3/keccak ./src/crypto/sha3/luffa ./src/crypto/sha3/shavite ./src/crypto/sha3/simd ./src/crypto/sha3/skein
PKG_CRYPTO = ./src/crypto/secp256k1 $(PKG_SHA3) ./src/crypto/x11
PKG_CORE = ./src/core ./src/core/types ./src/core/types/tx_io

PACKAGES =  $(PKG_CORE) $(PKG_CRYPTO) ./src/network/protocol ./src/utils ./src/wallet ./src/db

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

