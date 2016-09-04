#!/bin/bash

INPUT=database.influx.in

if [ -z "$1" ]
then
	echo "
$(basename $0) renders $INPUT with custom database name
and prints result to stdout.

Usage: $0 dbname
"
	exit 1
fi

cat $INPUT | DBNAME=$1 envsubst

