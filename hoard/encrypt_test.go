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
