package main

import (
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/v6/config"
	"github.com/monax/hoard/v6/stores"
	"github.com/monax/hoard/v6/stores/cloud"
	"github.com/monax/hoard/v6/stores/ipfs"
)

func StoreFromStorageConfig(storageConfig *config.Storage, logger log.Logger) (stores.NamedStore, error) {
	addressEncoding, err := stores.GetAddressEncoding(storageConfig.AddressEncoding)
	if err != nil {
		return nil, err
	}

	switch storageConfig.StorageType {
	case config.Memory, config.Unspecified:
		return stores.NewMemoryStore(), nil

	case config.Filesystem:
		fsConf := storageConfig.FileSystemConfig
		if fsConf == nil {
			return nil, errors.New("filesystem storage configuration must be " +
				"supplied to use the filesystem storage backend")
		}
		if fsConf.RootDirectory == "" {
			return nil, errors.New("rootDirectory key must be non-empty in " + "filesystem storage")

		}
		return stores.NewFileSystemStore(fsConf.RootDirectory, addressEncoding)

	case config.IPFS:
		ipfsConf := storageConfig.IPFSConfig
		if ipfsConf == nil {
			return nil, errors.New("IPFS storage configuration must be " +
				"supplied to use the filesystem storage backend")
		}
		if ipfsConf.RemoteAPI == "" {
			return nil, errors.New("http api url must be non-empty in " +
				"ipfs storage config")
		}
		return ipfs.NewStore(ipfsConf.RemoteAPI, addressEncoding)

	case config.AWS:
		awsConf := storageConfig.Cloud
		if awsConf == nil {
			return nil, errors.New("aws configuration must be supplied")
		}
		return cloud.NewStore(cloud.AWS, awsConf.Bucket, awsConf.Prefix, awsConf.Region, addressEncoding, logger)

	case config.Azure:
		azureConf := storageConfig.Cloud
		if azureConf == nil {
			return nil, errors.New("azure configuration must be supplied")
		}
		return cloud.NewStore(cloud.Azure, azureConf.Bucket, azureConf.Prefix, azureConf.Region, addressEncoding, logger)

	case config.GCP:
		gcpConf := storageConfig.Cloud
		if gcpConf == nil {
			return nil, errors.New("gcp configuration must be supplied")
		}
		return cloud.NewStore(cloud.GCP, gcpConf.Bucket, gcpConf.Prefix, gcpConf.Region, addressEncoding, logger)

	default:
		return nil, fmt.Errorf("did not recognise storage type '%s'",
			storageConfig.StorageType)
	}
}
