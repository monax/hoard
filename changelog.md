# Hoard Changelog
## Version 1.0.2
Improve success/failure logging of LoggingStore.

## Version 1.0.1
Add S3 integration test and include ca-certificates to Docker image so TLS (and S3) works.

## Version 1.0.0
Minor breaking change in that 'hoard init' becomes 'hoard config':
- 'hoard config' adds some niceties for printing JSON config for --env configuration source
- Added S3 'remote' credentials provider enabling credentials to be sourced from EC2 instance roles (note since [RemoteCredProvider()](https://github.com/aws/aws-sdk-go/blob/5a2026bfb28e86839f9fcc46523850319399006c/aws/defaults/defaults.go#L108) is used it also support ECS configuration via AWS_CONTAINER_CREDENTIALS_RELATIVE_URI and AWS_CONTAINER_CREDENTIALS_FULL_URI)


## Version 0.1.1
Include hoarctl in Docker image

## Version 0.1.0
Release adding environment config and docker image
- Adds --env flag to read JSON config from HOARD_JSON_CONFIG
- Add --json and --toml flags to 'hoard init' to generate JSON optionally
- Added alpine based docker image pushed on releases (that reads config from environment variable)


## Version 0.0.2
Bug fix release for FileSystemStorage:
- Switch to URL and filesystem compliant base64 alphabet so some addresses do not target non-existent directories
- Create root directory for FileSystemStorage if it does not exist


## Version 0.0.1
This is the first Hoard open source release and includes:
- Deterministic encryption scheme
- GRPC API for encryption, storage, and cleartext
- Memory, Filesystem, and S3 storage backends
- Configuration
- Hoar-Daemon hoard
- Hoar-Control hoarctl CLI


