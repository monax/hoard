#      ___           ___           ___           ___
#     /\  \         /\  \         /\  \         /\  \         _____
#     \:\  \       /::\  \       /::\  \       /::\  \       /::\  \
#      \:\  \     /:/\:\  \     /:/\:\  \     /:/\:\__\     /:/\:\  \
#  ___ /::\  \   /:/  \:\  \   /:/ /::\  \   /:/ /:/  /    /:/  \:\__\
# /\  /:/\:\__\ /:/__/ \:\__\ /:/_/:/\:\__\ /:/_/:/__/___ /:/__/ \:|__|
# \:\/:/  \/__/ \:\  \ /:/  / \:\/:/  \/__/ \:\/:::::/  / \:\  \ /:/  /
#  \::/__/       \:\  /:/  /   \::/__/       \::/~~/~~~~   \:\  /:/  /
#   \:\  \        \:\/:/  /     \:\  \        \:\~~\        \:\/:/  /
#    \:\__\        \::/  /       \:\__\        \:\__\        \::/  /
#     \/__/         \/__/         \/__/         \/__/         \/__/
#
# Hoard Makefile
#
# Requires go version 1.8 or later.
#
# To compile gRPC service also requires protobuf 3 and the protobuf go plugin.
# See http://www.grpc.io/docs/quickstart/go.html to get started.
#

SHELL := /bin/bash
GOFILES_NOVENDOR := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOPACKAGES_NOVENDOR := $(shell go list ./... | grep -v /vendor/)

# Install dependencies and also clear out vendor (we should do this in CI)

# to make sure we are not depending on any local changes to dependencies in
# vendor/
.PHONY: ensure_vendor
ensure_vendor:
	@rm -rf vendor
	@glide install

# to make sure we are not depending on any local changes to dependencies in
# vendor/
.PHONY: deps
deps: ensure_vendor
	@go get golang.org/x/tools/cmd/goimports
	@go get -u github.com/golang/protobuf/protoc-gen-go
	@go get -u github.com/Masterminds/glide

# Print version
.PHONY: version
version:
	@go run ./cmd/version/main.go

# Run goimports (also checks formatting) first display output first, then check for success
.PHONY: check
check:
	@goimports -l -d ${GOFILES_NOVENDOR}
	@goimports -l ${GOFILES_NOVENDOR} | read && echo && \
	echo "Your marmot has found a problem with the formatting style of the code."\
	 1>&2 && exit 1 || true

# Just fix it
.PHONY: fix
fix:
	@goimports -l -w ${GOFILES_NOVENDOR}

# Compile hoard.proto interface defition
.PHONY: build_protobuf
build_protobuf:
	@protoc -I ./core core/hoard.proto --go_out=plugins=grpc:core

# Build the hoard binary
.PHONY: build_hoard
build_hoard:
	@go build -o bin/hoard ./cmd/hoard

# Build the hoard binary
.PHONY: build_hoarctl
build_hoarctl:
	@go build -o bin/hoarctl ./cmd/hoarctl

# Build the hoard binaries
.PHONY: build_bin
build_bin:	build_hoard build_hoarctl

# Run tests
.PHONY:	test
test: check build_protobuf
	@go test ${GOPACKAGES_NOVENDOR}

# Run tests for developing (noisy)
.PHONY:	test_dev
test_dev: build_protobuf
	@go test -v ${GOPACKAGES_NOVENDOR}

# Build all the things
.PHONY: build
build:	build_protobuf build_bin

# Do all available tests and checks then build
.PHONY: build_ci
build_ci: ensure_vendor check test build
