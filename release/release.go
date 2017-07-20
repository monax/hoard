package release

import (
	"fmt"

	"github.com/pkg/errors"
)

// The purpose of this package is to capture version changes and change logs
// in a single location and use that data to generate releases and print
// changes to the command line

type Release struct {
	Version string
	Changes string
}

// The releases described by version string and changes, newest release first.
// The current release is taken to be the first release in the slice, and its
// version determines the single authoritative version for the next release.
//
// To cut a new release add a release to the front of this slice then run the
// release tagging script: ./scripts/tag_release.sh
var hoardReleases = []Release{
	{
		Version: "0.0.1",
		Changes: `This is the first Hoard open source release and includes:
- Deterministic encryption scheme
- GRPC API for encryption, storage, and cleartext
- Memory, Filesystem, and S3 storage backends
- Configuration
- Hoar-Daemon hoard
- Hoar-Control hoarctl CLI
`,
	},
}

// Checks that a sequence of releases are monotonically decreasing with each
// version being a simple major, minor, or patch bump of its successor in the
// slice
func AssertReleasesUniqueAndMonotonic(releases []Release) error {
	if len(releases) == 0 {
		return errors.New("At least one release must be defined")
	}
	maj, min, pat, err := ParseVersion(releases[0].Version)
	if err != nil {
		return err
	}
	for i := 1; i < len(releases); i++ {
		// The numbers of the lower version (expect descending sort)
		lMaj, lMin, lPat, err := ParseVersion(releases[i].Version)
		if err != nil {
			return err
		}
		// Check versions are consecutive
		if maj == lMaj+1 {
			// Major bump, so minor and patch versions must be reset
			if min != 0 || pat != 0 {
				return fmt.Errorf("Minor and patch versions must be reset to "+
					"0 after a major bump, but they are not in %s -> %s",
					releases[i].Version, releases[i-1].Version)
			}
		} else if maj == lMaj {
			// Same major number
			if min == lMin+1 {
				// Minor bump so patch version must be reset
				if pat != 0 {
					return fmt.Errorf("Patch version must be reset to "+
						"0 after a minor bump, but they are not in %s -> %s",
						releases[i].Version, releases[i-1].Version)
				}
			} else if min == lMin {
				// Same minor number so must be patch bump to be valid
				if pat != lPat+1 {
					return fmt.Errorf("Consecutive patch versions must be equal "+
						"or incremented by 1, but they are not in %s -> %s",
						releases[i].Version, releases[i-1].Version)
				}
			} else {
				return fmt.Errorf("Consecutive minor versions must be equal or "+
					"incremented by 1, but they are not in %s -> %s",
					releases[i].Version, releases[i-1].Version)
			}
		} else {
			return fmt.Errorf("Consecutive major versions must be equal or "+
				"incremented by 1, but they are not in  %s -> %s",
				releases[i].Version, releases[i-1].Version)
		}

		maj, min, pat = lMaj, lMin, lPat
	}
	return nil
}
