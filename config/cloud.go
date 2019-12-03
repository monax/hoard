package config

type Cloud struct {
	Bucket string
	Prefix string
	Region string
}

func NewDefaultCloud(cloud StorageType) *Storage {
	conf := NewDefaultStorage()
	conf.StorageType = cloud
	conf.Cloud = &Cloud{
		Bucket: "hoard",
		Prefix: "store",
		Region: "uk",
	}
	return conf
}
