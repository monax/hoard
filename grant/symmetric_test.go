package grant

import (
	"encoding/base64"
	"testing"

	"github.com/monax/hoard/v7/encryption"
	"github.com/monax/hoard/v7/reference"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
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

	_, err = SymmetricReferenceV0(grt, secret)
	assert.Error(t, err)
}

func TestSymmetricReferenceV0(t *testing.T) {
	// this fixture was taken from v6
	ciphertext, err := base64.StdEncoding.DecodeString("w07sfs2T8qjBkLahiTxG7eO1ex+xy1wacxTAyQUtJpbO8Vcv+UeoSNJc0dZY4Nds3XXs9npre0/EKHsxPNyvEzKwcpYf8IQlVFxRIU5Jm5wi5S5gzTrELHipzuVtPltGy3QUixKYtrxm3ykC6rI1kCd8eYZyRBkpnJeuL2HhEE6RT2FdUEGZthiyeZ0JTkDaln6Youm1Gp8G5UnSaQ==")
	require.NoError(t, err)
	secret := []byte("secret-passphrase")
	_, err = SymmetricReferenceV0(ciphertext, secret)
	require.NoError(t, err)
}
