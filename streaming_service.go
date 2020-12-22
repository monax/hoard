package hoard

import (
	"fmt"
	"io"

	"github.com/monax/hoard/v8/encryption"

	"github.com/monax/hoard/v8/protodet"
	"github.com/monax/hoard/v8/versions"

	"github.com/monax/hoard/v8/api"
	"github.com/monax/hoard/v8/grant"
	"github.com/monax/hoard/v8/reference"
	"github.com/monax/hoard/v8/stores"
)

// StreamingService provides the API implementation for Service without relying directly on the
// GRPC generated streaming types
type StreamingService struct {
	grantService GrantService
	chunkSize    int64
}

// Create a streaming service that will re-buffer any plaintext data in blocks of chunkSize. linker is a
// predicate used to decide whether PutSeal should store a LINK ref in a grant rather than an array of refs. It is
// passed the refs so it may decide to dynamically use LINK refs based on the number and size of references
func NewStreamingService(grantService GrantService, chunkSize int64) *StreamingService {
	return &StreamingService{
		grantService: grantService,
		chunkSize:    chunkSize,
	}
}

func (service *StreamingService) PutSeal(sendAndClose func(*grant.Grant) error, recv func() (*api.PlaintextAndGrantSpec, error)) error {
	first, err := recv()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("no messages sent before EOF")
		}
		return err
	}

	spec := first.GetGrantSpec()
	if spec == nil {
		return fmt.Errorf("grant spec expected in first message")
	}

	head := first.GetPlaintext().GetHead()

	var refs []*reference.Ref

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

	// TODO: it would be useful to be able to send Header.Data (i.e. metadata) as a trailer, this could then be
	//   normalised to the front of the refs array in the object pointed to by a LINK ref. This would allow things like
	//   file size, detected mime type, etc to be provided at the end of (or during) a stream. If a salt/chunksize is used
	//   we do  need that at the beginning of the stream though. The decision to explicitly only support first-message headers
	//   was to avoid unintended behaviour from multiple headers appearing in the stream. Another option is to provide another call
	//   that take a grant without a header and adds a header by creating a copy of the link ref with a header added.

	// Convert base refs into link ref(s) (usually a single unique link ref to allow for safe deletion of links)
	refs, err = link(refs, head.GetSalt(), spec.LinkNonce, service.grantService.Put)
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

		err = decode(data, ref.GetType(), service.grantService.Get, send, versions.LatestGrantVersion)
		if err != nil {
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
		if err == io.EOF {
			return fmt.Errorf("no messages sent before EOF")
		}
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

		err = decode(data, ref.GetType(), service.grantService.Get, send, versions.LatestGrantVersion)
		if err != nil {
			return err
		}
	}
}

// Encrypt data and return ciphertext
func (service *StreamingService) Encrypt(send func(*api.ReferenceAndCiphertext) error, recv func() (*api.Plaintext, error)) error {
	first, err := recv()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("no messages sent before EOF")
		}
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

		err = decode(data, refAndCiphertext.Reference.GetType(), service.grantService.Get, send, versions.LatestGrantVersion)
		if err != nil {
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
	var refs []*reference.Ref
	first, err := recv()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("no messages sent before EOF")
		}
		return err
	}

	spec := first.GetGrantSpec()
	if spec == nil {
		return fmt.Errorf("grant spec expected in first message")
	}

	for {
		if first.GetReference() != nil {
			ref := first.GetReference()
			first = nil
			refs = append(refs, ref)
		}
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
	// TODO: could provide a way to add Header metadata after the fact by re-linking some refs with header appended
	//  something like: refs ,err = link(append([]*reference.Ref{headerRef}, delink(refs)...), salt, linkNonce, service.grantService.Put)
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

// Put wrapped with dummy 'encrypt' signature to help with reuse
func (service *StreamingService) put(data, salt []byte) (*reference.Ref, []byte, error) {
	ref, err := service.grantService.Put(data, salt)
	return ref, nil, err
}

// Abstracts the handling of incoming plaintexts that is common between Encrypt, Put, and PutSeal
func encrypt(first *api.Plaintext,
	encrypt func(data []byte, salt []byte) (ref *reference.Ref, encryptedData []byte, err error),
	send func(ref *reference.Ref, encryptedData []byte) error,
	recv func() (*api.Plaintext, error),
	chunkSize int64) error {

	// Expect header to always be in first message if provided
	head := first.GetHead()
	if head != nil {
		// Use chunkSize if supplied
		if head.GetChunkSize() > 0 {
			chunkSize = head.ChunkSize
		}
		data, err := protodet.Marshal(head)
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
			// Consume any body that may have been in the first message
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

// Converts raw plaintext data to the wrapper type
// In the case of a HEADER ref type the plaintext is deserialised using the header type
// In the case of a LINK ref type the supplied get function is used to fetch additional plaintext data which are themselves each decoded
// Otherwise the data is returned as MustPlaintextFromRefs.Body
// The decoded plaintext(s) are then streamed as output via the supplied send function
func decode(data []byte, refType reference.Ref_RefType,
	get func(*reference.Ref) ([]byte, error),
	send func(*api.Plaintext) error, version int32) error {

	switch refType {
	case reference.Ref_HEADER:
		head := new(api.Header)
		err := protodet.Unmarshal(data, head)
		if err != nil {
			return err
		}
		return send(&api.Plaintext{Head: head})

	case reference.Ref_LINK:
		refs, err := reference.RefsFromPlaintext(data, version)
		if err != nil {
			return err
		}
		for _, ref := range refs {
			data, err := get(ref)
			if err != nil {
				return err
			}
			err = decode(data, ref.Type, get, send, version)
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

func link(refs []*reference.Ref, salt []byte, linkNonce []byte,
	put func(data, salt []byte) (*reference.Ref, error)) ([]*reference.Ref, error) {
	var err error
	// By default link refs use a unique nonce to allow them to be deletable unless the grant specifies otherwise
	if len(linkNonce) == 0 {
		linkNonce, err = encryption.NewNonce(encryption.NonceSize)
		if err != nil {
			return nil, fmt.Errorf("could not create nonce for LINK ref: %w", err)
		}
	}
	// Store refs as a plaintext document
	plaintext, err := reference.PlaintextFromRefs(refs, linkNonce)
	if err != nil {
		return nil, err
	}
	ref, err := put(plaintext, salt)
	if err != nil {
		return nil, err
	}
	// Mark this ref as a LINK so it will be followed during dereferencing
	ref.Type = reference.Ref_LINK
	return []*reference.Ref{ref}, nil
}
