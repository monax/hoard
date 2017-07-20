package release

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangelogForReleases(t *testing.T) {
	changelog, err := changelogForReleases(
		[]Release{
			{
				Version: "0.1.0",
				Changes: `Everything broken`,
			},
			{
				Version: "0.0.2",
				Changes: `Wonderful things were achieved`,
			},
			{
				Version: "0.0.1",
				Changes: `Marvelous advances were made`,
			},
		})
	assert.NoError(t, err)

	assert.Equal(t, `# Hoard Changelog
## Version 0.1.0
Everything broken

## Version 0.0.2
Wonderful things were achieved

## Version 0.0.1
Marvelous advances were made
`, changelog)
}
