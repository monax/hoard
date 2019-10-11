package helpers

import (
	"github.com/monax/hoard/v5/meta"
)

// -------------------------------------------------------
//
// types
//
// -------------------------------------------------------
type TestDocData struct {
	Type    string `json:"type"`
	RawData []byte `json:"data"`
}

type DocumentTest struct {
	Meta meta.Meta   `json:"meta"`
	Data TestDocData `json:"data"`
}

// -------------------------------------------------------
//
// constants
//
// -------------------------------------------------------
const FailureTemplate = `stuff {{ fail "You have failed" }} otherstuff`
