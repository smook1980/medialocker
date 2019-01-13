#!/usr/bin/env bash

MEDIALOCKER_DATE=`date -u +%y%m%d`
MEDIALOCKER_VERSION=`git describe --always`

if [[ -z $1 ]] || [[ -z $2 ]]; then
    echo "Please provide build mode and output file name" 1>&2
    exit 1
fi

if [[ $OS == "Windows_NT" ]]; then
    MEDIALOCKER_OS=win32
    if [[ $PROCESSOR_ARCHITEW6432 == "AMD64" ]]; then
        MEDIALOCKER_ARCH=amd64
    else
        if [[ $PROCESSOR_ARCHITECTURE == "AMD64" ]]; then
            MEDIALOCKER_ARCH=amd64
        fi
        if [[ $PROCESSOR_ARCHITECTURE == "x86" ]]; then
            MEDIALOCKER_ARCH=ia32
        fi
    fi
else
    MEDIALOCKER_OS=`uname -s`
    MEDIALOCKER_ARCH=`uname -p`
fi

if [[ $1 == "debug" ]]; then
    echo "Building development binary..."
  go build -ldflags "-X main.version=${MEDIALOCKER_DATE}-${MEDIALOCKER_VERSION}-${MEDIALOCKER_OS}-${MEDIALOCKER_ARCH}-DEBUG" -o $2 cmd/locker/locker.go
  du -h $2
  echo "Done."
else
    echo "Building production binary..."
  go build -ldflags "-s -w -X main.version=${MEDIALOCKER_DATE}-${MEDIALOCKER_VERSION}-${MEDIALOCKER_OS}-${MEDIALOCKER_ARCH}" -o $2 cmd/locker/locker.go
  du -h $2
  echo "Done."
fi
