#!/bin/sh

VERSION=$1
PROJECTNAME=$2
UUID=$3


if [ $4 = "c" ]; then
cd "./$PROJECTNAME/dump/$VERSION"
tar -jcvf ./$UUID.tar.bz2 ./$UUID.*
rm -f ./$UUID.txt
rm -f ./$UUID.zip
rm -f ./$UUID.txt.info
else
cd "./$PROJECTNAME/tencentdump/"
unzip ./$UUID.zip
unzip ./tomb.zip
rm -f ./cpuinfo.txt
rm -f ./log.txt
rm -f ./detail.txt
rm -f ./tomb.zip
rm -f ./$UUID.zip
path1=$(find ./ -name "tomb_*.txt" | head -n 1)
mkdir -p ../dump/$VERSION
mv ./$path1 ../dump/$VERSION/$UUID.txt
fi


