package config

type IPFSConfig struct {
	RemoteAPI string
}

func NewIPFSConfig(addressEncoding, host string) *Storage {
	return &Storage{
		StorageType:     IPFS,
		AddressEncoding: addressEncoding,
		IPFSConfig: &IPFSConfig{
			RemoteAPI: host,
		},
	}
}

func DefaultIPFSConfig() *Storage {
	return NewIPFSConfig(DefaultAddressEncodingName,
		"http://:5001",
	)
}
