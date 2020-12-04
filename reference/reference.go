package reference

import (
	"encoding/json"
	"fmt"
)

func New(address, secretKey, salt []byte, size int64) *Ref {
	if len(salt) == 0 {
		salt = nil
	}
	return &Ref{
		Address:   address,
		SecretKey: secretKey,
		Salt:      salt,
		Size_:     size,
	}
}

type Refs []*Ref

// Note the Salt here is different to the salt that may have been used to encrypt
// the data pointed to by the reference.
type refsWithNonce struct {
	Refs
	Nonce []byte `json:",omitempty"`
}

// Obtain the canonical plaintext for the Ref with an optional nonce that can be
// be used to salt the plaintext in order to obtain an unpredictable version of
// the plaintext for encryption purposes (i.e. for Grants). The nonce is
// discarded when read by FromPlaintext
func plaintext(wrapper interface{}) []byte {
	bs, err := json.Marshal(wrapper)
	if err != nil {
		panic(fmt.Errorf("did not expect an error when serialising reference, " +
			"error suppressed for security"))
	}
	return bs
}

func (refs Refs) Plaintext(nonce []byte) []byte {
	return plaintext(refsWithNonce{Refs: refs, Nonce: nonce})
}

func fromPlaintext(wrapper interface{}, plaintext []byte) {
	err := json.Unmarshal(plaintext, wrapper)
	if err != nil {
		panic(fmt.Errorf("did not expect an error when deserialising reference, " +
			"error suppressed for security"))
	}
}

func RepeatedFromPlaintext(plaintext []byte) Refs {
	wrapper := new(refsWithNonce)
	fromPlaintext(wrapper, plaintext)
	return wrapper.Refs
}
