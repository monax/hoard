package core

import (
	"crypto/sha256"
	"hash"

	"github.com/go-kit/kit/log"

	"code.monax.io/platform/hoard/core/encryption"
	"code.monax.io/platform/hoard/core/reference"
	"code.monax.io/platform/hoard/core/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// This is our top level API object providing library acting as a deterministic
// encrypted store and a grant issuer. It can be consumed as a Go library or as
// a GRPC service through grpcService which just plumbs this object into the
// hoard.proto interface.
type hoard struct {
	store  storage.ContentAddressedStore
	logger log.Logger
}

type DeterministicEncryptedStore interface {
	// Get encrypted data from underlying storage at address and decrypt it using
	// secretKey
	Get(ref *reference.Ref) (data []byte, err error)
	// Encrypt data and put it in underlying storage
	Put(data, salt []byte) (*reference.Ref, error)
	// Like Put, but just compute the reference
	Ref(data, salt []byte) (*reference.Ref, error)
	// Get the underlying ContentAddressedStore
	Store() storage.ContentAddressedStore
}

var _ DeterministicEncryptedStore = (*hoard)(nil)

func NewHoard(store storage.Store, logger log.Logger) DeterministicEncryptedStore {
	return &hoard{
		store:  storage.NewContentAddressedStore(makeAddresser(sha256.New()), store),
		logger: logger,
	}
}

// Gets encrypted blob
func (hrd *hoard) Get(ref *reference.Ref) ([]byte, error) {
	encryptedData, err := hrd.store.Get(ref.Address)
	if err != nil {
		return nil, err
	}

	// Some stores return nil/empty data for 'not found'. If the address is
	// non-zero // then by definition the encrypted data should be, so we
	// infer there is nothing stored at address
	if len(encryptedData) == 0 && len(ref.Address) != 0 {
		return nil, status.Errorf(codes.NotFound,
			"No data stored at address 0x%X", ref.Address)
	}

	data, err := encryption.Decrypt(ref.SecretKey, encryptedData, ref.Salt)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Encrypts data and stores it in underlying store and returns the address
func (hrd *hoard) Put(data, salt []byte) (*reference.Ref, error) {
	blob, err := encryption.Encrypt(data, salt)
	if err != nil {
		return nil, err
	}
	address, err := hrd.store.Put(blob.EncryptedData())
	if err != nil {
		return nil, err
	}
	return reference.New(address, blob.SecretKey(), salt), nil
}

func (hrd *hoard) Ref(data, salt []byte) (*reference.Ref, error) {
	blob, err := encryption.Encrypt(data, salt)
	if err != nil {
		return nil, err
	}
	address := hrd.store.Address(blob.EncryptedData())
	return reference.New(address, blob.SecretKey(), salt), nil

}

func (hrd *hoard) Store() storage.ContentAddressedStore {
	return hrd.store
}

// Close in hasher
func makeAddresser(hasher hash.Hash) func(data []byte) []byte {
	return func(data []byte) []byte {
		hasher.Reset()
		hasher.Write(data)
		return hasher.Sum(nil)
	}
}
