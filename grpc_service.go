package hoard

import (
	"context"

	"github.com/monax/hoard/grant"

	"github.com/monax/hoard/reference"
	"github.com/monax/hoard/storage"
)

// Here we implement the GRPC Hoard service. It should mostly be plumbing to
// a DeterministicEncryptedStore (for which hoard.hoard is the canonical example)
// and also to Grants.
type grpcService struct {
	des DeterministicEncryptedStore
	gs  GrantService
}

func NewHoardServer(des DeterministicEncryptedStore, gs GrantService) *grpcService {
	return &grpcService{
		des: des,
		gs:  gs,
	}
}

func (service *grpcService) Get(ctx context.Context, ref *reference.Ref) (*Plaintext, error) {
	data, err := service.des.Get(ref)
	if err != nil {
		return nil, err
	}

	return &Plaintext{
		Data: data,
		Salt: ref.Salt,
	}, nil
}

func (service *grpcService) Put(ctx context.Context, plaintext *Plaintext) (*reference.Ref, error) {
	return service.des.Put(plaintext.Data, plaintext.Salt)
}

func (service *grpcService) Encrypt(ctx context.Context, plaintext *Plaintext) (*ReferenceAndCiphertext, error) {
	ref, encryptedData, err := service.des.Encrypt(plaintext.Data, plaintext.Salt)
	if err != nil {
		return nil, err
	}

	return &ReferenceAndCiphertext{
		Reference: ref,
		Ciphertext: &Ciphertext{
			EncryptedData: encryptedData,
		},
	}, nil
}

func (service *grpcService) Decrypt(ctx context.Context, refAndCiphertext *ReferenceAndCiphertext) (*Plaintext, error) {
	data, err := service.des.Decrypt(refAndCiphertext.Reference, refAndCiphertext.Ciphertext.EncryptedData)
	if err != nil {
		return nil, err
	}
	return &Plaintext{
		Data: data,
		Salt: refAndCiphertext.Reference.Salt,
	}, nil
}

// StorageServer
func (service *grpcService) Push(ctx context.Context, ciphertext *Ciphertext) (*Address, error) {
	address, err := service.des.Store().Put(ciphertext.EncryptedData)
	if err != nil {
		return nil, err
	}
	return &Address{
		Address: address,
	}, nil
}

func (service *grpcService) Pull(ctx context.Context, address *Address) (*Ciphertext, error) {

	// Get from the underlying store
	encryptedData, err := service.des.Store().Get(address.Address)
	if err != nil {
		return nil, err
	}

	return &Ciphertext{
		EncryptedData: encryptedData,
	}, nil
}

func (service *grpcService) Stat(ctx context.Context, address *Address) (*storage.StatInfo, error) {
	statInfo, err := service.des.Store().Stat(address.Address)
	if err != nil {
		return nil, err
	}
	// For the master API we provide the address and the canonical
	// location in a StatInfo message
	statInfo.Address = address.Address
	statInfo.Location = service.des.Store().Location(address.Address)
	return statInfo, nil
}

// GrantServer

func (service *grpcService) Seal(ctx context.Context, arg *ReferenceAndGrantSpec) (*grant.Grant, error) {
	return service.gs.Seal(arg.Reference, arg.GrantSpec)
}

func (service *grpcService) Unseal(ctx context.Context, grt *grant.Grant) (*reference.Ref, error) {
	return service.gs.Unseal(grt)
}

func (service *grpcService) Reseal(ctx context.Context, arg *GrantAndGrantSpec) (*grant.Grant, error) {
	ref, err := service.gs.Unseal(arg.Grant)
	if err != nil {
		return nil, err
	}
	return service.gs.Seal(ref, arg.GrantSpec)
}

func (service *grpcService) PutSealed(ctx context.Context, arg *PlaintextAndGrantSpec) (*grant.Grant, error) {
	ref, err := service.des.Put(arg.Plaintext.Data, arg.Plaintext.Salt)
	if err != nil {
		return nil, err
	}
	return service.gs.Seal(ref, arg.GrantSpec)
}

func (service *grpcService) GetUnsealed(ctx context.Context, grt *grant.Grant) (*Plaintext, error) {
	ref, err := service.gs.Unseal(grt)
	if err != nil {
		return nil, err
	}
	return service.Get(ctx, ref)
}
