# Hoard

Hoard is a stateless, deterministically encrypted, content-addressed object store.

It is stateless in the sense of relying on storage backends for the actual persistence of objects. Currently supported backends are:

- Memory
- Filesystem
- [S3](https://aws.amazon.com/s3/)
- [GCS](https://cloud.google.com/storage/)
- [IPFS](https://ipfs.io/)


Planned storage backends are:

- BigchainDB (and IPDB)
- Tendermint

It encrypts deterministically (convergently) because it encrypts an object using the object's hash (SHA256) as the secret key (which can than be shared as a 'grant').

It is content-addressed because encrypted objects are stored at an address determined by the encrypted object's hash (SHA256 again).

![hoarding marmot](docs/images/hoard.jpg)

## Installing

Hoard should be go-gettable with:

```shell
# Install the Hoar-Daemon hoard:
go get github.com/monax/hoard/cmd/hoard

# Install the Hoar-Control hoarctl:
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

With no config file by default `hoard` will run a memory storage backend from which all objects will be lost when `hoard` is terminated. You can initialise a Hoard config by running one of:

```shell
# Initialise Hoard with memory backend
hoard config --init mem

# Initialise Hoard with filesystem backend
hoard config --init fs

# Initialise Hoard with S3 backend
hoard config --init s3

# Initialise Hoard with GCS backend
hoard config --init gcs

# Initialise Hoard with IPFS backend
hoard config --init ipfs
```

These will provide base configurations you can configure to meet your needs. The config is located by default in `$HOME/.config/hoard.toml` but you can specify a file with `hoard -c /path/to/config`. The XDG base directory specification is used to search for config.

You can interact with Hoard using `hoarctl`:

```shell
# Store an object:
ref=$(echo bar | hoarctl put)

# Retrieve 'bar' from its (deterministic) reference
echo $ref | hoarctl get

# Or get information about the object without decrypting
echo $ref | hoarctl stat

# This one-liner exercises the entire API:
echo foo | hoarctl put | hoarctl get | hoarctl put | hoarctl stat | hoarctl cat | hoarctl insert | hoarctl cat | hoarctl decrypt -k tbudgBSg+bHWHiHnlteNzN8TUvI80ygS9IULh4rklEw= | hoarctl encrypt
```

You can chop off segments of the final command to see the output of each intermediate command. It is contrived so that the outputs can be used as inputs for the next pipeline step. `hoarctl` either returns JSON references or raw bytes depending on the command. You may find the excellent [jq](https://stedolan.github.io/jq/) useful for working with single-line JSON files on the commandline.

## Config
Using the filesystem storage backend as an example (generated with `hoard init -o- fs`) you can configure Hoard with a file like:

```toml
# The listen address, also supported is "unix:///tmp/hoard.socket" for a unix domain socket
ListenAddress = "tcp://localhost:53431"

[Storage]
  StorageType = "filesystem"
  # One of: base64, base32, or hex (base 16)
  AddressEncoding = "base64"
  RootDirectory = "/home/silas/.local/share/hoard"

[Logging]
  LoggingType = "logfmt"
  # Removing "trace" from this array will reduce log output
  Channels = ["info", "trace"]
```

The default directory is `$HOME/.config/hoard.toml` or you can pass the file with `hoard -c`.

## Encryption scheme

Hoard implements an encryption scheme based off the SHA256 cryptographic hash function and the symmetric block cipher AES256-GCM (Galois Counter Mode is an authenticated mode of AES). It is an example of envelope encryption where an object is encrypted with a specific one-time key and where that secret key can itself be shared by encrypting it (asymmetrically or otherwise) and publishing it to a recipient. It is motivated by and possesses the following features:

- Usable on a publicly-accessible storage backend (strongly encrypted by AES256-GCM)
- Addressable by the hash of the ciphertext (the bytes of encrypted object) so can be ciphertext can de-duplicated and object existence can be queried without the secret key
- Ciphertext is recoverable from plaintext (if you have a copy of the plaintext you can check whether the ciphertext is stored and recover the secret key)
- Permits sharing of access grants in public via asymmetric encryption of the secret key and address (can be used to implement a decentralised equivalent of access control lists)
- Optional salt allows one to simulate random encryption keys or to add entropy to short objects

The encryption scheme relies on well-known cryptographic primitives. Given an object (just a sequence of bytes) and a salt the encryption proceeds as follows:

1. Compute the SHA256 of the object.
2. Prepend the salt (just arbitrary bytes) to the object.
3. Encrypt the salted object with AES256-GCM with an empty IV (this is okay since our keys are one-time since they are based on a semantically secure hash of the object) with authenticated additional data describing the presence of the salt (authenticated additional data is a feature of GCM).
4. Compute the SHA256 of the encrypted object.
5. Return the output of (1) as the `secretKey`, the output of (3) as the `encryptedData`, and the output of (4) as the `address`.

Given an encrypted object, secret key, and salt the decryption proceeds as follows:

1. Decrypt the encrypted object with AES256-GCM with an empty IV with `secretKey` and additional authenticated data associated with `salt`.
2. Trim the prefix `salt` from the output of (1).
3. Return the object as `data` from the output of (2).

### Security

By design this scheme is trivially vulnerable to known-plaintext attacks (if you know the plaintext you can find the key).

It is also vulnerable to a ciphertext-only attack where a rainbow table of ciphertexts may be computed where the space of variation in a set of plaintexts is small. This occurs when the plaintexts are short or when the plaintexts are mostly identical template with a small amount of variation (an 8-digit account number for example). This is only a problem if revealing that an object has been stored at all leaks sensitive information.

The scheme should not be vulnerable to general chosen-ciphertext attacks (CCA) through the use of Galois Counter Mode with AES that is CCA-secure.

If you want known-plaintext or ciphertext-only security you can provide a salt. The salt can be of any length and can be used a number of ways. If you are encrypting short plaintexts you can agree a sufficiently long salt to be pre-shared amongst a set of parties that adds entropy to the encryption so that a rainbow table of possible values cannot be feasibly built (semantic security of SHA256 means the salt will induce an unpredictable variation of ciphertexts). In this case amongst the parties in possession of the salt the advantages of deterministic encryption are preserved. Alternatively a random hash can be used for the salt, this will induce a different secret key and address each time you encrypt the same bytes. It is effectively the same as using a random key.

The encryption is furthermore vulnerable to the same timing and length attacks that to which AES is susceptible, but for most purposes these attacks are not usually considered an issue.

### Maturity

Hoard is still pre-release and pre-version, there may be breaking changes to the API at any time. Before version tag v1.0.0 changes minor version number changes may break the hoard.proto or hoarctl APIs (e.g. 0.3.4 -> 0.4.0) but patch number changes should leave it intact (e.g. 0.3.4 -> 0.3.5).

The cryptographic libraries used are standard Go libraries (and Go's NACL implementation) so should be of reasonable quality and are widely deployed. The encryption scheme is straight-forward and has an isolated implementation. However there may be bugs in the implementation.

## Specification

See [hoard.proto](protobuf/hoard.proto) for the protobuf3 definition of the API. Also see `hoarctl <CMD> -h` for full help on each sub-command.

## Clients

Hoard uses [GRPC](https://grpc.io/) for its API for which there is a wide range of client libraries available. You should be able to set up a client in any GRPC supported language with relative ease.

### Javascript

A Javascript client library can be found here: [hoard-js](https://github.com/monax/hoard-js).

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
