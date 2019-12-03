package grant

import (
	"github.com/monax/hoard/v7/reference"
)

// PlaintextGrant returns an encoded reference
func PlaintextGrant(ref *reference.Ref) []byte {
	return []byte(ref.Plaintext(nil))
}

// PlaintextReference decodes the grant
func PlaintextReference(ciphertext []byte) *reference.Ref {
	return reference.FromPlaintext(string(ciphertext))
}
