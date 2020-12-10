package reference

import (
	"bytes"
	"fmt"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/monax/hoard/v8/protodet"
	"github.com/monax/hoard/v8/versions"
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

// Obtain the canonical plaintext for the Ref with an optional nonce that can be used to make a particular
// array of refs unique, as is usually required for LINK refs
func PlaintextFromRefs(refs []*Ref, nonce []byte) ([]byte, error) {
	refsWithNonce := &RefsWithNonce{
		Refs:  refs,
		Nonce: nonce,
	}
	bs, err := protodet.Marshal(refsWithNonce)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling to plaintext, error supressed for security")
	}
	return bs, nil
}

func MustPlaintextFromRefs(refs []*Ref, nonce []byte) []byte {
	bs, err := PlaintextFromRefs(refs, nonce)
	if err != nil {
		panic(err)
	}
	return bs
}

func refsFromProtobuf(plaintext []byte) ([]*Ref, error) {
	wrapper := new(RefsWithNonce)
	err := protodet.Unmarshal(plaintext, wrapper)
	return wrapper.Refs, err
}

func refsFromJSON(plaintext []byte) ([]*Ref, error) {
	wrapper := new(RefsWithNonce)
	m := jsonpb.Unmarshaler{}
	err := m.Unmarshal(bytes.NewBuffer(plaintext), wrapper)
	return wrapper.Refs, err
}

func RefsFromPlaintext(plaintext []byte, version int32) (refs []*Ref, err error) {
	switch version {
	case 0, 1, 2:
		refs, err = refsFromJSON(plaintext)
		for _, ref := range refs {
			if ref.Version == versions.RefVersionIncorrectlyUsedToDenoteHeader {
				ref.Type = Ref_HEADER
			}
		}
	default:
		refs, err = refsFromProtobuf(plaintext)
	}
	if err != nil {
		err = fmt.Errorf("error while unmarshalling from plaintexst, error supressed for security")
	}
	return
}

func MustRefsFromPlaintext(plaintext []byte, version int32) []*Ref {
	refs, err := RefsFromPlaintext(plaintext, version)
	if err != nil {
		panic(err)
	}
	return refs
}
