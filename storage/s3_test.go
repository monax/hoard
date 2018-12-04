// +build integration

package storage

import (
	"testing"

	"context"

	"encoding/base32"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
)

func TestS3Store(t *testing.T) {
	bucket := "monax-hoard-test"
	prefix := "TestS3Store/"
	deletePrefix(bucket, prefix)
	s3s, err := NewS3Store(bucket, prefix, base32.StdEncoding, nil, nil)
	assert.NoError(t, err)
	testStore(t, s3s)
}

func deletePrefix(bucket, prefix string) {
	deleter := s3manager.NewBatchDelete(session.Must(Session(aws.NewConfig())))
	err := deleter.Delete(context.Background(),
		s3manager.NewDeleteListIterator(deleter.Client,
			&s3.ListObjectsInput{Bucket: &bucket, Prefix: &prefix}))
	if err != nil {
		panic(err)
	}
}
