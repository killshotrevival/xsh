# Makefile for building and managing the xsh application
PROJECT = bin/xsh

# Version of the appliction
VERSION = dev

# Detect the operating system
UNAME_S := $(shell uname -s)

ifeq ($(UNAME_S),Linux)
	os_build = build-linux
	os_suffix = linux
endif

ifeq ($(UNAME_S),Darwin)
	os_build = build-mac
	os_suffix = mac
endif

#### Build targets

build-mac:
	@echo "Building the app for mac OS"
	mkdir -p bin/
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X 'xsh/cmd.Version=${VERSION}'" -o ${PROJECT}-mac

build-linux:
	@echo "Building the app for linux OS"
	mkdir -p bin/	
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -extldflags '-static' -X 'xsh/cmd.Version=${VERSION}'" -o ${PROJECT}-linux

build: $(os_build)

build-all: build-linux build-mac
	@echo "Builing the app for all OS"

build-and-replace: $(os_build)
	@echo "Replacing the existing binary"
	mv bin/xsh-$(os_suffix) ~/.local/bin/xsh

clean:
	# Remove the binaries directory
	rm -rf ./bin/*

lint:
	@echo "Running Golang Lint..."
	golangci-lint run

unit-test:
	@echo "Running Unit tests"
	go test ./... -run=.

verify: unit-test lint 
	@echo "Code verification passed"

gendocs:
	@echo "Generating documentation"
	go run ./... gendocs   

put-host:
	go run ./... put h -i

get-hosts:
	go run ./... get h 

integration-test: $(os_build)
	@echo "CLI created successfully, starting integration test"
	bash test/integration.sh $(PWD)/bin/xsh-$(os_suffix) $(PWD)/test/test
	@echo "Integration test passed" 

	