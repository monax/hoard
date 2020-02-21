package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v7"
	"github.com/monax/hoard/v7/api"
	"github.com/monax/hoard/v7/grant"
)

// PutSeal encrypts and stores data then prints a grant
func (client *Client) PutSeal(cmd *cli.Cmd) {
	salt := addStringOpt(cmd, "salt", saltOpt)
	key := addStringOpt(cmd, "key", keyOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(*chunk)

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

		err = hoard.StreamFileFrom(os.Stdin, *chunk, func(data []byte) error {
			return putseal.Send(&api.PlaintextAndGrantSpec{
				Plaintext: &api.Plaintext{
					Body: data,
				},
			})
		})
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
		spec := grant.Spec{Plaintext: &grant.PlaintextSpec{}}
		if *key != "" {
			spec = grant.Spec{
				Plaintext: nil,
				Symmetric: &grant.SymmetricSpec{PublicID: *key},
			}
		}

		refs := readReferences()
		seal, err := client.grant.Seal(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		for _, ref := range refs {
			if err = seal.Send(&api.ReferenceAndGrantSpec{
				Reference: ref,
				GrantSpec: &spec,
			}); err != nil {
				fatalf("Error sending data: %v", err)
			}
		}
		if err = seal.CloseSend(); err != nil {
			fatalf("Error closing send: %v", err)
		}

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

		refs, err := hoard.ReceiveAllReferences(unseal)
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

		err = hoard.StreamFileTo(os.Stdout, func() ([]byte, error) {
			plaintext, err := unsealget.Recv()
			return plaintext.GetBody(), err
		})
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
