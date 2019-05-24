package grant

import (
	"fmt"

	"github.com/monax/hoard/v5/config"
	"github.com/monax/hoard/v5/reference"
)

// Seal this reference into a Grant as specified by Spec
func Seal(secret config.SecretsManager, ref *reference.Ref, spec *Spec) (*Grant, error) {
	grt := &Grant{Spec: spec}

	if s := spec.GetPlaintext(); s != nil {
		grt.EncryptedReference = PlaintextGrant(ref)
	} else if s := spec.GetSymmetric(); s != nil {
		secret, err := secret.Provider(s.PublicID)
		if err != nil {
			return nil, err
		}
		encRef, err := SymmetricGrant(ref, secret)
		if err != nil {
			return nil, err
		}
		grt.EncryptedReference = encRef
	} else if s := spec.GetOpenPGP(); s != nil {
		encRef, err := OpenPGPGrant(ref, s.PublicKey, secret.OpenPGP)
		if err != nil {
			return nil, err
		}
		grt.EncryptedReference = encRef
	} else {
		return nil, fmt.Errorf("grant type %v not recognised", s)
	}

	return grt, nil
}

// Unseal a Grant exposing its secret reference
func Unseal(secret config.SecretsManager, grt *Grant) (*reference.Ref, error) {
	if s := grt.Spec.GetPlaintext(); s != nil {
		return PlaintextReference(grt.EncryptedReference), nil
	} else if s := grt.Spec.GetSymmetric(); s != nil {
		secret, err := secret.Provider(s.PublicID)
		if err != nil {
			return nil, err
		}
		return SymmetricReference(grt.EncryptedReference, secret)
	} else if s := grt.Spec.GetOpenPGP(); s != nil {
		return OpenPGPReference(grt.EncryptedReference, secret.OpenPGP)
	} else {
		return nil, fmt.Errorf("grant type %v not recognised", s)
	}
}
