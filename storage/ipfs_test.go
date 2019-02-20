package storage

import (
	"net/http"
	"testing"

	"encoding/base64"

	"github.com/stretchr/testify/assert"
)

func TestIPFSStore(t *testing.T) {
	srv := &http.Server{Addr: ":5001"}
	go srv.ListenAndServe()
	inv, err := NewIPFSStore("http://:5001", base64.URLEncoding)
	assert.NoError(t, err)
	assert.Equal(t, "http://:5001/api/v0", inv.host)
	srv.Close()
}
