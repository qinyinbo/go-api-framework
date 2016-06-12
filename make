#!/bin/bash
if [ ! -f make ]; then
    echo 'make must be run within its container folder' 1>&2
    exit 1
fi

ROOT_DIR=`pwd`
APP=`basename $ROOT_DIR`

go build -o bin/$APP
