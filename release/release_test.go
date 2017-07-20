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
			Notes:   `Everything fixed`,
		},
		{
			Version: "2.1.0",
			Notes:   `Everything broken`,
		},
		{
			Version: "0.0.2",
			Notes:   `Wonderful things were achieved`,
		},
		{
			Version: "0.0.1",
			Notes:   `Marvelous advances were made`,
		},
	}
	err := AssertReleasesUniqueAndMonotonic(releases)
	assert.Error(t, err)

	releases = []Release{
		{
			Version: "2.1.1",
			Notes:   `Everything fixed`,
		},
		{
			Version: "2.1.0",
			Notes:   `Everything broken`,
		},
		{
			Version: "2.0.0",
			Notes:   `Wonderful things were achieved`,
		},
		{
			Version: "1.0.0",
			Notes:   `Wonderful things were achieved`,
		},
		{
			Version: "0.0.2",
			Notes:   `Wonderful things were achieved`,
		},
		{
			Version: "0.0.1",
			Notes:   `Marvelous advances were made`,
		},
	}
	err = AssertReleasesUniqueAndMonotonic(releases)
	assert.NoError(t, err)

	releases = []Release{
		{
			Version: "1.0.3",
			Notes:   `Wonderful things were achieved`,
		},
		{
			Version: "0.0.2",
			Notes:   `Wonderful things were achieved`,
		},
		{
			Version: "0.0.1",
			Notes:   `Marvelous advances were made`,
		},
	}
	err = AssertReleasesUniqueAndMonotonic(releases)
	assert.Error(t, err)

	releases = []Release{
		{
			Version: "0.1.3",
			Notes:   `Wonderful things were achieved`,
		},
		{
			Version: "0.0.2",
			Notes:   `Wonderful things were achieved`,
		},
		{
			Version: "0.0.1",
			Notes:   `Marvelous advances were made`,
		},
	}
	err = AssertReleasesUniqueAndMonotonic(releases)

	assert.Error(t, err)
	releases = []Release{
		{
			Version: "0.0.3",
			Notes:   `Wonderful things were achieved`,
		},
		{
			Version: "0.0.2",
			Notes:   `Wonderful things were achieved`,
		},
		{
			Version: "0.0.1",
			Notes:   `Marvelous advances were made`,
		},
	}
	err = AssertReleasesUniqueAndMonotonic(releases)
	assert.NoError(t, err)

	releases = []Release{
		{
			Version: "0.0.2",
			Notes:   `Wonderful things were achieved`,
		},
		{
			Version: "0.0.3",
			Notes:   `Wonderful things were achieved`,
		},
		{
			Version: "0.0.1",
			Notes:   `Marvelous advances were made`,
		},
	}
	err = AssertReleasesUniqueAndMonotonic(releases)
	assert.Error(t, err)
}
