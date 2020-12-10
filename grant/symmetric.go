package grant

import (
	"fmt"

	"github.com/monax/hoard/v8/encryption"
	"github.com/monax/hoard/v8/reference"
)

// SymmetricGrant encrypts the given reference based on a secret read from the provider store
func SymmetricGrant(refs []*reference.Ref, secret []byte) ([]byte, error) {
	if len(secret) < encryption.KeySize {
		return nil, fmt.Errorf("SymmetricGrant cannot encrypt with a secret of size < %d", encryption.KeySize)
	}
	// Generate AES nonce
	nonce, err := encryption.NewNonce(encryption.NonceSize)
	if err != nil {
		return nil, fmt.Errorf("SymmetricGrant failed to generate random nonce: %v", err)
	}
	// Encrypt reference with key and nonce
	plaintext, err := reference.PlaintextFromRefs(refs, nil)
	if err != nil {
		return nil, err
	}
	blob, err := encryption.Encrypt(plaintext, nonce, secret)
	if err != nil {
		return nil, fmt.Errorf("SymmetricGrant failed to encyrpt: %v", err)
	}

	// Store salt and nonce so we can re-derive key/decrypt later
	return encryption.Salinate(blob.EncryptedData, nonce), nil
}

func SymmetricReference(ciphertext, secret []byte, version int32) ([]*reference.Ref, error) {
	encryptedData, nonce := encryption.Desalinate(ciphertext, encryption.NonceSize)
	data, err := encryption.Decrypt(encryptedData, nonce, secret)
	if err != nil {
		return nil, fmt.Errorf("SymmetricReference failed to decrypt: %v", err)
	}
	return reference.RefsFromPlaintext(data, version)
}
