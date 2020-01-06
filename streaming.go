package hoard

import (
	"fmt"
	"io"

	"github.com/monax/hoard/v7/api"
	"github.com/monax/hoard/v7/grant"
	"github.com/monax/hoard/v7/meta"
	"github.com/monax/hoard/v7/reference"
)

type PlaintextReceiver interface {
	Recv() (*api.Plaintext, error)
}

func ReceivePlaintext(srv PlaintextReceiver) (*api.Plaintext, error) {
	accum := new(api.Plaintext)
	for {
		plaintext, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return accum, nil
			}

			return nil, err
		}

		err = consumePlaintext(accum, plaintext)
		if err != nil {
			return nil, err
		}
	}
}

func ReceiveAndPutPlaintext(srv PlaintextReceiver, obj ObjectService, ref *reference.Ref, salt []byte) error {
	for {
		accum, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		} else if len(accum.Salt) > 0 {
			return fmt.Errorf("received multiple salts but there can be at most one")
		}

		ref.Next, err = obj.Put(accum.Data, salt)
		if err != nil {
			return err
		}
		ref = ref.Next
	}
}

type PlaintextSender interface {
	Send(*api.Plaintext) error
}

func SendPlaintext(srv PlaintextSender, data, salt []byte, chunkSize int) error {
	err := srv.Send(&api.Plaintext{Salt: salt})
	if err != nil {
		return err
	}

	return sendChunks(data, chunkSize, func(chunk []byte) error {
		return srv.Send(&api.Plaintext{Data: chunk})
	})
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

func SendCiphertext(srv CiphertextSender, data []byte, chunkSize int) error {
	return sendChunks(data, chunkSize, func(chunk []byte) error {
		return srv.Send(&api.Ciphertext{EncryptedData: chunk})
	})
}

type PlaintextAndGrantSpecReceiver interface {
	Recv() (*api.PlaintextAndGrantSpec, error)
}

// Receive chunks of plaintext and spec and aggregate into complete objects
func ReceiveAndPutPlaintextAndGrantSpec(srv PlaintextAndGrantSpecReceiver, obj ObjectService, ref *reference.Ref, salt []byte) error {
	for {
		accum, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		} else if accum.GrantSpec != nil {
			return fmt.Errorf("received multiple grant specs but there can be at most one")
		} else if accum.Plaintext != nil && len(accum.Plaintext.Salt) > 0 {
			return fmt.Errorf("received multiple salts but there can be at most one")
		} else if accum.Plaintext == nil {
			accum.Plaintext = new(api.Plaintext)
		}

		ref.Next, err = obj.Put(accum.Plaintext.Data, salt)
		if err != nil {
			return err
		}
		ref = ref.Next
	}
}

type PlaintextAndGrantSpecSender interface {
	Send(*api.PlaintextAndGrantSpec) error
}

// Send some plaintext and spec to a service in chunks
func SendPlaintextAndGrantSpec(srv PlaintextAndGrantSpecSender, pgs *api.PlaintextAndGrantSpec, chunkSize int) error {
	err := srv.Send(&api.PlaintextAndGrantSpec{GrantSpec: pgs.GrantSpec, Plaintext: &api.Plaintext{Salt: pgs.Plaintext.Salt}})
	if err != nil {
		return err
	}

	return sendChunks(pgs.Plaintext.Data, chunkSize, func(chunk []byte) error {
		return srv.Send(&api.PlaintextAndGrantSpec{Plaintext: &api.Plaintext{Data: chunk}})
	})
}

type DocumentSender interface {
	Send(*api.PlaintextAndMeta) error
}

func SendDocument(srv DocumentSender, doc *meta.Document, salt []byte, chunkSize int) error {
	err := srv.Send(&api.PlaintextAndMeta{Meta: doc.Meta})
	if err != nil {
		return err
	}

	err = srv.Send(&api.PlaintextAndMeta{Plaintext: &api.Plaintext{Salt: salt}})
	if err != nil {
		return err
	}

	return sendChunks(doc.Data, chunkSize, func(chunk []byte) error {
		return srv.Send(&api.PlaintextAndMeta{Plaintext: &api.Plaintext{Data: chunk}})
	})
}

type DocumentAndGrantSender interface {
	Send(*api.PlaintextAndGrantSpecAndMeta) error
}

func SendDocumentAndGrantSpec(srv DocumentAndGrantSender, doc *meta.Document, salt []byte, spec *grant.Spec, chunkSize int) error {
	err := srv.Send(&api.PlaintextAndGrantSpecAndMeta{Meta: doc.Meta})
	if err != nil {
		return err
	}

	err = srv.Send(&api.PlaintextAndGrantSpecAndMeta{
		PlaintextAndGrantSpec: &api.PlaintextAndGrantSpec{GrantSpec: spec},
	})
	if err != nil {
		return err
	}

	err = srv.Send(&api.PlaintextAndGrantSpecAndMeta{
		PlaintextAndGrantSpec: &api.PlaintextAndGrantSpec{Plaintext: &api.Plaintext{Salt: salt}},
	})
	if err != nil {
		return err
	}

	return sendChunks(doc.Data, chunkSize, func(chunk []byte) error {
		return srv.Send(&api.PlaintextAndGrantSpecAndMeta{
			PlaintextAndGrantSpec: &api.PlaintextAndGrantSpec{Plaintext: &api.Plaintext{Data: chunk}},
		})
	})
}

type DocumentAndGrantReceiver interface {
	Recv() (*api.PlaintextAndGrantSpecAndMeta, error)
}

func ReceiveDocumentAndGrantSpec(srv DocumentAndGrantReceiver) (*api.PlaintextAndGrantSpecAndMeta, error) {
	accum := &api.PlaintextAndGrantSpecAndMeta{
		PlaintextAndGrantSpec: &api.PlaintextAndGrantSpec{Plaintext: &api.Plaintext{}},
	}
	for {
		d, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return accum, nil
			}
			return nil, err
		}

		// NOTE: for singular values we adopt the convention of accepting the first one
		if d.Meta != nil {
			if accum.Meta != nil {
				return nil, fmt.Errorf("received multiple document meta but there can be at most one")
			}
			accum.Meta = d.Meta
		}

		err = consumePlaintextAndGrantSpec(accum.PlaintextAndGrantSpec, d.PlaintextAndGrantSpec)
		if err != nil {
			return nil, err
		}
	}
}

type DocumentReceiver interface {
	Recv() (*api.PlaintextAndMeta, error)
}

func ReceiveDocument(srv DocumentReceiver) (*meta.Document, error) {
	accum := &api.PlaintextAndMeta{
		Plaintext: &api.Plaintext{},
	}

	for {
		d, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return &meta.Document{Meta: accum.Meta, Data: accum.Plaintext.Data}, nil
			}

			return nil, err
		}

		if d.Meta != nil {
			if accum.Meta != nil {
				return nil, fmt.Errorf("received multiple document meta but there can be at most one")
			}
			accum.Meta = d.Meta
		}

		err = consumePlaintext(accum.Plaintext, d.Plaintext)
		if err != nil {
			return nil, err
		}
	}
}

func consumePlaintextAndGrantSpec(accum, chunk *api.PlaintextAndGrantSpec) error {
	if chunk == nil {
		return nil
	}
	if chunk.GrantSpec != nil {
		if accum.GrantSpec != nil {
			return fmt.Errorf("received multiple grant specs but there can be at most one")
		}
		if len(accum.Plaintext.Data) > 0 {
			return fmt.Errorf("received grant spec after data but spec must come before all data chunks")
		}
		accum.GrantSpec = chunk.GrantSpec
	}
	return consumePlaintext(accum.Plaintext, chunk.Plaintext)
}

func consumePlaintext(accum, chunk *api.Plaintext) error {
	if chunk == nil {
		return nil
	}
	if len(chunk.Salt) > 0 {
		if len(accum.Salt) > 0 {
			return fmt.Errorf("received multiple salts but there can be at most one")
		}
		if len(accum.Data) > 0 {
			return fmt.Errorf("received salt after data but salt must come before all data chunks")
		}
		accum.Salt = chunk.Salt
	}
	accum.Data = append(accum.Data, chunk.Data...)
	return nil
}

func sendChunks(data []byte, chunkSize int, sender func(chunk []byte) error) error {
	var err error
	for i := 0; i < len(data); i += chunkSize {
		if i+chunkSize > len(data) {
			err = sender(data[i:])
		} else {
			err = sender(data[i : i+chunkSize])
		}
		if err != nil {
			return err
		}
	}
	return nil
}
