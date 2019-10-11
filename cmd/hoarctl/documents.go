package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v5"
	"github.com/monax/hoard/v5/grant"
	"github.com/monax/hoard/v5/meta"

	"github.com/h2non/filetype"
)

// Upload a document with metadata to hoard
func (client *Client) Upload(cmd *cli.Cmd) {
	chunk := addIntOpt(cmd, "chunk", chunkOpt, chunkSize)
	file := addStringOpt(cmd, "file", fileOpt)
	salt := addStringOpt(cmd, "salt", saltOpt)

	cmd.Action = func() {
		validateChunkSize(*chunk)

		file, err := openFile(file)
		if err != nil {
			fatalf("Error reading file: %v", err)
		}
		data := readData(file)

		kind, err := filetype.Match(data)
		if err != nil {
			fatalf("Error inspecting file: %v", err)
		}

		upload, err := client.documents.Upload(context.Background())
		if err != nil {
			fatalf("Error starting client: %v", err)
		}

		spec := &grant.Spec{Plaintext: &grant.PlaintextSpec{}}
		err = hoard.SendDocumentAndGrant(upload, &meta.Document{
			Meta: &meta.Meta{
				Name:     file.Name(),
				MimeType: kind.MIME.Type,
			},
			Data: data,
		}, parseSalt(salt), spec, *chunk)
		if err != nil {
			fatalf("Error sending data: %v", err)
		}

		grant, err := upload.CloseAndRecv()
		if err != nil {
			fatalf("Error closing client: %v", err)
		}

		fmt.Printf("%s\n", jsonString(grant))
	}
}

// Download a document and print the metadata or body
func (client *Client) Download(cmd *cli.Cmd) {
	cmd.Command("head", "get metadata", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			doc := client.getFile()
			fmt.Printf("%s\n", jsonString(doc.Meta))
		}
	})

	cmd.Command("body", "get data", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			doc := client.getFile()
			os.Stdout.Write(doc.Data)
		}
	})
}

func (client *Client) getFile() *meta.Document {
	grt := readGrant()
	download, err := client.documents.Download(context.Background(), grt)
	if err != nil {
		fatalf("Error starting client: %v", err)
	}

	doc, _, err := hoard.ReceiveDocument(download)
	if err != nil {
		fatalf("Error sending data: %v", err)
	}

	err = download.CloseSend()
	if err != nil {
		fatalf("Error closing client: %v", err)
	}

	return doc
}
