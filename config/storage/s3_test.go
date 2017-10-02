package storage

import (
	"testing"

	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/stretchr/testify/assert"
)

func TestDefaultS3Config(t *testing.T) {
	assertStorageConfigSerialisation(t, DefaultS3Config())
}
