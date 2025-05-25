# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- SQL executor for PostgreSQL database operations
- Transaction mode support (none, each, all) for SQL executor
- Multi-database example demonstrating complex setups
- Colorful ASCII art logo
- GitHub Actions release workflow with goreleaser
- Comprehensive shell completion with descriptions
- State management for resumable execution
- Dependency resolution between steps

### Changed
- Restructured CLI to follow Cobra best practices
- Changed ExecutorConfig to map[string]interface{} for flexibility
- Improved executor initialization to support custom configurations
- Updated installation documentation for binary releases

### Fixed
- Test compatibility with new executor configuration structure
- SQL file naming in multi-database example

## [0.1.0] - 2024-01-01

### Added
- Initial release
- Basic shell executor
- YAML-based execution plans
- State management
- CLI commands: execute, validate, status, reset
- Cross-platform support (Linux, macOS, Windows)

[Unreleased]: https://github.com/SphereStacking/plexr/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/SphereStacking/plexr/releases/tag/v0.1.0