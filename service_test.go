package hoard

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/v7/api"
	"github.com/monax/hoard/v7/config"
	"github.com/monax/hoard/v7/grant"
	"github.com/monax/hoard/v7/stores"
	"github.com/monax/hoard/v7/test/helpers"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestService(t *testing.T) {
	chunkSize := 67
	secrets := config.SecretsManager{
		Provider: func(secretID string) ([]byte, error) {
			return []byte(secretID + "shhhh"), nil
		},
	}
	hrd := NewHoard(stores.NewMemoryStore(), secrets, log.NewNopLogger())
	service := NewService(hrd, chunkSize)
	ctx := context.Background()
	err := helpers.RunWithTestServer(ctx, service, func(server *grpc.Server, conn *grpc.ClientConn) error {

		t.Run("Grants", func(t *testing.T) {
			t.Run("Streaming", func(t *testing.T) {
				data := []byte(helpers.LongText)
				client := api.NewGrantClient(conn)
				putStream, err := client.PutSeal(ctx)
				require.NoError(t, err)

				publicID := "code1"
				pgs := &api.PlaintextAndGrantSpec{
					Plaintext: &api.Plaintext{
						Data: data,
						Salt: []byte("celery"),
					},
					GrantSpec: &grant.Spec{
						Symmetric: &grant.SymmetricSpec{
							PublicID: publicID,
						},
					},
				}
				err = SendPlaintextAndGrantSpec(putStream, pgs, chunkSize)
				require.NoError(t, err)
				grt, err := putStream.CloseAndRecv()
				require.NoError(t, err)

				getStream, err := client.UnsealGet(ctx, grt)
				require.NoError(t, err)
				plaintextOut, err := ReceivePlaintext(getStream)
				require.NoError(t, err)
				require.Equal(t, pgs.Plaintext, plaintextOut)
			})
		})

		return nil
	})
	require.NoError(t, err)
}
