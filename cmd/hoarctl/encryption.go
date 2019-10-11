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

// Decrypt does what it says on the tin
func (client *Client) Decrypt(cmd *cli.Cmd) {
	secretKey := addStringOpt(cmd, "key", secretOpt)
	salt := addStringOpt(cmd, "salt", saltOpt)

	cmd.Action = func() {
		encryptedData := readData(os.Stdin)
		dec, err := client.encryption.Decrypt(context.Background(),
			&api.ReferenceAndCiphertext{
				Reference: &reference.Ref{
					SecretKey: readBase64(secretKey),
					Salt:      parseSalt(salt),
				},
				Ciphertext: &api.Ciphertext{
					EncryptedData: encryptedData,
				},
			})
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		data, _, err := hoard.ReceivePlaintext(dec)
		if err != nil {
			fatalf("Error receiving data: %v", err)
		}

		os.Stdout.Write(data)
	}
}

// Encrypt also does what it says on the tin
func (client *Client) Encrypt(cmd *cli.Cmd) {
	salt := addStringOpt(cmd, "salt", saltOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(*chunk)

		data := readData(os.Stdin)
		enc, err := client.encryption.Encrypt(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = hoard.SendPlaintext(enc, data, parseSalt(salt), *chunk)
		if err != nil {
			fatalf("Error sending data: %v", err)
		}

		refAndCiphertext, err := enc.CloseAndRecv()
		if err != nil {
			fatalf("Error closing client: %v", err)
		}

		os.Stdout.Write(refAndCiphertext.Ciphertext.EncryptedData)
	}
}

// Ref encrypts as above, but then packages the data in a ref
func (client *Client) Ref(cmd *cli.Cmd) {
	salt := addStringOpt(cmd, "salt", saltOpt)
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)

	cmd.Action = func() {
		validateChunkSize(*chunk)

		data := readData(os.Stdin)
		enc, err := client.encryption.Encrypt(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		err = hoard.SendPlaintext(enc, data, parseSalt(salt), *chunk)
		if err != nil {
			fatalf("Error sending data: %v", err)
		}

		refAndCiphertext, err := enc.CloseAndRecv()
		if err != nil {
			fatalf("Error closing client: %v", err)
		}

		fmt.Printf("%s\n", jsonString(refAndCiphertext.Reference))
	}
}
