package storage

import (
	"fmt"
)

type CloudConfig struct {
	Bucket string
	Prefix string
	Region string
}

func NewCloudConfig(encoding, cloud, bucket, prefix, region string) (*StorageConfig, error) {
	return &StorageConfig{
		StorageType:     StorageType(cloud),
		AddressEncoding: encoding,
		CloudConfig: &CloudConfig{
			Bucket: bucket,
			Prefix: prefix,
			Region: region,
		},
	}, nil
}

func DefaultCloudConfig(cloud string) *StorageConfig {
	conf, err := NewCloudConfig(DefaultAddressEncodingName,
		cloud,
		"hoard",
		"store",
		"uk",
	)
	if err != nil {
		panic(fmt.Errorf("could not generate example config: %s", err))
	}
	return conf
}
