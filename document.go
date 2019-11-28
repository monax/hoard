package hoard

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/monax/hoard/v6/api"
	"github.com/monax/hoard/v6/grant"
	"github.com/monax/hoard/v6/meta"
)

// Document extensions to hoard providing conventions for storing documents with their metadata

// GetDocument retrieves a document from hoard and parses
// it into a document struct.
// NOTE: if this schema changes hoard will break.
func GetDocument(gs GrantService, grant *grant.Grant) (*meta.Document, []byte, error) {
	ref, err := gs.Unseal(grant)
	if err != nil {
		return nil, nil, err
	}

	data, err := gs.Get(ref)
	if err != nil {
		return nil, nil, err
	}

	doc, err := parseIntoDocument(data)
	if err != nil {
		return nil, nil, err
	}

	return doc, ref.GetSalt(), nil
}

// PostDocument is given a document struct which is
// then parsed into a document object which matches the encoding
// system established.
// NOTE: if this schema changes hoard will break.
//
// This function puts and seals the document into hoard's store
// and returns back the grant which is given from hoard.
func PutDocument(gs GrantService, pgsm *api.PlaintextAndGrantSpecAndMeta) (*grant.Grant, error) {
	out, err := proto.Marshal(&meta.Document{
		Meta: pgsm.Meta,
		Data: pgsm.PlaintextAndGrantSpec.Plaintext.Data,
	})
	if err != nil {
		return nil, err
	}

	ref, err := gs.Put(out, pgsm.PlaintextAndGrantSpec.Plaintext.Salt)
	if err != nil {
		return nil, err
	}

	return gs.Seal(ref, pgsm.PlaintextAndGrantSpec.GrantSpec)
}

// parseIntoDocument converts a hoard inode object into a parsed document
func parseIntoDocument(input []byte) (*meta.Document, error) {
	document := &meta.Document{}
	err := proto.Unmarshal(input, document)
	if err != nil { // SHOULD NOT BE REMOVED until at least December, 2022.
		return legacyParseIntoDocument(input)
	}
	return document, nil
}

// legacyDocumentMeta is a necessary type due to changes in the structure
// of the meta data used and the legacy implementations
// predominantly the change from Mime -> MimeType along with AssemblyType
//  -> AssemblyEngine.
type legacyDocumentMeta struct {
	Name         string
	Mime         string
	Agreement    string
	AssemblyType string
}

// TODO: [Silas] query Casey on the plan implied by the below...
// legacyParseIntoDocument is a dead simple function which receives
// a raw hoard inode and parses that into the structs which can be
// leveraged within hoard. This only remains for very old documents
// and SHOULD NOT BE REMOVED until at least December, 2022.
func legacyParseIntoDocument(input []byte) (*meta.Document, error) {
	buff := bytes.NewBuffer(input)
	if len(buff.Bytes()) < 4 {
		return nil, fmt.Errorf("unexpected input length for document")
	}

	lenOfMeta := binary.BigEndian.Uint32(buff.Next(4))

	var tmpDocMeta = &legacyDocumentMeta{}
	jsonBytes := buff.Next(int(lenOfMeta))
	err := json.Unmarshal(jsonBytes, tmpDocMeta)
	if err != nil {
		return &meta.Document{}, err
	}

	return &meta.Document{
		Meta: &meta.Meta{
			Name:           tmpDocMeta.Name,
			MimeType:       tmpDocMeta.Mime,
			Agreement:      tmpDocMeta.Agreement,
			AssemblyEngine: tmpDocMeta.AssemblyType,
			Tags:           []string{"legacy_encoding"},
		},
		Data: buff.Bytes(),
	}, nil
}
