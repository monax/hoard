package storage

type FileSystemConfig struct {
	RootDirectory string
}

func NewFileSystemConfig(addressEncoding, rootDirectory string) *StorageConfig {
	return &StorageConfig{
		StorageType:     Filesystem,
		AddressEncoding: addressEncoding,
		FileSystemConfig: &FileSystemConfig{
			RootDirectory: rootDirectory,
		},
	}
}
