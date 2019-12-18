package grant

import (
	"encoding/base64"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/monax/hoard/v7/config"
	"github.com/monax/hoard/v7/encryption"
	"github.com/monax/hoard/v7/reference"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
)

func TestGrants(t *testing.T) {
	testRef := testReference()

	keyPrivate, err := ioutil.ReadFile("private.key.asc")
	assert.NoError(t, err)

	testPGP := config.OpenPGPSecret{
		PrivateID: "10449759736975846181",
		Data:      keyPrivate,
	}
	testSecrets := newSecretsManager(nil, &testPGP)

	plaintextSpec := Spec{Plaintext: &PlaintextSpec{}}
	plaintextGrant, err := Seal(testSecrets, testRef, &plaintextSpec)
	assert.NoError(t, err)
	assert.Equal(t, testRef.Address, reference.FromPlaintext(string(plaintextGrant.EncryptedReference)).Address)
	assert.Equal(t, testRef.SecretKey, reference.FromPlaintext(string(plaintextGrant.EncryptedReference)).SecretKey)
	plaintextRef, err := Unseal(testSecrets, plaintextGrant)
	assert.Equal(t, testRef, plaintextRef)

	// SymmetricGrant with empty provider
	symmetricSpec := Spec{Symmetric: &SymmetricSpec{PublicID: "test"}}
	symmetricGrant, err := Seal(testSecrets, testRef, &symmetricSpec)
	assert.Error(t, err)
	assert.Nil(t, symmetricGrant)

	secret := deriveSecret(t, []byte("sssshhhh"))

	// SymmetricGrant with correct provider
	testSecrets.Provider = func(_ string) (config.SymmetricSecret, error) {
		return config.SymmetricSecret{SecretKey: secret}, nil
	}
	symmetricGrant, err = Seal(testSecrets, testRef, &symmetricSpec)
	assert.NotNil(t, symmetricGrant)
	assert.NoError(t, err)
	symmetricRef, err := Unseal(testSecrets, symmetricGrant)
	assert.Equal(t, testRef, symmetricRef)
	assert.NoError(t, err)

	// OpenPGPGrant encrypt / decrypt with local keypair
	openpgpSpec := Spec{OpenPGP: &OpenPGPSpec{}}
	openpgpGrant, err := Seal(testSecrets, testRef, &openpgpSpec)
	assert.NoError(t, err)
	openpgpRef, err := Unseal(testSecrets, openpgpGrant)
	assert.Equal(t, testRef, openpgpRef)
	assert.NoError(t, err)
}

func mustDecodeString(str string) []byte {
	ciphertext, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return ciphertext
}

func TestUnsealV0Grant(t *testing.T) {
	secrets := newSecretsManager(map[string]string{
		"testing-id-1": strings.Repeat("A", encryption.KeySize),
		"testing-id-2": strings.Repeat("A", encryption.KeySize-1),
	}, nil)

	var params = []struct {
		id         string
		ciphertext string
	}{
		{
			"testing-id-1",
			"Rki+cOHZ1WClgLUx3/6AlP48p//fz8Y8hEbAqYsM2w/os1dQ+yViX6JPRI/BcJW7ebSmwzisnekowWjZ6w+Zpi7EFa52q8SXZOgg5Qi5RmAfHDpbbtQNGpLIQUrCIXaa/+6TKpiEKB67Vq+9OIhjtI1pThTPDyMGc6dBHx6P9d+zfALn4iAOPURWma93vjZKsJON6sU3YzHIc3+Gag==",
		},
		{
			"testing-id-2",
			"+WErtplQBsz3Uq+LTbyxEI1JMUDWBqJdHeFey3gSG/KOgnp55xRqDGa4bq/ByksQ1EOPjFSD3AwU/Zc2Z+1E1PhAizp+uhdbJvtHXbEL1x/Ox/zEBQ/x4ZI5cMxtiB0LtPWfAvWaA8OmHYZkvNnJ/zoD4Ch/TV4+Y8h7Q8dLipcsG6PEVNWvIW52W61XJUBQozf/iZOpx6dRcv4xwA==",
		},
	}

	for _, tt := range params {
		ciphertext, err := base64.StdEncoding.DecodeString(tt.ciphertext)
		require.NoError(t, err)

		_, err = Unseal(secrets, &Grant{
			Spec: &Spec{
				Symmetric: &SymmetricSpec{PublicID: tt.id},
			},
			EncryptedReference: ciphertext,
			Version:            0,
		})
		require.NoError(t, err)
	}
}

func newSecretsManager(secrets map[string]string, pgp *config.OpenPGPSecret) config.SecretsManager {
	return config.SecretsManager{
		Provider: func(id string) (config.SymmetricSecret, error) {
			return config.SymmetricSecret{
				SecretKey:  []byte(secrets[id]),
				Passphrase: secrets[id],
			}, nil
		},
		OpenPGP: pgp,
	}
}

func testReference() *reference.Ref {
	address := []byte{
		1, 2, 3, 4, 5, 6, 7, 1,
		1, 2, 3, 4, 5, 6, 7, 1,
		1, 2, 3, 4, 5, 6, 7, 1,
		1, 2, 3, 4, 5, 6, 7, 1,
	}
	secretKey := []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
	}
	return reference.New(address, secretKey, nil)
}

func deriveSecret(t *testing.T, data []byte) []byte {
	salt, err := encryption.NewNonce(encryption.NonceSize)
	assert.NoError(t, err)
	secret, err := encryption.DeriveSecretKey(data, salt)
	assert.NoError(t, err)
	return secret
}
