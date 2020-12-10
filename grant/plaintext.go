package grant

import (
	"github.com/monax/hoard/v8/reference"
)

// PlaintextGrant returns an encoded reference
func PlaintextGrant(refs []*reference.Ref) ([]byte, error) {
	return reference.PlaintextFromRefs(refs, nil)
}

func PlaintextReference(ciphertext []byte, version int32) ([]*reference.Ref, error) {
	return reference.RefsFromPlaintext(ciphertext, version)
}
