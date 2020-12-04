package grant

import (
	"testing"

	"github.com/monax/hoard/v8/encryption"
	"github.com/monax/hoard/v8/reference"
	"github.com/stretchr/testify/assert"
)

func TestSymmetricGrant(t *testing.T) {
	ref := reference.Refs{&reference.Ref{
		Address:   []byte("adddressss"),
		SecretKey: []byte("other secret"),
	}}
	grt, err := SymmetricGrant(ref, nil)
	assert.Error(t, err)
	assert.Nil(t, grt)

	secret := []byte("sshh")
	grt, err = SymmetricGrant(ref, secret)
	assert.Errorf(t, err, "SymmetricGrant cannot encrypt with a secret of size < %d", encryption.KeySize)
	assert.Nil(t, grt)

	secret = deriveSecret(t, []byte("sssshhhh"))
	grt, err = SymmetricGrant(ref, secret)
	assert.NoError(t, err)
	assert.NotNil(t, grt)
	refOut, err := SymmetricReferenceV2(grt, secret)
	assert.NoError(t, err)
	assert.Equal(t, ref, refOut)
}
