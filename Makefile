GO_SRC_FILES=$(shell find "cmd/gohooks" -name "*.go" | tr "\n" " ")

all: build check

.PHONY: build
build: modules
	go build -v -o bin/gohooks $(GO_SRC_FILES)

.PHONY: modules
modules:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	golangci-lint run