package grant

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/monax/hoard/v8/config"

	"github.com/monax/hoard/v8/reference"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

// OpenPGPGrant encrypts and signs a given reference
func OpenPGPGrant(refs []*reference.Ref, public string, keyring *config.OpenPGPSecret) ([]byte, error) {
	if keyring == nil {
		return nil, fmt.Errorf("cannot encrypt because no private key was provided")
	}

	buf := bytes.NewBuffer(nil)
	armorWriter, err := armor.Encode(buf, "PGP MESSAGE", nil)
	if err != nil {
		return nil, err
	}

	var to openpgp.EntityList
	if public != "" {
		// use public keyring
		if to, err = openpgp.ReadArmoredKeyRing(bytes.NewBufferString(public)); err != nil {
			return nil, fmt.Errorf("could not read public keyring: %s", err)
		}
	} else {
		// default to configured keyring
		if to, err = openpgp.ReadArmoredKeyRing(bytes.NewBuffer(keyring.Data)); err != nil {
			return nil, fmt.Errorf("could not read private keyring: %s", err)
		}
	}

	// read configured keyring
	keys, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer(keyring.Data))
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(keyring.PrivateID, 10, 64)
	if err != nil {
		return nil, err
	}

	// we can only sign with one key
	from := keys.KeysById(id)[0].Entity
	if from == nil {
		return nil, fmt.Errorf("signing identity not found")
	}

	plaintextWriter, err := openpgp.Encrypt(armorWriter, to, from, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not set up openpgp encryption: %s", err)
	}

	plaintext, err := reference.PlaintextFromRefs(refs, nil)
	if err != nil {
		return nil, err
	}
	_, err = plaintextWriter.Write(plaintext)
	if err != nil {
		return nil, err
	}

	plaintextWriter.Close()
	armorWriter.Close()
	return buf.Bytes(), nil
}

// openPGPReference decrypts a given grant
func openPGPReference(grant []byte, keyring *config.OpenPGPSecret) ([]byte, error) {
	if keyring == nil {
		return nil, fmt.Errorf("cannot decrypt because no private key was provided")
	}

	// read local keyring
	to, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer(keyring.Data))
	block, err := armor.Decode(bytes.NewBuffer(grant))
	if err != nil {
		return nil, err
	}

	// consume grant message
	messageReader, err := openpgp.ReadMessage(block.Body, to,
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

	// verify that client knows signer
	if !messageReader.IsSigned {
		return nil, errors.New("grant is not signed")
	}
	if keys := to.KeysById(messageReader.SignedByKeyId); len(keys) == 0 {
		return nil, errors.New("unknown message signature")
	}

	return bs, nil

}

func OpenPGPReference(grant []byte, keyring *config.OpenPGPSecret, version int32) ([]*reference.Ref, error) {
	data, err := openPGPReference(grant, keyring)
	if err != nil {
		return nil, err
	}

	return reference.RefsFromPlaintext(data, version)
}
