package reference

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// We more or may not want grants to be deterministic (as in byte-wise identical
// for the same reference). We provide the ability to salt a reference for the case when
// we expressly want to avoid them being deterministic. Furthermore the JSON
// spec doesn't specify a canonical field order. However golang has a pretty
// stable ordering for structs, though we probably shouldn't depend on it. In
// case we are this test is a canary that will alert us if that ever changes.
// If it is useful to have strictly deterministic grants then we should consider
// a canonical ordering. It might be useful in some circumstances to see that one
// reference is the same as an other without knowing the reference content (or having to
// decrypt it), but from this vantage point that case seems fairly marginal.
// If this test fails but TestGrantPlaintext passes, consider removing this test.
func TestReferencePlaintextDeterministic(t *testing.T) {
	assert.Equal(t,
		"{\"Address\":\"AQIDBAUGBwEBAgMEBQYHAQECAwQFBgcBAQIDBAUGBwE=\","+
			"\"SecretKey\":\"AQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIAQIDBAUGBwg=\""+
			"}",
		testReference(nil).Plaintext(nil))

	assert.Equal(t,
		"{\"Address\":\"AQIDBAUGBwEBAgMEBQYHAQECAwQFBgcBAQIDBAUGBwE=\","+
			"\"SecretKey\":\"AQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIAQIDBAUGBwg=\","+
			"\"Salt\":\"c2FsdA==\""+
			"}",
		testReference(([]byte)("salt")).Plaintext(nil))

	assert.Equal(t,
		"{\"Address\":\"AQIDBAUGBwEBAgMEBQYHAQECAwQFBgcBAQIDBAUGBwE=\","+
			"\"SecretKey\":\"AQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIAQIDBAUGBwg=\","+
			"\"Salt\":\"c2FsdA==\","+
			"\"Nonce\":\"bm9uY2U=\""+
			"}",
		testReference(([]byte)("salt")).Plaintext(([]byte)("nonce")))
}

func TestReferencePlaintext(t *testing.T) {
	ref := testReference(nil)
	assert.Equal(t, ref,
		FromPlaintext(ref.Plaintext(nil)))
	assert.Equal(t, ref,
		FromPlaintext(ref.Plaintext(([]byte)("nonce"))))
}

func testReference(salt []byte) *Ref {
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
	return New(address, secretKey, salt)
}
