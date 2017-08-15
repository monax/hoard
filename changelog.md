# Hoard Changelog
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


