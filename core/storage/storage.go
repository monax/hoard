package storage

import (
	"encoding/base64"

	"hash"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorAddressNotFound(address []byte) error {
	return status.Errorf(codes.NotFound, "No data stored at address %s",
		base64.StdEncoding.EncodeToString(address))
}

type Locator interface {
	// Provides a canonical external location for some data, typically a URI
	Location(address []byte) string
}

type StatInfo struct {
	Exists bool
	Size   uint64
}

type ReadStore interface {
	// Get data stored at address
	Get(address []byte) (data []byte, err error)
	// Get stats on file including existence
	Stat(address []byte) (*StatInfo, error)
}

type WriteStore interface {
	// Put data at address
	Put(address []byte, data []byte) ([]byte, error)
}

type Store interface {
	ReadStore
	WriteStore
	Locator
}

type NamedStore interface {
	// Human readable name describing the Store
	Name() string
	Store
}

type ContentAddressedStore interface {
	ReadStore
	Locator
	// Put the data at its address
	Put(data []byte) (address []byte, err error)
	// Get the address of some data without putting it at that address
	Address(data []byte) (address []byte)
}

type contentAddressedStore struct {
	// The addresser that derives an address from some data deterministically.
	// Generally we would expect addresser to be a good (enough) hash function for
	// the space of expected binary strings passed as data. It may also encode some
	// conventions around the location of the binary blob in addition to its role
	// as a hash function
	addresser func(data []byte) (address []byte)
	// The underlying store to store against
	store Store
}

func NewContentAddressedStore(addresser func([]byte) []byte, store Store) ContentAddressedStore {
	return &contentAddressedStore{
		addresser: addresser,
		store:     store,
	}
}

func (cas *contentAddressedStore) Address(data []byte) []byte {
	return cas.addresser(data)
}

func (cas *contentAddressedStore) Put(data []byte) ([]byte, error) {
	address := cas.addresser(data)
	address, err := cas.store.Put(address, data)
	return address, err
}

func (cas *contentAddressedStore) Get(address []byte) ([]byte, error) {
	return cas.store.Get(address)
}

func (cas *contentAddressedStore) Stat(address []byte) (*StatInfo, error) {
	return cas.store.Stat(address)
}

func (cas *contentAddressedStore) Location(address []byte) string {
	return cas.store.Location(address)
}

// Close in hasher
func MakeAddresser(hashProvider func() hash.Hash) func(data []byte) []byte {
	return func(data []byte) []byte {
		hasher := hashProvider()
		hasher.Write(data)
		return hasher.Sum(nil)
	}
}
