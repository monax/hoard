package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptionRoundTrip(t *testing.T) {
	plaintext := []byte("Hello this is a string")
	blob, err := Encrypt(plaintext, nil)
	assert.NoError(t, err)
	decryptedPlaintext, err := Decrypt(blob.SecretKey(), blob.EncryptedData(), nil)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decryptedPlaintext)
}

func TestEncryptedIdentity(t *testing.T) {
	plaintext := []byte("Identical plaintext should lead to identical blob")
	blob1, err := Encrypt(plaintext, nil)
	assert.NoError(t, err)
	blob2, err := Encrypt(plaintext, nil)
	assert.NoError(t, err)
	assert.Equal(t, blob1, blob2)
}

func TestSaltedEncryptionRoundTrip(t *testing.T) {
	plaintext := []byte("Hello this is a string")
	salt := []byte("salty like the sea")
	saltedBlob, err := Encrypt(plaintext, salt)
	assert.NoError(t, err)
	decryptedPlaintext, err := Decrypt(saltedBlob.SecretKey(),
		saltedBlob.EncryptedData(), salt)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decryptedPlaintext)

	// Should be error to decrypt salted an unsalted blob
	unsaltedBlob, err := Encrypt(plaintext, nil)
	_, err = Decrypt(unsaltedBlob.SecretKey(),
		unsaltedBlob.EncryptedData(), salt)
	assert.Error(t, err, "Should fail on salted decrypt of unsalted blob")

	// Conversely should be an error to decrypt normally a salted blob
	_, err = Decrypt(saltedBlob.SecretKey(), saltedBlob.EncryptedData(), nil)
	assert.Error(t, err, "Should fail on unsalted decrypt of salted blob")
}

func TestAdditionalDataForSalt(t *testing.T) {
	// This function may panic on marshalling so we try to cover it here
	assert.Nil(t, additionalDataForSalt(nil))
	assert.Nil(t, additionalDataForSalt([]byte("")))
	assert.Equal(t, "{\"SaltType\":\"prefix\",\"SaltLength\":21}",
		string(additionalDataForSalt([]byte("I _am_ a magical fish"))))
}
