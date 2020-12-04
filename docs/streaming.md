# Notes on streaming

Many of the API endpoints accept and/or return streams. Streaming messages are represented by protobuf messages that logically are either 'header type information' or 'chunk type information'. The header type messages must be within a single message and be the first message sent in the stream, the chunks can be repeated forming a stream of values.

Aa object of type `Plaintext` should have exactly one of `Head` or `Body` set on it. Similarly, a `PlaintextAndGrantSpec` should either have a `Spec` and `Plaintext.Header` (optionally) or `Plaintext.Body` set. An error is thrown if `Spec` or `Header` information are provided anywhere other than the first message in a stream. This allows us to normalise our storage and make sure `Header` info is accessible.
