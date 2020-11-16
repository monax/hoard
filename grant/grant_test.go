package grant

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/monax/hoard/v8/config"
	"github.com/monax/hoard/v8/encryption"
	"github.com/monax/hoard/v8/reference"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
)

func TestGrants(t *testing.T) {
	testRefs := testReferences()

	keyPrivate, err := ioutil.ReadFile("private.key.asc")
	assert.NoError(t, err)

	testPGP := config.OpenPGPSecret{
		PrivateID: "10449759736975846181",
		Data:      keyPrivate,
	}
	testSecrets := newSecretsManager(nil, &testPGP)

	plaintextSpec := Spec{Plaintext: &PlaintextSpec{}}
	plaintextGrant, err := Seal(testSecrets, testRefs, &plaintextSpec)
	assert.NoError(t, err)
	assert.Equal(t, testRefs[0].Address, reference.RepeatedFromPlaintext(string(plaintextGrant.EncryptedReferences))[0].Address)
	assert.Equal(t, testRefs[0].SecretKey, reference.RepeatedFromPlaintext(string(plaintextGrant.EncryptedReferences))[0].SecretKey)
	plaintextRef, err := Unseal(testSecrets, plaintextGrant)
	assert.Equal(t, testRefs, plaintextRef)

	// SymmetricGrant with empty provider
	symmetricSpec := Spec{Symmetric: &SymmetricSpec{PublicID: "test"}}
	symmetricGrant, err := Seal(testSecrets, testRefs, &symmetricSpec)
	assert.Error(t, err)
	assert.Nil(t, symmetricGrant)

	secret := deriveSecret(t, []byte("sssshhhh"))

	// SymmetricGrant with correct provider
	testSecrets.Provider = func(_ string) (config.SymmetricSecret, error) {
		return config.SymmetricSecret{SecretKey: secret}, nil
	}
	symmetricGrant, err = Seal(testSecrets, testRefs, &symmetricSpec)
	assert.NotNil(t, symmetricGrant)
	assert.NoError(t, err)
	symmetricRef, err := Unseal(testSecrets, symmetricGrant)
	assert.Equal(t, testRefs, symmetricRef)
	assert.NoError(t, err)

	// OpenPGPGrant encrypt / decrypt with local keypair
	openpgpSpec := Spec{OpenPGP: &OpenPGPSpec{}}
	openpgpGrant, err := Seal(testSecrets, testRefs, &openpgpSpec)
	assert.NoError(t, err)
	openpgpRef, err := Unseal(testSecrets, openpgpGrant)
	assert.Equal(t, testRefs, openpgpRef)
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
			EncryptedReferences: ciphertext,
			Version:             0,
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

func testReferences() reference.Refs {
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
	return reference.Refs{reference.New(address, secretKey, nil)}
}

func deriveSecret(t *testing.T, data []byte) []byte {
	salt, err := encryption.NewNonce(encryption.NonceSize)
	assert.NoError(t, err)
	secret, err := encryption.DeriveSecretKey(data, salt)
	assert.NoError(t, err)
	return secret
}

func TestUnmarshal(t *testing.T) {
	// the client library stores the grant with lowercase field names,
	// we expect the go server to correctly unmarshal this
	data := `{"spec":{"plaintext":{},"symmetric":null,"openpgp":null},"encryptedreferences":"eyJSZWZzIjpbeyJBZGRyZXNzIjoidDIzZjh1cTZsd3lJL2ZTTGJaMVJ2b3ZMYzFSSDMwWEk4cUlyUzBQZnljOD0iLCJTZWNyZXRLZXkiOiI0N0RFUXBqOEhCU2ErL1RJbVcrNUpDZXVRZVJrbTVOTXBKV1pHM2hTdUZVPSIsIlZlcnNpb24iOjF9LHsiQWRkcmVzcyI6Ii8rdWxUa0N6cFlnMnNQYVp0cVM4ZHljSkJMWTkzODd5WlBzdDhMWDVZTDA9IiwiU2VjcmV0S2V5IjoidGJ1ZGdCU2crYkhXSGlIbmx0ZU56TjhUVXZJODB5Z1M5SVVMaDRya2xFdz0ifV19","version":2}`
	grant := new(Grant)
	err := json.Unmarshal([]byte(data), grant)
	require.NoError(t, err)
	require.Equal(t, int32(2), grant.GetVersion())
	require.NotNil(t, grant.GetSpec().GetPlaintext())
}

const secret = `{
"PublicID": "uMfCSTU6UgNDi6lgapqM0tWwJdWXcT77",
"SecretKey": "a2gJTHGZZK1OHrMa5xdv2U3DB0AgvU78JHcB21rrnYk="
}`

func TestLongGrant(t *testing.T) {
	data := "{\"spec\":{\"symmetric\":{\"publicid\":\"uMfCSTU6UgNDi6lgapqM0tWwJdWXcT77\"}},\"encryptedreferences\":\"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\",\"version\":2}"

	sec := new(config.SymmetricSecret)
	err := json.Unmarshal([]byte(secret), sec)
	require.NoError(t, err)
	grant := new(Grant)
	err = json.Unmarshal([]byte(data), grant)
	require.NoError(t, err)
	secrets := config.SecretsManager{
		Provider: func(secretID string) (config.SymmetricSecret, error) {
			return *sec, nil
		},
		OpenPGP: nil,
	}
	refs, err := Unseal(secrets, grant)
	require.NoError(t, err)
	fmt.Println(refs)
}
