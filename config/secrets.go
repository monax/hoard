package config

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

// Secrets lists the configured secrets,
// Symmetric secrets are those local to the running daemon
// and OpenPGP identifies an entity in the given keyring
type Secrets struct {
	Symmetric []SymmetricSecret
	OpenPGP   *OpenPGPSecret
}

type SymmetricSecret struct {
	// An identifier for this secret that will be stored in the clear with the grant
	PublicID  string
	SecretKey string
}

type OpenPGPSecret struct {
	// A private (though not secret) identifier that points to a PGP keyring that this instance of hoard
	// will use to provide PGP grants
	PrivateID string
	File      string
	Data      []byte
}

type SecretsManager struct {
	Provider SymmetricProvider
	OpenPGP  *OpenPGPSecret
}

type SymmetricProvider func(secretID string) ([]byte, error)

// NoopSecretManager is an empty secret manager
var NoopSecretManager = SecretsManager{
	Provider: NoopSymmetricProvider,
	OpenPGP:  nil,
}

// NoopSymmetricProvider returns an empty provider
func NoopSymmetricProvider(_ string) ([]byte, error) {
	return nil, fmt.Errorf("no secrets provided to hoard")
}

// ProviderFromConfig creates a secret reader from a set of symmetric secrets
func NewSymmetricProvider(conf *Secrets, fromEnv bool) (SymmetricProvider, error) {
	if conf == nil || len(conf.Symmetric) == 0 {
		return NoopSymmetricProvider, nil
	}
	secs := make(map[string][]byte, len(conf.Symmetric))
	for _, s := range conf.Symmetric {
		if fromEnv {
			// sometimes we don't want to specify these in the config
			s.SecretKey = os.Getenv(s.PublicID)
		}
		secret, err := b64.StdEncoding.DecodeString(s.SecretKey)
		if err != nil {
			return nil, err
		}
		secs[s.PublicID] = secret
	}
	return func(id string) ([]byte, error) {
		if id == "" {
			return nil, fmt.Errorf("empty secret ID passed to provider")
		}
		if val, ok := secs[id]; ok {
			return val, nil
		}
		return nil, fmt.Errorf("could not find symmetric secret with ID '%s'", id)
	}, nil
}

// OpenPGPFromConfig reads a given PGP keyring
func NewOpenPGPSecret(conf *Secrets) *OpenPGPSecret {
	if conf == nil || conf.OpenPGP == nil {
		return nil
	}
	keyRing, err := ioutil.ReadFile(conf.OpenPGP.File)
	if err != nil {
		return nil
	}
	conf.OpenPGP.Data = keyRing
	return conf.OpenPGP
}
