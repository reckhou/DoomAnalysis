#!/bin/sh

VERSION=$1
FILENAME=$2
NDK=./android-ndk-r9d

cd "./dump/$VERSION"
../../tools/ndk-stack -sym ../../lib/ -dump $FILENAME > $FILENAME.info
