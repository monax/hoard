package storage

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type AddressEncoding interface {
	EncodeToString(address []byte) (addressString string)
	DecodeString(addressString string) (address []byte, err error)
}

func GetAddressEncoding(name string) (AddressEncoding, error) {
	switch name {
	case "base64", "":
		return base64.StdEncoding, nil
	case "base32":
		return base32.StdEncoding, nil
	case "base16", "hex":
		return NewAddressEncoding(hex.EncodeToString, hex.DecodeString), nil
	}
	return nil, fmt.Errorf("Could not find an address encoding named '%s'",
		name)
}

func NewAddressEncoding(encodeToString func([]byte) string,
	decodeString func(string) ([]byte, error)) *addressEncoding {
	return &addressEncoding{
		encodeToString: encodeToString,
		decodeString:   decodeString,
	}
}

var _ AddressEncoding = (*addressEncoding)(nil)

type addressEncoding struct {
	encodeToString func([]byte) string
	decodeString   func(string) ([]byte, error)
}

func (ae *addressEncoding) EncodeToString(address []byte) string {
	return ae.encodeToString(address)
}

func (ae *addressEncoding) DecodeString(addressString string) ([]byte, error) {
	return ae.decodeString(addressString)
}
