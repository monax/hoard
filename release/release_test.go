package release

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHoardReleasesUniqueAndMonotonic(t *testing.T) {
	err := AssertReleasesUniqueAndMonotonic(hoardReleases)
	assert.NoError(t, err)
}

func TestAssertReleasesUniqueAndMonotonic(t *testing.T) {
	releases := []Release{
		{
			Version: "2.1.1",
			Changes: `Everything fixed`,
		},
		{
			Version: "2.1.0",
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
	}
	err := AssertReleasesUniqueAndMonotonic(releases)
	assert.Error(t, err)

	releases = []Release{
		{
			Version: "2.1.1",
			Changes: `Everything fixed`,
		},
		{
			Version: "2.1.0",
			Changes: `Everything broken`,
		},
		{
			Version: "2.0.0",
			Changes: `Wonderful things were achieved`,
		},
		{
			Version: "1.0.0",
			Changes: `Wonderful things were achieved`,
		},
		{
			Version: "0.0.2",
			Changes: `Wonderful things were achieved`,
		},
		{
			Version: "0.0.1",
			Changes: `Marvelous advances were made`,
		},
	}
	err = AssertReleasesUniqueAndMonotonic(releases)
	assert.NoError(t, err)

	releases = []Release{
		{
			Version: "1.0.3",
			Changes: `Wonderful things were achieved`,
		},
		{
			Version: "0.0.2",
			Changes: `Wonderful things were achieved`,
		},
		{
			Version: "0.0.1",
			Changes: `Marvelous advances were made`,
		},
	}
	err = AssertReleasesUniqueAndMonotonic(releases)
	assert.Error(t, err)

	releases = []Release{
		{
			Version: "0.1.3",
			Changes: `Wonderful things were achieved`,
		},
		{
			Version: "0.0.2",
			Changes: `Wonderful things were achieved`,
		},
		{
			Version: "0.0.1",
			Changes: `Marvelous advances were made`,
		},
	}
	err = AssertReleasesUniqueAndMonotonic(releases)

	assert.Error(t, err)
	releases = []Release{
		{
			Version: "0.0.3",
			Changes: `Wonderful things were achieved`,
		},
		{
			Version: "0.0.2",
			Changes: `Wonderful things were achieved`,
		},
		{
			Version: "0.0.1",
			Changes: `Marvelous advances were made`,
		},
	}
	err = AssertReleasesUniqueAndMonotonic(releases)
	assert.NoError(t, err)

	releases = []Release{
		{
			Version: "0.0.2",
			Changes: `Wonderful things were achieved`,
		},
		{
			Version: "0.0.3",
			Changes: `Wonderful things were achieved`,
		},
		{
			Version: "0.0.1",
			Changes: `Marvelous advances were made`,
		},
	}
	err = AssertReleasesUniqueAndMonotonic(releases)
	assert.Error(t, err)
}
