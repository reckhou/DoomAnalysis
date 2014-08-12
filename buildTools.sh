#!/bin/sh

CURDIR=`pwd`
EXPORTPATH=$CURDIR/bin/client_tool
SRCPATH=$CURDIR/src/tools/main.go

echo "=== Running govet tools to check code validity ==="
go tool vet ./
echo "=== govet ends ==="

gofmt -w=true -tabwidth=2 -tabs=false $CURDIR
go build -v -o $EXPORTPATH $SRCPATH

EXPORTPATH=$CURDIR/bin/client_tool_linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o $EXPORTPATH $SRCPATH
