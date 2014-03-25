#!/bin/sh

CURDIR=`pwd`
EXPORTPATH=$CURDIR/bin/doomAnalysis
SRCPATH=$CURDIR/main.go

echo "=== Running govet tools to check code validity ==="
go tool vet ./
echo "=== govet ends ==="

gofmt -w=true -tabwidth=2 -tabs=false $CURDIR
go build -v -o $EXPORTPATH $SRCPATH

