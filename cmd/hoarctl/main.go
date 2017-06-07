package main

import (
	"os"

	"context"
	"encoding/hex"
	"io/ioutil"
	"strings"

	"fmt"

	"net"
	"time"

	"encoding/json"

	"code.monax.io/platform/hoard/cmd/shared"
	"code.monax.io/platform/hoard/hoard"
	"github.com/jawher/mow.cli"
	"google.golang.org/grpc"
)

func main() {
	hoarctlApp := cli.App("hoarctl",
		"Command line interface to the hoard daemon a content-addressed "+
			"deterministically encrypted blob storage system")

	dialURL := hoarctlApp.StringOpt("a address", "tcp://localhost:54193",
		"local address on which hoard is listening encoded as a URL with the "+
			"network protocol as the scheme, for example 'tcp://localhost:54192' "+
			"or 'unix:///tmp/hoard.sock'")

	saltString := hoarctlApp.StringOpt("s salt", "", "The salt "+
		"for put or get requests used for encryption and decryption. If begins "+
		"with 0x will be interpreted as a hex string.")

	var grpcClient hoard.HoardClient
	var salt []byte

	hoarctlApp.Before = func() {
		netProtocol, localAddress, err := shared.SplitListenURL(*dialURL)

		conn, err := grpc.Dial(*dialURL,
			grpc.WithInsecure(),
			// We have to bugger around with this so we can dial an arbitrary net.Conn
			grpc.WithDialer(func(string, time.Duration) (net.Conn, error) {
				return net.Dial(netProtocol, localAddress)
			}))

		if err != nil {
			shared.Fatalf("Could not dial hoard server on %s: %v", *dialURL, err)
		}
		//defer conn.Close()
		grpcClient = hoard.NewHoardClient(conn)
		salt = parseSalt(*saltString)
	}

	hoarctlApp.Command("put",
		"Put some data into encrypted data store and return a reference",
		func(cmd *cli.Cmd) {
			cmd.Action = func() {
				data, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					shared.Fatalf("Could read bytes to put from stdin: %v", err)
				}
				ref, err := grpcClient.Put(context.Background(),
					&hoard.Plaintext{
						Data: data,
						Salt: salt,
					})
				if err != nil {
					shared.Fatalf("Error putting file: %v", err)
				}
				fmt.Printf("%s\n", jsonString(ref))
			}
		})

	hoarctlApp.Run(os.Args)
}

func parseSalt(saltString string) []byte {
	if saltString == "" {
		return nil
	}
	if strings.HasPrefix(saltString, "0x") {
		salt, err := hex.DecodeString(strings.TrimPrefix(saltString, "0x"))
		if err != nil {
			shared.Fatalf("Could not decode '%s' as hex string: %s",
				saltString, err)
		}
		return salt
	}
	return ([]byte)(saltString)
}

func jsonString(v interface{}) string {
	bs, err := json.Marshal(v)
	if err != nil {
		shared.Fatalf("Could not serialise '%s' to json: %v", err)
	}
	return string(bs)

}
