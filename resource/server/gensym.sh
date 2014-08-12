#!/bin/sh

VERSION=$1
PROJECTNAME=$2
LIANYUN=$3
INPUTLIBNAME=$4
OUTPUTLIBNAME=$5


cd "./$PROJECTNAME/dump/$VERSION"
cp "../../lib/${VERSION}"_"$INPUTLIBNAME" ./"$INPUTLIBNAME"

../../../tools/dump_syms "$INPUTLIBNAME" > "$OUTPUTLIBNAME".sym
rm -f "$INPUTLIBNAME"

cd "../../lib"
touch "$VERSION".txt

