package grant

import (
	"crypto/rand"
	"fmt"

	"github.com/monax/hoard/encryption"
	"github.com/monax/hoard/reference"
	"golang.org/x/crypto/scrypt"
)

// We bump it a little from the 100ms for interactive logins rule: https://blog.filippo.io/the-scrypt-parameters/
const scryptSecurityWorkExponent = 16

// SymmetricGrant encrypts the given reference based on a secret read from the provider store
func SymmetricGrant(ref *reference.Ref, secret []byte) ([]byte, error) {
	// Generate scrypt salt
	salt := make([]byte, encryption.NonceSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("SymmetricGrant failed to generate random salt: %v", err)
	}
	// Derive key
	secretKey, err := DeriveSecretKey(secret, salt)
	if err != nil {
		return nil, fmt.Errorf("SymmetricGrant failed to derive secret key: %v", err)
	}
	// Generate AES nonce
	nonce := make([]byte, encryption.NonceSize)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("SymmetricGrant failed to generate random nonce: %v", err)
	}
	// Encrypt reference with key and nonce
	blob, err := encryption.Encrypt([]byte(ref.Plaintext(nil)), nonce, secretKey)
	if err != nil {
		return nil, fmt.Errorf("SymmetricGrant failed to encyrpt: %v", err)
	}
	// Store salt and nonce so we can re-derive key/decrypt later
	return encryption.Salinate(blob.EncryptedData, append(nonce, salt...)), nil
}

// SymmetricReference decrypts the given grant based on a secret read from the provider store
func SymmetricReference(ciphertext, secret []byte) (*reference.Ref, error) {
	// Extract nonce and salt stored with ciphertext
	encryptedData, nonceAndSalt := encryption.Desalinate(ciphertext, encryption.NonceSize+encryption.NonceSize)
	nonce, salt := nonceAndSalt[:encryption.NonceSize], nonceAndSalt[encryption.NonceSize:]
	// Re-derive key based on these
	secretKey, err := DeriveSecretKey(secret, salt)
	if err != nil {
		return nil, fmt.Errorf("SymmetricReference failed to derive secret key: %v", err)
	}
	// Decrypt
	data, err := encryption.Decrypt(encryptedData, nonce, secretKey)
	if err != nil {
		return nil, fmt.Errorf("SymmetricReference failed to decrypt: %v", err)
	}
	// Deserialise reference
	return reference.FromPlaintext(string(data)), nil
}

func DeriveSecretKey(secret, salt []byte) ([]byte, error) {
	return scrypt.Key(secret, salt, 1<<scryptSecurityWorkExponent, 8, 1, encryption.KeySize)
}
