GOPATH := $(shell pwd)

all: build

build:
	@mkdir -p bin
	@GOPATH=${GOPATH} go build -o bin/mdu src/main/mdu.go

test: build
	@GOPATH=${GOPATH} go test test/multithread_du_test.go

clean:
	@rm -rf bin
	@rm -rf test/tmp
