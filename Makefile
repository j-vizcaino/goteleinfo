#!/usr/bin/make -f

all: goconso

goconso: goconso.go
	go build goconso

clean:
	rm -f goconso

.PHONY: all clean
