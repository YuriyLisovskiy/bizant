#!/usr/bin/env bash

if hash git 2>/dev/null
then
    echo Installing submodules...
    git submodule init
    git submodule update
    echo Done.
else
    echo install.sh: can\'t install submodules$'\n'$'\t'reason: git is not installed, please install it and then try again
fi

if hash go 2>/dev/null
then
    echo Downloading Golang required libraries...
    go get github.com/stretchr/testify
    go get -u golang.org/x/crypto/...
    echo Done.
else
    echo install.sh: can\'t download go libraries$'\n'$'\t'reason: Golang is not installed, please install it and try again
fi
