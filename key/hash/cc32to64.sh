#!/usr/bin/env bash

fullcwd=`pwd`
basecwd=`basename $fullcwd`
if [ $basecwd != "hash" ]; then
	echo "Script must be run from go-functional-collections/fmap/hash" 1>&2
	exit 1
fi

cp val.go-32 val.go-64
perl -pi -e 's/32/64/g' val.go-64

ln -sf val.go-64 val.go
