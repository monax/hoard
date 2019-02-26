### Changed
- [PROTO] Renamed symmetric grant SecretID to PublicID
- [PROTO] Renamed openpgp grant ID to PrivateID

### Fixed
- [GRANTS] Throw an exception if symmetric secret for ID cannot be found

### Added
- [NODEJS] Added integration tests including test for symmetric secrets
- [GRANTS] Added openpgp grants example
- [CLI] Added ability to configure secrets on command line with hoard config <config> --secret
- 

