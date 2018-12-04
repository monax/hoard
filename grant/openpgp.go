package grant

import (
	"errors"
	"fmt"
	"strings"

	"bytes"
	"io"
	"io/ioutil"

	"github.com/monax/hoard/reference"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

const (
	OpenPGPGrantType = "HOARD OPENPGP GRANT"
)

func OpenPGPGrant(ref *reference.Ref, to []*openpgp.Entity,
	signed *openpgp.Entity) (string, error) {

	buf := bytes.NewBuffer(nil)
	armorWriter, err := armor.Encode(buf, OpenPGPGrantType, nil)
	if err != nil {
		return "", err
	}
	plaintextWriter, err := openpgp.Encrypt(armorWriter, to, signed, nil, nil)
	if err != nil {
		fmt.Errorf("Could not set up openpgp encryption: %s", err)
	}

	_, err = io.WriteString(plaintextWriter, ref.Plaintext(nil))
	if err != nil {
		return "", err
	}

	plaintextWriter.Close()
	armorWriter.Close()
	return buf.String(), nil
}

func OpenPGPGrantReference(grant string,
	keyRing openpgp.KeyRing) (*reference.Ref, error) {
	block, err := armor.Decode(strings.NewReader(grant))
	if err != nil {
		return nil, err
	}

	if block.Type != OpenPGPGrantType {
		return nil, errors.New("OpenPGP block does not appear to be a Hoard OpenPGP grant")
	}
	messageReader, err := openpgp.ReadMessage(block.Body, keyRing,
		func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
			return nil, errors.New("OpenPGPGrantReference does not support password prompting")
		}, nil)
	if err != nil {
		return nil, err
	}

	bs, err := ioutil.ReadAll(messageReader.UnverifiedBody)
	if err != nil {
		return nil, err
	}
	return reference.FromPlaintext(string(bs)), nil
}
