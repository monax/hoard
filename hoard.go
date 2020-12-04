package hoard

import (
	"crypto/sha256"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/v8/config"
	"github.com/monax/hoard/v8/encryption"
	"github.com/monax/hoard/v8/grant"
	"github.com/monax/hoard/v8/reference"
	"github.com/monax/hoard/v8/stores"
)

type EncryptionService interface {
	// Encrypt data and return it along with reference
	Encrypt(data, salt []byte) (ref *reference.Ref, encryptedData []byte, err error)
	// Encrypt data and return it along with reference
	Decrypt(ref *reference.Ref, encryptedData []byte) (data []byte, err error)
}

type ObjectService interface {
	EncryptionService
	// Get encrypted data from underlying storage at address and decrypt it
	Get(ref *reference.Ref) (data []byte, err error)
	// Encrypt data and put it in underlying storage
	Put(data, salt []byte) (*reference.Ref, error)
	// Delete underlying data obtained by address
	Delete(address []byte) error
	// Get the underlying ContentAddressedStore
	Store() stores.ContentAddressedStore
}

type GrantService interface {
	ObjectService
	// Seal a reference by encrypting it according to a grant spec
	Seal(refs reference.Refs, spec *grant.Spec) (*grant.Grant, error)
	// Unseal a grant by decrypting it and returning the reference
	Unseal(grt *grant.Grant) (reference.Refs, error)
}

// This is our top level API object providing library acting as a deterministic
// encrypted store and a grant issuer. It can be consumed as a Go library or as
// a GRPC service through grpcService which just plumbs this object into the
// hoard.proto interface.
type Hoard struct {
	name    string
	store   stores.ContentAddressedStore
	secrets config.SecretsManager
	logger  log.Logger
}

func NewHoard(store stores.NamedStore, secrets config.SecretsManager, logger log.Logger) *Hoard {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	return &Hoard{
		name: store.Name(),
		store: stores.NewContentAddressedStore(stores.MakeAddresser(sha256.New),
			stores.NewLoggingStore(stores.NewSyncStore(store), logger)),
		secrets: secrets,
		logger:  log.With(logger, "scope", "NewHoard"),
	}
}

func (hrd *Hoard) Name() string {
	return hrd.name
}

func (hrd *Hoard) Seal(refs reference.Refs, spec *grant.Spec) (*grant.Grant, error) {
	return grant.Seal(hrd.secrets, refs, spec)
}

func (hrd *Hoard) Unseal(grt *grant.Grant) (reference.Refs, error) {
	return grant.Unseal(hrd.secrets, grt)
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

// Encrypts data and storage it in underlying store and returns the address
func (hrd *Hoard) Put(data, salt []byte) (*reference.Ref, error) {
	blob, err := encryption.EncryptConvergent(data, salt)
	if err != nil {
		return nil, err
	}
	address, err := hrd.store.Put(blob.EncryptedData)
	if err != nil {
		return nil, err
	}
	return reference.New(address, blob.SecretKey, salt, int64(len(data))), nil
}

func (hrd *Hoard) Delete(address []byte) error {
	return hrd.store.Delete(address)
}

// Encrypt data and get reference
func (hrd *Hoard) Encrypt(data, salt []byte) (*reference.Ref, []byte, error) {
	blob, err := encryption.EncryptConvergent(data, salt)
	if err != nil {
		return nil, nil, err
	}
	address := hrd.store.Address(blob.EncryptedData)
	return reference.New(address, blob.SecretKey, salt, int64(len(data))), blob.EncryptedData, nil
}

// Decrypt data using reference
func (hrd *Hoard) Decrypt(ref *reference.Ref, encryptedData []byte) ([]byte, error) {
	data, err := encryption.DecryptConvergent(encryptedData, ref.Salt, ref.SecretKey)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (hrd *Hoard) Store() stores.ContentAddressedStore {
	return hrd.store
}
