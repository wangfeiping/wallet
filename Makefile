PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

APP_NAME=wallet

# TODO: Update the ldflags
ldflags = -X github.com/wangfeiping/wallet/wallet/version.Name=$(APP_NAME) \
	-X github.com/wangfeiping/wallet/wallet/version.Version=$(VERSION) \
	-X "github.com/wangfeiping/wallet/wallet/version.VersionCosmos=github.com/cosmos/cosmos-sdk v0.34.4-0.20200829041113-200e88ba075b" \
	-X "github.com/wangfeiping/wallet/wallet/version.VersionEthereum=github.com/ethereum/go-ethereum v1.9.18" \
	-X github.com/wangfeiping/wallet/wallet/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

all: build

build: go.sum
		go build $(BUILD_FLAGS) -o ./build/$(APP_NAME) ./wallet/cmd/walletcli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

# Uncomment when you have some tests
# test:
# 	@go test -mod=readonly $(PACKAGES)

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify
