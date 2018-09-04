Build Notes
===========
Some notes on how to build [Block Chain Go](https://github.com/YuriyLisovskiy/blockchain-go).

Before Build
---------------------
Build requires:
* The [Go programming language](https://golang.org), minimum version required: `go1.10`.
* [make](https://www.gnu.org/software/make/manual/make.html) - automation tool to build the project automatically.

#### Install dependencies
You can install required [Go](https://golang.org) libraries and submodules manually, see [manual-installation.md](manual-installation.md),
or use [dependencies.sh](../dependencies.sh) script:
```bash
$ make dependencies
```

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
$ make
```

This will build [blockchain-go](https://github.com/YuriyLisovskiy/blockchain-go) as well if the dependencies are met.

See [golang installation](https://golang.org/doc/install) for more info.

Go Libraries
---------------------
These libraries are required:

 Library     | Purpose          | Description
 ------------|------------------|----------------------
 crypto      | Cryptography     | Supplementary Go cryptography libraries
