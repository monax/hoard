package hoard

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/v8/api"
	"github.com/monax/hoard/v8/config"
	"github.com/monax/hoard/v8/encryption"
	"github.com/monax/hoard/v8/grant"
	"github.com/monax/hoard/v8/reference"
	"github.com/monax/hoard/v8/stores"
	"github.com/monax/hoard/v8/test/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func writeUIntBE(buffer []byte, value, offset, byteLength int64) error {
	slice := make([]byte, byteLength)

	buf := bytes.NewBuffer(slice)
	err := binary.Write(buf, binary.BigEndian, value)
	if err != nil {
		return err
	}

	slice = buf.Bytes()
	slice = slice[int64(len(slice))-byteLength : len(slice)]

	copy(buffer[offset:], slice)
	return nil
}

func TestService(t *testing.T) {
	chunkSize := int64(16)
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

				refs, err := ReceiveAllReferences(putStream.Recv)
				require.NoError(t, err)
				expected := int64(len(data))/chunkSize + 1
				if int64(len(data))%chunkSize > 0 {
					expected++
				}
				require.Equal(t, expected, int64(len(refs)))

				getStream, err := client.Get(ctx)
				require.NoError(t, err)
				for _, ref := range refs {
					err = getStream.Send(ref)
					require.NoError(t, err)
				}

				err = getStream.CloseSend()
				require.NoError(t, err)

				plaintext, err := ReceiveAllPlaintexts(getStream.Recv)
				require.NoError(t, err)
				require.Equal(t, data, plaintext.GetBody())
			})

			t.Run("LengthPrefix", func(t *testing.T) {
				const lengthPrefixByteLength = 4

				meta, err := json.Marshal(&struct {
					Name           string
					MimeType       string
					Tags           []string
					Agreement      string
					AssemblyEngine string
				}{
					Name:     "document",
					MimeType: ".docx",
				})
				require.NoError(t, err)

				buffer := make([]byte, lengthPrefixByteLength)
				err = writeUIntBE(buffer, int64(len(meta)), 0, lengthPrefixByteLength)
				require.NoError(t, err)

				data := make([]byte, 1000)
				msg := append(buffer, meta...)
				msg = append(msg, data...)

				ref, err := service.streaming.grantService.Put(msg, []byte{})
				require.NoError(t, err)

				client := api.NewCleartextClient(conn)
				getStream, err := client.Get(ctx)
				require.NoError(t, err)
				err = getStream.Send(ref)
				require.NoError(t, err)
				err = getStream.CloseSend()
				require.NoError(t, err)

				plaintext, err := ReceiveAllPlaintexts(getStream.Recv)
				require.NoError(t, err)
				require.Nil(t, plaintext.GetHead())
				body := plaintext.GetBody()
				size := binary.BigEndian.Uint32(body[:lengthPrefixByteLength])
				head := body[lengthPrefixByteLength : size+lengthPrefixByteLength]
				rest := body[size+lengthPrefixByteLength:]

				require.Equal(t, meta, head)
				require.Equal(t, data, rest)
			})

			t.Run("ChunkLarge", func(t *testing.T) {
				size := 124 * 124 * 10
				client := api.NewCleartextClient(conn)
				putStream, err := client.Put(ctx)
				require.NoError(t, err)

				bigBytes := make([]byte, size)
				bigBytes[333] = 23
				input := bytes.NewBuffer(bigBytes)

				var refs []*reference.Ref
				err = NewStreamer().
					WithChunkSize(512).
					WithInput(input).
					WithSend(
						func(chunk []byte) error {
							err := putStream.Send(&api.Plaintext{
								Body: chunk,
							})
							return err
						}).
					WithCloseSend(putStream.CloseSend).
					WithRecv(func() ([]byte, error) {
						ref, err := putStream.Recv()
						if err != nil {
							return nil, err
						}
						refs = append(refs, ref)
						// Hi, I'm Go, I don't have generics.
						return nil, nil
					}).
					Stream(context.Background())

				getStream, err := client.Get(ctx)
				require.NoError(t, err)

				output := new(bytes.Buffer)

				err = NewStreamer().
					WithSend(func(chunk []byte) error {
						if len(refs) == 0 {
							return io.EOF
						}
						ref := refs[0]
						refs = refs[1:]
						return getStream.Send(ref)
					}).
					WithRecv(func() ([]byte, error) {
						pt, err := getStream.Recv()
						if err != nil {
							return nil, err
						}
						return pt.Body, nil
					}).
					WithCloseSend(getStream.CloseSend).
					WithOutput(output).
					Stream(context.Background())

				require.NoError(t, err)

				require.True(t, bytes.Equal(bigBytes, output.Bytes()))
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
				plaintext, err := ReceiveAllPlaintexts(getStream.Recv)
				require.NoError(t, err)
				require.Equal(t, data, plaintext.GetBody())
			})
		})

		return nil
	})
	require.NoError(t, err)
}

func ReceiveAllPlaintexts(recv func() (*api.Plaintext, error)) (*api.Plaintext, error) {
	plaintext := new(api.Plaintext)

	for {
		pt, err := recv()
		if err != nil {
			if err == io.EOF {
				return plaintext, nil
			}

			return nil, err
		}

		plaintext.Body = append(plaintext.Body, pt.GetBody()...)
		if plaintext.Head == nil {
			plaintext.Head = pt.GetHead()
		}
	}
}

func ReceiveAllReferences(recv func() (*reference.Ref, error)) ([]*reference.Ref, error) {
	refs := make([]*reference.Ref, 0)

	for {
		ref, err := recv()
		if err != nil {
			if err == io.EOF {
				return refs, nil
			}

			return nil, err
		}

		refs = append(refs, ref)
	}
}
