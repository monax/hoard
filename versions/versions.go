package versions

// We use the Grant version as a kind of global version for the core protobuf types.
// The version is a simple counter and does not 'encode' a breaking change
// (the intention is by using the version number we support all non deprecated previous versions, mearning no change is 'breaking')
// Version history:
// 0: deprecated and removed
// 1: deprecated and removed
// 2: encrypted references array for streaming, non-derived keys, reference with version
// 3: reference Version -> Type, introduce LINK references, store plaintext data Size in reference
const LatestGrantVersion = 3

const RefVersionIncorrectlyUsedToDenoteHeader = 1
