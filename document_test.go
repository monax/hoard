package hoard

import (
	"testing"

	"github.com/monax/hoard/v7/meta"
	"github.com/monax/hoard/v7/test/helpers"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestObjectStructure(t *testing.T) {
	doc := helpers.ReadDocument(t, "markdownfile.md")
	data, err := proto.Marshal(doc)
	require.NoError(t, err)
	newDoc := &meta.Document{}
	err = proto.Unmarshal(data, newDoc)
	require.NoError(t, err)
}

func TestProperlyEncoding(t *testing.T) {
	doc := helpers.ReadDocument(t, "markdownfile.md")
	out, err := proto.Marshal(doc)
	require.NoError(t, err)

	doc1 := &meta.Document{}
	err = proto.Unmarshal(out, doc1)
	require.NoError(t, err)
	require.NotEqual(t, []byte{}, doc1.Data)
	require.Equal(t, "markdownfile.md", doc1.Meta.Name)
}

func TestProperlyDecoding(t *testing.T) {
	doc := helpers.ReadDocument(t, "markdownfile.md")
	input, err := proto.Marshal(doc)
	require.NoError(t, err)

	doc1, err := parseIntoDocument(input)
	require.NoError(t, err)
	require.Equal(t, 0, len(doc1.Meta.Tags))
	require.Equal(t, doc.Meta.Name, doc1.Meta.Name)
	require.Equal(t, doc.Meta.MimeType, doc1.Meta.MimeType)
	require.Equal(t, doc.Data, doc1.Data)
}

func TestImproperlyDecoding(t *testing.T) {
	_, err := parseIntoDocument([]byte(`11111`))
	require.Error(t, err)

	_, err = parseIntoDocument([]byte(`111`))
	require.Error(t, err)
}

// TODO: Remove this when...
// [Casey] We can remove that fallback if we ensure that every hoard object conforms to the new encoding system rather than the old.
func TestProperlyDecodingLegacy(t *testing.T) {
	docRaw, err := helpers.ReadFixture("rawObject.buffer")
	require.NoError(t, err)

	doc, err := legacyParseIntoDocument(docRaw)
	require.NoError(t, err)
	require.Equal(t, "ContractorProseTemplate.md", doc.Meta.Name)
	require.Equal(t, "6A0D529822970791495B7D239BF21A365F6237DE", doc.Meta.Agreement)
	require.Equal(t, "text/markdown", doc.Meta.MimeType)
	require.Equal(t, "scalia", doc.Meta.AssemblyEngine)
	require.Equal(t, 1, len(doc.Meta.Tags))
	require.Equal(t, "legacy_encoding", doc.Meta.Tags[0])
	require.NotEqual(t, 0, len(doc.Data))
}

// TODO: Remove this when...
// [Casey] We can remove that fallback if we ensure that every hoard object conforms to the new encoding system rather than the old.
func TestProperlyDecodingLegacyWithFallThru(t *testing.T) {
	docRaw, err := helpers.ReadFixture("rawObject.buffer")
	require.NoError(t, err)

	doc, err := parseIntoDocument(docRaw)
	require.NoError(t, err)
	require.Equal(t, "ContractorProseTemplate.md", doc.Meta.Name)
	require.Equal(t, "6A0D529822970791495B7D239BF21A365F6237DE", doc.Meta.Agreement)
	require.Equal(t, "text/markdown", doc.Meta.MimeType)
	require.Equal(t, "scalia", doc.Meta.AssemblyEngine)
	require.Equal(t, 1, len(doc.Meta.Tags))
	require.Equal(t, "legacy_encoding", doc.Meta.Tags[0])
	require.NotEqual(t, 0, len(doc.Data))
}

func TestLegacyParseOnError(t *testing.T) {
	docRaw, err := helpers.ReadFixture("rawObject.buffer.bad")
	require.NoError(t, err)

	_, err = legacyParseIntoDocument(docRaw)
	require.Error(t, err)
}
