package storage

type IPFSConfig struct {
	Protocol string
	Address  string
	Port     string
}

func NewIPFSConfig(addressEncoding, proto, address, port string) *StorageConfig {
	return &StorageConfig{
		StorageType:     IPFS,
		AddressEncoding: addressEncoding,
		IPFSConfig: &IPFSConfig{
			Protocol: proto,
			Address:  address,
			Port:     port,
		},
	}
}

func DefaultIPFSConfig() *StorageConfig {
	return NewIPFSConfig(DefaultAddressEncodingName,
		"https://",
		"127.0.0.1",
		"5001",
	)
}
