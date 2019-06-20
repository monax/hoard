// +build integration

package cloud

import (
	"context"
	"encoding/base32"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/monax/hoard/v5/stores"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func TestStoreGCS(t *testing.T) {
	bucket := "monax-hoard"
	prefix := "test-store"
	deleteGCSPrefix(bucket, prefix)
	store, err := NewStore(GCP, bucket, prefix, "", base32.StdEncoding, nil)
	assert.NoError(t, err)
	stores.RunTests(t, store)
}

func TestStoreS3(t *testing.T) {
	bucket := "monax-hoard-test"
	prefix := "TestS3Store/"
	deleteS3Prefix(bucket, prefix)
	store, err := NewStore(AWS, bucket, prefix, "", base32.StdEncoding, nil)
	assert.NoError(t, err)
	stores.RunTests(t, store)
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

	client, err := cloudstores.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		panic(err)
	}

	defer client.Close()
	bkt := client.Bucket(bucket)
	objs := bkt.Objects(ctx, &cloudstores.Query{Prefix: prefix})
	for obj, _ := objs.Next(); obj != nil; obj, _ = objs.Next() {
		bkt.Object(obj.Name).Delete(ctx)
	}
}
