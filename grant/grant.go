package grant

import (
	"fmt"

	"github.com/monax/hoard/v8/config"
	"github.com/monax/hoard/v8/reference"
)

// We use the Grant version as a kind of global version for the core protobuf types.
// The version is a simple counter and does not 'encode' a breaking change
// (the intention is by using the version number we support all non deprecated previous versions, mearning no change is 'breaking')
// Version history:
// 0: deprecated and removed
// 1: deprecated and removed
// 2: encrypted references array for streaming, non-derived keys, reference with version
// 3: reference Version -> Type, introduce LINK references, store plaintext data Size in reference
const LatestGrantVersion = 3

// Seal this reference into a Grant as specified by Spec
func Seal(secret config.SecretsManager, refs reference.Refs, spec *Spec) (*Grant, error) {
	grt := &Grant{Spec: spec, Version: LatestGrantVersion}

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
		default:
			return PlaintextReferenceV2(grt.EncryptedReferences), nil
		}

	} else if s := grt.Spec.GetSymmetric(); s != nil {
		secret, err := secret.Provider(s.PublicID)
		if err != nil {
			return nil, err
		}
		switch grt.GetVersion() {
		default:
			return SymmetricReferenceV2(grt.EncryptedReferences, secret.SecretKey)
		}

	} else if s := grt.Spec.GetOpenPGP(); s != nil {
		switch grt.GetVersion() {
		default:
			return OpenPGPReferenceV2(grt.EncryptedReferences, secret.OpenPGP)
		}

	} else {
		return nil, fmt.Errorf("grant type %v not recognised", s)
	}
}
