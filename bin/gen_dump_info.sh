#!/bin/sh

VERSION=$1
FILENAME=$2

cd "./dump/$VERSION"
../../tools/minidump_stackwalk $FILENAME symbols/ > "$FILENAME.info"

