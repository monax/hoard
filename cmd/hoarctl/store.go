package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v6"
	"github.com/monax/hoard/v6/api"
	"github.com/monax/hoard/v6/reference"
)

// Cat retrieves encrypted data from store
func (client *Client) Cat(cmd *cli.Cmd) {
	address := addStringOpt(cmd, "address", addrOpt)

	cmd.Action = func() {
		ref := readReference(address)
		pull, err := client.storage.Pull(context.Background(),
			&api.Address{Address: ref.Address})
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		data, err := hoard.ReceiveCiphertext(pull)
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}

		os.Stdout.Write(data)
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

		get, err := client.cleartext.Get(context.Background(), ref)
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		data, _, err := hoard.ReceivePlaintext(get)
		if err != nil {
			fatalf("Error retrieving data: %v", err)
		}
		os.Stdout.Write(data)
	}
}

// Insert data directly into store, preferably pre-encrypted
func (client *Client) Insert(cmd *cli.Cmd) {
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(*chunk)

		data := readData(os.Stdin)
		// If given address use it
		push, err := client.storage.Push(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = hoard.SendCiphertext(push, data, *chunk)
		if err != nil {
			fatalf("Error sending data: %v", err)
		}

		addr, err := push.CloseAndRecv()
		if err != nil {
			fatalf("Error closing client: %v", err)
		}

		fmt.Printf("%s\n", jsonString(addr))
	}
}

// Put encrypts data and stores it
func (client *Client) Put(cmd *cli.Cmd) {
	// TODO: check if salt is too big
	salt := addStringOpt(cmd, "salt", saltOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(*chunk)

		data := readData(os.Stdin)
		put, err := client.cleartext.Put(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = hoard.SendPlaintext(put, data, parseSalt(salt), *chunk)
		if err != nil {
			fatalf("Error sending data: %v", err)
		}

		ref, err := put.CloseAndRecv()
		if err != nil {
			fatalf("Error closing client: %v", err)
		}
		fmt.Printf("%s\n", jsonString(ref))
	}
}

// Delete removes the blob located at the provided address
func (client *Client) Delete(cmd *cli.Cmd) {
	address := addStringOpt(cmd, "address", addrOpt)

	cmd.Action = func() {
		ref := readReference(address)
		_, err := client.storage.Delete(context.Background(),
			&api.Address{
				Address: ref.Address,
			})
		if err != nil {
			fatalf("Error deleting data: %v", err)
		}
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
