package grant

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/monax/hoard/v4/reference"
)

func TestSymmetricGrant(t *testing.T) {
	ref := &reference.Ref{
		Address:   []byte("adddressss"),
		SecretKey: []byte("other secret"),
	}
	grt, err := SymmetricGrant(ref, nil)
	assert.Error(t, err)
	assert.Nil(t, grt)

	secret := []byte("sshh")
	grt, err = SymmetricGrant(ref, secret)
	assert.Errorf(t, err, "SymmetricGrant cannot encrypt with a secret of size < %d", minSecretSize)
	assert.Nil(t, grt)

	secret = []byte("sssshhhh")
	grt, err = SymmetricGrant(ref, secret)
	assert.NoError(t, err)
	assert.NotNil(t, grt)
	refOut, err := SymmetricReference(grt, secret)
	assert.NoError(t, err)
	assert.Equal(t, ref, refOut)
}
