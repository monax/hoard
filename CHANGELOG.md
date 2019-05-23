# [Monax Hoard](https://github.com/monax/hoard) Changelog
## [Unreleased]


## [4.0.0] - 2019-05-21
###Fixed
- [BUILD] Change hoard.pb.go to services/services.pb.go


## [3.2.1] - 2019-04-24
### Fixed
- [RELEASE] Push latest tag with version tag and perform release on CI


## [3.2.0] - 2019-04-23
### Fixed
- [MODULES] Add v3 to module declaration and update imports

### Removed
- [BUILD] Remove vendor and related scripting


## [3.1.0] - 2019-04-23
### Added
- [SERVER] Added Wait() function to wait until server is ready and ListenAddress for getting bound listen address (useful when using localhost:0 for a OS selected free port)


## [3.0.1] - 2019-03-01
### Added
- [CLI] Optional YAML configuration

## [3.0.0] - 2019-02-26
### Changed
- [PROTO] Renamed symmetric grant SecretID to PublicID
- [PROTO] Renamed openpgp grant ID to PrivateID

### Fixed
- [GRANTS] Throw an exception if symmetric secret for ID cannot be found

### Added
- [NODEJS] Added integration tests including test for symmetric secrets
- [GRANTS] Added openpgp grants example
- [CLI] Added ability to configure secrets on command line with hoard config <config> --secret

## [2.0.0] - 2019-02-21
### Changed
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
- [STORE] IPFS integration test

## [1.1.5] - 2018-10-17
Scripted integration tests, better makefile and ci configs, gcs creds read from env var.

## [1.1.4]
IPFS & GCP Support

## [1.1.3]
Just create new hasher each call of addresses - we only use SHA256 and this operation is cheap

## [1.1.2]
Upgrade all Go dependencies

## [1.1.1]
Bump docker image Alpine Linux version to 3.8 and Go to 1.10.3

## [1.1.0]
Fix unsafe concurrent access of hash.Hash function in makeAddresser with sync.Pool

## [1.0.2]
Improve success/failure logging of LoggingStore.

## [1.0.1]
Add S3 integration test and include ca-certificates to Docker image so TLS (and S3) works.

## [1.0.0]
Minor breaking change in that 'hoard init' becomes 'hoard config':
	- 'hoard config' adds some niceties for printing JSON config for --env configuration source
	- Added S3 'remote' credentials provider enabling credentials to be sourced from EC2 instance roles (note since [RemoteCredProvider()](https://github.com/aws/aws-sdk-go/blob/5a2026bfb28e86839f9fcc46523850319399006c/aws/defaults/defaults.go#L108) is used it also support ECS configuration via AWS_CONTAINER_CREDENTIALS_RELATIVE_URI and AWS_CONTAINER_CREDENTIALS_FULL_URI)

## [0.1.1]
Include hoarctl in Docker image

## [0.1.0]
Release adding environment config and docker image
	- Adds --env flag to read JSON config from HOARD_JSON_CONFIG
	- Add --json and --toml flags to 'hoard init' to generate JSON optionally
	- Added alpine based docker image pushed on releases (that reads config from environment variable)

## [0.0.2]
Bug fix release for FileSystemStorage:
	- Switch to URL and filesystem compliant base64 alphabet so some addresses do not target non-existent directories
	- Create root directory for FileSystemStorage if it does not exist

## [0.0.1]
This is the first Hoard open source release and includes:
	- Deterministic encryption scheme
	- GRPC API for encryption, storage, and cleartext
	- Memory, Filesystem, and S3 storage backends
	- Configuration
	- Hoar-Daemon hoard
	- Hoar-Control hoarctl CLI

[Unreleased]: https://github.com/monax/hoard/compare/v4.0.0...HEAD
[4.0.0]: https://github.com/monax/hoard/compare/v3.2.1...v4.0.0
[3.2.1]: https://github.com/monax/hoard/compare/v3.2.0...v3.2.1
[3.2.0]: https://github.com/monax/hoard/compare/v3.1.0...v3.2.0
[3.1.0]: https://github.com/monax/hoard/compare/v3.0.1...v3.1.0
[3.0.1]: https://github.com/monax/hoard/compare/v3.0.0...v3.0.1
[3.0.0]: https://github.com/monax/hoard/compare/v2.0.0...v3.0.0
[2.0.0]: https://github.com/monax/hoard/compare/v1.1.5...v2.0.0
[1.1.5]: https://github.com/monax/hoard/compare/v1.1.4...v1.1.5
[1.1.4]: https://github.com/monax/hoard/compare/v1.1.3...v1.1.4
[1.1.3]: https://github.com/monax/hoard/compare/v1.1.2...v1.1.3
[1.1.2]: https://github.com/monax/hoard/compare/v1.1.1...v1.1.2
[1.1.1]: https://github.com/monax/hoard/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/monax/hoard/compare/v1.0.2...v1.1.0
[1.0.2]: https://github.com/monax/hoard/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/monax/hoard/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/monax/hoard/compare/v0.1.1...v1.0.0
[0.1.1]: https://github.com/monax/hoard/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/monax/hoard/compare/v0.0.2...v0.1.0
[0.0.2]: https://github.com/monax/hoard/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/monax/hoard/commits/v0.0.1
