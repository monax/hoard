package main

import (
	"fmt"
	"net"
	"os"

	"code.monax.io/platform/hoard/cmd/shared"
	"code.monax.io/platform/hoard/hoard"
	"code.monax.io/platform/hoard/hoard/storage"
	"github.com/go-kit/kit/log"
	"github.com/jawher/mow.cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	hoardApp := cli.App("hoard",
		"A content-addressed deterministically encrypted blob storage system")
	listenURL := hoardApp.StringOpt("a address", "tcp://localhost:54193",
		"local address for hoard to listen on encoded as a URL with the "+
			"network protocol as the scheme, for example 'tcp://localhost:54192' "+
			"or 'unix:///tmp/hoard.sock'")
	// This string spec is parsed by mow.cli and has actual semantic significance
	// around optionality and ordering of options and arguments
	hoardApp.Spec = "[-a]"

	hoardApp.Action = func() {
		netProtocol, localAddress, err := shared.SplitListenURL(*listenURL)
		if err != nil {
			shared.Fatalf("Failed to split listen URL '%s': %v", *listenURL, err)
		}
		lis, err := net.Listen(netProtocol, localAddress)
		if err != nil {
			shared.Fatalf("Failed to create listener: %v", err)
		}
		grpcServer := grpc.NewServer()
		hrd := hoard.NewHoard(storage.NewMemoryStore(),
			log.NewLogfmtLogger(os.Stderr))
		hoard.RegisterHoardServer(grpcServer, hoard.NewHoardServer(hrd))
		// Register reflection service on gRPC server.
		reflection.Register(grpcServer)
		fmt.Fprintf(os.Stderr, "Running hoard at %s://%s...\n",
			netProtocol, localAddress)
		if err := grpcServer.Serve(lis); err != nil {
			shared.Fatalf("Failed to start GRPC Server: %v", err)
		}
	}

	hoardApp.Run(os.Args)
}
