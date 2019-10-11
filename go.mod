module github.com/monax/hoard/v5

go 1.13

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

require (
	cloud.google.com/go v0.37.0
	github.com/Azure/azure-storage-blob-go v0.6.0
	github.com/BurntSushi/toml v0.3.1
	github.com/OneOfOne/xxhash v1.2.5
	github.com/aws/aws-sdk-go v1.19.35
	github.com/cep21/xdgbasedir v0.0.0-20170329171747-21470bfc93b9
	github.com/eapache/channels v1.1.0
	github.com/go-kit/kit v0.9.0
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/go-stack/stack v1.8.0
	github.com/gogo/protobuf v1.3.0
	github.com/golang/protobuf v1.3.2
	github.com/h2non/filetype v1.0.10
	github.com/jawher/mow.cli v1.1.0
	github.com/monax/hoard v3.0.1+incompatible
	github.com/monax/relic v2.0.0+incompatible
	github.com/stretchr/testify v1.4.0
	gocloud.dev v0.13.0
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc
	golang.org/x/net v0.0.0-20191009170851-d66e71096ffb // indirect
	golang.org/x/oauth2 v0.0.0-20190517181255-950ef44c6e07
	golang.org/x/sys v0.0.0-20191009170203-06d7bd2c5f4f // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/api v0.2.0
	google.golang.org/genproto v0.0.0-20191009194640-548a555dbc03 // indirect
	google.golang.org/grpc v1.24.0
	gopkg.in/yaml.v2 v2.2.4
)
