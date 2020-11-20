# Hoard

![master](https://github.com/monax/hoard/workflows/master/badge.svg)
[![version](https://img.shields.io/github/tag/monax/hoard.svg)](https://github.com/monax/hoard/releases/latest)
[![GoDoc](https://godoc.org/github.com/hoard?status.png)](https://pkg.go.dev/github.com/monax/hoard)
[![license](https://img.shields.io/github/license/monax/hoard.svg)](LICENSE)
[![LoC](https://tokei.rs/b1/github/monax/hoard?category=lines)](https://github.com/monax/hoard)
[![codecov](https://codecov.io/gh/monax/hoard/branch/master/graph/badge.svg)](https://codecov.io/gh/monax/hoard)

Hoard is a stateless, deterministically encrypted, content-addressed object store.

![hoarding marmot](docs/images/hoard.jpg)

## Introduction
It convergently encrypts an object using its (SHA256) hash as the secret key (which can than be shared as a 'grant').
The address is then deterministically generated from the encrypted object's (SHA256) digest and allocated to the configured storage back-end:

- In-Memory
- Filesystem
- [AWS](https://aws.amazon.com/s3/)
- [GCP](https://cloud.google.com/storage/)
- [Azure](https://azure.microsoft.com/en-gb/services/storage/)
- [IPFS](https://ipfs.io/)

## Installing
Hoard should be go-gettable with:

```shell
# Install Hoard-Daemon:
go get github.com/monax/hoard/cmd/hoard

# Install Hoard-Control:
go get github.com/monax/hoard/cmd/hoarctl
```

## Usage
Hoard runs as a daemon providing a GRPC service to other clients including the command line client `hoarctl`. The purpose of the daemon is to read local secrets (such as PGP or other keys) and to configure itself to use a particular storage backend. You can run the daemon with:

```shell
# Run the daemon
hoard

# or with logging
hoard --logging
```

You can initialise a Hoard config by running one of:

```shell
# Initialise Hoard with memory backend
hoard config --init memory

# Initialise Hoard with filesystem backend
hoard config --init filesystem

# Initialise Hoard with AWS (S3) backend
hoard config --init aws

# Initialise Hoard with Azure backend
hoard config --init azure

# Initialise Hoard with GCP backend
hoard config --init gcp

# Initialise Hoard with IPFS backend
hoard config --init ipfs
```

These will provide base configurations you can configure to meet your needs. The config is located by default in `$HOME/.config/hoard.conf` but you can specify a file with `hoard -c /path/to/config`. The XDG base directory specification is used to search for config.

You can interact with Hoard using the go client `hoarctl`:

```shell
# Store an object:
ref=$(echo bar | hoarctl put)

# Retrieve 'bar' from its (deterministic) reference
echo $ref | hoarctl get

# Or get information about the object without decrypting
echo $ref | hoarctl stat

# This one-liner exercises the entire API:
echo foo | hoarctl put | hoarctl get | hoarctl putseal | hoarctl unsealget | hoarctl encrypt | hoarctl insert | hoarctl stat | hoarctl cat | hoarctl decrypt -k tbudgBSg+bHWHiHnlteNzN8TUvI80ygS9IULh4rklEw= | hoarctl ref | hoarctl seal | hoarctl reseal | hoarctl unseal | hoarctl get
```

You can chop off segments of the final command to see the output of each intermediate command. It is contrived so that the outputs can be used as inputs for the next pipeline step. `hoarctl` either returns JSON references or raw bytes depending on the command. You may find the excellent [jq](https://stedolan.github.io/jq/) useful for working with single-line JSON files on the command line.

## Config
Using the filesystem storage backend as an example (generated with `hoard init -o- fs`) you can configure Hoard with a file like:

```toml
# The listen address, also supported is "unix:///tmp/hoard.socket" for a unix domain socket
ListenAddress = "tcp://localhost:53431"

[Storage]
  StorageType = "filesystem"
  # One of: base64, base32, or hex (base 16)
  AddressEncoding = "base64"
  RootDirectory = "/home/user/.local/share/hoard"

[Logging]
  LoggingType = "logfmt"
  # Removing "trace" from this array will reduce log output
  Channels = ["info", "trace"]
```

The default directory is `$HOME/.config/hoard.toml` or you can pass the file with `hoard -c`.

## Specification
See [hoard.proto](protobuf/hoard.proto) for the protobuf3 definition of the API. Hoard uses [GRPC](https://grpc.io/) for its API for which there is a wide range of client libraries available. You should be able to set up a client in any GRPC supported language with relative ease. Also see `hoarctl <CMD> -h` for full help on each sub-command.

For more information on the design of Hoard please checkout our [documentation](docs/encryption.md).

## Building
To build Hoard you will need to have the following installed:
- Go (Version >= 1.11) (with $GOPATH/bin in $PATH)
- GNU make
- [Protocol Buffers 3](https://github.com/google/protobuf/releases/tag/v3.3.0)

Then, from the project root run:

```shell
# Install protobuf GRPC plugin, glide, and glide dependencies
make protobuf_deps
# Run checks, tests, and build binaries
make build && make install
```

## Javascript Client
A Javascript client library can be found here: [js](https://github.com/monax/hoard/tree/master/js).

Hoard-js is a fairly lightweight wrapper around the Hoard GRPC API. It mainly serves to abstract over the dynamic protobuf library and the static protobuf generation.

### Usage
First we need to have Hoard running. For development purposes this can be accomplished by:

```shell
go get github.com/monax/hoard/cmd/hoard 
# Run Hoard with logging
hoard --logging
```

Hoard will run with an in-memory store by default that will be discarded when it is shutdown, but will expose the same
interface as when using remote storage backends.

To interact with Hoard from Node see [example.js](example.js) for a self-contained example of how to use every method
of the API. To run use:

```shell
# Get dependencies
yarn install
# Run example
node example.js
```
