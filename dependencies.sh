#!/usr/bin/env bash

if hash go 2>/dev/null
then
    echo Downloading Golang required libraries...
    go get github.com/stretchr/testify
    go get -u golang.org/x/crypto/...
    echo Done.
else
    echo install.sh: can\'t download go libraries$'\n'$'\t'reason: Golang is not installed, please install it and try again
fi
