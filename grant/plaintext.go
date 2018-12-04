package grant

import (
	"github.com/monax/hoard/reference"
)

func PlaintextGrant(ref *reference.Ref) string {
	return SaltedPlaintextGrant(ref, nil)
}

func SaltedPlaintextGrant(ref *reference.Ref, salt []byte) string {
	return ref.Plaintext(salt)
}

func PlaintextGrantReference(grant string) *reference.Ref {
	return reference.FromPlaintext(grant)
}
