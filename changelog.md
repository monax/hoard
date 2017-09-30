# Hoard Changelog
## Version 0.1.1
Include hoarctl in Docker image

## Version 0.1.0
Release adding environment config and docker image
- Adds --env flag to read JSON config from HOARD_JSON_CONFIG
- Add --json and --toml flags to &#39;hoard init&#39; to generate JSON optionally
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


