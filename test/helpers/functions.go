package helpers

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func ReadFixture(filename string) ([]byte, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if ok {
		return ioutil.ReadFile(filepath.Join(filepath.Dir(thisFile), "..", "fixtures", filename))
	} else {
		return []byte{}, fmt.Errorf("cannot find the functions.go file in your filepath")
	}
}

// func ReadDocument(t *testing.T, filename string) *meta.Document {
// 	docRaw, err := ReadFixture(filename)
// 	require.NoError(t, err)
// 	doc := &meta.Document{
// 		Meta: &meta.Meta{
// 			Name: filename,
// 		},
// 		Data: docRaw,
// 	}

// 	switch filepath.Ext(filename) {
// 	case ".docx":
// 		doc.Meta.MimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
// 	case ".md":
// 		doc.Meta.MimeType = "text/markdown"
// 	case ".html":
// 		doc.Meta.MimeType = "text/html"
// 	case ".json":
// 		doc.Meta.MimeType = "application/json"
// 	case ".pdf":
// 		doc.Meta.MimeType = "application/pdf"
// 	case ".rtf":
// 		doc.Meta.MimeType = "application/rtf"
// 	case ".doc":
// 		doc.Meta.MimeType = "application/msword"
// 	case ".odt":
// 		doc.Meta.MimeType = "application/vnd.oasis.opendocument.text"
// 	}

// 	return doc
// }

func WriteResult(t *testing.T, data []byte) {
	_, thisFile, _, ok := runtime.Caller(0)
	if ok {
		err := ioutil.WriteFile(filepath.Join(filepath.Dir(thisFile), "..", "results", fmt.Sprintf("Result_%s.docx", t.Name())), data, 0644)
		require.NoError(t, err)
	}
}

func WriteResultWithExt(t *testing.T, data []byte, postfix string) {
	_, thisFile, _, ok := runtime.Caller(0)
	if ok {
		err := ioutil.WriteFile(filepath.Join(filepath.Dir(thisFile), "..", "results", fmt.Sprintf("Result_%s_%s", t.Name(), postfix)), data, 0644)
		require.NoError(t, err)
	}
}
