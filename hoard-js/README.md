## Hoard-js

Hoard-js is a fairly lightweight typescript wrapper around the Hoard GRPC API. It mainly serves to abstract over the dynamic
protobuf library and the static protobuf generation.

### Usage

First we need to run a Hoard server. For development purposes this can be accomplished by:

```shell
go get github.com/monax/hoard/cmd/hoard 
# Run in memory with secret-key and logging
hoard config memory -s testing-id:secret_pass | hoard --logging -c-
```

Hoard will run with an in-memory store by default that will be discarded when it is shutdown, but will expose the same
interface as when using remote storage backends.

To interact with Hoard see [examples.ts](src/examples.js) for a self-contained example of how to use every method
of the API. To run use:

```shell
# Get dependencies
npm install
# Run examples
npm test
```