package hoard

import (
	"context"

	"github.com/monax/hoard/v7/api"
	"github.com/monax/hoard/v7/grant"
	"github.com/monax/hoard/v7/reference"
	"github.com/monax/hoard/v7/stores"
)

// 1MiB
const MaxChunkSize = 1 << 20

// Here we implement the GRPC Hoard service. It should mostly be plumbing to
// a DeterministicEncryptedStore (for which hoard.hoard is the canonical example)
// and also to Grants.
type Service struct {
	grantService GrantService
	chunkSize    int
}

func NewService(grantService GrantService, chunkSize int) *Service {
	if chunkSize == 0 {
		chunkSize = MaxChunkSize
	}
	return &Service{
		grantService: grantService,
		chunkSize:    chunkSize,
	}
}

func (service *Service) Get(ref *reference.Ref, srv api.Cleartext_GetServer) error {
	data, err := service.grantService.Get(ref)
	if err != nil {
		return err
	}

	return SendPlaintext(srv, data, ref.Salt, service.chunkSize)
}

func (service *Service) Put(srv api.Cleartext_PutServer) error {
	plaintext, err := ReceivePlaintext(srv)
	if err != nil {
		return err
	}
	ref, err := service.grantService.Put(plaintext.Data, plaintext.Salt)
	if err != nil {
		return err
	}
	return srv.SendAndClose(ref)
}

func (service *Service) Encrypt(srv api.Encryption_EncryptServer) error {
	plaintext, err := ReceivePlaintext(srv)
	if err != nil {
		return err
	}

	ref, encryptedData, err := service.grantService.Encrypt(plaintext.Data, plaintext.Salt)
	if err != nil {
		return err
	}

	return srv.SendAndClose(&api.ReferenceAndCiphertext{
		Reference: ref,
		Ciphertext: &api.Ciphertext{
			EncryptedData: encryptedData,
		},
	})
}

func (service *Service) Decrypt(refAndCiphertext *api.ReferenceAndCiphertext, srv api.Encryption_DecryptServer) error {
	data, err := service.grantService.Decrypt(refAndCiphertext.Reference, refAndCiphertext.Ciphertext.EncryptedData)
	if err != nil {
		return err
	}

	return SendPlaintext(srv, data, refAndCiphertext.Reference.GetSalt(), service.chunkSize)
}

// StorageServer
func (service *Service) Push(srv api.Storage_PushServer) error {
	data, err := ReceiveCiphertext(srv)
	if err != nil {
		return err
	}

	address, err := service.grantService.Store().Put(data)
	if err != nil {
		return err
	}

	return srv.SendAndClose(&api.Address{
		Address: address,
	})
}

func (service *Service) Pull(address *api.Address, srv api.Storage_PullServer) error {
	// Get from the underlying store
	data, err := service.grantService.Store().Get(address.Address)
	if err != nil {
		return err
	}

	return SendCiphertext(srv, data, service.chunkSize)

}

func (service *Service) Delete(ctx context.Context, address *api.Address) (*api.Address, error) {
	return address, service.grantService.Store().Delete(address.Address)
}

func (service *Service) Stat(ctx context.Context, address *api.Address) (*stores.StatInfo, error) {
	statInfo, err := service.grantService.Store().Stat(address.Address)
	if err != nil {
		return nil, err
	}
	// For the master API we provide the address and the canonical
	// location in a StatInfo message
	statInfo.Address = address.Address
	statInfo.Location = service.grantService.Store().Location(address.Address)
	return statInfo, nil
}

// GrantServer

func (service *Service) Seal(ctx context.Context, arg *api.ReferenceAndGrantSpec) (*grant.Grant, error) {
	return service.grantService.Seal(arg.Reference, arg.GrantSpec)
}

func (service *Service) Unseal(ctx context.Context, grt *grant.Grant) (*reference.Ref, error) {
	return service.grantService.Unseal(grt)
}

func (service *Service) Reseal(ctx context.Context, arg *api.GrantAndGrantSpec) (*grant.Grant, error) {
	ref, err := service.grantService.Unseal(arg.Grant)
	if err != nil {
		return nil, err
	}
	return service.grantService.Seal(ref, arg.GrantSpec)
}

func (service *Service) PutSeal(srv api.Grant_PutSealServer) error {
	pgs, err := ReceivePlaintextAndGrantSpec(srv)
	if err != nil {
		return err
	}

	ref, err := service.grantService.Put(pgs.Plaintext.Data, pgs.Plaintext.Salt)
	if err != nil {
		return err
	}

	grt, err := service.grantService.Seal(ref, pgs.GrantSpec)
	if err != nil {
		return err
	}

	return srv.SendAndClose(grt)
}

func (service *Service) UnsealGet(grt *grant.Grant, srv api.Grant_UnsealGetServer) error {
	ref, err := service.grantService.Unseal(grt)
	if err != nil {
		return err
	}

	data, err := service.grantService.Get(ref)
	if err != nil {
		return err
	}

	return SendPlaintext(srv, data, ref.Salt, service.chunkSize)
}

func (service *Service) UnsealDelete(ctx context.Context, grt *grant.Grant) (*api.Address, error) {
	ref, err := service.grantService.Unseal(grt)
	if err != nil {
		return nil, err
	}
	return service.Delete(ctx, &api.Address{Address: ref.Address})
}

func (service *Service) Download(grt *grant.Grant, srv api.Document_DownloadServer) error {
	doc, salt, err := GetDocument(service.grantService, grt)
	if err != nil {
		return err
	}

	return SendDocument(srv, doc, salt, service.chunkSize)
}

func (service *Service) Upload(srv api.Document_UploadServer) error {
	pgsm, err := ReceiveDocumentAndGrantSpec(srv)
	if err != nil {
		return err
	}

	grt, err := PutDocument(service.grantService, pgsm)
	if err != nil {
		return err
	}

	return srv.SendAndClose(grt)
}
