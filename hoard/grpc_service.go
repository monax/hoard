package hoard

import (
	"code.monax.io/platform/hoard/hoard/reference"
	"code.monax.io/platform/hoard/hoard/storage"
	"golang.org/x/net/context"
)

// Here we implement the GRPC Hoard service. It should mostly be plumbing to
// a DeterministicEncryptedStore (for which hoard.hoard is the canonical example)
// and also to Grants.
type grpcService struct {
	des DeterministicEncryptedStore
}

//service Hoard {
//rpc Get (Reference) returns (Plaintext);
//rpc Put (Plaintext) returns (Reference);
//rpc Ref (Plaintext) returns (Reference);
//rpc Stat (Address) returns (StatInfo);
//}
var _ HoardServer = (*grpcService)(nil)

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

func (service *grpcService) Ref(ctx context.Context,
	plaintext *Plaintext) (*Reference, error) {

	ref, err := service.des.Ref(plaintext.Data, plaintext.Salt)
	if err != nil {
		return nil, err
	}

	return protobufRef(ref), nil
}

func (service *grpcService) Stat(ctx context.Context,
	address *Address) (*StatInfo, error) {

	statInfo, err := service.des.Store().Stat(address.Address)
	if err != nil {
		return nil, err
	}

	return protobufStatInfo(statInfo), nil
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
