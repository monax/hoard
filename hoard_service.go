package hoard

import (
	"context"
	"io"

	"github.com/monax/hoard/v6/api"
	"github.com/monax/hoard/v6/grant"
	"github.com/monax/hoard/v6/reference"
	"github.com/monax/hoard/v6/stores"
)

// Here we implement the GRPC Hoard service. It should mostly be plumbing to
// a DeterministicEncryptedStore (for which hoard.hoard is the canonical example)
// and also to Grants.
type hoardService struct {
	gs GrantService
	cs int
}

func NewHoardServer(gs GrantService, chunkSize int) *hoardService {
	return &hoardService{
		gs: gs,
		cs: chunkSize,
	}
}

func (service *hoardService) Get(ref *reference.Ref, srv api.Cleartext_GetServer) error {
	data, err := service.gs.Get(ref)
	if err != nil {
		return err
	}

	return SendPlaintext(srv, data, ref.Salt, service.cs)
}

func (service *hoardService) Put(srv api.Cleartext_PutServer) error {
	data, salt, err := ReceivePlaintext(srv)
	if err != nil {
		return err
	}
	ref, err := service.gs.Put(data, salt)
	if err != nil {
		return err
	}
	return srv.SendAndClose(ref)
}

func (service *hoardService) Encrypt(srv api.Encryption_EncryptServer) error {
	data, salt, err := ReceivePlaintext(srv)
	if err != nil {
		return err
	}

	ref, encryptedData, err := service.gs.Encrypt(data, salt)
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

func (service *hoardService) Decrypt(refAndCiphertext *api.ReferenceAndCiphertext, srv api.Encryption_DecryptServer) error {
	data, err := service.gs.Decrypt(refAndCiphertext.Reference, refAndCiphertext.Ciphertext.EncryptedData)
	if err != nil {
		return err
	}

	return SendPlaintext(srv, data, refAndCiphertext.Reference.GetSalt(), service.cs)
}

// StorageServer
func (service *hoardService) Push(srv api.Storage_PushServer) error {
	data, err := ReceiveCiphertext(srv)
	if err != nil {
		return err
	}

	address, err := service.gs.Store().Put(data)
	if err != nil {
		return err
	}

	return srv.SendAndClose(&api.Address{
		Address: address,
	})
}

func (service *hoardService) Pull(address *api.Address, srv api.Storage_PullServer) error {
	// Get from the underlying store
	data, err := service.gs.Store().Get(address.Address)
	if err != nil {
		return err
	}

	return SendCiphertext(srv, data, service.cs)

}

func (service *hoardService) Delete(ctx context.Context, addr *api.Address) (*api.Address, error) {
	return addr, service.gs.Store().Delete(addr.Address)
}

func (service *hoardService) Stat(ctx context.Context, address *api.Address) (*stores.StatInfo, error) {
	statInfo, err := service.gs.Store().Stat(address.Address)
	if err != nil {
		return nil, err
	}
	// For the master API we provide the address and the canonical
	// location in a StatInfo message
	statInfo.Address = address.Address
	statInfo.Location = service.gs.Store().Location(address.Address)
	return statInfo, nil
}

// GrantServer

func (service *hoardService) Seal(ctx context.Context, arg *api.ReferenceAndGrantSpec) (*grant.Grant, error) {
	return service.gs.Seal(arg.Reference, arg.GrantSpec)
}

func (service *hoardService) Unseal(ctx context.Context, grt *grant.Grant) (*reference.Ref, error) {
	return service.gs.Unseal(grt)
}

func (service *hoardService) Reseal(ctx context.Context, arg *api.GrantAndGrantSpec) (*grant.Grant, error) {
	ref, err := service.gs.Unseal(arg.Grant)
	if err != nil {
		return nil, err
	}
	return service.gs.Seal(ref, arg.GrantSpec)
}

func (service *hoardService) PutSeal(srv api.Grant_PutSealServer) error {
	data, salt, spec, err := ReceivePlaintextAndGrantSpec(srv)
	if err != nil {
		return err
	}

	ref, err := service.gs.Put(data, salt)
	if err != nil {
		return err
	}

	grant, err := service.gs.Seal(ref, spec)
	if err != nil {
		return err
	}

	return srv.SendAndClose(grant)
}

func (service *hoardService) UnsealGet(grt *grant.Grant, srv api.Grant_UnsealGetServer) error {
	ref, err := service.gs.Unseal(grt)
	if err != nil {
		return err
	}

	data, err := service.gs.Get(ref)
	if err != nil {
		return err
	}

	return SendPlaintext(srv, data, ref.Salt, service.cs)
}

func (service *hoardService) UnsealDelete(ctx context.Context, grt *grant.Grant) (*api.Address, error) {
	ref, err := service.gs.Unseal(grt)
	if err != nil {
		return nil, err
	}
	return service.Delete(ctx, &api.Address{Address: ref.Address})
}

type PlaintextReceiver interface {
	Recv() (*api.Plaintext, error)
}

func ReceivePlaintext(srv PlaintextReceiver) ([]byte, []byte, error) {
	var data, salt []byte
	for {
		c, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return data, salt, nil
			}

			return nil, nil, err
		}

		switch x := c.GetInput().(type) {
		case *api.Plaintext_Salt:
			salt = x.Salt
		case *api.Plaintext_Data:
			data = append(data, x.Data...)
		}
	}
}

type PlaintextSender interface {
	Send(*api.Plaintext) error
}

func SendPlaintext(srv PlaintextSender, data, salt []byte, cs int) error {
	out := new(api.Plaintext)
	out.Input = &api.Plaintext_Salt{Salt: salt}
	err := srv.Send(out)
	if err != nil {
		return err
	}

	for i := 0; i < len(data); i += cs {
		if i+cs > len(data) {
			out.Input = &api.Plaintext_Data{Data: data[i:len(data)]}
		} else {
			out.Input = &api.Plaintext_Data{Data: data[i : i+cs]}
		}
		if err := srv.Send(out); err != nil {
			return err
		}
	}

	return nil
}

type CiphertextReceiver interface {
	Recv() (*api.Ciphertext, error)
}

func ReceiveCiphertext(srv CiphertextReceiver) ([]byte, error) {
	var data []byte
	for {
		c, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return data, nil
			}

			return nil, err
		}

		data = append(data, c.EncryptedData...)
	}
}

type CiphertextSender interface {
	Send(*api.Ciphertext) error
}

func SendCiphertext(srv CiphertextSender, data []byte, cs int) error {
	out := new(api.Ciphertext)
	for i := 0; i < len(data); i += cs {
		if i+cs > len(data) {
			out.EncryptedData = data[i:len(data)]
		} else {
			out.EncryptedData = data[i : i+cs]
		}
		if err := srv.Send(out); err != nil {
			return err
		}
	}

	return nil
}

type PlaintextAndGrantSpecReceiver interface {
	Recv() (*api.PlaintextAndGrantSpec, error)
}

func ReceivePlaintextAndGrantSpec(srv PlaintextAndGrantSpecReceiver) ([]byte, []byte, *grant.Spec, error) {
	spec := new(grant.Spec)
	var data, salt []byte
	for {
		g, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return data, salt, spec, nil
			}

			return nil, nil, nil, err
		}

		switch x := g.GetInput().(type) {
		case *api.PlaintextAndGrantSpec_Plaintext:
			switch y := x.Plaintext.GetInput().(type) {
			case *api.Plaintext_Salt:
				salt = y.Salt
			case *api.Plaintext_Data:
				data = append(data, y.Data...)
			}
		case *api.PlaintextAndGrantSpec_GrantSpec:
			spec = x.GrantSpec
		}
	}
}

type PlaintextAndGrantSpecSender interface {
	Send(*api.PlaintextAndGrantSpec) error
}

func SendPlaintextAndGrantSpec(srv PlaintextAndGrantSpecSender, spec *grant.Spec, data, salt []byte, cs int) error {
	out := new(api.PlaintextAndGrantSpec)
	out.Input = &api.PlaintextAndGrantSpec_GrantSpec{GrantSpec: spec}
	if err := srv.Send(out); err != nil {
		return err
	}

	if len(salt) > 0 {
		out.Input = &api.PlaintextAndGrantSpec_Plaintext{
			Plaintext: &api.Plaintext{
				Input: &api.Plaintext_Salt{Salt: salt},
			},
		}
		if err := srv.Send(out); err != nil {
			return err
		}
	}

	for i := 0; i < len(data); i += cs {
		if i+cs > len(data) {
			out.Input = &api.PlaintextAndGrantSpec_Plaintext{
				Plaintext: &api.Plaintext{
					Input: &api.Plaintext_Data{Data: data[i:len(data)]},
				},
			}
		} else {
			out.Input = &api.PlaintextAndGrantSpec_Plaintext{
				Plaintext: &api.Plaintext{
					Input: &api.Plaintext_Data{Data: data[i : i+cs]},
				},
			}
		}
		if err := srv.Send(out); err != nil {
			return err
		}
	}

	return nil
}
