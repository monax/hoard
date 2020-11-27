package hoard

import (
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/monax/hoard/v8/api"
	"github.com/monax/hoard/v8/grant"
	"github.com/monax/hoard/v8/reference"
)

type PlaintextSender interface {
	Send(*api.Plaintext) error
}

// SendPlaintext gets the plaintext for a given reference and sends it to the client
func SendPlaintext(data []byte, chunkSize int, srv PlaintextSender, version int32) error {
	if version == defaultRefVersionForHeader {
		head := new(api.Header)
		err := proto.Unmarshal(data, head)
		if err != nil {
			return err
		}
		return srv.Send(&api.Plaintext{Head: head})
	}

	return sendChunks(data, chunkSize, func(chunk []byte) error {
		return srv.Send(&api.Plaintext{
			Body: chunk,
		})
	})
}

func receiveReferencesAndGrantSpec(srv api.Grant_SealServer) (reference.Refs, *grant.Spec, error) {
	var refs reference.Refs
	spec, err := consumeHeadFromReferenceAndGrantSpec(srv.Recv())
	if err != nil {
		return nil, nil, err
	}

	for {
		refAndSpec, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return refs, spec, nil
			}
			return nil, nil, err
		}

		if s := refAndSpec.GrantSpec; s != nil {
			return nil, nil, fmt.Errorf("received multiple grant specs but there can be at most one")
		}

		refs = append(refs, refAndSpec.Reference)
	}
}

func consumeHeadFromReferenceAndGrantSpec(rgs *api.ReferenceAndGrantSpec, err error) (*grant.Spec, error) {
	if err != nil {
		return nil, err
	}

	spec := rgs.GetGrantSpec()
	if spec == nil {
		return nil, fmt.Errorf("grant spec expected in first message")
	}

	return spec, nil
}

func consumeHeadFromPlaintextAndGrantSpec(ptgs *api.PlaintextAndGrantSpec, err error) (*grant.Spec, *api.Header, error) {
	head, err := consumeHeadFromPlaintext(ptgs.GetPlaintext(), err)
	if err != nil {
		return nil, nil, err
	}

	spec := ptgs.GetGrantSpec()
	if spec == nil {
		return nil, nil, fmt.Errorf("grant spec expected in first message")
	}

	return spec, head, nil
}

func consumeHeadFromPlaintext(pt *api.Plaintext, err error) (*api.Header, error) {
	if err != nil {
		return nil, err
	} else if len(pt.GetBody()) > 0 {
		return nil, fmt.Errorf("no data expected in first message")
	}

	return pt.GetHead(), nil
}

func consumeBodyFromPlaintextAndGrantSpec(ptgs *api.PlaintextAndGrantSpec, acc []byte, chunkSize int, cb func([]byte) error) ([]byte, error) {
	leftover, err := consumeBodyFromPlaintext(ptgs.GetPlaintext(), acc, chunkSize, cb)
	if err != nil {
		return nil, err
	} else if ptgs.GetGrantSpec() != nil {
		return nil, fmt.Errorf("received multiple grant specs but there can be at most one")
	}

	return leftover, nil
}

func consumeBodyFromPlaintext(pt *api.Plaintext, acc []byte, chunkSize int, cb func([]byte) error) ([]byte, error) {
	if pt.GetHead() != nil {
		return nil, fmt.Errorf("received header in a Plaintext frame other than the first")
	}

	// acc is the lefto
	limit := chunkSize - len(acc)
	data := pt.GetBody()

	if len(data) < limit {
		return append(acc, data...), nil
	} else if len(data) == limit {
		chunk := make([]byte, chunkSize)
		copy(chunk, append(acc, data...))
		return nil, cb(chunk)
	}

	// TODO: investigate dirty write
	chunk := make([]byte, chunkSize)
	copy(chunk, acc)
	copy(chunk[len(acc):], data[:limit])

	leftover := make([]byte, len(data)-limit)
	copy(leftover, data[limit:])

	if err := cb(chunk); err != nil {
		return nil, err
	}

	index := 0
	for index = 0; index < len(leftover)-chunkSize; index += chunkSize {
		next := make([]byte, chunkSize)
		copy(next, leftover[index:index+chunkSize])
		if err := cb(next); err != nil {
			return nil, err
		}
	}

	next := make([]byte, len(leftover[index:]))
	copy(next, leftover[index:])
	return next, nil
}

func ReceiveAllPlaintexts(cli interface {
	Recv() (*api.Plaintext, error)
}) (*api.Plaintext, error) {
	plaintext := new(api.Plaintext)

	for {
		pt, err := cli.Recv()
		if err != nil {
			if err == io.EOF {
				return plaintext, nil
			}

			return nil, err
		}

		plaintext.Body = append(plaintext.Body, pt.GetBody()...)
		if plaintext.Head == nil {
			plaintext.Head = pt.GetHead()
		}
	}
}

func ReceiveAllReferences(cli interface {
	Recv() (*reference.Ref, error)
}) (reference.Refs, error) {
	refs := make(reference.Refs, 0)

	for {
		ref, err := cli.Recv()
		if err != nil {
			if err == io.EOF {
				return refs, nil
			}

			return nil, err
		}

		refs = append(refs, ref)
	}
}

func ReceiveAllAddresses(cli interface {
	Recv() (*api.Address, error)
}) ([]*api.Address, error) {
	addrs := make([]*api.Address, 0)

	for {
		addr, err := cli.Recv()
		if err != nil {
			if err == io.EOF {
				return addrs, nil
			}

			return nil, err
		}

		addrs = append(addrs, addr)
	}
}

// StreamFileFrom provides a convenience wrapper over an io.Reader
func StreamFileFrom(reader io.Reader, chunkSize int, sender func(chunk []byte) error) error {
	out := make([]byte, chunkSize)
	for {
		n, err := reader.Read(out)
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		if err = sender(out[:n]); err != nil {
			return err
		}
	}
}

// StreamFileTo provides a convenience wrapper over an io.Writer
func StreamFileTo(writer io.Writer, receiver func() ([]byte, error)) error {
	for {
		chunk, err := receiver()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		n, err := writer.Write(chunk)
		if err != nil {
			return err
		} else if n != len(chunk) {
			return fmt.Errorf("failed to write data")
		}
	}
}

// ReadStream returns an interface when it is non-nil
func ReadStream(receiver func() (interface{}, error)) (interface{}, error) {
	for {
		chunk, err := receiver()
		if err != nil {
			return nil, err
		} else if chunk != nil {
			return chunk, nil
		}
	}
}

func sendChunks(data []byte, chunkSize int, sender func(chunk []byte) error) error {
	var err error
	for len(data) > chunkSize {
		err = sender(data[:chunkSize])
		if err != nil {
			return err
		}
		data = data[chunkSize:]
	}
	if len(data) == 0 {
		return nil
	}
	return sender(data)
}
