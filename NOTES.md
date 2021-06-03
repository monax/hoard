### Fixed
- [JS] Streaming functions in JS client would swallow all GRPC errors and instead throw on a null exception on getHead for the first frame of messages, now we wait for error message and reject with that message

### Added
- [JS] Convenience methods for serialising and deserialising grants to base64 so grants can be treated as opaque identifiers

