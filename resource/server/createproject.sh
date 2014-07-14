#!/bin/sh

PROJECTNAME=$1
PROJETYPE_1=$2

mkdir -p "$PROJECTNAME/lib"
mkdir -p "$PROJECTNAME/dump"

if [ $2 = "a" ]; then
mkdir -p $PROJECTNAME"_java/lib"
mkdir -p $PROJECTNAME"_java/dump"
mkdir -p $PROJECTNAME"_js/lib"
mkdir -p $PROJECTNAME"_js/dump"
fi

