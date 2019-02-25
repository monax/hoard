package grant

import (
	"io/ioutil"
	"testing"

	"github.com/monax/hoard/config/secrets"
	"github.com/stretchr/testify/assert"
)

func TestOpenPGPGrant(t *testing.T) {
	testRef := testReference()

	keyPublic, err := ioutil.ReadFile("public.key.asc")
	assert.NoError(t, err)
	keyPrivate, err := ioutil.ReadFile("private.key.asc")
	assert.NoError(t, err)

	testPGP := secrets.OpenPGPSecret{
		ID:   "10449759736975846181",
		Data: []byte(keyPrivate),
	}

	// Create grant from public
	grant, err := OpenPGPGrant(testRef, []byte(keyPublic), &testPGP)
	assert.NoError(t, err)

	// Try to read reference from grant
	ref, err := OpenPGPReference(grant, &testPGP)
	assert.NoError(t, err)
	assert.EqualValues(t, testRef, ref)

	// Create grant from private
	grant, err = OpenPGPGrant(testRef, []byte(keyPrivate), &testPGP)
	assert.NoError(t, err)

	// Try to read reference from grant
	ref, err = OpenPGPReference(grant, &testPGP)
	assert.NoError(t, err)
	assert.EqualValues(t, testRef, ref)

	// Create grant from signer
	grant, err = OpenPGPGrant(testRef, nil, &testPGP)
	assert.NoError(t, err)

	// Try to read reference from grant
	ref, err = OpenPGPReference(grant, &testPGP)
	assert.NoError(t, err)
	assert.EqualValues(t, testRef, ref)

	ref, err = OpenPGPReference(grant, nil)
	assert.Errorf(t, err, "hoard is not currently configured to use openpgp")

	grant, err = OpenPGPGrant(testRef, []byte(keyPublic), nil)
	assert.Errorf(t, err, "hoard is not currently configured to use openpgp")
}
