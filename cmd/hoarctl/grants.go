package main

import (
	"context"
	"fmt"
	"os"

	"github.com/monax/hoard/v8/reference"

	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v8"
	"github.com/monax/hoard/v8/api"
	"github.com/monax/hoard/v8/grant"
)

// PutSeal encrypts and stores data then prints a grant
func (client *Client) PutSeal(cmd *cli.Cmd) {
	salt := addStringOpt(cmd, "salt", saltOpt)
	key := addStringOpt(cmd, "key", keyOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(int64(*chunk))

		spec := &grant.Spec{Plaintext: &grant.PlaintextSpec{}}
		if *key != "" {
			spec = &grant.Spec{
				Plaintext: nil,
				Symmetric: &grant.SymmetricSpec{PublicID: *key},
			}
		}

		putseal, err := client.grant.PutSeal(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = putseal.Send(&api.PlaintextAndGrantSpec{
			Plaintext: &api.Plaintext{
				Head: &api.Header{
					Salt: parseSalt(salt),
				},
			},
			GrantSpec: spec,
		})
		if err != nil {
			fatalf("Error sending head: %v", err)
		}

		err = hoard.NewStreamer().WithChunkSize(int64(*chunk)).WithInput(os.Stdin).WithSend(func(data []byte) error {
			return putseal.Send(&api.PlaintextAndGrantSpec{
				Plaintext: &api.Plaintext{
					Body: data,
				},
			})
		}).Stream(context.Background())
		if err != nil {
			fatalf("Error sending body: %v", err)
		}

		if err = putseal.CloseSend(); err != nil {
			fatalf("Error closing send: %v", err)
		}

		grt, err := putseal.CloseAndRecv()
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}

		fmt.Printf("%s\n", jsonString(grt))
	}
}

// Seal reads encrypted data then prints a grant
func (client *Client) Seal(cmd *cli.Cmd) {
	key := addStringOpt(cmd, "key", keyOpt)

	cmd.Action = func() {
		spec := &grant.Spec{Plaintext: &grant.PlaintextSpec{}}
		if *key != "" {
			spec = &grant.Spec{
				Plaintext: nil,
				Symmetric: &grant.SymmetricSpec{PublicID: *key},
			}
		}

		seal, err := client.grant.Seal(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		// Send spec
		err = seal.Send(&api.ReferenceAndGrantSpec{
			GrantSpec: spec,
		})
		if err != nil {
			fatalf("Could not send spec to seal: %v", err)
		}

		err = hoard.NewStreamer().
			WithSend(readReferences(func(ref *reference.Ref) error {
				return seal.Send(&api.ReferenceAndGrantSpec{
					Reference: ref,
				})
			})).
			WithCloseSend(seal.CloseSend).
			Stream(context.Background())

		grt, err := seal.CloseAndRecv()
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}

		fmt.Printf("%s\n", jsonString(grt))
	}
}

// Reseal reads a grant then prints a new grant
func (client *Client) Reseal(cmd *cli.Cmd) {
	key := addStringOpt(cmd, "key", keyOpt)

	cmd.Action = func() {
		prev := readGrant()
		next := grant.Spec{Plaintext: &grant.PlaintextSpec{}}

		if *key != "" {
			next = grant.Spec{
				Plaintext: nil,
				Symmetric: &grant.SymmetricSpec{PublicID: *key},
			}
		}

		grt, err := client.grant.Reseal(context.Background(),
			&api.GrantAndGrantSpec{
				Grant:     prev,
				GrantSpec: &next,
			})

		if err != nil {
			fatalf("Error resealing data: %v", err)
		}
		fmt.Printf("%s\n", jsonString(grt))
	}
}

// Unseal reads a grant then prints the original reference
func (client *Client) Unseal(cmd *cli.Cmd) {
	cmd.Action = func() {
		grt := readGrant()
		unseal, err := client.grant.Unseal(context.Background(), grt)
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		if err = unseal.CloseSend(); err != nil {
			fatalf("Error closing send: %v", err)
		}

		refs := []*reference.Ref{}
		err = hoard.NewStreamer().
			WithRecv(recvReferences(&refs, unseal.Recv)).
			Stream(context.Background())

		if err != nil {
			fatalf("Error receiving data: %v", err)
		}

		fmt.Printf("%s\n", jsonString(refs))
	}
}

// UnsealGet reads a grant, decrypts and prints the stored data
func (client *Client) UnsealGet(cmd *cli.Cmd) {
	cmd.Action = func() {
		grt := readGrant()

		unsealget, err := client.grant.UnsealGet(context.Background(), grt)
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = hoard.NewStreamer().
			WithRecv(func() ([]byte, error) {
				plaintext, err := unsealget.Recv()
				if err != nil {
					return nil, err
				}
				return plaintext.GetBody(), nil
			}).
			WithOutput(os.Stdout).
			Stream(context.Background())
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}
	}
}

// UnsealDelete reads a grant and deletes the encrypted data
func (client *Client) UnsealDelete(cmd *cli.Cmd) {
	cmd.Action = func() {
		grt := readGrant()
		_, err := client.grant.UnsealDelete(context.Background(), grt)
		if err != nil {
			fatalf("Error unsealing data: %v", err)
		}
	}
}
