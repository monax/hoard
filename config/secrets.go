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
	Symmetric []*SymmetricSecret
	OpenPGP   *OpenPGPSecret
}

type SymmetricSecret struct {
	// An identifier for this secret that will be stored in the clear with the grant
	PublicID string
	// We expect this to be base64 encoded
	SecretKey SecretKey
	// Needed for backwards compatability
	Passphrase string
}

// SecretKey allows us to encode yaml and toml as base64
type SecretKey []byte

// MarshalText should fulfil most serialization interfaces to ensure that the
// secret key in the config is always base64 encoded
func (sec SecretKey) MarshalText() ([]byte, error) {
	data := b64.StdEncoding.EncodeToString(sec)
	return []byte(data), nil
}

func (sec *SymmetricSecret) UnmarshalTOML(in interface{}) error {
	if sec == nil {
		sec = new(SymmetricSecret)
	}

	data, _ := in.(map[string]interface{})
	sec.PublicID, _ = data["PublicID"].(string)
	sec.Passphrase, _ = data["Passphrase"].(string)

	secret, _ := data["SecretKey"].(string)
	key, err := b64.StdEncoding.DecodeString(secret)
	sec.SecretKey = key
	return err
}

func (sec *SymmetricSecret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if sec == nil {
		sec = new(SymmetricSecret)
	}

	secret := &struct {
		PublicID   string
		SecretKey  string
		Passphrase string
	}{}
	if err := unmarshal(secret); err != nil {
		return err
	}
	sec.PublicID = secret.PublicID
	sec.Passphrase = secret.Passphrase
	key, err := b64.StdEncoding.DecodeString(secret.SecretKey)
	sec.SecretKey = key
	return err
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

type SymmetricProvider func(secretID string) (SymmetricSecret, error)

// NoopSecretManager is an empty secret manager
var NoopSecretManager = SecretsManager{
	Provider: NoopSymmetricProvider,
	OpenPGP:  nil,
}

// NoopSymmetricProvider returns an empty provider
func NoopSymmetricProvider(_ string) (SymmetricSecret, error) {
	return SymmetricSecret{}, fmt.Errorf("no secrets provided to hoard")
}

// ProviderFromConfig creates a secret reader from a set of symmetric secrets
func NewSymmetricProvider(conf *Secrets, fromEnv bool) (SymmetricProvider, error) {
	if conf == nil || len(conf.Symmetric) == 0 {
		return NoopSymmetricProvider, nil
	}
	secs := make(map[string]SymmetricSecret, len(conf.Symmetric))
	for _, s := range conf.Symmetric {
		if fromEnv {
			// sometimes we don't want to specify these in the config
			secret := os.Getenv(s.PublicID)
			s.SecretKey = []byte(secret)
			s.Passphrase = secret
		}
		secs[s.PublicID] = SymmetricSecret{
			Passphrase: s.Passphrase,
			SecretKey:  s.SecretKey,
		}
	}
	return func(id string) (SymmetricSecret, error) {
		if id == "" {
			return SymmetricSecret{}, fmt.Errorf("empty secret ID passed to provider")
		}
		if val, ok := secs[id]; ok {
			return val, nil
		}
		return SymmetricSecret{}, fmt.Errorf("could not find symmetric secret with ID '%s'", id)
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
