#!/bin/sh

CURDIR=`pwd`
EXPORTPATH=$CURDIR/bin/doomAnalysis_osx
SRCPATH=$CURDIR/main.go

echo "=== Running govet tools to check code validity ==="
go tool vet ./
echo "=== govet ends ==="

gofmt -w=true -tabwidth=2 -tabs=false $CURDIR

EXPORTPATH=$CURDIR/bin/doomAnalysis_linux_amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o $EXPORTPATH $SRCPATH
