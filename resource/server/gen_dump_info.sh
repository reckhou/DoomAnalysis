#!/bin/sh

VERSION=$1
FILENAME=$2
PROJECTNAME=$3

cd "./$PROJECTNAME/dump/$VERSION"
../../../tools/minidump_stackwalk $FILENAME symbols/ > "$FILENAME.info"

