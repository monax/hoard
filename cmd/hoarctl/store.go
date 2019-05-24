package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v5/api"
	"github.com/monax/hoard/v5/reference"
)

// Cat retrieves encrypted data from store
func (client *Client) Cat(cmd *cli.Cmd) {
	address := addStringOpt(cmd, "address", addrOpt)

	cmd.Action = func() {
		ref := readReference(address)
		ciphertext, err := client.storage.Pull(context.Background(),
			&api.Address{Address: ref.Address})
		if err != nil {
			fatalf("Error querying data: %v", err)
		}
		os.Stdout.Write(ciphertext.EncryptedData)
	}
}

// Get retrieves and decrypts data from store
func (client *Client) Get(cmd *cli.Cmd) {
	address := addStringOpt(cmd, "address", addrOpt)
	secretKey := addStringOpt(cmd, "key", secretOpt)
	salt := addStringOpt(cmd, "salt", saltOpt)

	cmd.Action = func() {
		// If given address then try to read reference from arguments and option
		ref := readReference(address)
		if ref.SecretKey == nil {
			if secretKey == nil || *secretKey == "" {
				fatalf("A secret key must be provided in order to decrypt")
			}
			ref = &reference.Ref{
				Address:   readBase64(address),
				SecretKey: readBase64(secretKey),
				Salt:      parseSalt(salt),
			}
		}
		plaintext, err := client.cleartext.Get(context.Background(), ref)
		if err != nil {
			fatalf("Error retrieving data: %v", err)
		}
		os.Stdout.Write(plaintext.Data)
	}
}

// Insert data directly into store, preferably pre-encrypted
func (client *Client) Insert(cmd *cli.Cmd) {
	cmd.Action = func() {
		data := readData()
		// If given address use it
		address, err := client.storage.Push(context.Background(),
			&api.Ciphertext{EncryptedData: data})
		if err != nil {
			fatalf("Error querying data: %v", err)
		}
		fmt.Printf("%s\n", jsonString(address))
	}
}

// Put encrypts data and stores it
func (client *Client) Put(cmd *cli.Cmd) {
	salt := addStringOpt(cmd, "salt", saltOpt)

	cmd.Action = func() {
		data := readData()
		ref, err := client.cleartext.Put(context.Background(),
			&api.Plaintext{
				Data: data,
				Salt: parseSalt(salt),
			})
		if err != nil {
			fatalf("Error storing data: %v", err)
		}
		fmt.Printf("%s\n", jsonString(ref))
	}
}

// Stat retrieves info about the stored data
func (client *Client) Stat(cmd *cli.Cmd) {
	address := addStringOpt(cmd, "address", addrOpt)

	cmd.Action = func() {
		ref := readReference(address)
		statInfo, err := client.storage.Stat(context.Background(),
			&api.Address{Address: ref.Address})
		if err != nil {
			fatalf("Error querying data: %v", err)
		}
		fmt.Printf("%s\n", jsonString(statInfo))
	}
}
