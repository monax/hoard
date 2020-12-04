package hoard

import (
	"fmt"
	"io"

	"github.com/monax/hoard/v8/encryption"

	"github.com/golang/protobuf/proto"
	"github.com/monax/hoard/v8/api"
	"github.com/monax/hoard/v8/grant"
	"github.com/monax/hoard/v8/reference"
	"github.com/monax/hoard/v8/stores"
)

type Linker func(refs reference.Refs, head *api.Header, put func(data, salt []byte) (*reference.Ref, error)) (reference.Refs, error)

// StreamingService provides the API implementation for Service without relying directly on the
// GRPC generated streaming types
type StreamingService struct {
	grantService GrantService
	chunkSize    int64
	linker       Linker
}

// Create a streaming service that will re-buffer any plaintext data in blocks of chunkSize. linker is a
// predicate used to decide whether PutSeal should store a LINK ref in a grant rather than an array of refs. It is
// passed the refs so it may decide to dynamically use LINK refs based on the number and size of references
func NewStreamingService(grantService GrantService, chunkSize int64, linker Linker) *StreamingService {
	return &StreamingService{
		grantService: grantService,
		chunkSize:    chunkSize,
		linker:       linker,
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

	err = encrypt(first.GetPlaintext(), service.put,
		func(ref *reference.Ref, encryptedData []byte) error {
			refs = append(refs, ref)
			return nil
		},
		func() (*api.Plaintext, error) {
			ptgs, err := recv()
			if err != nil {
				return nil, err
			}
			return ptgs.GetPlaintext(), nil
		}, service.chunkSize)

	// Convert base refs into link ref(s) (usually a single unique link ref to allow for safe deletion of links)
	refs, err = service.linker(refs, head, service.grantService.Put)
	if err != nil {
		return fmt.Errorf("could not link refs: %w", err)
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
	first, err := recv()
	if err != nil {
		return err
	}

	err = encrypt(first, service.put, func(ref *reference.Ref, _ []byte) error { return send(ref) },
		recv, service.chunkSize)

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

	err = encrypt(first, service.grantService.Encrypt, func(ref *reference.Ref, encryptedData []byte) error {
		return send(&api.ReferenceAndCiphertext{
			Reference: ref,
			Ciphertext: &api.Ciphertext{
				EncryptedData: encryptedData,
			},
		})
	}, recv, service.chunkSize)

	if err != nil {
		return fmt.Errorf("Could not encrypt data: %w", err)
	}
	return nil
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

// Put wrapped with dummy 'encrypt' signature to help with reuse
func (service *StreamingService) put(data, salt []byte) (*reference.Ref, []byte, error) {
	ref, err := service.grantService.Put(data, salt)
	return ref, nil, err
}

// Abstracts common Plaintext input flow handling
func encrypt(first *api.Plaintext,
	encrypt func(data []byte, salt []byte) (ref *reference.Ref, encryptedData []byte, err error),
	send func(ref *reference.Ref, encryptedData []byte) error,
	recv func() (*api.Plaintext, error), chunkSize int64) error {

	head := first.GetHead()
	if head != nil {
		// Use chunkSize if supplied
		if head.GetChunkSize() > 0 {
			chunkSize = head.ChunkSize
		}
		data, err := proto.Marshal(head)
		if err != nil {
			return err
		}

		ref, encryptedData, err := encrypt(data, head.GetSalt())
		if err != nil {
			return err
		}
		ref.Type = reference.Ref_HEADER

		err = send(ref, encryptedData)
		if err != nil {
			return err
		}
	}
	// Truncate to max chunkSize
	if chunkSize > MaxChunkSize {
		chunkSize = MaxChunkSize
	}
	return CopyChunked(
		func(chunk []byte) error {
			ref, encryptedData, err := encrypt(chunk, head.GetSalt())
			if err != nil {
				return err
			}
			return send(ref, encryptedData)
		},
		func() ([]byte, error) {
			// If and only if there is no head
			if first.GetBody() != nil {
				body := first.GetBody()
				first = nil
				return body, nil
			}
			plaintext, err := recv()
			if err != nil {
				return nil, err
			}
			return plaintext.Body, nil
		},
		chunkSize)
}

func defaultLinker(refs reference.Refs, head *api.Header, put func(data, salt []byte) (*reference.Ref, error)) (reference.Refs, error) {
	// Always link refs and always use a unique nonce to allow them to be deletable
	nonce, err := encryption.NewNonce(encryption.NonceSize)
	if err != nil {
		return nil, fmt.Errorf("could not create nonce for LINK ref: %w", err)
	}
	// Store refs as a plaintext document
	ref, err := put(refs.Plaintext(nonce), head.GetSalt())
	if err != nil {
		return nil, err
	}
	// Mark this ref as a LINK so it will be followed during dereferencing
	ref.Type = reference.Ref_LINK
	return reference.Refs{ref}, nil
}
