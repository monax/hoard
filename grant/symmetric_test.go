package grant

import (
	"testing"

	"github.com/monax/hoard/reference"
	"github.com/stretchr/testify/require"
)

func TestSymmetricGrant(t *testing.T) {
	secret := []byte("sshh")
	ref := &reference.Ref{
		Address:   []byte("adddressss"),
		SecretKey: []byte("other secret"),
	}
	grt, err := SymmetricGrant(ref, secret)
	require.NoError(t, err)
	refOut, err := SymmetricReference(grt, secret)
	require.NoError(t, err)
	require.Equal(t, ref, refOut)
}
