This breaking changes refactors the exported API to make it possible to have a much more minimal import tree. Not all storage backends are imported when depending just on api (containing protobuf generated code) or on hoard/v5 root package which allows, for instance, importing the root package to run an in memory test server without all the storage backend dependencies.

### Changed
- Renamed services package to api
- Move services.NewHoardServer to hoard.NewServer
- Renamed storage package to stores
- Made ipfs and cloud their own subpackages to avoid massive import tree

