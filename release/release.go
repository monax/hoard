package release

import (
	"fmt"

	"errors"
)

// The purpose of this package is to capture version changes and change logs
// in a single location and use that data to generate releases and print
// changes to the command line

type Release struct {
	Version string
	Notes   string
}

// The releases described by version string and changes, newest release first.
// The current release is taken to be the first release in the slice, and its
// version determines the single authoritative version for the next release.
//
// To cut a new release add a release to the front of this slice then run the
// release tagging script: ./scripts/tag_release.sh
var hoardReleases = []Release{
	{
		Version: "1.1.1",
		Notes:   `Bump docker image Alpine Linux version to 3.8 and Go to 1.10.3`,
	},
	{
		Version: "1.1.0",
		Notes:   `Fix unsafe concurrent access of hash.Hash function in makeAddresser with sync.Pool`,
	},
	{
		Version: "1.0.2",
		Notes:   `Improve success/failure logging of LoggingStore.`,
	},
	{
		Version: "1.0.1",
		Notes:   `Add S3 integration test and include ca-certificates to Docker image so TLS (and S3) works.`,
	},
	{
		Version: "1.0.0",
		Notes: `Minor breaking change in that 'hoard init' becomes 'hoard config':
- 'hoard config' adds some niceties for printing JSON config for --env configuration source
- Added S3 'remote' credentials provider enabling credentials to be sourced from EC2 instance roles (note since [RemoteCredProvider()](https://github.com/aws/aws-sdk-go/blob/5a2026bfb28e86839f9fcc46523850319399006c/aws/defaults/defaults.go#L108) is used it also support ECS configuration via AWS_CONTAINER_CREDENTIALS_RELATIVE_URI and AWS_CONTAINER_CREDENTIALS_FULL_URI)
`,
	},
	{
		Version: "0.1.1",
		Notes:   `Include hoarctl in Docker image`,
	},
	{
		Version: "0.1.0",
		Notes: `Release adding environment config and docker image
- Adds --env flag to read JSON config from HOARD_JSON_CONFIG
- Add --json and --toml flags to 'hoard init' to generate JSON optionally
- Added alpine based docker image pushed on releases (that reads config from environment variable)
`,
	},
	{
		Version: "0.0.2",
		Notes: `Bug fix release for FileSystemStorage:
- Switch to URL and filesystem compliant base64 alphabet so some addresses do not target non-existent directories
- Create root directory for FileSystemStorage if it does not exist
`,
	},
	{
		Version: "0.0.1",
		Notes: `This is the first Hoard open source release and includes:
- Deterministic encryption scheme
- GRPC API for encryption, storage, and cleartext
- Memory, Filesystem, and S3 storage backends
- Configuration
- Hoar-Daemon hoard
- Hoar-Control hoarctl CLI
`,
	},
}

func Version() string {
	return hoardReleases[0].Version
}

func Notes() string {
	return hoardReleases[0].Notes
}

// Checks that a sequence of releases are monotonically decreasing with each
// version being a simple major, minor, or patch bump of its successor in the
// slice
func AssertReleasesUniqueAndMonotonic(releases []Release) error {
	if len(releases) == 0 {
		return errors.New("at least one release must be defined")
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
				return fmt.Errorf("minor and patch versions must be reset to "+
					"0 after a major bump, but they are not in %s -> %s",
					releases[i].Version, releases[i-1].Version)
			}
		} else if maj == lMaj {
			// Same major number
			if min == lMin+1 {
				// Minor bump so patch version must be reset
				if pat != 0 {
					return fmt.Errorf("patch version must be reset to "+
						"0 after a minor bump, but they are not in %s -> %s",
						releases[i].Version, releases[i-1].Version)
				}
			} else if min == lMin {
				// Same minor number so must be patch bump to be valid
				if pat != lPat+1 {
					return fmt.Errorf("consecutive patch versions must be equal "+
						"or incremented by 1, but they are not in %s -> %s",
						releases[i].Version, releases[i-1].Version)
				}
			} else {
				return fmt.Errorf("consecutive minor versions must be equal or "+
					"incremented by 1, but they are not in %s -> %s",
					releases[i].Version, releases[i-1].Version)
			}
		} else {
			return fmt.Errorf("consecutive major versions must be equal or "+
				"incremented by 1, but they are not in  %s -> %s",
				releases[i].Version, releases[i-1].Version)
		}

		maj, min, pat = lMaj, lMin, lPat
	}
	return nil
}
