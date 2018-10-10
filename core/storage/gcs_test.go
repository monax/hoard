// +build integration

package storage

import (
	"testing"

	"encoding/base32"

	"github.com/stretchr/testify/assert"
)

func TestGCSStore(t *testing.T) {
	bucket := "monax-hoard-test"
	prefix := "TestGCSStore/"
	gcss, err := NewGCSStore(bucket, prefix, base32.StdEncoding, nil)
	assert.NoError(t, err)
	testStore(t, gcss)
}
