package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/monax/hoard/v8"

	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v8/grant"
	"github.com/monax/hoard/v8/reference"
)

// extra cli options
func addStringOpt(cmd *cli.Cmd, arg, desc string) *string {
	opt := cmd.StringOpt(fmt.Sprintf("%s %s", string(arg[0]), arg), "", desc)
	cmd.Spec += fmt.Sprintf("[-%s | --%s]", string(arg[0]), arg)
	return opt
}

func addIntOpt(cmd *cli.Cmd, arg, desc string, def int) *int {
	opt := cmd.IntOpt(fmt.Sprintf("%s %s", string(arg[0]), arg), def, desc)
	cmd.Spec += fmt.Sprintf("[-%s | --%s]", string(arg[0]), arg)
	return opt
}

func validateChunkSize(chunkSize int64) {
	if chunkSize == 0 {
		fatalf("Chunk size cannot be 0")
	} else if chunkSize > hoard.MaxChunkSize {
		fatalf("Chunk size cannot be greater than 4Mb")
	}
}

func parseSalt(saltString *string) []byte {
	if saltString == nil {
		return nil
	}
	saltBytes, err := base64.StdEncoding.DecodeString(*saltString)
	if err == nil {
		return saltBytes
	}
	return ([]byte)(*saltString)
}

func jsonString(v interface{}) string {
	bs, err := json.Marshal(v)
	if err != nil {
		fatalf("Could not serialise '%s' to json: %v", err)
	}
	return string(bs)

}

func readReferences(send func(ref *reference.Ref) error) func(chunk []byte) error {
	decoder := json.NewDecoder(os.Stdin)
	return func(chunk []byte) error {
		refs := new([]*reference.Ref)
		err := decoder.Decode(refs)
		if err != nil {
			return err
		}
		for _, ref := range *refs {
			err := send(ref)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
func recvReferences(refs *[]*reference.Ref, recv func() (*reference.Ref, error)) func() ([]byte, error) {
	return func() ([]byte, error) {
		ref, err := recv()
		if err != nil {
			return nil, err
		}
		*refs = append(*refs, ref)
		return nil, nil
	}
}

func readGrant() *grant.Grant {
	grt := new(grant.Grant)
	err := json.NewDecoder(os.Stdin).Decode(grt)
	if err != nil {
		fatalf("Could not read grant from STDIN: %v", err)
	}
	return grt
}

func readBase64(base64String *string) []byte {
	if base64String == nil {
		return nil
	}
	secretKeyBytes, err := base64.StdEncoding.DecodeString(*base64String)
	if err != nil {
		fatalf("Could not decode '%s' as base64-encoded string", base64String)
	}
	return secretKeyBytes
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
