package grant

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/monax/hoard/config/secrets"
	"github.com/monax/hoard/reference"
)

func TestGrants(t *testing.T) {
	testRef := testReference()

	keyPrivate, err := ioutil.ReadFile("private.key.asc")
	assert.NoError(t, err)

	testPGP := secrets.OpenPGPSecret{
		PrivateID: "10449759736975846181",
		Data:      []byte(keyPrivate),
	}

	testSecrets := secrets.Manager{
		Provider: func(_ string) []byte {
			return nil
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

	symmetricSpec := Spec{Symmetric: &SymmetricSpec{PublicID: "test"}}
	symmetricGrant, err := Seal(testSecrets, testRef, &symmetricSpec)
	assert.NoError(t, err)
	symmetricRef, err := Unseal(testSecrets, symmetricGrant)
	assert.Equal(t, testRef, symmetricRef)
	assert.NoError(t, err)

	asymmetricSpec := Spec{OpenPGP: &OpenPGPSpec{}}
	asymmetricGrant, err := Seal(testSecrets, testRef, &asymmetricSpec)
	assert.NoError(t, err)
	asymmetricRef, err := Unseal(testSecrets, asymmetricGrant)
	assert.Equal(t, testRef, asymmetricRef)
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
