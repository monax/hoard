package config

type IPFSConfig struct {
	RemoteAPI string
}

func NewIPFSConfig(addressEncoding, host string) *Storage {
	conf := NewDefaultStorage()
	conf.StorageType = IPFS
	conf.AddressEncoding = addressEncoding
	conf.IPFSConfig = &IPFSConfig{
		RemoteAPI: host,
	}
	return conf
}

func NewDefaultIPFSConfig() *Storage {
	return NewIPFSConfig(DefaultAddressEncodingName,
		"http://:5001",
	)
}
