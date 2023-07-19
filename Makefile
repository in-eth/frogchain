#!/usr/bin/make -f
PACKAGES=$(shell go list ./...)
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf:1.9.0

###############################################################################
###                           Install                                       ###
###############################################################################
install: go.sum
		@echo "--> Installing frogchaind"
		@go install ./cmd/frogchaind

install-debug: go.sum
	go build -gcflags="all=-N -l" ./cmd/frogchaind

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

test:
	@go test -mod=readonly $(PACKAGES) -cover

lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify

build:
	@mkdir -p build
	@go build -o build/frogchaind ./cmd/frogchaind

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.11.6
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=main

CMT_URL             = https://raw.githubusercontent.com/cometbft/cometbft/v0.37.0/proto/tendermint
GOGO_PROTO_URL      = https://raw.githubusercontent.com/regen-network/protobuf/cosmos
ICS23_PROTO_URL     = https://raw.githubusercontent.com/cosmos/ics23/v0.7.1
SDK_PROTO_URL 		= https://raw.githubusercontent.com/cosmos/cosmos-sdk/v0.47.0-rc3/proto/cosmos

CMT_CRYPTO_TYPES    = third_party/proto/tendermint/crypto
CMT_ABCI_TYPES      = third_party/proto/tendermint/abci
CMT_TYPES           = third_party/proto/tendermint/types
CMT_VERSION         = third_party/proto/tendermint/version
CMT_LIBS            = third_party/proto/tendermint/libs/bits
CMT_P2P             = third_party/proto/tendermint/p2p

SDK_QUERY 			= third_party/proto/cosmos/base/query/v1beta1
SDK_BASE 			= third_party/proto/cosmos/base/v1beta1
SDK_UPGRADE			= third_party/proto/cosmos/upgrade

GOGO_PROTO_TYPES    = third_party/proto/gogoproto
ICS23_TYPES         = third_party/proto/confio

proto-update-deps:
	@mkdir -p $(GOGO_PROTO_TYPES)
	@curl -sSL $(GOGO_PROTO_URL)/gogoproto/gogo.proto > $(GOGO_PROTO_TYPES)/gogo.proto

	@mkdir -p $(SDK_QUERY)
	@curl -sSL $(SDK_PROTO_URL)/base/query/v1beta1/pagination.proto > $(SDK_QUERY)/pagination.proto

	@mkdir -p $(SDK_BASE)
	@curl -sSL $(SDK_PROTO_URL)/base/v1beta1/coin.proto > $(SDK_BASE)/coin.proto

	@mkdir -p $(SDK_UPGRADE)
	@curl -sSL $(SDK_PROTO_URL)/upgrade/v1beta1/upgrade.proto > $(SDK_UPGRADE)/v1beta1/upgrade.proto

## Importing of tendermint protobuf definitions currently requires the
## use of `sed` in order to build properly with cosmos-sdk's proto file layout
## (which is the standard Buf.build FILE_LAYOUT)
## Issue link: https://github.com/tendermint/tendermint/issues/5021
	@mkdir -p $(CMT_TYPES)
	@curl -sSL $(CMT_URL)/types/types.proto > $(CMT_TYPES)/types.proto
	@curl -sSL $(CMT_URL)/types/validator.proto > $(CMT_TYPES)/validator.proto

	@mkdir -p $(CMT_VERSION)
	@curl -sSL $(CMT_URL)/version/types.proto > $(CMT_VERSION)/types.proto

	@mkdir -p $(CMT_LIBS)
	@curl -sSL $(CMT_URL)/libs/bits/types.proto > $(CMT_LIBS)/types.proto

	@mkdir -p $(CMT_CRYPTO_TYPES)
	@curl -sSL $(CMT_URL)/crypto/proof.proto > $(CMT_CRYPTO_TYPES)/proof.proto
	@curl -sSL $(CMT_URL)/crypto/keys.proto > $(CMT_CRYPTO_TYPES)/keys.proto

	@mkdir -p $(ICS23_TYPES)
	@curl -sSL $(ICS23_PROTO_URL)/proofs.proto > $(ICS23_TYPES)/proofs.proto

## insert go package option into proofs.proto file
## Issue link: https://github.com/confio/ics23/issues/32
	@sed -i '4ioption go_package = "github.com/cosmos/ics23/go";' $(ICS23_TYPES)/proofs.proto

.PHONY: proto-all proto-gen proto-gen-any proto-swagger-gen proto-format proto-lint proto-check-breaking proto-update-deps

###############################################################################
###                                Initialize                               ###
###############################################################################

init-hermes: kill-dev install 
	@echo "Initializing both blockchains..."
	./network/init.sh
	./network/start.sh
	@echo "Initializing relayer..." 
	./network/hermes/restore-keys.sh
	./network/hermes/create-conn.sh

start: 
	@echo "Starting up test network"
	./network/start.sh

start-hermes:
	./network/hermes/start.sh

kill-dev:
	@echo "Killing frogchaind and removing previous data"
	-@rm -rf ./data
	-@killall frogchaind 2>/dev/null
