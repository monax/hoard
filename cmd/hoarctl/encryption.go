package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v8"
	"github.com/monax/hoard/v8/api"
	"github.com/monax/hoard/v8/reference"
)

// Decrypt does what it says on the tin
func (client *Client) Decrypt(cmd *cli.Cmd) {
	secretKey := addStringOpt(cmd, "key", secretOpt)
	salt := addStringOpt(cmd, "salt", saltOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(int64(*chunk))

		dec, err := client.encryption.Decrypt(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = hoard.NewStreamer().WithChunkSize(int64(*chunk)).WithInput(os.Stdin).
			WithSend(
				func(data []byte) error {
					return dec.Send(&api.ReferenceAndCiphertext{
						Reference: &reference.Ref{
							SecretKey: readBase64(secretKey),
							Salt:      parseSalt(salt),
						},
						Ciphertext: &api.Ciphertext{
							EncryptedData: data,
						},
					})
				}).
			WithCloseSend(dec.CloseSend).
			WithRecv(func() ([]byte, error) {
				plaintext, err := dec.Recv()
				return plaintext.GetBody(), err
			}).WithOutput(os.Stdout).
			Stream(context.Background())

		if err != nil {
			fatalf("Error streaming data: %v", err)
		}

	}
}

// Encrypt also does what it says on the tin
func (client *Client) Encrypt(cmd *cli.Cmd) {
	salt := addStringOpt(cmd, "salt", saltOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(int64(*chunk))

		enc, err := client.encryption.Encrypt(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = enc.Send(&api.Plaintext{Head: &api.Header{Salt: parseSalt(salt)}})
		if err != nil {
			fatalf("Could not send encryption head: %v", err)
		}

		// Throw away the encrypted header
		_, err = enc.Recv()
		if err != nil {
			fatalf("Could not receive (and discard) header")
		}

		err = hoard.NewStreamer().
			WithChunkSize(int64(*chunk)).
			WithInput(os.Stdin).
			WithSend(
				func(data []byte) error {
					return enc.Send(&api.Plaintext{Body: data})
				},
			).
			WithCloseSend(enc.CloseSend).
			WithRecv(func() ([]byte, error) {
				refAndCiphertext, err := enc.Recv()
				if err != nil {
					return nil, err
				}
				encryptedData := refAndCiphertext.GetCiphertext().GetEncryptedData()
				return encryptedData, nil
			}).
			WithOutput(os.Stdout).
			Stream(context.Background())
		if err != nil {
			fatalf("Error sending body: %v", err)
		}

	}
}

// Ref encrypts as above, but then reads the reference
func (client *Client) Ref(cmd *cli.Cmd) {
	salt := addStringOpt(cmd, "salt", saltOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(int64(*chunk))

		enc, err := client.encryption.Encrypt(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = enc.Send(&api.Plaintext{Head: &api.Header{Salt: parseSalt(salt)}})
		if err != nil {
			fatalf("Error sending head: %v", err)
		}

		var refs reference.Refs

		err = hoard.NewStreamer().
			WithChunkSize(int64(*chunk)).
			WithInput(os.Stdin).
			WithSend(func(data []byte) error {
				return enc.Send(&api.Plaintext{Body: data})
			}).
			WithCloseSend(enc.CloseSend).
			WithRecv(
				func() ([]byte, error) {
					refAndCiphertext, err := enc.Recv()
					if err != nil {
						return nil, err
					}
					refs = append(refs, refAndCiphertext.GetReference())
					return nil, nil
				}).Stream(context.Background())
		if err != nil {
			fatalf("Error streaming data: %v", err)
		}

		fmt.Printf("%s\n", jsonString(refs))
	}
}
