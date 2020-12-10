package project

import (
	"github.com/monax/relic"
)

// Can be used to set the commit hash version of the binary at build time with:
// `go build -ldflags "-X github.com/monax/hoard/project.commit=$(git rev-parse --short HEAD)" ./cmd/hoard`

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

// Use below as template for change notes, delete empty sections but keep order
/*
### Security
### Changed
### Fixed
### Added
### Removed
### Deprecated
*/

// The releases described by version string and changes, newest release first.
// The current release is taken to be the first release in the slice, and its
// version determines the single authoritative version for the next release.
//
// To cut a new release add a release to the front of this slice then run the
// release tagging script: ./scripts/tag_release.sh
var History relic.ImmutableHistory = relic.NewHistory("Monax Hoard", "https://github.com/monax/hoard").
	MustDeclareReleases("9.0.0",
		`### Added
- [Hoard] Grants store the refs to their chunks in a new LINK ref type that is followed during dereferencing (Get, UnsealGet, Decrypt). Version 4 grants _always_ store a single LINK ref. This LINK ref is guaranteed to be unique to the Grant (and therefore grants are now unique). This means UnsealDelete can be safely called without the risk of deleting data still referenced by other grants. This also means grants are not linear in the number of chunks used to store them which keeps grants in constant size.
- [Hoard] Refs now store the size of the plaintext data stored behind them. This allows for easier random access and predictable downloads.

### Changed
- [JS] Move to pure-js @grpc/grpc-js library
- [JS] Expose more usable methods from the client (breaking)
- [JS] Support streaming versions of calls taking BytesLike
- [Hoard] Make default storage ChunkSize 3 MiB
- [Hoard] Body can be sent in same message as MustPlaintext Header (but Header data will be normalised out into first message on storage and retrieval)


### Fixed
- [Hoard] Blocking read-then-write write-then-read usage
- [Hoard] Unnecessary copying for streams
- [Hoard] Encrypt endpoints not chunking
- [Cloud] Stat now explicitly checks for NotFound error for Exists flag, and throws other errors

`,
		"8.2.3 - 2020-03-27",
		`### Added
- ReadHeader will stop once it gets the head
`,
		"8.2.2 - 2020-03-25",
		`### Fixed
- Chunk plaintext out if data too big
`,
		"8.2.1 - 2020-03-23",
		`### Fixed
- Convergent encryption nonce is now specified
`,
		"8.2.0 - 2020-03-20",
		`### Changed
- Grant json now uses lowercase field names for compatability with client lib
`,
		"8.1.0 - 2020-03-09",
		`### Changed
- Header now specifies arbitrary data payload
`,
		"8.0.2 - 2020-03-09",
		`### Fixed
- NPM package.json now includes proto defs and README
`,
		"8.0.1 - 2020-03-09",
		`### Fixed
- NPM publish via auth token
`,

		"8.0.0 - 2020-02-24",
		`### Changed
- Service now encrypts per chunk of plaintext
- Grants (v2) contain an array of references

### Removed
- Document service - metadata now in header
`,
		"7.2.0 - 2019-12-17",
		`### Fixed
- Symmetric grants are now versioned to preserve backwards compatability
`,
		"7.1.0 - 2019-12-12",
		`### Changed
- [CMD] Secret keys can now be loaded from the environment

### Fixed
- [ENCRYPTION] Secret keys are no longer derived at runtime due to high scrypt memory overhead
`,
		"7.0.0 - 2019-12-02",
		`This release makes some changes to the Hoard protobuf and service that are backwards compatible for clients - Hoard v6 clients should work with Hoard v7 but hoard-js v7 will not work entirely correctly with Hoard v6 due to removal of oneof.

### Changed
- [API] Drop use of oneof in protobuf files - allow singleton fields to be sent with streamable fields
- [API] Enforce that we only receive exactly one salt and grant spec in streams and that they come first
- [NODEJS] Expose streaming promise client-side API to take advantage of streaming rather than loading entire file into buffer

### Fixed
- Ignoring Spec if Salt present in single message
`,
		"6.0.0 - 2019-10-11",
		`### Added
- Document streaming service
`,
		"5.1.0 - 2019-09-08",
		`### Added
- Ability to delete files located at address
- Stream all files to overcome GRPC message limit
`,
		"5.0.1 - 2019-06-20",
		`### Fixed
- JS client - v5, npm publish
`,
		"5.0.0 - 2019-05-24",
		`This breaking changes refactors the exported API to make it possible to have a much more minimal import tree. Not all storage backends are imported when depending just on api (containing protobuf generated code) or on hoard/v5 root package which allows, for instance, importing the root package to run an in memory test server without all the storage backend dependencies.

### Changed
- Renamed services package to api
- Move services.NewHoardServer to hoard.NewServer
- Renamed storage package to stores
- Made ipfs and cloud their own subpackages to avoid massive import tree
`,
		"4.0.0 - 2019-05-21",
		`### Fixed
- [BUILD] Change hoard.pb.go to services/api.pb.go
`,
		"3.2.1 - 2019-04-24",
		`### Fixed
- [RELEASE] Push latest tag with version tag and perform release on CI
`,
		"3.2.0 - 2019-04-23",
		`### Fixed
- [MODULES] Add v3 to module declaration and update imports

### Removed
- [BUILD] Remove vendor and related scripting
`,
		"3.1.0 - 2019-04-23",
		`### Added
- [SERVER] Added Wait() function to wait until server is ready and ListenAddress for getting bound listen address (useful when using localhost:0 for a OS selected free port)
`,
		"3.0.1 - 2019-03-01",
		`### Added
- [CLI] Optional YAML configuration`,

		"3.0.0 - 2019-02-26",
		`### Changed
- [PROTO] Renamed symmetric grant SecretID to PublicID
- [PROTO] Renamed openpgp grant ID to PrivateID

### Fixed
- [GRANTS] Throw an exception if symmetric secret for ID cannot be found

### Added
- [NODEJS] Added integration tests including test for symmetric secrets
- [GRANTS] Added openpgp grants example
- [CLI] Added ability to configure secrets on command line with hoard config <config> --secret`,
		"2.0.0 - 2019-02-21",
		`### Changed
- [PROTO] Upper case field names in protobuf
- [PROTO] Used gogoproto for types
- [STORE] Go-cloud project for s3 store backend
- [STORE] Minimized IPFS configuration

### Added
- [STORE] Initial support for Azure backend thanks to go-cloud
- [GRANTS] Interface and GRPC service
- [GRANTS] Symmetric AES-GCM-based grants
- [GRANTS] Asymmetric OpenPGP support
- [GRANTS] Go-client (hoarctl) tooling
- [NODEJS] Support for using the grant service + examples
- [STORE] IPFS integration test`,

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
	- Added S3 'remote' credentials provider enabling credentials to be sourced from EC2 instance roles (note since [RemoteCredProvider()](https://github.com/aws/aws-sdk-go/blob/5a2026bfb28e86839f9fcc46523850319399006c/aws/defaults/defaults.go#L108) is used it also support ECS configuration via AWS_CONTAINER_CREDENTIALS_RELATIVE_URI and AWS_CONTAINER_CREDENTIALS_FULL_URI)`,

		"0.1.1",
		`Include hoarctl in Docker image`,

		"0.1.0",
		`Release adding environment config and docker image
	- Adds --env flag to read JSON config from HOARD_JSON_CONFIG
	- Add --json and --toml flags to 'hoard init' to generate JSON optionally
	- Added alpine based docker image pushed on releases (that reads config from environment variable)`,

		"0.0.2",
		`Bug fix release for FileSystemStorage:
	- Switch to URL and filesystem compliant base64 alphabet so some addresses do not target non-existent directories
	- Create root directory for FileSystemStorage if it does not exist`,

		"0.0.1",
		`This is the first Hoard open source release and includes:
	- Deterministic encryption scheme
	- GRPC API for encryption, storage, and cleartext
	- Memory, Filesystem, and S3 storage backends
	- Configuration
	- Hoar-Daemon hoard
	- Hoar-Control hoarctl CLI`,
	)
