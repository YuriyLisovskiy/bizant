BINARY = blockchain

ARCH_i386 = 386
ARCH_arm = arm
ARCH_amd64 = amd64
ARCH_arm64 = arm64

FLAGS = main.go

all: test target

depends:
	go get github.com/stretchr/testify
	go get ./...

target:
	go build -o bin/${BINARY} ${FLAGS}

cross: clean linux windows darwin freebsd

linux: linux-i386 linux-arm linux-amd64 linux-arm64

darwin: darwin-i386 darwin-amd64

windows: windows-i386 windows-amd64

freebsd: freebsd-i386 freebsd-arm freebsd-amd64


linux-i386:
	GOOS=linux GOARCH=${ARCH_i386} go build -o bin/linux/i${ARCH_i386}/${BINARY} ${FLAGS}

linux-arm:
	GOOS=linux GOARCH=${ARCH_arm} go build -o bin/linux/${ARCH_arm}/${BINARY} ${FLAGS}

linux-amd64:
	GOOS=linux GOARCH=${ARCH_amd64} go build -o bin/linux/${ARCH_amd64}/${BINARY} ${FLAGS}

linux-arm64:
	GOOS=linux GOARCH=${ARCH_arm64} go build -o bin/linux/${ARCH_arm64}/${BINARY} ${FLAGS}


darwin-i386:
	GOOS=darwin GOARCH=${ARCH_i386} go build -o bin/darwin/i${ARCH_i386}/${BINARY} ${FLAGS}

darwin-amd64:
	GOOS=darwin GOARCH=${ARCH_amd64} go build -o bin/darwin/${ARCH_amd64}/${BINARY} ${FLAGS}


windows-i386:
	GOOS=windows GOARCH=${ARCH_i386} go build -o bin/windows/i${ARCH_i386}/${BINARY}.exe ${FLAGS}

windows-amd64:
	GOOS=windows GOARCH=${ARCH_amd64} go build -o bin/windows/${ARCH_amd64}/${BINARY}.exe ${FLAGS}


freebsd-i386:
	GOOS=freebsd GOARCH=${ARCH_i386} go build -o bin/freebsd/i${ARCH_i386}/${BINARY}.out ${FLAGS}

freebsd-arm:
	GOOS=freebsd GOARCH=${ARCH_arm} go build -o bin/freebsd/${ARCH_arm}/${BINARY}.out ${FLAGS}

freebsd-amd64:
	GOOS=freebsd GOARCH=${ARCH_amd64} go build -o bin/freebsd/${ARCH_amd64}/${BINARY}.out ${FLAGS}


test:
	go test ./src/cli
	go test ./src/consensus
	go test ./src/db
	go test ./src/network/protocol
	go test ./src/network/services
	go test ./src/network/static
	go test ./src/network/util
	go test ./src/primitives
	go test ./src/primitives/tx_io
	go test ./src/utils
	go test ./src/wallet

clean:
	-rm -rf bin/

renewchain:
	make build
	cp Genesis.db BlockChain_3000.db
	cp Genesis.db BlockChain_3001.db
	clear
	./bin/linux/amd64/blockchain startnode
