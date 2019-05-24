package config

import (
	"fmt"
)

type Cloud struct {
	Bucket string
	Prefix string
	Region string
}

func NewCloud(encoding, cloud, bucket, prefix, region string) (*Storage, error) {
	return &Storage{
		StorageType:     StorageType(cloud),
		AddressEncoding: encoding,
		Cloud: &Cloud{
			Bucket: bucket,
			Prefix: prefix,
			Region: region,
		},
	}, nil
}

func DefaultCloud(cloud string) *Storage {
	conf, err := NewCloud(DefaultAddressEncodingName,
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
