#!/usr/bin/make -f

SRC=$(shell find . -type f -name '*.go' ! -name '*_test.go')
SRC_TESTS=$(shell find . -path './vendor/*' -prune -o -name '*_test.go' -print)
TEST_DIRS=$(sort $(dir $(SRC_TESTS)))

all: deps goteleinfo

deps:
	glide install

goteleinfo: $(SRC)
	go build goteleinfo

test: $(SRC_TESTS)
	go test -v $(TEST_DIRS)

clean:
	rm -f goteleinfo

.PHONY: all clean test

	
