package reference

import (
	"encoding/json"
	"fmt"
)

func New(address, secretKey, salt []byte) *Ref {
	if len(salt) == 0 {
		salt = nil
	}
	return &Ref{
		Address:   address,
		SecretKey: secretKey,
		Salt:      salt,
	}
}

type Refs []*Ref

// Note the Salt here is different to the salt that may have been used to encrypt
// the data pointed to by the reference.
type refWithNonce struct {
	*Ref
	Nonce []byte `json:",omitempty"`
}

type refsWithNonce struct {
	Refs
	Nonce []byte `json:",omitempty"`
}

// Obtain the canonical plaintext for the Ref with an optional nonce that can be
// be used to salt the plaintext in order to obtain an unpredictable version of
// the plaintext for encryption purposes (i.e. for Grants). The nonce is
// discarded when read by FromPlaintext
func plaintext(wrapper interface{}) string {
	bs, err := json.Marshal(wrapper)
	if err != nil {
		panic(fmt.Errorf("did not expect an error when serialising reference, " +
			"error suppressed for security"))
	}
	return string(bs)
}

func (ref *Ref) Plaintext(nonce []byte) string {
	return plaintext(refWithNonce{Ref: ref, Nonce: nonce})

}

func (refs Refs) Plaintext(nonce []byte) string {
	return plaintext(refsWithNonce{Refs: refs, Nonce: nonce})
}

func fromPlaintext(wrapper interface{}, plaintext string) {
	err := json.Unmarshal([]byte(plaintext), wrapper)
	if err != nil {
		panic(fmt.Errorf("did not expect an error when deserialising reference, " +
			"error suppressed for security"))
	}
}

func FromPlaintext(plaintext string) *Ref {
	wrapper := new(refWithNonce)
	fromPlaintext(wrapper, plaintext)
	return wrapper.Ref
}

func RepeatedFromPlaintext(plaintext string) Refs {
	wrapper := new(refsWithNonce)
	fromPlaintext(wrapper, plaintext)
	return wrapper.Refs
}
