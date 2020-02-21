// +build integration

package cloud

import (
	"context"
	"encoding/base32"
	"fmt"
	"os"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/monax/hoard/v7/stores"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func TestStoreGCS(t *testing.T) {
	bucket := "monax-hoard"
	prefix := "test-store"
	err := deleteGCSPrefix(bucket, prefix)
	require.NoError(t, err)
	store, err := NewStore(GCP, bucket, prefix, "", base32.StdEncoding, nil)
	assert.NoError(t, err)
	stores.RunTests(t, store)
}

func deleteGCSPrefix(bucket, prefix string) error {
	ctx := context.Background()
	serviceKey := os.Getenv(GcloudServiceKeyEnvVar)
	if len(serviceKey) == 0 {
		return fmt.Errorf("environment variable %s not set", GcloudServiceKeyEnvVar)
	}
	creds, err := google.CredentialsFromJSON(ctx, []byte(serviceKey), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return fmt.Errorf("could not parse gcloud credentials: %v", err)
	}

	client, err := storage.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return err
	}

	defer client.Close()
	bkt := client.Bucket(bucket)
	objs := bkt.Objects(ctx, &storage.Query{Prefix: prefix})
	for obj, _ := objs.Next(); obj != nil; obj, _ = objs.Next() {
		bkt.Object(obj.Name).Delete(ctx)
	}
	return nil
}
