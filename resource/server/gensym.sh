#!/bin/sh

VERSION=$1
PROJECTNAME=$2
LIANYUN=$3
INPUTLIBNAME=$4
OUTPUTLIBNAME=$5


cd "./$PROJECTNAME/dump/$VERSION"
cp "../../lib/${VERSION}"_"$INPUTLIBNAME" ./"$OUTPUTLIBNAME"

../../../tools/dump_syms "$OUTPUTLIBNAME" > "$OUTPUTLIBNAME".sym
rm -f "$OUTPUTLIBNAME"

cd "../../lib"
touch "$VERSION".txt

