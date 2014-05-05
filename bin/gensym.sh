#!/bin/sh

VERSION=$1
PROJECTNAME=$2
LIANYUN=$3

if [ $3 = "sxda_tr" ]; then
cd "./$PROJECTNAME/dump/$VERSION"
cp "../../lib/sxda_tr/${VERSION}"_libgame.so ./libgame.so
elif [ $3 = "sxda_kr" ]; then
cd "./$PROJECTNAME/dump/$VERSION"
cp "../../lib/sxda_kr/${VERSION}"_libgame.so ./libgame.so
else
cd "./$PROJECTNAME/dump/$VERSION"
cp "../../lib/${VERSION}"_libgame.so ./libgame.so
fi

../../../tools/dump_syms libgame.so > libgame.so.sym
rm -f libgame.so

cd "../../lib"
touch "$VERSION".txt

