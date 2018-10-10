package storage

import (
	"fmt"
)

type GCSConfig struct {
	GCSBucket string
	GCSPrefix string
}

func NewGCSConfig(addressEncoding, gcsBucket, gcsPrefix string) (*StorageConfig, error) {
	return &StorageConfig{
		StorageType:     GCS,
		AddressEncoding: addressEncoding,
		GCSConfig: &GCSConfig{
			GCSBucket: gcsBucket,
			GCSPrefix: gcsPrefix,
		},
	}, nil
}

func DefaultGCSConfig() *StorageConfig {
	gcsc, err := NewGCSConfig(DefaultAddressEncodingName,
		"monax-hoard-test",
		"store",
	)
	if err != nil {
		panic(fmt.Errorf("could not generate example config: %s", err))
	}
	return gcsc
}
