package secrets

type SecretsConfig struct {
	Secrets []Secret
}

type Secret struct {
	ID         string
	Passphrase string
}
