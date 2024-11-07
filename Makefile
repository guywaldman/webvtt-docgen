OUT=./build
SRC=.
SRC_CLI=$(SRC)/main.go
BINARY=webvtt-docgen

.PHONY: run
run:
	@go run $(SRC_CLI)

.PHONY: build
build:
	@go build  -o $(OUT)/$(BINARY) $(SRC_CLI)

.PHONY: test
test:
	@go test ./... -v

.PHONY: lint
lint:
	@golangci-lint run ./...