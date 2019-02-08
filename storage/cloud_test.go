// +build integration

package storage

import (
	"encoding/base32"
	"os"
	"testing"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/google"
)

func TestStoreGCS(t *testing.T) {
	bucket := "monax-hoard"
	prefix := "test-store"
	deleteGCSPrefix(bucket, prefix)
	store, err := NewCloudStore(GCP, bucket, prefix, "", base32.StdEncoding, nil)
	assert.NoError(t, err)
	testStore(t, store)
}

func TestStoreS3(t *testing.T) {
	bucket := "monax-hoard-test"
	prefix := "TestS3Store/"
	deleteS3Prefix(bucket, prefix)
	store, err := NewCloudStore(AWS, bucket, prefix, "", base32.StdEncoding, nil)
	assert.NoError(t, err)
	testStore(t, store)
}

func deleteS3Prefix(bucket, prefix string) {
	deleter := s3manager.NewBatchDelete(session.Must(session.New(aws.NewConfig()), nil))
	err := deleter.Delete(context.Background(),
		s3manager.NewDeleteListIterator(deleter.Client,
			&s3.ListObjectsInput{Bucket: &bucket, Prefix: &prefix}))
	if err != nil {
		panic(err)
	}
}

func deleteGCSPrefix(bucket, prefix string) {
	ctx := context.Background()
	creds, err := google.CredentialsFromJSON(ctx, []byte(os.Getenv("GCLOUD_SERVICE_KEY")), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		panic(err)
	}

	client, err := storage.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		panic(err)
	}

	defer client.Close()
	bkt := client.Bucket(bucket)
	objs := bkt.Objects(ctx, &storage.Query{Prefix: prefix})
	for obj, _ := objs.Next(); obj != nil; obj, _ = objs.Next() {
		bkt.Object(obj.Name).Delete(ctx)
	}
}
