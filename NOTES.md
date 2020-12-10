This is a fairly major release and the client APIs change significantly. However Grant v2s are still supported and the protobuf API is backwards-compatible.

### Added
- [Hoard] Grants store the refs to their chunks in a new LINK ref type that is followed during dereferencing (Get, UnsealGet, Decrypt). Version 4 grants _always_ store a single LINK ref. This LINK ref is guaranteed to be unique to the Grant (and therefore grants are now unique). This means UnsealDelete can be safely called without the risk of deleting data still referenced by other grants. This also means grants are not linear in the number of chunks used to store them which keeps grants in constant size.
- [Hoard] Refs now store the size of the plaintext data stored behind them. This allows for easier random access and predictable downloads.
- [Hoard] Added regression test to check grant-to-plaintext compatibility between versions.

### Changed
- [JS] Move to pure-js @grpc/grpc-js library
- [JS] Expose more usable methods from the client (breaking)
- [JS] Support streaming versions of calls taking BytesLike
- [Hoard] Make default storage ChunkSize 3 MiB
- [Hoard] Body can be sent in same message as MustPlaintext Header (but Header data will be normalised out into first message on storage and retrieval)
- [Hoard] Grants now encode their references using Protobuf rather than JSON (backwards compatible with grant V2)


### Fixed
- [Hoard] Blocking read-then-write write-then-read usage
- [Hoard] Unnecessary copying for streams
- [Hoard] Encrypt endpoints not chunking
- [Cloud] Stat now explicitly checks for NotFound error for Exists flag, and throws other errors


