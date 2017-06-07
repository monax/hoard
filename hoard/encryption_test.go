package hoard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptionRoundTrip(t *testing.T) {
	plaintext := []byte("Hello this is a string")
	blob, err := Encrypt(plaintext)
	assert.NoError(t, err)
	decryptedPlaintext, err := Decrypt(blob.SecretKey(), blob.EncryptedData())
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decryptedPlaintext)
}

func TestEncryptedIdentity(t *testing.T) {
	plaintext := []byte("Identical plaintext should lead to identical blob")
	blob1, err := Encrypt(plaintext)
	assert.NoError(t, err)
	blob2, err := Encrypt(plaintext)
	assert.NoError(t, err)
	assert.Equal(t, blob1, blob2)
}

func TestSaltedEncryptionRoundTrip(t *testing.T) {
	plaintext := []byte("Hello this is a string")
	salt := []byte("salty like the sea")
	saltedBlob, err := EncryptSalted(plaintext, salt)
	assert.NoError(t, err)
	decryptedPlaintext, err := DecryptSalted(saltedBlob.SecretKey(), salt,
		saltedBlob.EncryptedData())
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decryptedPlaintext)

	// Should be error to decrypt salted an unsalted blob
	unsaltedBlob, err := Encrypt(plaintext)
	_, err = DecryptSalted(unsaltedBlob.SecretKey(), salt,
		unsaltedBlob.EncryptedData())
	assert.Error(t, err, "Should fail on salted decrypt of unsalted blob")

	// Conversely should be an error to decrypt normally a salted blob
	_, err = Decrypt(saltedBlob.SecretKey(), saltedBlob.EncryptedData())
	assert.Error(t, err, "Should fail on unsalted decrypt of salted blob")
}

func TestAdditionalDataForSalt(t *testing.T) {
	// This function may panic on marshalling so we try to cover it here
	assert.Equal(t, "{\"SaltType\":\"prefix\",\"SaltLength\":0}",
		string(additionalDataForSalt(nil)))
	assert.Equal(t, "{\"SaltType\":\"prefix\",\"SaltLength\":0}",
		string(additionalDataForSalt([]byte(""))))
	assert.Equal(t, "{\"SaltType\":\"prefix\",\"SaltLength\":21}",
		string(additionalDataForSalt([]byte("I _am_ a magical fish"))))
}
