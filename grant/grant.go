package grant

import (
	"fmt"

	"github.com/monax/hoard/reference"
)

// Interface to secret backend
type SecretProvider func(secretID string) []byte

// Seal this reference into a Grant as specified by Spec
func Seal(secret SecretProvider, ref *reference.Ref, spec *Spec) (*Grant, error) {
	grt := &Grant{Spec: spec}

	switch s := spec.GetValue().(type) {
	case PlaintextSpec:
		grt.EncryptedReference = PlaintextGrant(ref)
	case SymmetricSpec:
		encRef, err := SymmetricGrant(ref, secret(s.SecretID))
		if err != nil {
			return nil, err
		}
		grt.EncryptedReference = encRef
	case OpenPGPSpec:
		//OpenPGPGrant(ref, spec.Data)
	default:
		return nil, fmt.Errorf("grant type %v not recognised", s)
	}
	return grt, nil
}

// Unseal a Grant exposing its secret reference
func Unseal(secret SecretProvider, grt *Grant) (*reference.Ref, error) {
	switch s := grt.Spec.GetValue().(type) {
	case PlaintextSpec:
		return PlaintextReference(grt.EncryptedReference), nil
	case SymmetricSpec:
		return SymmetricReference(grt.EncryptedReference, secret(s.SecretID))
	case OpenPGPSpec:
		//OpenPGPGrant(ref, spec.Data)
	}
	return nil, fmt.Errorf("grant type %v not recognised", grt.Spec)
}
