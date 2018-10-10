package project

import (
	"github.com/monax/relic"
)

// Can be used to set the commit hash version of the binary at build time with:
// `go build -ldflags "-X github.com/hyperledger/burrow/project.commit=$(git rev-parse --short HEAD)" ./cmd/burrow`

var commit = ""
var date = ""

func Commit() string {
	return commit
}

func FullVersion() string {
	version := History.CurrentVersion().String()
	if commit != "" {
		return version + "+commit." + commit
	}
	return version
}

// The releases described by version string and changes, newest release first.
// The current release is taken to be the first release in the slice, and its
// version determines the single authoritative version for the next release.
//
// To cut a new release add a release to the front of this slice then run the
// release tagging script: ./scripts/tag_release.sh
var History relic.ImmutableHistory = relic.NewHistory("Monax Hoard", "https://github.com/monax/hoard").
	MustDeclareReleases("0.2.0 - 2018-10-09",
		`### Changed
- IPFS & GCP Support
`,
	)
