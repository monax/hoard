package grant

import (
	"github.com/monax/hoard/v7/reference"
)

// PlaintextGrant returns an encoded reference
func PlaintextGrant(refs reference.Refs) []byte {
	return []byte(refs.Plaintext(nil))
}

// PlaintextReferenceV0 decodes the grant
func PlaintextReferenceV0(ciphertext []byte) reference.Refs {
	return reference.Refs{reference.FromPlaintext(string(ciphertext))}
}

func PlaintextReferenceV1(ciphertext []byte) reference.Refs {
	return PlaintextReferenceV0(ciphertext)
}

func PlaintextReferenceV2(ciphertext []byte) reference.Refs {
	return reference.RepeatedFromPlaintext(string(ciphertext))
}
