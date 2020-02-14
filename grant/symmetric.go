package grant

import (
	"fmt"

	"github.com/monax/hoard/v7/encryption"
	"github.com/monax/hoard/v7/reference"
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
	blob, err := encryption.Encrypt([]byte(ref.Plaintext(nil)), nonce, secret)
	if err != nil {
		return nil, fmt.Errorf("SymmetricGrant failed to encyrpt: %v", err)
	}

	// Store salt and nonce so we can re-derive key/decrypt later
	return encryption.Salinate(blob.EncryptedData, nonce), nil
}

// SymmetricReferenceV0 decrypts the given grant based on a passphrase read from the provider store
// TODO: deprecate after migration due to high memory overhead of scrypt
func SymmetricReferenceV0(ciphertext, secret []byte) (reference.Refs, error) {
	// Extract nonce and salt stored with ciphertext
	encryptedData, nonceAndSalt := encryption.Desalinate(ciphertext, encryption.NonceSize+encryption.NonceSize)
	nonce, salt := nonceAndSalt[:encryption.NonceSize], nonceAndSalt[encryption.NonceSize:]
	// Re-derive key based on these
	secretKey, err := encryption.DeriveSecretKey(secret, salt)
	if err != nil {
		return nil, fmt.Errorf("SymmetricReferenceV0 failed to derive secret key: %v", err)
	}
	// Decrypt
	data, err := encryption.Decrypt(encryptedData, nonce, secretKey)
	if err != nil {
		return nil, fmt.Errorf("SymmetricReferenceV0 failed to decrypt: %v", err)
	}
	// Deserialise reference
	return reference.Refs{reference.FromPlaintext(string(data))}, nil
}

// SymmetricReferenceV1 decrypts the given grant based on a secret read from the provider store
func SymmetricReferenceV1(ciphertext, secret []byte) (reference.Refs, error) {
	encryptedData, nonce := encryption.Desalinate(ciphertext, encryption.NonceSize)
	data, err := encryption.Decrypt(encryptedData, nonce, secret)
	if err != nil {
		return nil, fmt.Errorf("SymmetricReferenceV1 failed to decrypt: %v", err)
	}
	return reference.Refs{reference.FromPlaintext(string(data))}, nil
}

func SymmetricReferenceV2(ciphertext, secret []byte) (reference.Refs, error) {
	encryptedData, nonce := encryption.Desalinate(ciphertext, encryption.NonceSize)
	data, err := encryption.Decrypt(encryptedData, nonce, secret)
	if err != nil {
		return nil, fmt.Errorf("SymmetricReferenceV1 failed to decrypt: %v", err)
	}
	return reference.RepeatedFromPlaintext(string(data)), nil
}
