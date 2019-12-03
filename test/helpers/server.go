package helpers

import (
	"context"
	"net"

	"github.com/monax/hoard/v7/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type service interface {
	api.CleartextServer
	api.EncryptionServer
	api.StorageServer
	api.GrantServer
	api.DocumentServer
}

// Provided with a HoardService executes runner in the context of a client-server connection over test buffer connection
// with all the hoard GRPC services registered on service
func RunWithTestServer(ctx context.Context, service service,
	runner func(server *grpc.Server, conn *grpc.ClientConn) error) error {

	grpcServer := grpc.NewServer()

	api.RegisterCleartextServer(grpcServer, service)
	api.RegisterEncryptionServer(grpcServer, service)
	api.RegisterStorageServer(grpcServer, service)
	api.RegisterGrantServer(grpcServer, service)
	api.RegisterDocumentServer(grpcServer, service)

	const bufferSize = 1 << 20
	l := bufconn.Listen(bufferSize)
	errCh := make(chan error)
	go func() {
		errCh <- grpcServer.Serve(l)
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return l.Dial() }),
		grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	// Do stuff while server runs
	err = runner(grpcServer, conn)
	if err != nil {
		return err
	}

	err = l.Close()
	if err != nil {
		return err
	}

	err = <-errCh
	if err != nil && err.Error() != "closed" {
		return err
	}
	return nil
}
