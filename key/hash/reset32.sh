#!/usr/bin/env bash

fullcwd=`pwd`
basecwd=`basename $fullcwd`
if [ $basecwd != "hash" ]; then
	echo "Script must be run from go-functional-collections/fmap/hash" 1>&2
	exit 1
fi


ln -sf val.go-32 val.go
