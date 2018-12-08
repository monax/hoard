package hoard

import (
	"crypto/sha256"

	"github.com/monax/hoard/grant"

	"github.com/go-kit/kit/log"

	"github.com/monax/hoard/encryption"
	"github.com/monax/hoard/reference"
	"github.com/monax/hoard/storage"
)

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

type GrantService interface {
	// Seal a reference by encrypting it according to a grant spec
	Seal(ref *reference.Ref, spec *grant.Spec) (*grant.Grant, error)
	// Unseal a grant by decrypting it and returning the reference
	Unseal(grt *grant.Grant) (*reference.Ref, error)
}

// This is our top level API object providing library acting as a deterministic
// encrypted store and a grant issuer. It can be consumed as a Go library or as
// a GRPC service through grpcService which just plumbs this object into the
// hoard.proto interface.
type Hoard struct {
	name           string
	store          storage.ContentAddressedStore
	secretProvider grant.SecretProvider
	logger         log.Logger
}

func NewHoard(store storage.NamedStore, secretProvider grant.SecretProvider, logger log.Logger) *Hoard {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	return &Hoard{
		name: store.Name(),
		store: storage.NewContentAddressedStore(storage.MakeAddresser(sha256.New),
			storage.NewLoggingStore(storage.NewSyncStore(store), logger)),
		secretProvider: secretProvider,
		logger:         log.With(logger, "scope", "NewHoard"),
	}
}

func (hrd *Hoard) Name() string {
	return hrd.name
}

func (hrd *Hoard) Seal(ref *reference.Ref, spec *grant.Spec) (*grant.Grant, error) {
	return grant.Seal(hrd.secretProvider, ref, spec)
}

func (hrd *Hoard) Unseal(grt *grant.Grant) (*reference.Ref, error) {
	return grant.Unseal(hrd.secretProvider, grt)
}

// Gets encrypted blob
func (hrd *Hoard) Get(ref *reference.Ref) ([]byte, error) {
	encryptedData, err := hrd.store.Get(ref.Address)
	if err != nil {
		return nil, err
	}

	data, err := encryption.DecryptConvergent(encryptedData, ref.Salt, ref.SecretKey)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Encrypts data and stores it in underlying store and returns the address
func (hrd *Hoard) Put(data, salt []byte) (*reference.Ref, error) {
	blob, err := encryption.EncryptConvergent(data, salt)
	if err != nil {
		return nil, err
	}
	address, err := hrd.store.Put(blob.EncryptedData)
	if err != nil {
		return nil, err
	}
	return reference.New(address, blob.SecretKey, salt), nil
}

// Encrypt data and get reference
func (hrd *Hoard) Encrypt(data, salt []byte) (*reference.Ref, []byte, error) {
	blob, err := encryption.EncryptConvergent(data, salt)
	if err != nil {
		return nil, nil, err
	}
	address := hrd.store.Address(blob.EncryptedData)
	return reference.New(address, blob.SecretKey, salt), blob.EncryptedData, nil

}

// Decrypt data using reference
func (hrd *Hoard) Decrypt(ref *reference.Ref, encryptedData []byte) ([]byte, error) {
	data, err := encryption.DecryptConvergent(encryptedData, ref.Salt, ref.SecretKey)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (hrd *Hoard) Store() storage.ContentAddressedStore {
	return hrd.store
}
