package storage

import (
	"fmt"
)

type GCSConfig struct {
	GCSBucket string
}

func NewGCSConfig(addressEncoding, gcsBucket string) (*StorageConfig, error) {
	return &StorageConfig{
		StorageType:     GCS,
		AddressEncoding: addressEncoding,
		GCSConfig: &GCSConfig{
			GCSBucket: gcsBucket,
		},
	}, nil
}

func DefaultGCSConfig() *StorageConfig {
	gcsc, err := NewGCSConfig(DefaultAddressEncodingName,
		"monax-hoard-test",
	)
	if err != nil {
		panic(fmt.Errorf("could not generate example config: %s", err))
	}
	return gcsc
}
