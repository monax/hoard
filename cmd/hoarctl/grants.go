package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v5"
	"github.com/monax/hoard/v5/api"
	"github.com/monax/hoard/v5/grant"
)

// PutSeal encrypts and stores data then prints a grant
func (client *Client) PutSeal(cmd *cli.Cmd) {
	salt := addStringOpt(cmd, "salt", saltOpt)
	key := addStringOpt(cmd, "key", keyOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(*chunk)

		spec := grant.Spec{Plaintext: &grant.PlaintextSpec{}}
		if *key != "" {
			spec = grant.Spec{
				Plaintext: nil,
				Symmetric: &grant.SymmetricSpec{PublicID: *key},
			}
		}

		data := readData(os.Stdin)
		seal, err := client.grant.PutSeal(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = hoard.SendPlaintextAndGrantSpec(seal, &spec, data, parseSalt(salt), *chunk)
		if err != nil {
			fatalf("Error sending data: %v", err)
		}

		grant, err := seal.CloseAndRecv()
		if err != nil {
			fatalf("Error closing client: %v", err)
		}

		fmt.Printf("%s\n", jsonString(grant))
	}
}

// Seal reads encrypted data then prints a grant
func (client *Client) Seal(cmd *cli.Cmd) {
	address := addStringOpt(cmd, "address", addrOpt)
	key := addStringOpt(cmd, "key", keyOpt)

	cmd.Action = func() {
		spec := grant.Spec{Plaintext: &grant.PlaintextSpec{}}
		if *key != "" {
			spec = grant.Spec{
				Plaintext: nil,
				Symmetric: &grant.SymmetricSpec{PublicID: *key},
			}
		}

		ref := readReference(address)
		seal, err := client.grant.Seal(context.Background(),
			&api.ReferenceAndGrantSpec{
				Reference: ref,
				GrantSpec: &spec,
			},
		)

		if err != nil {
			fatalf("Error sealing data: %v", err)
		}
		fmt.Printf("%s\n", jsonString(seal))
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

		ref, err := client.grant.Reseal(context.Background(),
			&api.GrantAndGrantSpec{
				Grant:     prev,
				GrantSpec: &next,
			})

		if err != nil {
			fatalf("Error resealing data: %v", err)
		}
		fmt.Printf("%s\n", jsonString(ref))
	}
}

// Unseal reads a grant then prints the original reference
func (client *Client) Unseal(cmd *cli.Cmd) {
	cmd.Action = func() {
		grt := readGrant()
		ref, err := client.grant.Unseal(context.Background(), grt)
		if err != nil {
			fatalf("Error unsealing data: %v", err)
		}
		fmt.Printf("%s\n", jsonString(ref))
	}
}

// UnsealGet reads a grant, decrypts and prints the stored data
func (client *Client) UnsealGet(cmd *cli.Cmd) {
	cmd.Action = func() {
		grt := readGrant()

		unseal, err := client.grant.UnsealGet(context.Background(), grt)
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		data, _, err := hoard.ReceivePlaintext(unseal)
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}

		os.Stdout.Write(data)
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
		return
	}
}
