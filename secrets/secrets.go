package secrets

import (
	"github.com/monax/hoard/config/secrets"
	"github.com/monax/hoard/grant"
)

func NoopSecretProvider(_ string) []byte {
	return nil
}

func SecretProviderFromConfig(conf *secrets.SecretsConfig) grant.SecretProvider {
	if conf == nil {
		return NoopSecretProvider
	}
	secs := make(map[string][]byte, len(conf.Secrets))
	for _, s := range conf.Secrets {
		secs[s.ID] = []byte(s.Passphrase)
	}
	return func(id string) []byte {
		return secs[id]
	}
}
