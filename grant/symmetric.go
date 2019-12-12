package grant

import (
	"fmt"

	"github.com/monax/hoard/v7/encryption"
	"github.com/monax/hoard/v7/reference"
)

// SymmetricGrant encrypts the given reference based on a secret read from the provider store
func SymmetricGrant(ref *reference.Ref, secret []byte) ([]byte, error) {
	if len(secret) < encryption.KeySize {
		return nil, fmt.Errorf("SymmetricGrant cannot encrypt with a secret of size < %d", encryption.KeySize)
	}
	// Generate AES nonce
	nonce, err := encryption.NewNonce(encryption.NonceSize)
	if err != nil {
		return nil, fmt.Errorf("SymmetricGrant failed to generate random nonce: %v", err)
	}
	// Encrypt reference with key and nonce
	blob, err := encryption.Encrypt([]byte(ref.Plaintext(nil)), nonce, secret)
	if err != nil {
		return nil, fmt.Errorf("SymmetricGrant failed to encyrpt: %v", err)
	}

	// we previously stored a key-derivation salt as part of a prefix to the ciphertext,
	// we no longer use this value but we must add padding to preserve backward compatibility
	padding := make([]byte, encryption.NonceSize)

	// Store salt and nonce so we can re-derive key/decrypt later
	return encryption.Salinate(blob.EncryptedData, append(nonce, padding...)), nil
}

// SymmetricReference decrypts the given grant based on a secret read from the provider store
func SymmetricReference(ciphertext, secret []byte) (*reference.Ref, error) {
	// Extract nonce and salt stored with ciphertext
	encryptedData, nonceAndSalt := encryption.Desalinate(ciphertext, encryption.NonceSize+encryption.NonceSize)

	// we previously stored a key-derivation salt in the latter half of nonceAndSalt,
	// we no longer use this but we still store empty space there to preserve backward compatibility
	nonce := nonceAndSalt[:encryption.NonceSize]

	// Decrypt
	data, err := encryption.Decrypt(encryptedData, nonce, secret)
	if err != nil {
		return nil, fmt.Errorf("SymmetricReference failed to decrypt: %v", err)
	}
	// Deserialise reference
	return reference.FromPlaintext(string(data)), nil
}
