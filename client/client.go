package client

import (
	"context"
	"fmt"
	"io"

	"github.com/monax/hoard/v8/api"
	"github.com/monax/hoard/v8/grant"
	"github.com/monax/hoard/v8/streamer"
	"google.golang.org/grpc"
)

type Client struct {
	grant api.GrantClient
}

type PlaintextStream struct {
	Head    *api.Header
	writeTo func(w io.Writer) (n int64, err error)
	closer  func() error
}

func (p *PlaintextStream) Close() error {
	return p.closer()
}

func (p *PlaintextStream) WriteTo(w io.Writer) (n int64, err error) {
	return p.writeTo(w)
}

func (p *PlaintextStream) GetHead() *api.Header {
	if p == nil {
		return nil
	}
	return p.Head
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{
		grant: api.NewGrantClient(conn),
	}
}

func (c Client) PutSeal(ctx context.Context,
	spec *grant.Spec,
	header *api.Header,
	plaintextReader io.Reader,
	opts ...grpc.CallOption) (*grant.Grant, error) {

	stream, err := c.grant.PutSeal(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("PutSeal: could not establish stream: %w", err)
	}
	defer stream.CloseSend()
	err = stream.Send(&api.PlaintextAndGrantSpec{
		Plaintext: &api.Plaintext{
			Head: header,
		},
		GrantSpec: spec,
	})
	if err != nil {
		return nil, fmt.Errorf("PutSeal: could not send grant spec: %w", err)
	}
	err = streamer.New().WithInput(plaintextReader).
		WithSend(func(chunk []byte) error {
			return stream.Send(&api.PlaintextAndGrantSpec{
				Plaintext: &api.Plaintext{
					Body: chunk,
				},
			})
		}).Stream(ctx)
	if err != nil {
		return nil, fmt.Errorf("PutSeal: could not read and send plaintext: %w", err)
	}
	grt, err := stream.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("PutSeal: could not get grant from stream: %w", err)
	}
	return grt, err
}

func (c Client) UnsealGet(ctx context.Context, grt *grant.Grant,
	opts ...grpc.CallOption) (*PlaintextStream, error) {
	stream, err := c.grant.UnsealGet(ctx, grt, opts...)
	if err != nil {
		return nil, fmt.Errorf("UnsealGet: could not establish stream: %w", err)
	}

	first, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("UnsealGet: not get first frame from stream: %w", err)
	}
	head := first.GetHead()

	return &PlaintextStream{
		Head:   head,
		closer: stream.CloseSend,
		writeTo: func(plaintextWriter io.Writer) (int64, error) {
			defer stream.CloseSend()
			_, n, err := streamer.New().WithRecv(func() ([]byte, error) {
				defer func() { first = nil }()
				if first.GetBody() != nil {
					return first.GetBody(), nil
				}
				plaintext, err := stream.Recv()
				if err != nil {
					return nil, err
				}
				return plaintext.Body, nil
			}).WithOutput(plaintextWriter).StreamCount(ctx)
			if err != nil {
				return n, fmt.Errorf("UnsealGet: could not receive and write plaintext: %w", err)
			}
			return 0, nil
		},
	}, nil
}

func (c Client) Seal(ctx context.Context, opts ...grpc.CallOption) (api.Grant_SealClient, error) {
	panic("implement me")
}

func (c Client) Unseal(ctx context.Context, in *grant.Grant, opts ...grpc.CallOption) (api.Grant_UnsealClient, error) {
	panic("implement me")
}

func (c Client) Reseal(ctx context.Context, in *api.GrantAndGrantSpec, opts ...grpc.CallOption) (*grant.Grant, error) {
	panic("implement me")
}

func (c Client) UnsealDelete(ctx context.Context, in *grant.Grant, opts ...grpc.CallOption) (api.Grant_UnsealDeleteClient, error) {
	panic("implement me")
}
