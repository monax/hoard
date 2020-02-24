package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
	hoard "github.com/monax/hoard/v8"
	"github.com/monax/hoard/v8/api"
)

// Cat retrieves encrypted data from store
func (client *Client) Cat(cmd *cli.Cmd) {
	cmd.Action = func() {
		refs := readReferences()

		pull, err := client.storage.Pull(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		for _, ref := range refs {
			if err = pull.Send(&api.Address{Address: ref.Address}); err != nil {
				fatalf("Error sending data: %v", err)
			}
		}
		if err = pull.CloseSend(); err != nil {
			fatalf("Error closing send: %v", err)
		}

		err = hoard.StreamFileTo(os.Stdout, func() ([]byte, error) {
			ciphertext, err := pull.Recv()
			return ciphertext.GetEncryptedData(), err
		})
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}
	}
}

// Get retrieves and decrypts data from store
func (client *Client) Get(cmd *cli.Cmd) {
	cmd.Action = func() {
		refs := readReferences()

		get, err := client.cleartext.Get(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		for _, ref := range refs {
			if err = get.Send(ref); err != nil {
				fatalf("Error sending data: %v", err)
			}
		}
		if err = get.CloseSend(); err != nil {
			fatalf("Error closing send: %v", err)
		}

		err = hoard.StreamFileTo(os.Stdout, func() ([]byte, error) {
			plaintext, err := get.Recv()
			return plaintext.GetBody(), err
		})
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}
	}
}

// Insert data directly into store, preferably pre-encrypted
func (client *Client) Insert(cmd *cli.Cmd) {
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(*chunk)

		// If given address use it
		push, err := client.storage.Push(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = hoard.StreamFileFrom(os.Stdin, *chunk, func(data []byte) error {
			return push.Send(&api.Ciphertext{EncryptedData: data})
		})
		if err != nil {
			fatalf("Error sending data: %v", err)
		}
		if err = push.CloseSend(); err != nil {
			fatalf("Error closing send: %v", err)
		}

		addrs, err := hoard.ReceiveAllAddresses(push)
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}

		fmt.Printf("%s\n", jsonString(addrs))
	}
}

// Put encrypts data and stores it
func (client *Client) Put(cmd *cli.Cmd) {
	// TODO: check if salt is too big
	salt := addStringOpt(cmd, "salt", saltOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(*chunk)

		put, err := client.cleartext.Put(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = put.Send(&api.Plaintext{Head: &api.Header{Salt: parseSalt(salt)}})
		if err != nil {
			fatalf("Error sending head: %v", err)
		}

		err = hoard.StreamFileFrom(os.Stdin, *chunk, func(data []byte) error {
			return put.Send(&api.Plaintext{Body: data})
		})
		if err != nil {
			fatalf("Error sending body: %v", err)
		}

		if err = put.CloseSend(); err != nil {
			fatalf("Error closing send: %v", err)
		}

		refs, err := hoard.ReceiveAllReferences(put)
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}

		fmt.Printf("%s\n", jsonString(refs))
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
