# Release Notes

## Version

v1.1.0

## Release Date

2026-07-30

## Summary

- 27 features
- 3 bug fixes
- 14 documentation
- 2 refactoring
- 1 ci/cd

---

## Features

- Add history command
-  **makefile:**Add link target for Windows symlink live testing
-  **cmd:**Create plugins/ and rules/ directories on install
-  **cmd:**Implement rule add/remove/export/import commands
-  **runner:**Register custom rules and pass through unknown skip names
-  **config:**Add CustomRules field and Add/Remove/RulesDir helpers
-  **cmd:**Copy plugin descriptor to plugins dir on install, remove on uninstall
-  **runner:**Accept configPath parameter for plugin resolution
-  **plugins:**Resolve relative plugin paths against config dir
-  **config:**Add ConfigDir and PluginsDir helpers
-  **cmd:**Resolve absolute path on plugins install and show in list
-  **plugins:**Load plugin descriptor files from disk
-  **config:**Replace PluginEntry.Rules with Path to descriptor file
-  **cmd:**Add plugins uninstall command and --disabled flag
-  **runner:**Simplify plugin policy registration
-  **policy:**Add CustomPolicy for YAML-defined custom rules
-  **plugins:**Rewrite loader to create policies from YAML rules
-  **config:**Replace PluginEntry.Path with Rules and add CustomRuleDef type
-  **cmd:**Implement plugins install and list commands
-  **runner:**Wire plugin policies into engine
-  **plugins:**Implement plugin loading from config
-  **config:**Add plugin config support
-  **cmd:**Add rule skip subcommand
-  **runner:**Wire skip list with auto-clear on success
-  **engine:**Add skip list support to policy engine
-  **git:**Add local git config helpers
- Add initial roadmap documentation with built-in policies and CLI commands

## Bug Fixes

- Persist skip list across all hooks instead of auto-clearing
- Use global cfgFile in history command for proper config path handling
-  **makefile:**Update binary installation path and improve symbolic link creation

## Documentation

- Update docs to include history command
- Update roadmap to mark history command as completed
- Add log format design for history subcommand
- Add history subcommand to v1 roadmap
-  **examples:**Add single rule example for rule import
- Document rule add/remove/export/import commands
- Update ROADMAP.md and CONTRIBUTING.md for YAML plugin system
- Rename GUIDE.md to USER-GUIDE.md and document all commands
-  **examples:**Use path-based plugin entries in examples
-  **examples:**Update plugin and config examples to new YAML format
- Mark rule skip as implemented in roadmap
- Add rule skip to Commands at a Glance in README
- Update roadmap with built-in policies and CLI commands formatting
- Add v1.0.1 release notes

## Refactoring

-  **history:**Switch to global history file instead of per-repo hashed files, update Query/Clear accordingly
-  **config:**Update ConfigDir and HistoryDir to use DefaultConfigPath when configPath is empty

## CI/CD

- Add caching for faster build and golangci-lint config

## Contributors

- Mark Wayne Buncaras Menorca (47 commits)

## Full Changelog

v1.0.1...v1.1.0

