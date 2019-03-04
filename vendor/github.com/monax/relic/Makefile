SHELL := /bin/bash
GO_FILES := $(shell go list -f "{{.Dir}}" ./...)

### Formatting, linting and vetting
# Run goimports (also checks formatting) first display output first, then check for success
.PHONY: check
check:
	@go get golang.org/x/tools/cmd/goimports
	@goimports -l -d ${GO_FILES}
	@goimports -l ${GO_FILES} | read && echo && \
	echo "Your marmot has found a problem with the formatting style of the code."\
	 1>&2 && exit 1 || true

# Just fix it
.PHONY: fix
fix:
	@goimports -l -w ${GO_FILES}

# test burrow
.PHONY: test
test: check docs
	@go test ./...

### Release and versioning
.PHONY: version
version:
	@go run ./project/cmd/version/main.go

# Generate full changelog of all release notes
CHANGELOG.md: ./project/releases.go history.go
	@go run ./project/cmd/changelog/main.go > CHANGELOG.md

# Generated release notes for this version
NOTES.md: ./project/releases.go history.go
	@go run ./project/cmd/notes/main.go > NOTES.md

.PHONY: docs
docs: CHANGELOG.md NOTES.md

# Tag a release a push it
.PHONY: tag_release
tag_release: test check docs
	@scripts/tag_release.sh
