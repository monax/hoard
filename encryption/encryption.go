package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// There are issues with the proof of security with other nonce sizes so we only use 0 (for convergent) or 12
const NonceSize = 12
const KeySize = 32

type BlockCipherMaker func(key []byte) (cipher.Block, error)

type Blob struct {
	SecretKey     []byte
	EncryptedData []byte
}

type Args struct {
	BlockCipherMaker BlockCipherMaker
	Nonce            []byte
	SecretKey        []byte
	AdditionalData   []byte
}

// Encrypt data convergently by using a securely generated deterministic
// key that is a hash of the plaintext (data/blob) itself. Allows for
// deduplication of ciphertexts and recovery of keys from plaintext alone.

// Deterministically encrypt data using a supplied salt to produce a
// distinguished the encrypted result that will have a different content hash,
// secret key, and address than the same data encrypted with a different salt
// (or not salt). Can be used to watermark a copy of a blob shared with a
// particular party or to hide the fact a certain plaintext is stored.
func EncryptConvergent(data, salt []byte) (*Blob, error) {
	// The SHA 256 hasher will be used to generate the secret key for AES. Since
	// the AES cipher is parameterised by the length of the secret key in this case
	// with the 32 byte key from SHA 256 we will get a AES 256 block cipher.
	// First hash the plaintext securely, we will use its hash as a key
	hasher := sha256.New()
	hasher.Write(data)
	secretKey := hasher.Sum(nil)
	return encrypt(Salinate(data, salt), Args{
		BlockCipherMaker: aes.NewCipher,
		SecretKey:        secretKey,
		Nonce:            nil,
		AdditionalData:   additionalDataForSalt(salt),
	})
}

// Decrypt data that was deterministically encrypted with the provided salt
func DecryptConvergent(encryptedData, salt, secretKey []byte) ([]byte, error) {
	data, err := decrypt(encryptedData, Args{
		BlockCipherMaker: aes.NewCipher,
		SecretKey:        secretKey,
		Nonce:            nil,
		AdditionalData:   additionalDataForSalt(salt),
	})
	if err != nil {
		return nil, err
	}
	plaintext, _ := Desalinate(data, len(salt))
	return plaintext, nil
}

// Encrypt data using random nonce
func Encrypt(data, nonce, secretKey []byte) (*Blob, error) {
	return encrypt(data, Args{
		BlockCipherMaker: aes.NewCipher,
		SecretKey:        secretKey,
		Nonce:            nonce,
	})
}

func Decrypt(encryptedData, nonce, secretKey []byte) ([]byte, error) {
	return decrypt(encryptedData, Args{
		BlockCipherMaker: aes.NewCipher,
		SecretKey:        secretKey,
		Nonce:            nonce,
	})
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
func encrypt(plaintext []byte, args Args) (*Blob, error) {
	blockCipher, err := args.BlockCipherMaker(args.SecretKey)
	if err != nil {
		return nil, err
	}
	// We can operate without a nonce because we are using a one-time key (the
	// secure hash of the data) that will be not used for other messages (blobs)
	// so IV/key pair is unique
	gcmCipher, err := cipher.NewGCMWithNonceSize(blockCipher, len(args.Nonce))
	if err != nil {
		return nil, err
	}

	// Encrypt authenticated with Galois Counter mode
	// TODO: consider storing contract address relating to blob in additional data
	ciphertext := gcmCipher.Seal(nil, args.Nonce, plaintext, args.AdditionalData)

	return &Blob{
		SecretKey:     args.SecretKey,
		EncryptedData: ciphertext,
	}, nil
}

// Decrypt ciphertext encrypted in Galois Counter Mode over the block cipher
// provided by blockCipherMaker assuming a one-time key and so no nonce as
// would be encrypted by encryptConvergent (though would work for any one-time
// key case). Salt is used as GCM additional authenticated data.
func decrypt(ciphertext []byte, args Args) ([]byte, error) {
	// Construct the underlying block cipher
	blockCipher, err := args.BlockCipherMaker(args.SecretKey)
	if err != nil {
		return nil, err
	}

	gcmCipher, err := cipher.NewGCMWithNonceSize(blockCipher, len(args.Nonce))
	if err != nil {
		return nil, err
	}

	return gcmCipher.Open(nil, args.Nonce, ciphertext, args.AdditionalData)
}

func Salinate(data, salt []byte) []byte {
	return append(data, salt...)
}

func Desalinate(data []byte, length int) (desalted []byte, salt []byte) {
	saltStart := len(data) - length
	return data[:saltStart], data[saltStart:]
}

// Provides additional authenticated data to fix context of our salting procedure
// using this means if we try to decrypt an unsalted message with a salt or visa
// versa we will get an error decrypting.
func additionalDataForSalt(salt []byte) []byte {
	if len(salt) == 0 {
		return nil
	}
	additionalData := struct {
		SaltType   string
		SaltLength int
	}{
		SaltType:   "prefix",
		SaltLength: len(salt),
	}
	jsonBytes, err := json.Marshal(additionalData)
	if err != nil {
		// We control this struct, we can exhaustively test so shouldn't panic
		panic(fmt.Errorf("could not marshal additional data describing "+
			"salting procedure: %#v, error: %s", jsonBytes, err))
	}
	return jsonBytes
}
