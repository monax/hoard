# [Relic](https://github.com/monax/relic) Changelog
## [2.0.0] - 2018-08-15
### Changed
- Versions must start from 0.0.1 (0.0.0 is reserved for unreleased)
- Default changelog format follows https://keepachangelog.com/en/1.0.0/
- NewHistory takes second parameter for project URL
- Dropped getters from Version since already passed by value so immutable

### Added
- Optional top (most recent) release may be provided with empty Version with (via empty string in DeclareReleases) whereby its notes will be listed under 'Unreleased'
- Optional date can be appended to version using the exact format <major.minor.patch - YYYY-MM-DD> e.g. '5.4.2 - 2018-08-14'
- Default changelog format footnote references standard github compare links to see commits between version tags


## [1.1.0]
Add ImmutableHistory and tweak suggested usage docs

## [1.0.1]
Documentation fixes and typos

## [1.0.0]
Minor improvements:
- Rename DeclareReleases to DeclareReleases (breaking API change)
- Add sample snippet to readme
- Sign version tags


## [0.0.1]
First release of Relic extracted from various initial projects, it can:
- Generate changelogs
- Print the current version
- Ensure valid semantic version numbers


[2.0.0]: https://github.com/monax/relic/compare/v1.1.0...v2.0.0
[1.1.0]: https://github.com/monax/relic/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/monax/relic/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/monax/relic/compare/v0.0.1...v1.0.0
[0.0.1]: https://github.com/monax/relic/commits/v0.0.1
