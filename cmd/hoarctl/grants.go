package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v4/grant"
	"github.com/monax/hoard/v4/services"
)

// PutSeal encrypts and stores data then prints a grant
func (client *Client) PutSeal(cmd *cli.Cmd) {
	salt := addStringOpt(cmd, "salt", saltOpt)
	key := addStringOpt(cmd, "key", keyOpt)

	cmd.Action = func() {
		var seal *grant.Grant
		var err error

		spec := grant.Spec{Plaintext: &grant.PlaintextSpec{}}
		if *key != "" {
			spec = grant.Spec{
				Plaintext: nil,
				Symmetric: &grant.SymmetricSpec{PublicID: *key},
			}
		}

		data := readData()
		seal, err = client.grant.PutSeal(context.Background(),
			&services.PlaintextAndGrantSpec{
				Plaintext: &services.Plaintext{
					Data: data,
					Salt: parseSalt(salt),
				},
				GrantSpec: &spec,
			},
		)

		if err != nil {
			fatalf("Error sealing data: %v", err)
		}
		fmt.Printf("%s\n", jsonString(seal))
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
			&services.ReferenceAndGrantSpec{
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
			&services.GrantAndGrantSpec{
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

		plaintext, err := client.grant.UnsealGet(context.Background(), grt)
		if err != nil {
			fatalf("Error unsealing data: %v", err)
		}
		os.Stdout.Write(plaintext.Data)
		return
	}
}
