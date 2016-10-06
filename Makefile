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

arm: $(SRC)
	env GOOS=linux GOARCH=arm go build goteleinfo
	@mv goteleinfo goteleinfo.arm

clean:
	rm -f goteleinfo

.PHONY: all clean test

