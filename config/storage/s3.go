package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

type S3Config struct {
	Bucket                   string
	Prefix                   string
	Region                   string
	CredentialsProviderChain []*CredentialsProviderConfig
}

type CredentialsProviderConfig struct {
	Provider string
	*credentials.EnvProvider
	*credentials.SharedCredentialsProvider
	*StaticProviderConfig
}

// Almost the same a credentials.Value
type StaticProviderConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

func (spc *StaticProviderConfig) Provider() credentials.Provider {
	return &credentials.StaticProvider{
		Value: credentials.Value{
			AccessKeyID:     spc.AccessKeyID,
			SecretAccessKey: spc.SecretAccessKey,
			SessionToken:    spc.SessionToken,
		},
	}
}

func NewS3Config(addressEncoding, s3Bucket, s3Prefix, region string,
	providers ...credentials.Provider) (*StorageConfig, error) {

	cpConfigs := make([]*CredentialsProviderConfig, 0, len(providers))

	for _, provider := range providers {
		cpConfig, err := ProviderConfig(provider)
		if err != nil {
			return nil, err
		}
		cpConfigs = append(cpConfigs, cpConfig)
	}
	return &StorageConfig{
		StorageType:     S3,
		AddressEncoding: addressEncoding,
		S3Config: &S3Config{
			Bucket: s3Bucket,
			Prefix: s3Prefix,
			Region: region,
			CredentialsProviderChain: cpConfigs,
		},
	}, nil
}

func AWSCredentialsFromConfig(cpcs []*CredentialsProviderConfig) (*credentials.Credentials, error) {
	providers := make([]credentials.Provider, 0, len(cpcs))
	for _, cpc := range cpcs {
		switch cpc.Provider {
		case credentials.EnvProviderName:
			providers = append(providers, cpc.EnvProvider)
		case credentials.SharedCredsProviderName:
			providers = append(providers, cpc.SharedCredentialsProvider)
		case credentials.StaticProviderName:
			providers = append(providers, cpc.StaticProviderConfig.Provider())
		default:
			return nil, fmt.Errorf("Credential provider named '%s' could not "+
				"be found", cpc.Provider)
		}
	}

	var creds *credentials.Credentials

	if len(providers) > 0 {
		creds = credentials.NewChainCredentials(providers)
	}

	return creds, nil
}

func ProviderConfig(provider credentials.Provider) (*CredentialsProviderConfig, error) {
	switch p := provider.(type) {
	case *credentials.EnvProvider:
		return &CredentialsProviderConfig{
			Provider:    credentials.EnvProviderName,
			EnvProvider: p,
		}, nil
	case *credentials.SharedCredentialsProvider:
		return &CredentialsProviderConfig{
			Provider:                  credentials.SharedCredsProviderName,
			SharedCredentialsProvider: p,
		}, nil
	case *credentials.StaticProvider:
		return &CredentialsProviderConfig{
			Provider: credentials.StaticProviderName,
			StaticProviderConfig: &StaticProviderConfig{
				AccessKeyID:     p.AccessKeyID,
				SecretAccessKey: p.SecretAccessKey,
				SessionToken:    p.SessionToken,
			},
		}, nil
	default:
		return nil, fmt.Errorf("Credential provide %s is not recognised", p)
	}
}
