# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.1] - 2025-05-26

### Added
- Progress display system for future TUI migration
- Enhanced visual feedback during execution

### Documentation
- Updated documentation for v0.1.0 release

## [0.1.0] - 2025-05-25

### Added
- Initial release with core functionality
- Shell executor for running scripts and commands
- SQL executor with PostgreSQL support
- Transaction mode support (none, each, all) for SQL executor
- YAML-based execution plans with dependency resolution
- State management for resumable execution
- CLI commands: execute, validate, status, reset
- Environment variable expansion in configurations
- Platform-specific file selection
- Colorful ASCII art logo with NO_COLOR support
- Multi-database example demonstrating complex setups
- GitHub Actions release workflow with goreleaser
- Comprehensive shell completion with descriptions
- Cross-platform support (Linux, macOS, Windows)

### Changed
- Restructured CLI to follow Cobra best practices
- Changed ExecutorConfig to map[string]interface{} for flexibility
- Improved executor initialization to support custom configurations

### Fixed
- Test compatibility with new executor configuration structure
- SQL file naming in multi-database example

[Unreleased]: https://github.com/SphereStacking/plexr/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/SphereStacking/plexr/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/SphereStacking/plexr/releases/tag/v0.1.0