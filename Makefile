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

# Gets implicit default GOPATH if not set
GOPATH?=$(shell go env GOPATH)
BIN_PATH?=${GOPATH}/bin

SHELL := /bin/bash
REPO := $(shell pwd)
GOFILES := $(shell find . -name '*.pb.go' -prune -o -not -path './vendor/*' -type f -name '*.go' -print)

# Protobuf generated go files
PROTO_FILES = $(shell find . -path ./js -prune -o -path ./node_modules -prune -o -type f -name '*.proto' -print)
PROTO_GO_FILES = $(patsubst %.proto, %.pb.go, $(PROTO_FILES))
PROTO_GO_FILES_REAL = $(shell find . -type f -name '*.pb.go' -print)
PROTO_TS_FILES = $(patsubst %.proto, %.pb.ts, $(PROTO_FILES))

HOARD_TS_PATH = ./js
PROTO_GEN_TS_PATH = ${HOARD_TS_PATH}/proto
PROTOC_GEN_TS_PATH = ${HOARD_TS_PATH}/node_modules/.bin/protoc-gen-ts
PROTOC_GEN_GRPC_PATH= ${HOARD_TS_PATH}/node_modules/.bin/grpc_tools_node_protoc_plugin

GO_BUILD_ARGS = -ldflags "-extldflags '-static' -X $(shell go list)/project.commit=$(shell cat commit_hash.txt) -X $(shell go list)/project.date=$(shell date '+%Y-%m-%d')"

export DOCKER_HUB := quay.io
export DOCKER_REPO := $(DOCKER_HUB)/monax/hoard
export BUILD_IMAGE := $(DOCKER_REPO):build

## Release and Versioning

VERSION := $(shell go run ./project/cmd/version/main.go)

## print version
.PHONY: version
version:
	@echo $(VERSION)

# Formatting, linting and vetting

## check the code for style standards; currently enforces go formatting.
.PHONY: check
check:
	@echo "Checking code for formatting style compliance."
	@gofmt -l -d ${GOFILES}
	@gofmt -l ${GOFILES} | read && echo && echo "Your marmot has found a problem with the formatting style of the code." 1>&2 && exit 1 || true

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

%.pb.ts: %.proto
	protoc -I protobuf \
		--plugin="protoc-gen-ts=${PROTOC_GEN_TS_PATH}" \
		--plugin=protoc-gen-grpc=${PROTOC_GEN_GRPC_PATH} \
		--js_out="import_style=commonjs,binary:${PROTO_GEN_TS_PATH}" \
		--ts_out="service=grpc-node,mode=grpc-js:${PROTO_GEN_TS_PATH}" \
		--grpc_out="grpc_js:${PROTO_GEN_TS_PATH}" $<


.PHONY: protobuf
protobuf: ${PROTO_GO_FILES} ${PROTO_TS_FILES}
	rsync -r .gopath/github.com/monax/hoard/v8/ ./
	rm -rf .gopath

.PHONY: clean_protobuf
clean_protobuf:
	@rm -f $(PROTO_GO_FILES_REAL)

.PHONY: npm_install
npm_install:
	@cd ${HOARD_TS_PATH} && npm install

.PHONY: protobuf_deps
protobuf_deps:
	@go get -u github.com/gogo/protobuf/protoc-gen-gogo
	@cd ${HOARD_TS_PATH} && npm install grpc-tools
	@cd ${HOARD_TS_PATH} && npm install ts-protoc-gen

## build the hoard binary
.PHONY: build_hoard
build_hoard: commit_hash
	go build $(GO_BUILD_ARGS) -o bin/hoard ./cmd/hoard

## build the hoard binary
.PHONY: build_hoarctl
build_hoarctl: commit_hash
	go build $(GO_BUILD_ARGS) -o bin/hoarctl ./cmd/hoarctl

.PHONY: install
install: build_hoarctl build_hoard
	mkdir -p ${BIN_PATH}
	install -T ${REPO}/bin/hoarctl ${BIN_PATH}/hoarctl
	install -T ${REPO}/bin/hoard ${BIN_PATH}/hoard

## build all targets in github.com/monax/hoard
.PHONY: build
build: check build_hoard build_hoarctl

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
	@scripts/bin_wrapper.sh go test -v ./... ${GO_TEST_ARGS}

.PHONY: test_js
test_js: build install npm_install
	@scripts/test_js.sh

## run tests including integration tests
.PHONY:	test_integration
test_integration: check
	@go test -v -tags integration ./...
	@scripts/integration/test_gcp.sh
	@scripts/integration/test_ipfs.sh

# Clean Up

## clean removes the target folder containing build artefacts
.PHONY: clean
clean:
	-rm -r ./bin

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

HELM_PATH?=helm/package
HELM_PACKAGE=$(HELM_PATH)/hoard-$(VERSION).tgz
ARCH?=linux-amd64

# Note --set flag currently needs helm 3 version < 3.0.3 https://github.com/helm/helm/issues/3141 - but hopefully they will reintroduce support
bin/helm:
	@echo Downloading helm...
	mkdir -p bin
	curl https://get.helm.sh/helm-v3.0.2-$(ARCH).tar.gz | tar xvzO $(ARCH)/helm > bin/helm && chmod +x bin/helm

.PHONY: helm_deps
helm_deps: bin/helm
	@bin/helm repo add --username "$(CM_USERNAME)" --password "$(CM_PASSWORD)" chartmuseum $(CM_URL)

.PHONY: helm_test
helm_test: bin/helm
	bin/helm dep up helm/hoard
	bin/helm lint helm/hoard

helm_package: $(HELM_PACKAGE)

$(HELM_PACKAGE): helm_test bin/helm
	bin/helm package helm/hoard \
		--version "$(VERSION)" \
		--app-version "$(VERSION)" \
		--set "image.tag=$(VERSION)" \
		--dependency-update \
		--destination helm/package

.PHONY: helm_push
helm_push: helm_package
	@echo pushing helm chart...
	@curl -u ${CM_USERNAME}:${CM_PASSWORD} \
		--data-binary "@$(HELM_PACKAGE)" $(CM_URL)/api/charts
