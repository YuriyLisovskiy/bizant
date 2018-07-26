Go Libraries
---------------------
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

Submodules
---------------------
Use `git` to clone all submodules:
```bash
$ git submodule init
$ git submodule update
```
