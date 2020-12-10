package reference

import (
	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/gogo/protobuf/proto"
	"github.com/monax/hoard/v8/versions"
	"github.com/stretchr/testify/assert"
	"testing"
)

// References should be serialised using a deterministic method, a non-deterministic salt or nonce provides a
// mechanism to make an individual ref or an array of refs non-deterministic
func TestReferencePlaintextDeterministic(t *testing.T) {
	cases := []struct {
		name string
		refs []*Ref
		nonce []byte
	}{
		{
			name: "UnsaltedNoNonce",
			refs:  testRefs(nil),
		},
		{
			name: "SaltedNoNonce",
			refs: testRefs([]byte("salt")),
		},
		{
			name: "SaltedNonce",
			refs: testRefs([]byte("salt")),
			nonce: []byte("nonce"),

		},
		{
			name: "RepeatedSaltedNonce",
			refs: append(testRefs([]byte("salt1")), testRefs([]byte("salt2"))...),
			nonce: []byte("nonce"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cupaloy.SnapshotT(t, string(MustPlaintextFromRefs(c.refs, c.nonce)))
		})
	}
}

func TestReferencePlaintext(t *testing.T) {
	refs := testRefs(nil)
	assertRefsEqual(t, refs,
		MustRefsFromPlaintext(MustPlaintextFromRefs(refs, nil), versions.LatestGrantVersion))
	assertRefsEqual(t, refs,
		MustRefsFromPlaintext(MustPlaintextFromRefs(refs, ([]byte)("nonce")), versions.LatestGrantVersion))
}

func testRefs(salt []byte) []*Ref {
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
	return []*Ref{New(address, secretKey, salt, 1024)}
}

func assertRefsEqual(t *testing.T, as, bs []*Ref) {
	for i, ref := range as {
		assert.True(t, proto.Equal(ref, bs[i]))
	}
}
