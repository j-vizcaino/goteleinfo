#!/usr/bin/make -f

SRC=$(shell find . -type f -name '*.go' ! -name '*_test.go')
SRC_TESTS=$(shell find . -path './vendor/*' -prune -o -name '*_test.go' -print)
TEST_DIRS=$(sort $(dir $(SRC_TESTS)))

all: goconso

goconso: $(SRC)
	go build goconso

test: $(SRC_TESTS)
	for d in $(TEST_DIRS); do cd $$d; echo "*** Running tests in $$d..."; go test -v; cd - > /dev/null; done

clean:
	rm -f goconso

.PHONY: all clean test $(TESTS)


	
