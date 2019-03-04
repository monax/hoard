### Changed
- Versions must start from 0.0.1 (0.0.0 is reserved for unreleased)
- Default changelog format follows https://keepachangelog.com/en/1.0.0/
- NewHistory takes second parameter for project URL
- Dropped getters from Version since already passed by value so immutable

### Added
- Optional top (most recent) release may be provided with empty Version with (via empty string in DeclareReleases) whereby its notes will be listed under 'Unreleased'
- Optional date can be appended to version using the exact format <major.minor.patch - YYYY-MM-DD> e.g. '5.4.2 - 2018-08-14'
- Default changelog format footnote references standard github compare links to see commits between version tags

