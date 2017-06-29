# Hoard

Hoard is an encrypted blob store for the Monax platform

![hoarding marmot](docs/images/hoard.jpg)

## Building

To build Hoard you will need to have the following installed:
- The Go language (with $GOPATH/bin in $PATH)
- GNU make
- [Protocol Buffers 3](https://github.com/google/protobuf/releases/tag/v3.3.0)

Then, from the project root run:

```shell
# Install protobuf GRPC plugin, glide, and glide dependencies
make deps
# Run checks, tests, and build binaries
make build
```