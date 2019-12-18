package grant

import (
	"fmt"

	"github.com/monax/hoard/v7/config"
	"github.com/monax/hoard/v7/reference"
)

const defaultGrantVersion = 1

// Seal this reference into a Grant as specified by Spec
func Seal(secret config.SecretsManager, ref *reference.Ref, spec *Spec) (*Grant, error) {
	grt := &Grant{Spec: spec, Version: defaultGrantVersion}

	if s := spec.GetPlaintext(); s != nil {
		grt.EncryptedReference = PlaintextGrant(ref)
	} else if s := spec.GetSymmetric(); s != nil {
		secret, err := secret.Provider(s.PublicID)
		if err != nil {
			return nil, err
		}
		encRef, err := SymmetricGrant(ref, secret.SecretKey)
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
		switch grt.GetVersion() {
		case 0:
			return SymmetricReferenceV0(grt.EncryptedReference, []byte(secret.Passphrase))
		default:
			return SymmetricReferenceV1(grt.EncryptedReference, secret.SecretKey)
		}
	} else if s := grt.Spec.GetOpenPGP(); s != nil {
		return OpenPGPReference(grt.EncryptedReference, secret.OpenPGP)
	} else {
		return nil, fmt.Errorf("grant type %v not recognised", s)
	}
}
