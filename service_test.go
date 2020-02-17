package hoard

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/v7/api"
	"github.com/monax/hoard/v7/config"
	"github.com/monax/hoard/v7/encryption"
	"github.com/monax/hoard/v7/grant"
	"github.com/monax/hoard/v7/stores"
	"github.com/monax/hoard/v7/test/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestService(t *testing.T) {
	chunkSize := 67
	salt, err := encryption.NewNonce(encryption.NonceSize)
	assert.NoError(t, err)
	secret, err := encryption.DeriveSecretKey([]byte("shhhh"), salt)
	assert.NoError(t, err)
	secrets := config.SecretsManager{
		Provider: func(secretID string) (config.SymmetricSecret, error) {
			return config.SymmetricSecret{SecretKey: secret}, nil
		},
	}
	hrd := NewHoard(stores.NewMemoryStore(), secrets, log.NewNopLogger())
	service := NewService(hrd, chunkSize)
	ctx := context.Background()
	err = helpers.RunWithTestServer(ctx, service, func(server *grpc.Server, conn *grpc.ClientConn) error {

		t.Run("Cleartext", func(t *testing.T) {
			t.Run("Streaming", func(t *testing.T) {
				data := make([]byte, 1000)
				salt := []byte("celery")

				client := api.NewCleartextClient(conn)
				putStream, err := client.Put(ctx)
				require.NoError(t, err)
				err = putStream.Send(&api.Plaintext{Head: &api.Header{Salt: salt}})
				require.NoError(t, err)
				for _, b := range data {
					err = putStream.Send(&api.Plaintext{Body: []byte{b}})
					require.NoError(t, err)
				}
				err = putStream.CloseSend()
				require.NoError(t, err)

				refs, err := ReceiveAllReferences(putStream)
				require.NoError(t, err)
				expected := len(data)/chunkSize + 1
				if len(data)%chunkSize > 0 {
					expected++
				}
				require.Equal(t, expected, len(refs))

				getStream, err := client.Get(ctx)
				require.NoError(t, err)
				for _, ref := range refs {
					err = getStream.Send(ref)
					require.NoError(t, err)
				}

				err = getStream.CloseSend()
				require.NoError(t, err)

				plaintext, err := ReceiveAllPlaintexts(getStream)
				require.NoError(t, err)
				require.Equal(t, data, plaintext.GetBody())
			})
		})

		t.Run("Grants", func(t *testing.T) {
			t.Run("Streaming", func(t *testing.T) {
				data := []byte(helpers.LongText)
				salt := []byte("celery")
				publicID := "code1"
				gs := &grant.Spec{
					Symmetric: &grant.SymmetricSpec{
						PublicID: publicID,
					},
				}

				client := api.NewGrantClient(conn)
				putStream, err := client.PutSeal(ctx)
				require.NoError(t, err)

				err = putStream.Send(&api.PlaintextAndGrantSpec{
					Plaintext: &api.Plaintext{Head: &api.Header{Salt: salt}},
					GrantSpec: gs,
				})
				require.NoError(t, err)
				err = putStream.Send(&api.PlaintextAndGrantSpec{Plaintext: &api.Plaintext{Body: data}})
				require.NoError(t, err)

				grt, err := putStream.CloseAndRecv()
				require.NoError(t, err)

				getStream, err := client.UnsealGet(ctx, grt)
				require.NoError(t, err)
				plaintext, err := ReceiveAllPlaintexts(getStream)
				require.NoError(t, err)
				require.Equal(t, data, plaintext.GetBody())
			})
		})

		return nil
	})
	require.NoError(t, err)
}
