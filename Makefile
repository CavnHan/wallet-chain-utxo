SHELL := /bin/bash

GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGSSTRING +=-X main.GitVersion=$(GITVERSION)
LDFLAGS :=-ldflags "$(LDFLAGSSTRING)"

wallet-chain-utxo:
	env GO111MODULE=on go build $(LDFLAGS)
.PHONY: wallet-chain-utxo

clean:
	rm wallet-chain-utxo

test:
	go test -v ./...

lint:
	golangci-lint run ./...