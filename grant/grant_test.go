package grant

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/monax/hoard/v8/versions"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"

	"github.com/monax/hoard/v8/config"
	"github.com/monax/hoard/v8/encryption"
	"github.com/monax/hoard/v8/reference"
	"github.com/stretchr/testify/assert"
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
	refs := reference.MustRefsFromPlaintext(plaintextGrant.EncryptedReferences, versions.LatestGrantVersion)
	require.NotEmpty(t, refs)
	assert.Equal(t, testRefs[0].Address, refs[0].Address)
	assert.Equal(t, testRefs[0].SecretKey, refs[0].SecretKey)
	plaintextRef, err := Unseal(testSecrets, plaintextGrant)
	assertRefsEqual(t, testRefs, plaintextRef)

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
	assertRefsEqual(t, testRefs, symmetricRef)
	assert.NoError(t, err)

	// OpenPGPGrant encrypt / decrypt with local keypair
	openpgpSpec := Spec{OpenPGP: &OpenPGPSpec{}}
	openpgpGrant, err := Seal(testSecrets, testRefs, &openpgpSpec)
	assert.NoError(t, err)
	openpgpRef, err := Unseal(testSecrets, openpgpGrant)
	assertRefsEqual(t, testRefs, openpgpRef)
	assert.NoError(t, err)
}

func assertRefsEqual(t *testing.T, as, bs []*reference.Ref) {
	for i, ref := range as {
		assert.True(t, proto.Equal(ref, bs[i]))
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

func testReferences() []*reference.Ref {
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
	return []*reference.Ref{reference.New(address, secretKey, nil, 1024)}
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
