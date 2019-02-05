package grant

import (
	"github.com/monax/hoard/reference"
)

func PlaintextGrant(ref *reference.Ref) []byte {
	return []byte(ref.Plaintext(nil))
}

func PlaintextReference(ciphertext []byte) *reference.Ref {
	return reference.FromPlaintext(string(ciphertext))
}
