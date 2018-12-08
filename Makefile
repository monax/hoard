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
REPO := $(shell pwd)
GOFILES_NOVENDOR := $(shell find . -path ./vendor -prune -o -name '*.pb.go' -prune -o -type f -name '*.go' -print)
PACKAGES_NOVENDOR := $(shell go list ./...)

# Protobuf generated go files
PROTO_FILES = $(shell find . -path ./vendor -prune -o -path ./hoard-js/node_modules -prune -o -type f -name '*.proto' -print)
PROTO_GO_FILES = $(patsubst %.proto, %.pb.go, $(PROTO_FILES))
PROTO_GO_FILES_REAL = $(shell find . -path ./vendor -prune -o -type f -name '*.pb.go' -print)

OS_ARCHS := "linux/arm linux/386 linux/amd64 darwin/386 darwin/amd64 windows/386 windows/amd64"
DIST := "dist"
GOX_OUTPUT := "$DIST/{{.Dir}}_{{.OS}}_{{.Arch}}"
BUILD_IMAGE := "quay.io/monax/hoard:build"


# Formatting, linting and vetting

## check the code for style standards; currently enforces go formatting.
.PHONY: check
check:
	@echo "Checking code for formatting style compliance."
	@goimports -l -d ${GOFILES_NOVENDOR}
	@goimports -l ${GOFILES_NOVENDOR} | read && echo && echo "Your marmot has found a problem with the formatting style of the code." 1>&2 && exit 1 || true

## just fix it
.PHONY: fix
fix:
	@goimports -l -w ${GOFILES_NOVENDOR}

## lint installs golint and prints recommendations for coding style.
lint:
	@echo "Running lint checks."
	go get -u github.com/golang/lint/golint
	@for file in $(GOFILES_NOVENDOR); do \
		echo; \
		golint --set_exit_status $${file}; \
	done

# Dependency Management

## erase vendor wipes the full vendor directory
.PHONY: erase_vendor
erase_vendor:
	rm -rf ${REPO}/vendor/

## install vendor uses dep to install vendored dependencies
.PHONY: reinstall_vendor
reinstall_vendor: erase_vendor
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure -v

## delete the vendor directly and pull back using dep lock and constraints file
## will exit with an error if the working directory is not clean (any missing files or new
## untracked ones)
.PHONY: ensure_vendor
ensure_vendor: reinstall_vendor
	@scripts/is_checkout_dirty.sh


# Building

## output commit_hash but only if we have the git repo (e.g. not in docker build)
.PHONY: commit_hash
commit_hash:
	@git status &> /dev/null && scripts/commit_hash.sh > commit_hash.txt || true

# Protobuffing

## compile hoard.proto interface definition
%.pb.go: %.proto
	protoc -I protobuf -I vendor $< --gogo_out=plugins=grpc:${GOPATH}/src

.PHONY: protobuf
protobuf: $(PROTO_GO_FILES)

.PHONY: clean_protobuf
clean_protobuf:
	@rm -f $(PROTO_GO_FILES_REAL)

.PHONY: protobuf_deps
protobuf_deps:
	@go get -u github.com/gogo/protobuf/protoc-gen-gogo

## build the hoard binary
.PHONY: build_hoard
build_hoard:
	@go build -o bin/hoard ./cmd/hoard

## build the hoard binary
.PHONY: build_hoarctl
build_hoarctl:
	@go build -o bin/hoarctl ./cmd/hoarctl

.PHONY: install
install:
	@go install ./cmd/hoard
	@go install ./cmd/hoarctl

## build all targets in github.com/monax/hoard
.PHONY: build
build:	check build_hoard build_hoarctl protobuf

.PHONY: docker_build
docker_build: check commit_hash
	@scripts/build_tool.sh

## build binaries for all architectures
.PHONY: build_dist
build_dist:	protobuf
	@goreleaser --rm-dist --skip-publish --skip-validate


# Testing

.PHONY:	test
test: check protobuf
	@scripts/bin_wrapper.sh go test -v ./... ${GOPACKAGES_NOVENDOR}

## run tests including integration tests
.PHONY:	test_integration
test_integration: check protobuf
	@go test -v -tags integration ./... ${GOPACKAGES_NOVENDOR}
	@integration/test_gcp.sh
	@integration/test_aws.sh

# Clean Up

## clean removes the target folder containing build artefacts
.PHONY: clean
clean:
	-rm -r ./bin

## Release and Versioning

## print version
.PHONY: version
version:
	@go run ./project/cmd/version/main.go

## generate full changelog of all release notes
CHANGELOG.md: project/history.go project/cmd/changelog/main.go
	@go run ./project/cmd/changelog/main.go > CHANGELOG.md

## generated release note for this version
NOTES.md: project/history.go project/cmd/notes/main.go
	@go run ./project/cmd/notes/main.go > NOTES.md

.PHONY: docs
docs: CHANGELOG.md NOTES.md

## tag the current HEAD commit with the current release defined in
## ./release/release.go
.PHONY: tag_release
tag_release: test check docs build
	@scripts/tag_release.sh

.PHONY: release
release: docs check test docker_build
	@scripts/is_checkout_dirty.sh || (echo "checkout is dirty so not releasing!" && exit 1)
	@scripts/release.sh

.PHONY: release_dev
release_dev: test docker_build
	@scripts/release_dev.sh

.PHONY: build_ci_image
build_ci_image:
	docker build -t ${CI_IMAGE} -f ./.circleci/Dockerfile .

.PHONY: push_ci_image
push_ci_image: build_ci_image
	docker push ${CI_IMAGE}
