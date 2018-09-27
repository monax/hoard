// +build integration

package storage

import (
	"testing"

	"encoding/base32"

	"github.com/stretchr/testify/assert"
)

func TestGCSStore(t *testing.T) {
	bucket := "monax-hoard-test"
	gcss, err := NewGCSStore(bucket, base32.StdEncoding, nil, nil)
	assert.NoError(t, err)
	testStore(t, gcss)
}
