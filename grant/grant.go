package grant

import (
	"fmt"

	"github.com/monax/hoard/v8/versions"

	"github.com/monax/hoard/v8/config"
	"github.com/monax/hoard/v8/reference"
)

// Seal this reference into a Grant as specified by Spec
func Seal(secret config.SecretsManager, refs []*reference.Ref, spec *Spec) (*Grant, error) {
	grt := &Grant{Spec: spec, Version: versions.LatestGrantVersion}

	var err error
	if s := spec.GetPlaintext(); s != nil {
		grt.EncryptedReferences, err = PlaintextGrant(refs)
		if err != nil {
			return nil, err
		}
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
func Unseal(secret config.SecretsManager, grt *Grant) ([]*reference.Ref, error) {
	// invert version switch, deal with old JSON references
	if s := grt.Spec.GetPlaintext(); s != nil {
		return PlaintextReference(grt.EncryptedReferences, grt.GetVersion())

	}
	if s := grt.Spec.GetSymmetric(); s != nil {
		secret, err := secret.Provider(s.PublicID)
		if err != nil {
			return nil, err
		}
		return SymmetricReference(grt.EncryptedReferences, secret.SecretKey, grt.GetVersion())
	}
	if s := grt.Spec.GetOpenPGP(); s != nil {
		return OpenPGPReference(grt.EncryptedReferences, secret.OpenPGP, grt.GetVersion())
	}
	return nil, fmt.Errorf("grant type not recognised")
}
