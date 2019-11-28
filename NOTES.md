This release makes some changes to the Hoard protobuf and service that are backwards compatible for clients - Hoard v6 clients should work with Hoard v7 but hoard-js v7 will not work entirely correctly with Hoard v6 due to removal of oneof.

### Changed
- [API] Drop use of oneof in protobuf files - allow singleton fields to be sent with streamable fields
- [API] Enforce that we only receive exactly one salt and grant spec in streams and that they come first
- [NODEJS] Expose streaming promise client-side API to take advantage of streaming rather than loading entire file into buffer

### Fixed
- Ignoring Spec if Salt present in single message

