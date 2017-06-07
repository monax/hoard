package hoard

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"hash"
)

type BlockCipherMaker func(key []byte) (cipher.Block, error)

type EncryptedBlob interface {
	SecretKey() []byte
	Address() []byte
	EncryptedData() []byte
}

type encryptedBlob struct {
	secretKey     []byte
	address       []byte
	encryptedData []byte
}

func (blob *encryptedBlob) SecretKey() []byte {
	return blob.secretKey
}

func (blob *encryptedBlob) Address() []byte {
	return blob.address
}

func (blob *encryptedBlob) EncryptedData() []byte {
	return blob.encryptedData
}

// Encrypt stuff convergently. That is, using a securely generated deterministic
// key that is a hash of the plaintext (data/blob) itself. Allows for
// deduplication of ciphertexts and recovery of keys from plaintext alone.

// Deterministically encrypt data under hoard
func Encrypt(data []byte) (EncryptedBlob, error) {
	// We'll use sha256 like IPFS, and AES-256 (using hash as key)
	return encryptConvergent(sha256.New(), aes.NewCipher, data)
}

// Decrypt deterministically encrypted data under hoard
func Decrypt(secretKey, encryptedData []byte) ([]byte, error) {
	return decryptConvergent(aes.NewCipher, secretKey, encryptedData)
}

// Encrypt plaintext convergently by using a secure hash of the plaintext as the
// secret key to the block cipher produced by blockCipherMaker that will be used
// in Galois Counter Mode as a stream cipher with salt used as additional
// authenticated data.
//
// Note that this deterministic encryption is by design not secure under a chosen
// plaintext attack. However it can be used in this mode by prefixing a random,
// secret, or unique salt to the plaintext itself.
//
// The idea of this encryption modality is to make the secret key, ciphertext,
// and ultimately storage address recoverable from a copy of plaintext. If a
// semantically secure hash and block cipher are used then this does not leak
// information to a chosen plaintext attacker.
//
// However this is vulnerable to a 'guess the missing information' attack;
// if most of the plaintext is known then it is possible to brute force some
// remaining portion (such as an account number) to query whether a particular
// blob is stored. We actually want this behaviour to deduplicate and locate
// encrypted blobs. However if you want to distinguish copies of a plaintext or
// hide them add a random salt as above.
func encryptConvergent(hasher hash.Hash, blockCipherMaker BlockCipherMaker,
	plaintext []byte) (EncryptedBlob, error) {

	// First hash the plaintext securely, we will use its hash as a key
	hasher.Write(plaintext)
	secretKey := hasher.Sum(nil)
	blockCipher, err := blockCipherMaker(secretKey)
	if err != nil {
		return nil, err
	}
	// We can operate without a nonce because we are using a one-time key (the
	// secure hash of the data) that will be not used for other messages (blobs)
	// so IV/key pair is unique
	gcmCipher, err := cipher.NewGCMWithNonceSize(blockCipher, 0)
	if err != nil {
		return nil, err
	}

	// Encrypt authenticated with Galois Counter mode
	ciphertext := gcmCipher.Seal(nil, nil, plaintext, nil)

	// Hash the ciphertext to get the canonical content address
	hasher.Reset()
	hasher.Write(ciphertext)
	address := hasher.Sum(nil)

	return &encryptedBlob{
		secretKey:     secretKey,
		address:       address,
		encryptedData: ciphertext,
	}, nil
}

// Decrypt ciphertext encrypted in Galois Counter Mode over the block cipher
// provided by blockCipherMaker assuming a one-time key and so no nonce as
// would be encrypted by encryptConvergent (though would work for any one-time
// key case). Salt is used as GCM additional authenticated data.
func decryptConvergent(blockCipherMaker BlockCipherMaker, secretKey,
	ciphertext []byte) ([]byte, error) {
	// Construct the underlying block cipher
	blockCipher, err := blockCipherMaker(secretKey)
	if err != nil {
		return nil, err
	}

	gcmCipher, err := cipher.NewGCMWithNonceSize(blockCipher, 0)
	if err != nil {
		return nil, err
	}

	return gcmCipher.Open(nil, nil, ciphertext, nil)
}
