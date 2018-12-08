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
	if commit != "" && date != "" {
		return version + "+commit." + commit + "+date." + date
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
	MustDeclareReleases("",
		`Changed:
- Switched to upper case field names in protobuf
- Used gogoproto for types

Added:
- Grant interface and GRPC service
- Symmetric AES-GCM-based grants
`,
		"1.1.5 - 2018-10-17",
		`Scripted integration tests, better makefile and ci configs, gcs creds read from env var.`,

		"1.1.4",
		`IPFS & GCP Support`,

		"1.1.3",
		`Just create new hasher each call of addresses - we only use SHA256 and this operation is cheap`,

		"1.1.2",
		`Upgrade all Go dependencies`,

		"1.1.1",
		`Bump docker image Alpine Linux version to 3.8 and Go to 1.10.3`,

		"1.1.0",
		`Fix unsafe concurrent access of hash.Hash function in makeAddresser with sync.Pool`,

		"1.0.2",
		`Improve success/failure logging of LoggingStore.`,

		"1.0.1",
		`Add S3 integration test and include ca-certificates to Docker image so TLS (and S3) works.`,

		"1.0.0",
		`Minor breaking change in that 'hoard init' becomes 'hoard config':
	- 'hoard config' adds some niceties for printing JSON config for --env configuration source
	- Added S3 'remote' credentials provider enabling credentials to be sourced from EC2 instance roles (note since [RemoteCredProvider()](https://github.com/aws/aws-sdk-go/blob/5a2026bfb28e86839f9fcc46523850319399006c/aws/defaults/defaults.go#L108) is used it also support ECS configuration via AWS_CONTAINER_CREDENTIALS_RELATIVE_URI and AWS_CONTAINER_CREDENTIALS_FULL_URI)
`,
		"0.1.1",
		`Include hoarctl in Docker image`,

		"0.1.0",
		`Release adding environment config and docker image
	- Adds --env flag to read JSON config from HOARD_JSON_CONFIG
	- Add --json and --toml flags to 'hoard init' to generate JSON optionally
	- Added alpine based docker image pushed on releases (that reads config from environment variable)
`,
		"0.0.2",
		`Bug fix release for FileSystemStorage:
	- Switch to URL and filesystem compliant base64 alphabet so some addresses do not target non-existent directories
	- Create root directory for FileSystemStorage if it does not exist
	`,
		"0.0.1",
		`This is the first Hoard open source release and includes:
	- Deterministic encryption scheme
	- GRPC API for encryption, storage, and cleartext
	- Memory, Filesystem, and S3 storage backends
	- Configuration
	- Hoar-Daemon hoard
	- Hoar-Control hoarctl CLI
	`,
	)
