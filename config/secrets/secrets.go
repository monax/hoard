package secrets

import "io/ioutil"

type SecretsConfig struct {
	Secrets []SymmetricSecret
	OpenPGP *OpenPGPSecret
}

type SymmetricSecret struct {
	ID         string
	Passphrase string
}

type OpenPGPSecret struct {
	ID   uint64
	File string
	Data []byte
}

type SymmetricProvider func(secretID string) []byte

type Manager struct {
	Provider SymmetricProvider
	OpenPGP  *OpenPGPSecret
}

func NoopSymmetricProvider(_ string) []byte {
	return nil
}

var NoopSecretManager = Manager{
	Provider: NoopSymmetricProvider,
	OpenPGP:  nil,
}

func ProviderFromConfig(conf *SecretsConfig) SymmetricProvider {
	if conf == nil {
		return NoopSymmetricProvider
	}
	secs := make(map[string][]byte, len(conf.Secrets))
	for _, s := range conf.Secrets {
		secs[s.ID] = []byte(s.Passphrase)
	}
	return func(id string) []byte {
		return secs[id]
	}
}

func OpenPGPFromConfig(conf *SecretsConfig) *OpenPGPSecret {
	if conf == nil {
		return nil
	}
	keyRing, err := ioutil.ReadFile(conf.OpenPGP.File)
	if err != nil {
		return nil
	}
	conf.OpenPGP.Data = keyRing
	return conf.OpenPGP
}
