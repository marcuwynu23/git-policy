# Release Notes

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

---

## [Unreleased]

### Added
- `history` command to track git-policy runs across all repositories
- History configuration options (`history.enabled`, `history.maxRecords`)
- History filters: `--limit`, `--repo`, `--status`, `--clear`
- Global history storage in `~/.config/git-policy/history/history.jsonl`

### Changed
- Updated `ConfigDir` in internal config to use default config path when configPath is empty

### Deprecated
- 

### Removed
- 

### Fixed
- History command now properly uses the global config path specified via `--config`

### Security
- 

---

## [1.0.0] - YYYY-MM-DD

### Added
- Initial release of the project
- Core features implemented:
  - Feature A
  - Feature B

### Changed
- 

### Fixed
- 

### Security
- 

---

## Release Guidelines

### Versioning
This project follows **Semantic Versioning (SemVer)**:
- **MAJOR**: incompatible API changes
- **MINOR**: backwards-compatible features
- **PATCH**: backwards-compatible bug fixes

---

## Notes

- Include links to issues or PRs when possible:
  - Example: Fixed login crash ([#42](https://example.com/issues/42))
- Highlight breaking changes clearly under a "Breaking Changes" section if needed.
- Keep entries concise and user-focused.

---

## Contributors

Thanks to everyone who contributed to this release:

- Name or GitHub handle