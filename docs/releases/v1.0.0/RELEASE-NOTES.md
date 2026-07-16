# Release Notes

## Version

v1.0.0

## Release Date

2026-07-16

---

## Features

- Add initial changelog for project documentation
- Add comprehensive user guide for git-policy
- Add example config, sample hooks, and test data
- Add dedicated packages for sync.
- Add dedicated packages for plugins.
- Add dedicated packages for  branch.
- Add dedicated packages for commitmsg.
- Add dedicated packages for files.
- Add dedicated packages for secrets.
-  **cmd:**Add CLI commands (install, uninstall, run, doctor, validate, version, rule)
-  **runner:**Wire Git context into policy engine
-  **hook:**Add global Git hook installer and uninstaller
-  **engine:**Add policy execution engine with disabled-rule support
-  **policy:**Implement Policy interface and 6 built-in rules
-  **git:**Add Git CLI wrappers
- Add logger and path utilities
-  **config:**Add YAML config loading and defaults
- Scaffold Go module with Makefile
- Initial commit
- Add GitHub Actions workflow for release process

## Bug Fixes

- Rename release job to build and update artifact handling in workflow
- Update build command to correctly set version variable in release workflow
- Update Makefile to handle binary extension and versioning
- Update version to development stage

## Documentation

- Add release notes for version 1.0.0
- Remove link to PLAN.md from README.md and clarify roadmap section
- Update license badge to Apache 2.0 and add limitations and community sections in README.md
- Update LICENSE to Apache License 2.0 with detailed terms and conditions
- Update README.md to correct project name and enhance command examples
- Enhance README.md with additional alignment and badge formatting
- Add comprehensive README.md to outline git-policy features and usage
- Create RELEASE-NOTES.md to document project changes and guidelines
- Add SECURITY.md to outline security policy and supported versions
- Add MIT License to the project
- Add pull request template to standardize contributions
- Add FUNDING.yml to support project funding options
- Add issue templates for bug reports and feature requests
- Add Code of Conduct to establish community standards and expectations
- Add contributing guidelines and project structure documentation
- Add SUPPORT.md with documentation, issue reporting, discussions, and security guidelines
- Update links in Community section to remove .github prefix

## Refactoring

- Remove unused commit message retrieval in Run function

## CI/CD

- Rename build job to release and consolidate artifact building steps
- Add GitHub Actions workflow for testing Go code

## Maintenance

- Add .gitignore file to exclude build artifacts and IDE settings

## Contributors

- Mark Wayne Buncaras Menorca (45 commits)
- Mark Wayne Menorca (4 commits)

