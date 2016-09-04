#!/usr/bin/make -f

SRC=$(shell find . -type f -name '*.go' ! -name '*_test.go')
SRC_TESTS=$(shell find . -path './vendor/*' -prune -o -name '*_test.go' -print)
TEST_DIRS=$(sort $(dir $(SRC_TESTS)))

all: goconso

goconso: $(SRC)
	go build goconso

test: $(SRC_TESTS)
	go test -v $(TEST_DIRS)

clean:
	rm -f goconso

.PHONY: all clean test $(TESTS)


	
