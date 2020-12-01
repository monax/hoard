package grant

import (
	"fmt"

	"github.com/monax/hoard/v8/encryption"
	"github.com/monax/hoard/v8/reference"
)

// SymmetricGrant encrypts the given reference based on a secret read from the provider store
func SymmetricGrant(ref reference.Refs, secret []byte) ([]byte, error) {
	if len(secret) < encryption.KeySize {
		return nil, fmt.Errorf("SymmetricGrant cannot encrypt with a secret of size < %d", encryption.KeySize)
	}
	// Generate AES nonce
	nonce, err := encryption.NewNonce(encryption.NonceSize)
	if err != nil {
		return nil, fmt.Errorf("SymmetricGrant failed to generate random nonce: %v", err)
	}
	// Encrypt reference with key and nonce
	blob, err := encryption.Encrypt(ref.Plaintext(nil), nonce, secret)
	if err != nil {
		return nil, fmt.Errorf("SymmetricGrant failed to encyrpt: %v", err)
	}

	// Store salt and nonce so we can re-derive key/decrypt later
	return encryption.Salinate(blob.EncryptedData, nonce), nil
}

func SymmetricReferenceV2(ciphertext, secret []byte) (reference.Refs, error) {
	encryptedData, nonce := encryption.Desalinate(ciphertext, encryption.NonceSize)
	data, err := encryption.Decrypt(encryptedData, nonce, secret)
	if err != nil {
		return nil, fmt.Errorf("SymmetricReferenceV2 failed to decrypt: %v", err)
	}
	return reference.RepeatedFromPlaintext(data), nil
}
