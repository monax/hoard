package storage

import (
	"testing"

	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/stretchr/testify/assert"
)

func TestAWSCredentialsFromConfig(t *testing.T) {
	chain := []*CredentialsProviderConfig{
		{
			Provider: EnvProviderName,
		},
	}

	value := credentials.Value{
		AccessKeyID:     "access-key-1",
		SecretAccessKey: "secret-access-key-1",
		ProviderName:    credentials.EnvProviderName,
	}

	err := os.Setenv("AWS_ACCESS_KEY", value.AccessKeyID)
	assert.NoError(t, err)
	err = os.Setenv("AWS_SECRET_KEY", value.SecretAccessKey)
	assert.NoError(t, err)
	creds, err := AWSCredentialsFromChain(chain)
	assert.NoError(t, err)
	valueOut, err := creds.Get()
	assert.NoError(t, err)
	assert.Equal(t, value, valueOut)
}

func TestDefaultS3Config(t *testing.T) {
	assertStorageConfigSerialisation(t, DefaultS3Config())
}
