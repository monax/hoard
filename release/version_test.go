package release

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	// Check parsed version numbers match the version string
	assert.Equal(t, Version(), fmt.Sprintf("%d.%d.%d",
		Major(), Minor(), Patch()))
}

func TestParseVersion(t *testing.T) {
	maj, min, pat, err := ParseVersion("23.255.1")
	assert.NoError(t, err)
	assert.Equal(t, uint8(23), maj)
	assert.Equal(t, uint8(255), min)
	assert.Equal(t, uint8(1), pat)

	maj, min, pat, err = ParseVersion("2312.3.1")
	assert.Error(t, err)

	maj, min, pat, err = ParseVersion("231.256.1")
	assert.Error(t, err)

	maj, min, pat, err = ParseVersion("231.3.5645")
	assert.Error(t, err)
}
