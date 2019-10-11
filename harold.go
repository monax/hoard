package hoard

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/monax/hoard/v6/grant"
	"github.com/monax/hoard/v6/meta"
)

// Harold is an extension to the hoard interface to allow
// stateful retrieval of document objects

// GetDocument retrieves a document from hoard and parses
// it into a document struct which harold will use to perform
// it's assembly functionality. During the document parsing
// sequence we parse the document based upon the API
// encoding schema.
// NOTE: if this schema changes harold will break.
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

// PostDocument is given a harold document struct which is
// then parsed into a document object which matches the encoding
// system established by Monax API. NOTE that if Monax API changes
// this encoding scheme this will break everything for harold.
//
// This function puts and seals the document into hoard's store
// and returns back the grant which is given from hoard.
func PostDocument(gs GrantService, document *meta.Document, spec *grant.Spec, salt []byte) (*grant.Grant, error) {
	out, err := proto.Marshal(document)
	if err != nil {
		return nil, err
	}

	ref, err := gs.Put(out, salt)
	if err != nil {
		return nil, err
	}

	return gs.Seal(ref, spec)
}

// PutDocument currently mirrors PostDocument in that we do not
// have hoard functionality to remove documents. However, in the
// future when we do have the ability to remove documents then
// the PutDocument function will remove the given inputGrant in
// favor of the input document by saving the input document into
// hoard and removing the input grant. It returns the given grant.
func PutDocument(gs GrantService, document *meta.Document, spec *grant.Spec, salt []byte, replace *grant.Grant) (*grant.Grant, error) {
	ref, err := gs.Unseal(replace)
	if err != nil {
		return nil, err
	}

	err = gs.Delete(ref.GetAddress())
	if err != nil {
		return nil, err
	}

	return PostDocument(gs, document, spec, salt)
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
// of the meta data used by  harold and the legacy implementations
// predominantly the change from Mime -> MimeType along with AssemblyType
//  -> AssemblyEngine.
type legacyDocumentMeta struct {
	Name         string
	Mime         string
	Agreement    string
	AssemblyType string
}

// legacyParseIntoDocument is a dead simple function which receives
// a raw hoard inode and parses that into the structs which can be
// leveraged within harold. This only remains for very old documents
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
