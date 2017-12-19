#!/usr/bin/env bash

fullcwd=`pwd`
basecwd=`basename $fullcwd`
if [ $basecwd != "hash" ]; then
	echo "Script must be run from go-functional-collections/fmap/hash" 1>&2
	exit 1
fi

cp hashval.go-32 hashval.go-64
perl -pi -e 's/32/64/g' hashval.go-64

ln -sf hashval.go-64 hashval.go

