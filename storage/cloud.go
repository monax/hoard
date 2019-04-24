package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/v3/logging"
	"github.com/monax/hoard/v3/logging/structure"
	"gocloud.dev/blob"
	"gocloud.dev/blob/azureblob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/blob/s3blob"
	"gocloud.dev/gcp"
	"golang.org/x/oauth2/google"
)

type CloudType string

const (
	AWS   CloudType = "aws"
	Azure CloudType = "azure"
	GCP   CloudType = "gcp"
)

type cloudStore struct {
	back     context.Context
	blob     *blob.Bucket
	bucket   string
	prefix   string
	encoding AddressEncoding
	logger   log.Logger
}

func NewCloudStore(cloud CloudType, bucket, prefix, region string, addrenc AddressEncoding, logger log.Logger) (*cloudStore, error) {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	var conn *blob.Bucket
	var err error
	ctx := context.Background()

	switch cloud {
	case AWS:
		// AWS_ACCESS_KEY_ID
		// AWS_SECRET_ACCESS_KEY
		awsConf := &aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewEnvCredentials(),
		}
		client := session.Must(session.NewSession(awsConf))
		conn, err = s3blob.OpenBucket(ctx, client, bucket, nil)

	case Azure:
		// AZURE_STORAGE_ACCOUNT_NAME
		// AZURE_STORAGE_ACCOUNT_KEY
		accountName, err := azureblob.DefaultAccountName()
		if err != nil {
			return nil, err
		}
		accountKey, err := azureblob.DefaultAccountKey()
		if err != nil {
			return nil, err
		}
		credential, err := azureblob.NewCredential(accountName, accountKey)
		if err != nil {
			return nil, err
		}
		p := azureblob.NewPipeline(credential, azblob.PipelineOptions{})
		conn, err = azureblob.OpenBucket(ctx, p, accountName, bucket, nil)

	case GCP:
		creds, err := google.CredentialsFromJSON(ctx, []byte(os.Getenv("GCLOUD_SERVICE_KEY")), "https://www.googleapis.com/auth/cloud-platform")
		if err != nil {
			return nil, err
		}
		client, err := gcp.NewHTTPClient(gcp.DefaultTransport(), gcp.CredentialsTokenSource(creds))
		if err != nil {
			return nil, err
		}
		conn, err = gcsblob.OpenBucket(ctx, client, bucket, nil)
	}

	if err != nil {
		return nil, err
	}

	prefix = strings.TrimRight(prefix, "/")
	inv := &cloudStore{
		back:     ctx,
		blob:     conn,
		bucket:   bucket,
		prefix:   prefix,
		encoding: addrenc,
		logger: logging.TraceLogger(log.With(logger,
			structure.ComponentKey, "storage")),
	}
	inv.logger = log.With(inv.logger, "store_name", inv.Name())
	return inv, nil
}

func (inv *cloudStore) Put(address, data []byte) ([]byte, error) {
	writer, err := inv.blob.NewWriter(inv.back, fmt.Sprintf("%s/%s", inv.prefix, inv.encode(address)), nil)
	if err != nil {
		return nil, err
	}

	n, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	inv.logger.Log("method", "Put",
		"location", inv.Location,
		"encoded_address", inv.encode(address),
		"uploaded_bytes", n)

	return address, nil
}

func (inv *cloudStore) Get(address []byte) ([]byte, error) {
	reader, err := inv.blob.NewReader(inv.back, fmt.Sprintf("%s/%s", inv.prefix, inv.encode(address)), nil)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	io.Copy(&buf, reader)

	inv.logger.Log("method", "Get",
		"encoded_address", inv.encode(address),
		"downloaded_bytes", reader.Size())

	err = reader.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (inv *cloudStore) Stat(address []byte) (*StatInfo, error) {
	reader, err := inv.blob.NewReader(inv.back, fmt.Sprintf("%s/%s", inv.prefix, inv.encode(address)), nil)
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

	inv.logger.Log("method", "Stat",
		"encoded_address", inv.encode(address))
	return &StatInfo{
		Exists: true,
		Size:   uint64(n),
	}, nil
}

func (inv *cloudStore) Location(address []byte) string {
	return fmt.Sprintf("gs://%s/%s", inv.bucket,
		inv.encode(address))
}

func (inv *cloudStore) Name() string {
	return fmt.Sprintf("gcpStore[bucket=%s]", inv.bucket)
}

func (inv *cloudStore) encode(address []byte) string {
	return inv.encoding.EncodeToString(address)
}
