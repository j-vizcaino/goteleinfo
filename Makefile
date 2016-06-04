#!/usr/bin/make -f

all: goconso

goconso: goconso.go
	go build goconso

test:
	go test -v

clean:
	rm -f goconso

.PHONY: all clean test
