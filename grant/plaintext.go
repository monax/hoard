package grant

import (
	"github.com/monax/hoard/v8/reference"
)

// PlaintextGrant returns an encoded reference
func PlaintextGrant(refs reference.Refs) []byte {
	return refs.Plaintext(nil)
}

func PlaintextReferenceV2(ciphertext []byte) reference.Refs {
	return reference.RepeatedFromPlaintext(ciphertext)
}
