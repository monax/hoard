package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v7"
	"github.com/monax/hoard/v7/api"
	"github.com/monax/hoard/v7/reference"
)

// Decrypt does what it says on the tin
func (client *Client) Decrypt(cmd *cli.Cmd) {
	secretKey := addStringOpt(cmd, "key", secretOpt)
	salt := addStringOpt(cmd, "salt", saltOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(*chunk)

		encryptedData := readData(os.Stdin)
		dec, err := client.encryption.Decrypt(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = hoard.StreamFileFrom(os.Stdin, *chunk, func(data []byte) error {
			return dec.Send(&api.ReferenceAndCiphertext{
				Reference: &reference.Ref{
					SecretKey: readBase64(secretKey),
					Salt:      parseSalt(salt),
				},
				Ciphertext: &api.Ciphertext{
					EncryptedData: encryptedData,
				},
			})
		})
		if err != nil {
			fatalf("Error sending data: %v", err)
		}

		err = hoard.StreamFileTo(os.Stdout, func() ([]byte, error) {
			plaintext, err := dec.Recv()
			return plaintext.GetBody(), err
		})
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}
	}
}

// Encrypt also does what it says on the tin
func (client *Client) Encrypt(cmd *cli.Cmd) {
	salt := addStringOpt(cmd, "salt", saltOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(*chunk)

		enc, err := client.encryption.Encrypt(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = enc.Send(&api.Plaintext{Head: &api.Header{Salt: parseSalt(salt)}})
		if err != nil {
			fatalf("Error sending head: %v", err)
		}

		err = hoard.StreamFileFrom(os.Stdin, *chunk, func(data []byte) error {
			return enc.Send(&api.Plaintext{Body: data})
		})
		if err != nil {
			fatalf("Error sending body: %v", err)
		}

		err = hoard.StreamFileTo(os.Stdout, func() ([]byte, error) {
			refAndCiphertext, err := enc.Recv()
			return refAndCiphertext.GetCiphertext().GetEncryptedData(), err
		})
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}
	}
}

// Ref encrypts as above, but then reads the reference
func (client *Client) Ref(cmd *cli.Cmd) {
	salt := addStringOpt(cmd, "salt", saltOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(*chunk)

		enc, err := client.encryption.Encrypt(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = enc.Send(&api.Plaintext{Head: &api.Header{Salt: parseSalt(salt)}})
		if err != nil {
			fatalf("Error sending head: %v", err)
		}

		err = hoard.StreamFileFrom(os.Stdin, *chunk, func(data []byte) error {
			return enc.Send(&api.Plaintext{Body: data})
		})
		if err != nil {
			fatalf("Error sending body: %v", err)
		}

		ref, err := hoard.ReadStream(func() (interface{}, error) {
			refAndCiphertext, err := enc.Recv()
			return refAndCiphertext.GetReference(), err
		})
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}

		fmt.Printf("%s\n", jsonString(ref))
	}
}
