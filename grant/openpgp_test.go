package grant

import (
	"io/ioutil"
	"testing"

	"github.com/monax/hoard/v8/versions"

	"github.com/monax/hoard/v8/config"

	"github.com/stretchr/testify/assert"
)

func TestOpenPGPGrant(t *testing.T) {
	testRefs := testReferences()

	keyPublic, err := ioutil.ReadFile("public.key.asc")
	assert.NoError(t, err)
	keyPrivate, err := ioutil.ReadFile("private.key.asc")
	assert.NoError(t, err)

	testPGP := config.OpenPGPSecret{
		PrivateID: "10449759736975846181",
		Data:      keyPrivate,
	}

	// Create grant from public
	grant, err := OpenPGPGrant(testRefs, string(keyPublic), &testPGP)
	assert.NoError(t, err)

	// Try to read reference from grant
	ref, err := OpenPGPReference(grant, &testPGP, versions.LatestGrantVersion)
	assert.NoError(t, err)
	assertRefsEqual(t, testRefs, ref)

	// Create grant from private
	grant, err = OpenPGPGrant(testRefs, string(keyPrivate), &testPGP)
	assert.NoError(t, err)

	// Try to read reference from grant
	ref, err = OpenPGPReference(grant, &testPGP, versions.LatestGrantVersion)
	assert.NoError(t, err)
	assertRefsEqual(t, testRefs, ref)

	// Create grant from signer
	grant, err = OpenPGPGrant(testRefs, "", &testPGP)
	assert.NoError(t, err)

	// Try to read reference from grant
	ref, err = OpenPGPReference(grant, &testPGP, versions.LatestGrantVersion)
	assert.NoError(t, err)
	assertRefsEqual(t, testRefs, ref)

	grant, err = OpenPGPGrant(testRefs, string(keyPublic), nil)
	assert.Errorf(t, err, "hoard is not currently configured to use openpgp")
}
