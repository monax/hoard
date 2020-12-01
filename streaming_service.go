package hoard

import (
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/monax/hoard/v8/api"
	"github.com/monax/hoard/v8/grant"
	"github.com/monax/hoard/v8/reference"
	"github.com/monax/hoard/v8/stores"
)

// StreamingService provides the API implementation for Service without relying directly on the
// GRPC generated streaming types
type StreamingService struct {
	grantService      GrantService
	chunkSize         int
	useLinkRefInGrant func(refs []*reference.Ref) bool
}

// Create a streaming service that will re-buffer any plaintext data in blocks of chunkSize. useLinkRefInGrant is a
// predicate used to decide whether PutSeal should store a LINK ref in a grant rather than an array of refs. It is
// passed the refs so it may decide to dynamically use LINK refs based on the number and size of references
func NewStreamingService(grantService GrantService, chunkSize int, useLinkRefInGrant func(refs []*reference.Ref) bool) *StreamingService {
	return &StreamingService{
		grantService:      grantService,
		chunkSize:         chunkSize,
		useLinkRefInGrant: useLinkRefInGrant,
	}
}

func (service *StreamingService) PutSeal(sendAndClose func(*grant.Grant) error, recv func() (*api.PlaintextAndGrantSpec, error)) error {
	first, err := recv()
	if err != nil {
		return err
	}
	spec := first.GetGrantSpec()
	if spec == nil {
		return fmt.Errorf("grant spec expected in first message")
	}

	head := first.GetPlaintext().GetHead()

	var refs reference.Refs
	if head != nil {
		ref, err := headerReference(head, service.grantService.Put)
		if err != nil {
			return err
		}
		refs = append(refs, ref)
	}

	err = CopyChunked(
		func(chunk []byte) error {
			ref, err := service.grantService.Put(chunk, head.GetSalt())
			if err != nil {
				return err
			}
			refs = append(refs, ref)
			return nil
		},
		func() ([]byte, error) {
			if first != nil {
				body := first.GetPlaintext().GetBody()
				first = nil
				return body, nil
			}
			ptgs, err := recv()
			if err != nil {
				return nil, err
			}
			return ptgs.GetPlaintext().GetBody(), nil
		},
		service.chunkSize)

	// Should we store the refs themselves in hoard and seal a LINK ref into the grant
	if service.useLinkRefInGrant(refs) {
		ref, err := service.grantService.Put(refs.Plaintext(nil), head.GetSalt())
		if err != nil {
			return err
		}
		ref.Type = reference.Ref_LINK
		refs = reference.Refs{ref}
	}

	// Now send the grant
	grt, err := service.grantService.Seal(refs, spec)
	if err != nil {
		return err
	}

	return sendAndClose(grt)
}

// UnsealGet decrypts and gets plaintext associated with a grant
func (service *StreamingService) UnsealGet(grt *grant.Grant, send func(*api.Plaintext) error) error {
	refs, err := service.grantService.Unseal(grt)
	if err != nil {
		return err
	}

	for _, ref := range refs {
		data, err := service.grantService.Get(ref)
		if err != nil {
			return err
		}

		if err = decodePlaintext(data, ref.GetType(), service.grantService.Get, send); err != nil {
			return err
		}
	}
	return nil
}

// UnsealDelete gets the references stored in a grant and deletes them
func (service *StreamingService) UnsealDelete(grt *grant.Grant, send func(address *api.Address) error) error {
	refs, err := service.grantService.Unseal(grt)
	if err != nil {
		return err
	}

	for _, ref := range refs {
		err := service.grantService.Store().Delete(ref.Address)
		if err != nil {
			return err
		}
		if err = send(&api.Address{Address: ref.Address}); err != nil {
			return err
		}
	}
	return nil
}

// Put encrypted data in the store
func (service *StreamingService) Put(send func(*reference.Ref) error, recv func() (*api.Plaintext, error)) error {
	pt, err := recv()
	if err != nil {
		return err
	}

	head := pt.GetHead()
	if head != nil {
		ref, err := headerReference(head, service.grantService.Put)
		if err != nil {
			return err
		}
		err = send(ref)
		if err != nil {
			return err
		}
	}
	err = CopyChunked(
		func(chunk []byte) error {
			ref, err := service.grantService.Put(chunk, head.GetSalt())
			if err != nil {
				return err
			}
			return send(ref)
		},
		func() ([]byte, error) {
			// When Header is omitted first message should contain b
			if pt.GetBody() != nil {
				body := pt.Body
				pt = nil
				return body, nil
			}
			plaintext, err := recv()
			if err != nil {
				return nil, err
			}
			return plaintext.Body, nil
		},
		service.chunkSize)

	if err != nil {
		return fmt.Errorf("Put: could not put plaintexts: %w", err)
	}

	return nil
}

// Get decrypted data from the store
func (service *StreamingService) Get(send func(*api.Plaintext) error, recv func() (*reference.Ref, error)) error {
	for {
		ref, err := recv()
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

		if err = decodePlaintext(data, ref.GetType(), service.grantService.Get, send); err != nil {
			return err
		}
	}
}

// Encrypt data and return ciphertext
func (service *StreamingService) Encrypt(send func(*api.ReferenceAndCiphertext) error, recv func() (*api.Plaintext, error)) error {
	first, err := recv()
	if err != nil {
		return err
	}
	head := first.GetHead()

	if head != nil {
		var encryptedData []byte
		ref, err := headerReference(head, func(data, salt []byte) (ref *reference.Ref, err error) {
			ref, encryptedData, err = service.grantService.Encrypt(data, salt)
			return ref, err
		})
		if err != nil {
			return fmt.Errorf("Encrypt: could not encode header: %w", err)
		}

		err = send(&api.ReferenceAndCiphertext{
			Reference: ref,
			Ciphertext: &api.Ciphertext{
				EncryptedData: encryptedData,
			},
		})
		if err != nil {
			return err
		}
	}

	var body []byte
	for {
		if first.GetBody() != nil {
			body = first.GetBody()
			first = nil
		} else {
			plaintext, err := recv()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			body = plaintext.GetBody()
		}

		ref, encryptedData, err := service.grantService.Encrypt(body, head.GetSalt())
		if err != nil {
			return err
		}

		if err = send(&api.ReferenceAndCiphertext{
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
func (service *StreamingService) Decrypt(send func(*api.Plaintext) error, recv func() (*api.ReferenceAndCiphertext, error)) error {
	for {
		refAndCiphertext, err := recv()
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

		if err = decodePlaintext(data, refAndCiphertext.Reference.GetType(), service.grantService.Get, send); err != nil {
			return err
		}
	}
}

// StorageServer

// Push ciphertext directly to store
func (service *StreamingService) Push(send func(*api.Address) error, recv func() (*api.Ciphertext, error)) error {
	for {
		ciphertext, err := recv()
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

		if err = send(&api.Address{Address: addr}); err != nil {
			return err
		}
	}
}

// Pull gets ciphertext directly from the store
func (service *StreamingService) Pull(send func(*api.Ciphertext) error, recv func() (*api.Address, error)) error {
	for {
		addr, err := recv()
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

		if err = send(&api.Ciphertext{EncryptedData: data}); err != nil {
			return err
		}
	}
}

// GrantServer

// Seal puts refs in a shareable grant
func (service *StreamingService) Seal(sendAndClose func(*grant.Grant) error, recv func() (*api.ReferenceAndGrantSpec, error)) error {
	var refs reference.Refs
	rgs, err := recv()
	if err != nil {
		return err
	}

	spec := rgs.GetGrantSpec()
	if spec == nil {
		return fmt.Errorf("grant spec expected in first message")
	}

	for {
		refAndSpec, err := recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if s := refAndSpec.GrantSpec; s != nil {
			return fmt.Errorf("received multiple grant specs but there can be at most one")
		}

		refs = append(refs, refAndSpec.Reference)
	}

	grt, err := service.grantService.Seal(refs, spec)
	if err != nil {
		return err
	}

	return sendAndClose(grt)
}

// Unseal gets the refs stored in a grant
func (service *StreamingService) Unseal(grt *grant.Grant, send func(*reference.Ref) error) error {
	refs, err := service.grantService.Unseal(grt)
	if err != nil {
		return err
	}

	for _, ref := range refs {
		if err = send(ref); err != nil {
			return err
		}
	}

	return nil
}

// Reseal changes how the references in a grant are stored
func (service *StreamingService) Reseal(arg *api.GrantAndGrantSpec) (*grant.Grant, error) {
	refs, err := service.grantService.Unseal(arg.Grant)
	if err != nil {
		return nil, err
	}
	return service.grantService.Seal(refs, arg.GrantSpec)
}

func (service *StreamingService) Stat(address *api.Address) (*stores.StatInfo, error) {
	statInfo, err := service.grantService.Store().Stat(address.Address)
	if err != nil {
		return nil, err
	}
	// provide the address and the canonical location
	statInfo.Address = address.Address
	statInfo.Location = service.grantService.Store().Location(address.Address)
	return statInfo, nil
}

func (service *StreamingService) Delete(address []byte) error {
	return service.grantService.Store().Delete(address)
}

// Converts raw plaintext data to the Plaintext wrapper type
// In the case of a HEADER ref type the plaintext is deserialised using the header type
// In the case of a LINK ref type the supplied get function is used to fetch additional plaintext data which are themselves each decoded
// Otherwise the data is returned as Plaintext.Body
// The decoded plaintext(s) are then streamed as output via the supplied send function
func decodePlaintext(data []byte, refType reference.Ref_RefType, get func(*reference.Ref) ([]byte, error),
	send func(*api.Plaintext) error) error {

	switch refType {
	case reference.Ref_HEADER:
		head := new(api.Header)
		err := proto.Unmarshal(data, head)
		if err != nil {
			return err
		}
		return send(&api.Plaintext{Head: head})

	case reference.Ref_LINK:
		refs := reference.RepeatedFromPlaintext(data)
		for _, ref := range refs {
			data, err := get(ref)
			if err != nil {
				return err
			}
			err = decodePlaintext(data, ref.Type, get, send)
			if err != nil {
				return err
			}
		}
		return nil

	default:
		return send(&api.Plaintext{
			Body: data,
		})
	}
}

func headerReference(head *api.Header, getRef func(data, salt []byte) (*reference.Ref, error)) (*reference.Ref, error) {
	data, err := proto.Marshal(head)
	if err != nil {
		return nil, err
	}

	ref, err := getRef(data, head.GetSalt())
	if err != nil {
		return nil, err
	}
	ref.Type = reference.Ref_HEADER

	return ref, nil
}
