package config

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/monax/hoard/v8/encryption"
	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestSecretKeyMarshal(t *testing.T) {
	salt := make([]byte, encryption.NonceSize)
	key, err := encryption.DeriveSecretKey([]byte("hello"), salt)
	assert.NoError(t, err)

	secret := SecretKey(key)
	data, err := secret.MarshalText()
	assert.NoError(t, err)
	expected := "bFQ+wRhNaOgC4fNcliGFaZ5Xr3wOywYJZP1eqj6SDCk="
	assert.Equal(t, expected, string(data))

	inSecret := new(SymmetricSecret)
	inSecret.SecretKey = secret
	outSecret := new(SymmetricSecret)

	data, err = json.Marshal(inSecret)
	assert.NoError(t, err)
	assert.Equal(t, "{\"PublicID\":\"\",\"SecretKey\":\"bFQ+wRhNaOgC4fNcliGFaZ5Xr3wOywYJZP1eqj6SDCk=\",\"Passphrase\":\"\"}", string(data))
	err = json.Unmarshal(data, outSecret)
	assert.NoError(t, err)
	assert.Equal(t, key, []byte(outSecret.SecretKey))

	data, err = yaml.Marshal(inSecret)
	assert.NoError(t, err)
	assert.Equal(t, "publicid: \"\"\nsecretkey: bFQ+wRhNaOgC4fNcliGFaZ5Xr3wOywYJZP1eqj6SDCk=\npassphrase: \"\"\n", string(data))
	err = yaml.Unmarshal(data, outSecret)
	assert.NoError(t, err)
	assert.Equal(t, key, []byte(outSecret.SecretKey))

	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	err = encoder.Encode(inSecret)
	assert.NoError(t, err)
	assert.Equal(t, "PublicID = \"\"\nSecretKey = \"bFQ+wRhNaOgC4fNcliGFaZ5Xr3wOywYJZP1eqj6SDCk=\"\nPassphrase = \"\"\n", buf.String())
	err = toml.Unmarshal(buf.Bytes(), outSecret)
	assert.NoError(t, err)
	assert.Equal(t, key, []byte(outSecret.SecretKey))

	err = toml.Unmarshal([]byte("SecretKey = \"bFQ+wRhNaOgC4fNcliGFaZ5Xr3wOywYJZP1eqj6SDCk=\"\n"), outSecret)
	assert.NoError(t, err)
	assert.Equal(t, key, []byte(outSecret.SecretKey))
	assert.Equal(t, "", outSecret.PublicID)
	assert.Equal(t, "", outSecret.Passphrase)

	err = toml.Unmarshal([]byte("PublicID = \"\"\nSecretKey = \"badkey=\"\n"), outSecret)
	assert.Error(t, err)
}
