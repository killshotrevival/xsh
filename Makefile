# Makefile for building and managing the xsh application
PROJECT = bin/xsh

# Version of the appliction
VERSION = dev


build-mac:
	@echo "Building the app for mac OS"
	mkdir -p bin/
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X 'xsh/cmd.Version=${VERSION}'" -o ${PROJECT}-mac

build-linux:
	@echo "Building the app for linux OS"
	mkdir -p bin/	
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -extldflags '-static' -X 'xsh/cmd.Version=${VERSION}'" -o ${PROJECT}-linux

build: build-mac build-linux
	@echo "Builing for mac and linux arch"

clean:
	# Remove the binaries directory
	rm -rf ./bin/*

lint:
	@echo "Running Golang Lint..."
	golangci-lint run

test:
	@echo "Running Unit tests"
	go test ./... -run=.

verify: test lint 
	@echo "Code verification passed"

gendocs:
	@echo "Generating documentation"
	go run ./... gendocs   

put-host:
	go run ./... put h -i

get-hosts:
	go run ./... get h  