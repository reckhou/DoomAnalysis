#!/bin/sh

VERSION=$1
PROJECTNAME=$2
UUID=$3

cd "./$PROJECTNAME/dump/$VERSION"
if [ $4 = "c" ]; then
tar -jcvf ./$UUID.tar.bz2 ./$UUID.*
rm -f ./$UUID.log
rm -f ./$UUID.txt
rm -f ./$UUID.txt.info
rm -f ./$UUID.txt.ndk.info
rm -f ./$UUID.txt.ndk
else
tar -xjvf $UUID.tar.bz2
rm -f ./$UUID.tar.bz2
fi

