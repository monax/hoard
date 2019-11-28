package hoard

import (
	"io"
	"testing"

	"github.com/monax/hoard/v6/api"
	"github.com/monax/hoard/v6/grant"
	"github.com/monax/hoard/v6/meta"
	"github.com/monax/hoard/v6/test/helpers"
	"github.com/stretchr/testify/require"
)

type duplexer struct {
	plaintextGrantSpecs chan *api.PlaintextAndGrantSpecAndMeta
	sent                int
	received            int
}

func newDuplexer() *duplexer {
	return &duplexer{plaintextGrantSpecs: make(chan *api.PlaintextAndGrantSpecAndMeta, 100)}
}

func (d *duplexer) Recv() (*api.PlaintextAndGrantSpecAndMeta, error) {
	msg := <-d.plaintextGrantSpecs
	if msg == nil {
		return nil, io.EOF
	}
	d.received++
	return msg, nil
}

func (d *duplexer) Send(msg *api.PlaintextAndGrantSpecAndMeta) error {
	d.plaintextGrantSpecs <- msg
	if msg != nil {
		d.sent++
	}
	return nil
}

func (d *duplexer) Close() {
	d.plaintextGrantSpecs <- nil
}

func TestReceiveDocumentAndGrant(t *testing.T) {
	d := newDuplexer()
	spec := &grant.Spec{Plaintext: &grant.PlaintextSpec{}}
	chunkSize := 100
	data := []byte(helpers.LongText)
	salt := []byte("me hearties")

	doc := &meta.Document{
		Meta: &meta.Meta{
			Name:     "Storia dei Musulmani di Sicilia ",
			MimeType: "text/plain",
		},
		Data: data,
	}

	err := SendDocumentAndGrantSpec(d, doc, salt, spec, chunkSize)
	require.NoError(t, err)

	// meta + spec + salt + data chunks
	msgCount := 3 + (len(data)+chunkSize-1)/chunkSize
	require.Equal(t, msgCount, d.sent)
	// Signal EOF in duplexer
	d.Close()
	pgsm, err := ReceiveDocumentAndGrantSpec(d)
	require.Equal(t, msgCount, d.received)
	require.NoError(t, err)
	require.Equal(t, doc, &meta.Document{
		Meta: pgsm.Meta,
		Data: pgsm.PlaintextAndGrantSpec.Plaintext.Data,
	})
	require.Equal(t, salt, pgsm.PlaintextAndGrantSpec.Plaintext.Salt)
	require.Equal(t, spec, pgsm.PlaintextAndGrantSpec.GrantSpec)
}
