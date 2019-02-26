package grant

import (
	"fmt"

	"github.com/monax/hoard/config/secrets"
	"github.com/monax/hoard/reference"
)

// Seal this reference into a Grant as specified by Spec
func Seal(secret secrets.Manager, ref *reference.Ref, spec *Spec) (*Grant, error) {
	grt := &Grant{Spec: spec}

	switch s := spec.GetValue().(type) {
	case *PlaintextSpec:
		grt.EncryptedReference = PlaintextGrant(ref)
	case *SymmetricSpec:
		encRef, err := SymmetricGrant(ref, secret.Provider(s.PublicID))
		if err != nil {
			return nil, err
		}
		grt.EncryptedReference = encRef
	case *OpenPGPSpec:
		encRef, err := OpenPGPGrant(ref, s.PublicKey, secret.OpenPGP)
		if err != nil {
			return nil, err
		}
		grt.EncryptedReference = encRef
	default:
		return nil, fmt.Errorf("grant type %v not recognised", s)
	}
	return grt, nil
}

// Unseal a Grant exposing its secret reference
func Unseal(secret secrets.Manager, grt *Grant) (*reference.Ref, error) {
	switch s := grt.Spec.GetValue().(type) {
	case *PlaintextSpec:
		return PlaintextReference(grt.EncryptedReference), nil
	case *SymmetricSpec:
		return SymmetricReference(grt.EncryptedReference, secret.Provider(s.PublicID))
	case *OpenPGPSpec:
		return OpenPGPReference(grt.EncryptedReference, secret.OpenPGP)
	}
	return nil, fmt.Errorf("grant type %v not recognised", grt.Spec)
}
