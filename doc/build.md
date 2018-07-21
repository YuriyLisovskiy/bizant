Build Notes
===========
Some notes on how to build [Block Chain Go](https://github.com/YuriyLisovskiy/blockchain-go).

Build
---------------------

```bash
$ make
```
See [`supported-platforms.md`](supported-platforms.md) for more info.

Build binaries for all supported platforms:
```bash
$ make cross
```

This will build [blockchain-go](https://github.com/YuriyLisovskiy/blockchain-go) as well if the dependencies are met.

Minimum version of Go language required: `go1.10`.

See [golang installation](https://golang.org/doc/install) for more info.

Go Libraries
---------------------
These libraries are required:

 Library     | Purpose          | Description
 ------------|------------------|----------------------
 crypto      | Cryptography     | Supplementary Go cryptography libraries
 bolt        | Bolt DB          | An embedded key/value database

#### Cryptography:

The easiest way to install is to run
```bash
$ go get -u golang.org/x/crypto/...
```
You can also manually git clone the repository to `$GOPATH/src/golang.org/x/crypto`.

#### Bolt DB:

This will retrieve the library and install the bolt command line utility into your `GOBIN` path:
```bash
$ go get github.com/boltdb/bolt/...
```
