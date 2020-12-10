package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/monax/hoard/v8/stores"

	"github.com/monax/hoard/v8/reference"

	cli "github.com/jawher/mow.cli"
	hoard "github.com/monax/hoard/v8"
	"github.com/monax/hoard/v8/api"
)

// Cat retrieves encrypted data from store
func (client *Client) Cat(cmd *cli.Cmd) {
	cmd.Action = func() {
		pull, err := client.storage.Pull(context.Background())
		decoder := json.NewDecoder(os.Stdin)
		err = hoard.NewStreamer().WithSend(func(chunk []byte) error {
			refs := new([]*reference.Ref)
			err := decoder.Decode(refs)
			if err != nil {
				return err
			}
			for _, ref := range *refs {
				err := pull.Send(&api.Address{Address: ref.Address})
				if err != nil {
					return err
				}
			}
			return nil
		}).WithCloseSend(pull.CloseSend).WithRecv(func() ([]byte, error) {
			ciphertext, err := pull.Recv()
			if err != nil {
				return nil, err
			}
			return ciphertext.GetEncryptedData(), err
		}).WithOutput(os.Stdout).Stream(context.Background())
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}
	}
}

// Get retrieves and decrypts data from store
func (client *Client) Get(cmd *cli.Cmd) {
	cmd.Action = func() {
		get, err := client.cleartext.Get(context.Background())
		err = hoard.NewStreamer().
			WithSend(readReferences(get.Send)).
			WithCloseSend(get.CloseSend).
			WithRecv(func() ([]byte, error) {
				plaintext, err := get.Recv()
				if err != nil {
					return nil, err
				}
				return plaintext.GetBody(), err
			}).WithOutput(os.Stdout).Stream(context.Background())
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}
	}
}

// Insert data directly into store, preferably pre-encrypted
func (client *Client) Insert(cmd *cli.Cmd) {
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(int64(*chunk))

		// If given address use it
		push, err := client.storage.Push(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		var addresses []*api.Address

		err = hoard.NewStreamer().WithChunkSize(int64(*chunk)).
			WithInput(os.Stdin).
			WithSend(func(data []byte) error {
				return push.Send(&api.Ciphertext{EncryptedData: data})
			}).
			WithCloseSend(push.CloseSend).
			WithRecv(func() ([]byte, error) {
				address, err := push.Recv()
				if err != nil {
					return nil, err
				}
				addresses = append(addresses, address)
				return nil, nil
			}).
			Stream(context.Background())
		if err != nil {
			fatalf("Error sending data: %v", err)
		}

		fmt.Printf("%s\n", jsonString(addresses))
	}
}

// Put encrypts data and stores it
func (client *Client) Put(cmd *cli.Cmd) {
	// TODO: check if salt is too big
	salt := addStringOpt(cmd, "salt", saltOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(int64(*chunk))

		put, err := client.cleartext.Put(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = put.Send(&api.Plaintext{Head: &api.Header{Salt: parseSalt(salt)}})
		if err != nil {
			fatalf("Error sending head: %v", err)
		}

		refs := []*reference.Ref{}
		err = hoard.NewStreamer().WithChunkSize(int64(*chunk)).
			WithInput(os.Stdin).
			WithSend(func(data []byte) error {
				return put.Send(&api.Plaintext{Body: data})
			}).
			WithCloseSend(put.CloseSend).
			WithRecv(recvReferences(&refs, put.Recv)).
			Stream(context.Background())

		if err != nil {
			fatalf("Error sending body: %v", err)
		}

		fmt.Printf("%s\n", jsonString(refs))
	}
}

// Delete removes the blob located at the provided address
func (client *Client) Delete(cmd *cli.Cmd) {
	cmd.Action = func() {
		err := hoard.NewStreamer().WithSend(readReferences(func(ref *reference.Ref) error {
			_, err := client.storage.Delete(context.Background(),
				&api.Address{
					Address: ref.Address,
				})
			return err
		})).Stream(context.Background())
		if err != nil {
			fatalf("Error deleting data: %v", err)
		}

	}
}

// Stat retrieves info about the stored data
func (client *Client) Stat(cmd *cli.Cmd) {
	cmd.Action = func() {
		var statInfos []*stores.StatInfo
		err := hoard.NewStreamer().WithSend(readReferences(func(ref *reference.Ref) error {
			statInfo, err := client.storage.Stat(context.Background(),
				&api.Address{Address: ref.Address})
			if err != nil {
				return err
			}
			statInfos = append(statInfos, statInfo)
			return nil
		})).Stream(context.Background())

		fmt.Printf("%s\n", jsonString(statInfos))
		if err != nil {
			fatalf("Error querying blobs: %v", err)
		}
	}
}
