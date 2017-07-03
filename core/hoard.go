package core

import (
	"crypto/sha256"
	"hash"

	"github.com/go-kit/kit/log"

	"encoding/base64"

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

type DeterministicEncryptor interface {
	// Encrypt data and return it along with reference
	Encrypt(data, salt []byte) (ref *reference.Ref, encryptedData []byte, err error)
	// Encrypt data and return it along with reference
	Decrypt(ref *reference.Ref, encryptedData []byte) (data []byte, err error)
}

type DeterministicEncryptedStore interface {
	DeterministicEncryptor
	// Get encrypted data from underlying storage at address and decrypt it using
	// secretKey
	Get(ref *reference.Ref) (data []byte, err error)
	// Encrypt data and put it in underlying storage
	Put(data, salt []byte) (*reference.Ref, error)
	// Get the underlying ContentAddressedStore
	Store() storage.ContentAddressedStore
}

var _ DeterministicEncryptedStore = (*hoard)(nil)
var _ DeterministicEncryptor = (*hoard)(nil)

func NewHoard(store storage.Store, logger log.Logger) DeterministicEncryptedStore {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &hoard{
		store: storage.NewContentAddressedStore(makeAddresser(sha256.New()),
			storage.NewLoggingStore(store, logger)),
		logger: log.With(logger, "scope", "NewHoard"),
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
			"No data stored at address %s",
			base64.StdEncoding.EncodeToString(ref.Address))
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

// Encrypt data and get reference
func (hrd *hoard) Encrypt(data, salt []byte) (*reference.Ref, []byte, error) {
	blob, err := encryption.Encrypt(data, salt)
	if err != nil {
		return nil, nil, err
	}
	address := hrd.store.Address(blob.EncryptedData())
	return reference.New(address, blob.SecretKey(), salt), blob.EncryptedData(), nil

}

// Decrypt data using reference
func (hrd *hoard) Decrypt(ref *reference.Ref, encryptedData []byte) ([]byte, error) {
	data, err := encryption.Decrypt(ref.SecretKey, encryptedData, ref.Salt)
	if err != nil {
		return nil, err
	}
	return data, nil
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
