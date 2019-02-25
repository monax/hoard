package secrets

import (
	"io/ioutil"
)

// SecretsConfig lists the configured secrets,
// Symmetric secrets are those local to the running daemon
// and OpenPGP identifies an entity in the given keyring
type SecretsConfig struct {
	Symmetric []SymmetricSecret
	OpenPGP   OpenPGPSecret
}

type SymmetricSecret struct {
	ID         string
	Passphrase string
}

type OpenPGPSecret struct {
	ID   string
	File string
	Data []byte
}

type Manager struct {
	Provider SymmetricProvider
	OpenPGP  *OpenPGPSecret
}

type SymmetricProvider func(secretID string) []byte

// NoopSecretManager is an empty secret manager
var NoopSecretManager = Manager{
	Provider: NoopSymmetricProvider,
	OpenPGP:  nil,
}

// NoopSymmetricProvider returns an empty provider
func NoopSymmetricProvider(_ string) []byte {
	return nil
}

// ProviderFromConfig creates a secret reader from a set of symmetric secrets
func ProviderFromConfig(conf *SecretsConfig) SymmetricProvider {
	if conf == nil || len(conf.Symmetric) == 0 {
		return NoopSymmetricProvider
	}
	secs := make(map[string][]byte, len(conf.Symmetric))
	for _, s := range conf.Symmetric {
		secs[s.ID] = []byte(s.Passphrase)
	}
	return func(id string) []byte {
		return secs[id]
	}
}

// OpenPGPFromConfig reads a given PGP keyring
func OpenPGPFromConfig(conf *SecretsConfig) *OpenPGPSecret {
	if conf == nil || conf.OpenPGP.File == "" {
		return nil
	}
	keyRing, err := ioutil.ReadFile(conf.OpenPGP.File)
	if err != nil {
		return nil
	}
	conf.OpenPGP.Data = keyRing
	return &conf.OpenPGP
}
