package hoard

import (
	"context"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/monax/hoard/v8/api"
	"github.com/monax/hoard/v8/grant"
	"github.com/monax/hoard/v8/reference"
	"github.com/monax/hoard/v8/stores"
)

const defaultRefVersionForHeader = 1

// MaxChunkSize = 1MiB
const MaxChunkSize = 1 << 20

// Service implements the GRPC Hoard service. It should mostly be plumbing to
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

// Get decrypted data from the store
func (service *Service) Get(srv api.Cleartext_GetServer) error {
	for {
		ref, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		data, err := service.grantService.Get(ref)
		if err != nil {
			return err
		}

		if err = SendPlaintext(data, service.chunkSize, srv, ref.GetVersion()); err != nil {
			return err
		}
	}
}

// Put encrypted data in the store
func (service *Service) Put(srv api.Cleartext_PutServer) error {
	head, err := consumeHeadFromPlaintext(srv.Recv())
	if err != nil {
		return err
	}

	if err = service.putHeader(head, func(ref *reference.Ref) error {
		return srv.Send(ref)
	}); err != nil {
		return err
	}

	var accumulator []byte
	for {
		plaintext, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return service.putPlaintext(accumulator, head.GetSalt(), func(ref *reference.Ref) error {
					return srv.Send(ref)
				})
			}

			return err
		}

		accumulator, err = consumeBodyFromPlaintext(plaintext, accumulator, service.chunkSize, func(chunk []byte) error {
			return service.putPlaintext(chunk, head.GetSalt(), func(ref *reference.Ref) error {
				return srv.Send(ref)
			})
		})
		if err != nil {
			return err
		}
	}
}

// Encrypt data and return ciphertext
func (service *Service) Encrypt(srv api.Encryption_EncryptServer) error {
	head, err := consumeHeadFromPlaintext(srv.Recv())
	if err != nil {
		return err
	}

	if err = service.encHeader(head, func(rct *api.ReferenceAndCiphertext) error {
		return srv.Send(rct)
	}); err != nil {
		return err
	}

	for {
		plaintext, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		ref, encryptedData, err := service.grantService.Encrypt(plaintext.GetBody(), head.GetSalt())
		if err != nil {
			return err
		}

		if err = srv.Send(&api.ReferenceAndCiphertext{
			Reference: ref,
			Ciphertext: &api.Ciphertext{
				EncryptedData: encryptedData,
			},
		}); err != nil {
			return err
		}
	}
}

// Decrypt ciphertext and return plaintext
func (service *Service) Decrypt(srv api.Encryption_DecryptServer) error {
	for {
		refAndCiphertext, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		data, err := service.grantService.Decrypt(refAndCiphertext.Reference, refAndCiphertext.Ciphertext.EncryptedData)
		if err != nil {
			return err
		}

		if err = SendPlaintext(data, service.chunkSize, srv, refAndCiphertext.Reference.GetVersion()); err != nil {
			return err
		}
	}
}

// StorageServer

// Push ciphertext directly to store
func (service *Service) Push(srv api.Storage_PushServer) error {
	for {
		ciphertext, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		addr, err := service.grantService.Store().Put(ciphertext.EncryptedData)
		if err != nil {
			return err
		}

		if err = srv.Send(&api.Address{Address: addr}); err != nil {
			return err
		}
	}
}

// Pull gets ciphertext directly from the store
func (service *Service) Pull(srv api.Storage_PullServer) error {
	for {
		addr, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		data, err := service.grantService.Store().Get(addr.Address)
		if err != nil {
			return err
		}

		if err = srv.Send(&api.Ciphertext{EncryptedData: data}); err != nil {
			return err
		}
	}
}

// Delete removes the data located at the address
func (service *Service) Delete(ctx context.Context, address *api.Address) (*api.Address, error) {
	return address, service.grantService.Store().Delete(address.Address)
}

// Stat checks the data stored at the given address
func (service *Service) Stat(ctx context.Context, address *api.Address) (*stores.StatInfo, error) {
	statInfo, err := service.grantService.Store().Stat(address.Address)
	if err != nil {
		return nil, err
	}
	// provide the address and the canonical location
	statInfo.Address = address.Address
	statInfo.Location = service.grantService.Store().Location(address.Address)
	return statInfo, nil
}

// GrantServer

// Seal puts refs in a shareable grant
func (service *Service) Seal(srv api.Grant_SealServer) error {
	refs, spec, err := receiveReferencesAndGrantSpec(srv)
	if err != nil {
		return err
	}

	grt, err := service.grantService.Seal(refs, spec)
	if err != nil {
		return err
	}

	return srv.SendAndClose(grt)
}

// Unseal gets the refs stored in a grant
func (service *Service) Unseal(grt *grant.Grant, srv api.Grant_UnsealServer) error {
	refs, err := service.grantService.Unseal(grt)
	if err != nil {
		return err
	}

	for _, ref := range refs {
		if err = srv.Send(ref); err != nil {
			return err
		}
	}

	return nil
}

// Reseal changes how the references in a grant are stored
func (service *Service) Reseal(ctx context.Context, arg *api.GrantAndGrantSpec) (*grant.Grant, error) {
	refs, err := service.grantService.Unseal(arg.Grant)
	if err != nil {
		return nil, err
	}
	return service.grantService.Seal(refs, arg.GrantSpec)
}

// PutSeal encrypts and seals plaintext
func (service *Service) PutSeal(srv api.Grant_PutSealServer) error {
	spec, head, err := consumeHeadFromPlaintextAndGrantSpec(srv.Recv())
	if err != nil {
		return err
	}

	var refs reference.Refs
	if err = service.putHeader(head, func(ref *reference.Ref) error {
		refs = append(refs, ref)
		return nil
	}); err != nil {
		return err
	}

	var accumulator []byte
	for {
		ptgs, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				err := service.putPlaintext(accumulator, head.GetSalt(), func(ref *reference.Ref) error {
					refs = append(refs, ref)
					return nil
				})
				if err != nil {
					return err
				}

				grt, err := service.grantService.Seal(refs, spec)
				if err != nil {
					return err
				}
				return srv.SendAndClose(grt)
			}

			return err
		}

		accumulator, err = consumeBodyFromPlaintextAndGrantSpec(ptgs, accumulator, service.chunkSize, func(chunk []byte) error {
			return service.putPlaintext(chunk, head.GetSalt(), func(ref *reference.Ref) error {
				refs = append(refs, ref)
				return nil
			})
		})
		if err != nil {
			return err
		}
	}
}

// UnsealGet decrypts and gets plaintext associated with a grant
func (service *Service) UnsealGet(grt *grant.Grant, srv api.Grant_UnsealGetServer) error {
	refs, err := service.grantService.Unseal(grt)
	if err != nil {
		return err
	}

	for _, ref := range refs {
		data, err := service.grantService.Get(ref)
		if err != nil {
			return err
		}

		if err = SendPlaintext(data, service.chunkSize, srv, ref.GetVersion()); err != nil {
			return err
		}
	}
	return nil
}

// UnsealDelete gets the references stored in a grant and deletes them
func (service *Service) UnsealDelete(grt *grant.Grant, srv api.Grant_UnsealDeleteServer) error {
	refs, err := service.grantService.Unseal(grt)
	if err != nil {
		return err
	}

	for _, ref := range refs {
		addr, err := service.Delete(srv.Context(), &api.Address{Address: ref.Address})
		if err != nil {
			return err
		}
		if err = srv.Send(addr); err != nil {
			return err
		}
	}
	return nil
}

func (service *Service) encHeader(head *api.Header, cb func(*api.ReferenceAndCiphertext) error) error {
	data, err := proto.Marshal(head)
	if err != nil {
		return err
	}

	ref, encryptedData, err := service.grantService.Encrypt(data, head.GetSalt())
	if err != nil {
		return err
	}
	ref.Version = defaultRefVersionForHeader

	return cb(&api.ReferenceAndCiphertext{
		Reference: ref,
		Ciphertext: &api.Ciphertext{
			EncryptedData: encryptedData,
		},
	})
}

func (service *Service) putHeader(head *api.Header, cb func(*reference.Ref) error) error {
	data, err := proto.Marshal(head)
	if err != nil {
		return err
	}

	ref, err := service.grantService.Put(data, head.GetSalt())
	if err != nil {
		return err
	}
	ref.Version = defaultRefVersionForHeader

	return cb(ref)
}

func (service *Service) putPlaintext(data, salt []byte, cb func(*reference.Ref) error) error {
	if data == nil {
		return nil
	}

	ref, err := service.grantService.Put(data, salt)
	if err != nil {
		return err
	}

	return cb(ref)
}
