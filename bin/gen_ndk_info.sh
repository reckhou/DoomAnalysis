#!/bin/sh

VERSION=$1
FILENAME=$2
PROJECTNAME=$3
NDK=./android-ndk-r9d

cd "./$PROJECTNAME/dump/$VERSION"
../../../tools/ndk-stack -sym ../../lib/ -dump $FILENAME > $FILENAME.info
