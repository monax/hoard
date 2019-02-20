package storage

type IPFSConfig struct {
	RemoteAPI string
}

func NewIPFSConfig(addressEncoding, host string) *StorageConfig {
	return &StorageConfig{
		StorageType:     IPFS,
		AddressEncoding: addressEncoding,
		IPFSConfig: &IPFSConfig{
			RemoteAPI: host,
		},
	}
}

func DefaultIPFSConfig() *StorageConfig {
	return NewIPFSConfig(DefaultAddressEncodingName,
		"http://:5001",
	)
}
