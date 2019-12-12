package grant

import (
	"io/ioutil"
	"testing"

	"github.com/monax/hoard/v7/config"
	"github.com/monax/hoard/v7/encryption"
	"github.com/monax/hoard/v7/reference"
	"github.com/stretchr/testify/assert"
)

func TestGrants(t *testing.T) {
	testRef := testReference()

	keyPrivate, err := ioutil.ReadFile("private.key.asc")
	assert.NoError(t, err)

	testPGP := config.OpenPGPSecret{
		PrivateID: "10449759736975846181",
		Data:      keyPrivate,
	}

	testSecrets := config.SecretsManager{
		Provider: func(_ string) ([]byte, error) {
			return nil, nil
		},
		OpenPGP: &testPGP,
	}

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

	secret, _ := computeSecretAndSalt(t, []byte("sssshhhh"))

	// SymmetricGrant with correct provider
	testSecrets.Provider = func(_ string) ([]byte, error) {
		return secret, nil
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

func computeSecretAndSalt(t *testing.T, data []byte) ([]byte, []byte) {
	salt, err := encryption.NewNonce(encryption.NonceSize)
	assert.NoError(t, err)
	secret, err := encryption.DeriveSecretKey(data, salt)
	assert.NoError(t, err)
	return secret, salt
}
