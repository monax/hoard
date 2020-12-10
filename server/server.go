package server

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/v8"
	"github.com/monax/hoard/v8/api"
	"github.com/monax/hoard/v8/config"
	"github.com/monax/hoard/v8/logging"
	"github.com/monax/hoard/v8/logging/loggers"
	"github.com/monax/hoard/v8/stores"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	listenURL  string
	listener   net.Listener
	hoard      *hoard.Hoard
	chunk      int64
	grpcServer *grpc.Server
	ready      chan struct{}
	logger     log.Logger
}

func New(listenURL string, store stores.NamedStore, secretManager config.SecretsManager, chunkSize int64, logger log.Logger) *Server {
	return &Server{
		listenURL: listenURL,
		hoard:     hoard.NewHoard(store, secretManager, logger),
		chunk:     chunkSize,
		ready:     make(chan struct{}),
		logger:    logger,
	}
}

func (serv *Server) Serve() error {
	netProtocol, localAddress, err := SplitListenURL(serv.listenURL)
	if err != nil {
		return fmt.Errorf("failed to split listen URL '%s': %v", serv.listenURL, err)
	}
	serv.listener, err = net.Listen(netProtocol, localAddress)
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}
	serv.grpcServer = grpc.NewServer()
	if serv.logger == nil {
		serv.logger = log.NewNopLogger()
	} else {
		serv.logger = loggers.Compose(logging.WithMetadata, loggers.NonBlockingLogger,
			loggers.VectorValuedLogger)(serv.logger)
	}

	logging.InfoMsg(serv.logger, "Initialising Hoard server",
		"store_name", serv.hoard.Name())

	hoardService := hoard.NewService(serv.hoard, serv.chunk)
	api.RegisterCleartextServer(serv.grpcServer, hoardService)
	api.RegisterEncryptionServer(serv.grpcServer, hoardService)
	api.RegisterStorageServer(serv.grpcServer, hoardService)
	api.RegisterGrantServer(serv.grpcServer, hoardService)

	// Register reflection service on gRPC server.
	reflection.Register(serv.grpcServer)
	// Announce ready
	close(serv.ready)
	err = serv.grpcServer.Serve(serv.listener)
	if err != nil {
		return fmt.Errorf("failed to start GRPC Server: %v", err)
	}
	return nil
}

func (serv *Server) ListenAddress() net.Addr {
	return serv.listener.Addr()
}

// Wait until server is listening or context is done
func (serv *Server) Wait(ctx context.Context) error {
	select {
	case <-serv.ready:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (serv *Server) Stop() {
	serv.grpcServer.Stop()
}

func SplitListenURL(listenOn string) (string, string, error) {
	// net.Listen does not want a parsed url.URL so it seems to make more sense
	// just to do a dumb split here to support the various networks
	listenParts := strings.Split(listenOn, "://")
	if len(listenParts) != 2 {
		return "", "", fmt.Errorf("expected a Go net.Listen URL of the form "+
			"'<net>://<laddr>', but got: '%s'", listenOn)
	}
	if listenParts[0] == "" {
		return "", "", fmt.Errorf("expected the URL scheme to be present, "+
			"but got '%s'", listenOn)
	}
	if listenParts[1] == "" {
		return "", "", fmt.Errorf("expected the URL host to be present, "+
			"but got '%s'", listenOn)
	}
	return listenParts[0], listenParts[1], nil
}
