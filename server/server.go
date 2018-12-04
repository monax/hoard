package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/monax/hoard"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/logging"
	"github.com/monax/hoard/logging/loggers"
	"github.com/monax/hoard/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	listenURL  string
	store      storage.NamedStore
	grpcServer *grpc.Server
	logger     log.Logger
}

func New(listenURL string, store storage.NamedStore, logger log.Logger) *server {
	return &server{
		listenURL: listenURL,
		store:     store,
		logger:    logger,
	}
}

func (serv *server) Serve() error {
	netProtocol, localAddress, err := SplitListenURL(serv.listenURL)
	if err != nil {
		return fmt.Errorf("Failed to split listen URL '%s': %v", serv.listenURL, err)
	}
	listener, err := net.Listen(netProtocol, localAddress)
	if err != nil {
		return fmt.Errorf("Failed to create listener: %v", err)
	}
	serv.grpcServer = grpc.NewServer()
	if serv.logger == nil {
		serv.logger = log.NewNopLogger()
	} else {
		serv.logger = loggers.Compose(logging.WithMetadata, loggers.NonBlockingLogger,
			loggers.VectorValuedLogger)(serv.logger)
	}

	logging.InfoMsg(serv.logger, "Initialising Hoard server",
		"store_name", serv.store.Name())
	hoardServer := hoard.NewHoardServer(hoard.NewHoard(serv.store, serv.logger))

	hoard.RegisterCleartextServer(serv.grpcServer, hoardServer)
	hoard.RegisterEncryptionServer(serv.grpcServer, hoardServer)
	hoard.RegisterStorageServer(serv.grpcServer, hoardServer)
	// Register reflection service on gRPC server.
	reflection.Register(serv.grpcServer)
	err = serv.grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("Failed to start GRPC Server: %v", err)
	}
	return nil
}

func (serv *server) Stop() {
	serv.grpcServer.Stop()
}

func SplitListenURL(listenOn string) (string, string, error) {
	// net.Listen does not want a parsed url.URL so it seems to make more sense
	// just to do a dumb split here to support the various networks
	listenParts := strings.Split(listenOn, "://")
	if len(listenParts) != 2 {
		return "", "", fmt.Errorf("Expected a Go net.Listen URL of the form "+
			"'<net>://<laddr>', but got: '%s'", listenOn)
	}
	if listenParts[0] == "" {
		return "", "", fmt.Errorf("Expected the URL scheme to be present, "+
			"but got '%s'", listenOn)
	}
	if listenParts[1] == "" {
		return "", "", fmt.Errorf("Expected the URL host to be present, "+
			"but got '%s'", listenOn)
	}
	return listenParts[0], listenParts[1], nil
}
