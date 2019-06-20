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

export GO111MODULE := on

SHELL := /bin/bash
REPO := $(shell pwd)
GOFILES := $(shell find . -name '*.pb.go' -prune -o -not -path './vendor/*' -type f -name '*.go' -print)

# Protobuf generated go files
PROTO_FILES = $(shell find . -path ./hoard-js -prune -o -path ./node_modules -prune -o -type f -name '*.proto' -print)
PROTO_GO_FILES = $(patsubst %.proto, %.pb.go, $(PROTO_FILES))
PROTO_GO_FILES_REAL = $(shell find . -type f -name '*.pb.go' -print)

export DOCKER_HUB := quay.io
export DOCKER_REPO := $(DOCKER_HUB)/monax/hoard
export BUILD_IMAGE := $(DOCKER_REPO):build

# Formatting, linting and vetting

## check the code for style standards; currently enforces go formatting.
.PHONY: check
check:
	@echo "Checking code for formatting style compliance."
	@goimports -l -d ${GOFILES}
	@goimports -l ${GOFILES} | read && echo && echo "Your marmot has found a problem with the formatting style of the code." 1>&2 && exit 1 || true

## just fix it
.PHONY: fix
fix:
	@goimports -l -w ${GOFILES}

## lint installs golint and prints recommendations for coding style.
.PHONY: lint
lint:
	@echo "Running lint checks."
	@for file in $(GOFILES); do \
		echo; \
		golint --set_exit_status $${file}; \
	done

# Building

## output commit_hash but only if we have the git repo (e.g. not in docker build)
.PHONY: commit_hash
commit_hash:
	@git status &> /dev/null && scripts/commit_hash.sh > commit_hash.txt || true

# Protobuffing

## compile hoard.proto interface definition
%.pb.go: %.proto
	@mkdir -p .gopath
	protoc -I protobuf $< --gogo_out=plugins=grpc:.gopath

.PHONY: protobuf
protobuf: $(PROTO_GO_FILES)
	rsync -r .gopath/github.com/monax/hoard/v5/ ./
	rm -rf .gopath

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
build:	check build_hoard build_hoarctl

.PHONY: docker_build
docker_build: commit_hash
	@scripts/build_tool.sh

## build binaries for all architectures
.PHONY: build_dist
build_dist:
	@goreleaser --rm-dist --skip-publish --skip-validate

# Testing

.PHONY:	test
test: check
	@scripts/bin_wrapper.sh go test -v ./...

.PHONY: test_js
test_js: build install
	$(eval HID := $(shell hoard config memory -s test:secret_pass | hoard -c- &> /dev/null & echo $$!))
	npm test
	kill ${HID}

## run tests including integration tests
.PHONY:	test_integration
test_integration: check
	@go test -v -tags integration ./...
	@scripts/integration/test_aws.sh
	@scripts/integration/test_gcp.sh
	@scripts/integration/test_ipfs.sh

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
## ./project/history.go and push to remote to trigger actual release
.PHONY: tag_release
tag_release: test check docs build
	@scripts/tag_release.sh

## To be run by CI to effect actual release
.PHONY: release
release:
	@scripts/is_checkout_dirty.sh || (echo "checkout is dirty so not releasing!" && exit 1)
	@scripts/release.sh

.PHONY: build_ci_image
build_ci_image:
	docker build -t ${BUILD_IMAGE} -f ./.circleci/Dockerfile .

.PHONY: push_ci_image
push_ci_image: build_ci_image
	docker push ${BUILD_IMAGE}
