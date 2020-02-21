package grant

import (
	"fmt"

	"github.com/monax/hoard/v7/config"
	"github.com/monax/hoard/v7/reference"
)

const defaultGrantVersion = 2

// Seal this reference into a Grant as specified by Spec
func Seal(secret config.SecretsManager, refs reference.Refs, spec *Spec) (*Grant, error) {
	grt := &Grant{Spec: spec, Version: defaultGrantVersion}

	if s := spec.GetPlaintext(); s != nil {
		grt.EncryptedReferences = PlaintextGrant(refs)
	} else if s := spec.GetSymmetric(); s != nil {
		secret, err := secret.Provider(s.PublicID)
		if err != nil {
			return nil, err
		}
		encRef, err := SymmetricGrant(refs, secret.SecretKey)
		if err != nil {
			return nil, err
		}
		grt.EncryptedReferences = encRef
	} else if s := spec.GetOpenPGP(); s != nil {
		encRef, err := OpenPGPGrant(refs, s.PublicKey, secret.OpenPGP)
		if err != nil {
			return nil, err
		}
		grt.EncryptedReferences = encRef
	} else {
		return nil, fmt.Errorf("grant type %v not recognised", s)
	}

	return grt, nil
}

// Unseal a Grant exposing its secret reference
func Unseal(secret config.SecretsManager, grt *Grant) (reference.Refs, error) {
	if s := grt.Spec.GetPlaintext(); s != nil {
		switch grt.GetVersion() {
		case 0:
			return PlaintextReferenceV0(grt.EncryptedReferences), nil
		case 1:
			return PlaintextReferenceV1(grt.EncryptedReferences), nil
		default:
			return PlaintextReferenceV2(grt.EncryptedReferences), nil
		}

	} else if s := grt.Spec.GetSymmetric(); s != nil {
		secret, err := secret.Provider(s.PublicID)
		if err != nil {
			return nil, err
		}
		switch grt.GetVersion() {
		case 0:
			return SymmetricReferenceV0(grt.EncryptedReferences, []byte(secret.Passphrase))
		case 1:
			return SymmetricReferenceV1(grt.EncryptedReferences, secret.SecretKey)
		default:
			return SymmetricReferenceV2(grt.EncryptedReferences, secret.SecretKey)
		}

	} else if s := grt.Spec.GetOpenPGP(); s != nil {
		switch grt.GetVersion() {
		case 0:
			return OpenPGPReferenceV0(grt.EncryptedReferences, secret.OpenPGP)
		case 1:
			return OpenPGPReferenceV1(grt.EncryptedReferences, secret.OpenPGP)
		default:
			return OpenPGPReferenceV2(grt.EncryptedReferences, secret.OpenPGP)
		}

	} else {
		return nil, fmt.Errorf("grant type %v not recognised", s)
	}
}
