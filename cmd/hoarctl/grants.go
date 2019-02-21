package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard"
	"github.com/monax/hoard/grant"
)

// PutSeal encrypts and stores data then prints a grant
func (client *Client) PutSeal(cmd *cli.Cmd) {
	salt := addOpt(cmd, "salt", saltOpt, "").(*string)
	key := addOpt(cmd, "key", saltOpt, "").(*string)

	cmd.Action = func() {

		var seal *grant.Grant
		var err error

		spec := grant.Spec{Plaintext: &grant.PlaintextSpec{}}
		if key != nil {
			spec = grant.Spec{
				Plaintext: nil,
				Symmetric: &grant.SymmetricSpec{SecretID: string(parseSalt(*key))},
			}
		}

		data := readData()
		seal, err = client.grant.PutSeal(context.Background(),
			&hoard.PlaintextAndGrantSpec{
				Plaintext: &hoard.Plaintext{
					Data: data,
					Salt: parseSalt(*salt),
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
	address := addOpt(cmd, "address", addrOpt, "").(*string)
	key := addOpt(cmd, "key", saltOpt, "").(*string)

	cmd.Action = func() {
		spec := grant.Spec{Plaintext: &grant.PlaintextSpec{}}
		if key != nil {
			spec = grant.Spec{
				Plaintext: nil,
				Symmetric: &grant.SymmetricSpec{SecretID: string(parseSalt(*key))},
			}
		}

		ref := readReference(address)
		seal, err := client.grant.Seal(context.Background(),
			&hoard.ReferenceAndGrantSpec{
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
	salt := addOpt(cmd, "salt", saltOpt, "").(*string)

	cmd.Action = func() {
		prev := readGrant()
		next := grant.Spec{Plaintext: &grant.PlaintextSpec{}}

		if salt != nil {
			next = grant.Spec{
				Plaintext: nil,
				Symmetric: &grant.SymmetricSpec{SecretID: string(parseSalt(*salt))},
			}
		}

		ref, err := client.grant.Reseal(context.Background(),
			&hoard.GrantAndGrantSpec{
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
