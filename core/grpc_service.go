package core

import (
	"github.com/monax/hoard/core/reference"
	"github.com/monax/hoard/core/storage"
	"golang.org/x/net/context"
)

// Here we implement the GRPC Hoard service. It should mostly be plumbing to
// a DeterministicEncryptedStore (for which hoard.hoard is the canonical example)
// and also to Grants.
type grpcService struct {
	des DeterministicEncryptedStore
}

type HoardServer interface {
	CleartextServer
	EncryptionServer
	StorageServer
}

func NewHoardServer(des DeterministicEncryptedStore) HoardServer {
	return &grpcService{
		des: des,
	}
}

func (service *grpcService) Get(ctx context.Context,
	ref *Reference) (*Plaintext, error) {

	data, err := service.des.Get(hoardRef(ref))
	if err != nil {
		return nil, err
	}

	return &Plaintext{
		Data: data,
		Salt: ref.Salt,
	}, nil
}

func (service *grpcService) Put(ctx context.Context,
	plaintext *Plaintext) (*Reference, error) {

	ref, err := service.des.Put(plaintext.Data, plaintext.Salt)
	if err != nil {
		return nil, err
	}

	return protobufRef(ref), nil
}

func (service *grpcService) Encrypt(ctx context.Context,
	plaintext *Plaintext) (*ReferenceAndCiphertext, error) {

	ref, encryptedData, err := service.des.Encrypt(plaintext.Data, plaintext.Salt)
	if err != nil {
		return nil, err
	}

	return &ReferenceAndCiphertext{
		Reference: protobufRef(ref),
		Ciphertext: &Ciphertext{
			EncryptedData: encryptedData,
		},
	}, nil
}

func (service *grpcService) Decrypt(ctx context.Context,
	refAndCiphertext *ReferenceAndCiphertext) (*Plaintext, error) {
	data, err := service.des.Decrypt(hoardRef(refAndCiphertext.Reference),
		refAndCiphertext.Ciphertext.EncryptedData)
	if err != nil {
		return nil, err
	}
	return &Plaintext{
		Data: data,
		Salt: refAndCiphertext.Reference.Salt,
	}, nil
}

// StorageServer
func (service *grpcService) Push(ctx context.Context,
	ciphertext *Ciphertext) (*Address, error) {
	address, err := service.des.Store().Put(ciphertext.EncryptedData)
	if err != nil {
		return nil, err
	}
	return &Address{
		Address: address,
	}, nil
}

func (service *grpcService) Pull(ctx context.Context,
	address *Address) (*Ciphertext, error) {

	// Get from the underlying store
	encryptedData, err := service.des.Store().Get(address.Address)
	if err != nil {
		return nil, err
	}

	return &Ciphertext{
		EncryptedData: encryptedData,
	}, nil
}

func (service *grpcService) Stat(ctx context.Context,
	address *Address) (*StatInfo, error) {

	statInfo, err := service.des.Store().Stat(address.Address)
	if err != nil {
		return nil, err
	}

	pbStatInfo := protobufStatInfo(statInfo)
	// For the master API we provide the address and the canonical
	// location in a StatInfo message
	pbStatInfo.Address = address.Address
	pbStatInfo.Location = service.des.Store().Location(address.Address)
	return pbStatInfo, nil
}

// From bitter experience it is better to decouple your serialisation types
// from your object in-memory object model because they change for different
// reasons So we bite the bullet and map between protobuf and hoard objects.

func hoardRef(ref *Reference) *reference.Ref {
	return reference.New(ref.Address, ref.SecretKey, ref.Salt)
}

func protobufRef(ref *reference.Ref) *Reference {
	return &Reference{
		Address:   ref.Address,
		SecretKey: ref.SecretKey,
		Salt:      ref.Salt,
	}
}

func protobufStatInfo(statInfo *storage.StatInfo) *StatInfo {
	return &StatInfo{
		Exists: statInfo.Exists,
		Size:   statInfo.Size,
	}
}
