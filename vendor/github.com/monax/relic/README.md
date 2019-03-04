# Relic
Relic is a library to help with versioning your projects by storing release metadata 
and versions as code.

## Purpose
Relic allows you define your project version history in a declarative style by
defining a `History` object somewhere in your project whose methods allow you to
declare releases defined by a version number and release note. It ensures your releases
have monotonically increasing unique versions. 

Relic can generate the current version and a complete [changelog](CHANGELOG.md) using this information.

Relic can be used by CI systems to output version numbers for tagging artefacts and automatically pushing releases,
such as [goreleaser](https://github.com/goreleaser/goreleaser).

By keeping the changelog with the version they are synchronised and you are reminded to produce 
the changelog.

## Usage
```go
// Add file to your project in which to record your projects revision history
package project

import (
	"fmt"
	"text/template"
	"github.com/monax/relic"
)

// Create a global variable in which to store your project history.
// MustDeclareReleases allows you to declare your releases by specifying a version and release note
// for each release. To add a new release just insert it at the top.
var History relic.ImmutableHistory = relic.NewHistory("Relic", "https://github.com/monax/relic").
	MustDeclareReleases(
		"2.0.0 - 2018-08-15",
		`### Changed
- Versions must start from 0.0.1 (0.0.0 is reserved for unreleased)
- Default changelog format follows https://keepachangelog.com/en/1.0.0/
- NewHistory takes second parameter for project URL
- Dropped getters from Version since already passed by value so immutable

### Added
- Optional top (most recent) release may be provided with empty Version with (via empty string in DeclareReleases) whereby its notes will be listed under 'Unreleased'
- Optional date can be appended to version using the exact format <major.minor.patch - YYYY-MM-DD> e.g. '5.4.2 - 2018-08-14'
- Default changelog format footnote references standard github compare links to see commits between version tags
`,
		"1.1.0",
		`Add ImmutableHistory and tweak suggested usage docs`,
		"1.0.1",
		`Documentation fixes and typos`,
		"1.0.0",
		`Minor improvements:
- Rename DeclareReleases to DeclareReleases (breaking API change)
- Add sample snippet to readme
- Sign version tags
`,
		"0.0.1",
		`First release of Relic extracted from various initial projects, it can:
- Generate changelogs
- Print the current version
- Ensure valid semantic version numbers
`)

func PrintReleaseInfo() {
	// Print the current version
	fmt.Printf("%s (Version: %v)\n", History.Project(), History.CurrentVersion().String())
	// Print the complete changelog 
	fmt.Println(History.Changelog())
	// Get specific release
	release, err := History.Release("0.0.1")
	if err != nil {
		panic(err)
	}
	// Print major version of release
	fmt.Printf("Release Major version: %v", release.Version.Major)
}

// You can also define histories with a custom template
var ProjectWithTemplate = relic.NewHistory("Test Project", "https://github.com/test/project").
		WithChangelogTemplate(template.Must(template.New("tests").
			Parse("{{range .Releases}}{{$.Name}} (v{{.Version}}): {{.Notes}}\n{{end}}"))).
		MustDeclareReleases(
			// Releases may optionally have a date which is included in the changelog
			"0.1.0 - 2016-07-12",
			"Basic functionality",
			"0.0.2",
			"Build scripts",
			"0.0.1",
			"Proof of concept",
		)

```

See Relic's own [`project` package](project/releases.go) and [Makefile](Makefile) for suggested usage within a project.

## Dependencies
Go standard library and tooling plus Make and Bash for builds (but not required).