Build Notes
===========
Some notes on how to build [Block Chain Go](https://github.com/YuriyLisovskiy/blockchain-go).

Testing
---------------------
Run tests:
```bash
$ make test
```

Build
---------------------
Build binary for target platform:
```bash
$ make target
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

Submodules
---------------------
Module      | Purpose                  | Description
------------|--------------------------|----------------------
secp256k1   | secp256k1 elliptic curve | Optimized C library for EC operations on curve secp256k1

#### Install dependencies
You can install required go libraries and submodules manually, see [manual-installation.md](manual-installation.md),
or use [install-deps.sh](../install-deps.sh) script:
```bash
$ make dependencies
```
