package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/gcsblob"
	"github.com/google/go-cloud/gcp"
	"github.com/monax/hoard/core/logging"
	"github.com/monax/hoard/core/logging/structure"
)

type gcsStore struct {
	back            context.Context
	gcpGCS          *blob.Bucket
	gcsBucket       string
	addressEncoding AddressEncoding
	logger          log.Logger
}

func NewGCSStore(gcsBucket string, addressEncoding AddressEncoding,
	logger log.Logger) (*gcsStore, error) {

	if logger == nil {
		logger = log.NewNopLogger()
	}

	ctx := context.Background()
	// obtain default GCP credentials from Cloud Platform scope
	creds, err := gcp.DefaultCredentials(ctx)
	gcsClient, err := gcp.NewHTTPClient(gcp.DefaultTransport(), gcp.CredentialsTokenSource(creds))
	gcpSession, err := gcsblob.OpenBucket(ctx, gcsBucket, gcsClient)
	if err != nil {
		return nil, err
	}
	gcss := &gcsStore{
		back:            ctx,
		gcpGCS:          gcpSession,
		gcsBucket:       gcsBucket,
		addressEncoding: addressEncoding,
		logger: logging.TraceLogger(log.With(logger,
			structure.ComponentKey, "storage")),
	}
	gcss.logger = log.With(gcss.logger, "store_name", gcss.Name())
	return gcss, nil
}

func (gcss *gcsStore) Put(address, data []byte) ([]byte, error) {
	writer, err := gcss.gcpGCS.NewWriter(gcss.back, gcss.encode(address), nil)
	if err != nil {
		return address, err
	}

	n, err := writer.Write(data)
	if err != nil {
		return address, err
	}

	if err = writer.Close(); err != nil {
		return address, err
	}

	gcss.logger.Log("method", "Put",
		"location", gcss.Location,
		"encoded_address", gcss.encode(address),
		"uploaded_bytes", n)

	return address, err
}

func (gcss *gcsStore) Get(address []byte) ([]byte, error) {
	reader, err := gcss.gcpGCS.NewReader(gcss.back, gcss.encode(address))
	if err != nil {
		return address, err
	}

	var buf bytes.Buffer
	io.Copy(&buf, reader)

	gcss.logger.Log("method", "Get",
		"encoded_address", gcss.encode(address),
		"downloaded_bytes", reader.Size())

	err = reader.Close()
	if err != nil {
		return address, err
	}
	return buf.Bytes(), nil
}

func (gcss *gcsStore) Stat(address []byte) (*StatInfo, error) {
	reader, err := gcss.gcpGCS.NewReader(gcss.back, gcss.encode(address))
	if err != nil {
		return &StatInfo{
			Exists: false,
		}, nil
	}

	n := reader.Size()
	err = reader.Close()
	if err != nil {
		return nil, err
	}

	gcss.logger.Log("method", "Stat",
		"encoded_address", gcss.encode(address))
	return &StatInfo{
		Exists: true,
		Size:   uint64(n),
	}, nil
}

func (gcss *gcsStore) Location(address []byte) string {
	return fmt.Sprintf("gs://%s/%s", gcss.gcsBucket,
		gcss.Key(address))
}

func (gcss *gcsStore) Key(address []byte) string {
	return fmt.Sprintf("%s", gcss.encode(address))
}

func (gcss *gcsStore) Name() string {
	return fmt.Sprintf("gcsStore[bucket=%s]", gcss.gcsBucket)
}

func (gcss *gcsStore) encode(address []byte) string {
	return gcss.addressEncoding.EncodeToString(address)
}
