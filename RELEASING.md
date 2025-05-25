# Releasing Plexr

This document describes the release process for Plexr.

## Prerequisites

- Write access to the repository
- GPG key for signing commits (optional but recommended)
- [goreleaser](https://goreleaser.com/) installed locally (for testing)

## Release Process

### 1. Prepare the Release

1. Ensure all changes are merged to `main` branch
2. Run tests to ensure everything is working:
   ```bash
   make test
   make lint
   ```

3. Update the version in code if needed (though goreleaser will handle this)

4. Update CHANGELOG.md with release notes (optional, as goreleaser generates them)

### 2. Create a Release Tag

```bash
# For a new version (e.g., v1.0.0)
git tag -a v1.0.0 -m "Release v1.0.0"

# For a pre-release (e.g., v1.0.0-rc.1)
git tag -a v1.0.0-rc.1 -m "Release v1.0.0-rc.1"

# Push the tag
git push origin v1.0.0
```

### 3. GitHub Actions Will Automatically:

1. Run all tests
2. Build binaries for all platforms:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)
3. Create a GitHub Release with:
   - Changelog generated from commit messages
   - Binary archives with checksums
   - Installation instructions

### 4. Verify the Release

1. Check the [Releases page](https://github.com/SphereStacking/plexr/releases)
2. Download and test a binary
3. Verify `go install` works:
   ```bash
   go install github.com/SphereStacking/plexr/cmd/plexr@v1.0.0
   ```

## Release Naming Convention

- Production releases: `v1.0.0`
- Pre-releases: `v1.0.0-rc.1`, `v1.0.0-beta.1`
- Development builds: Automatically created by goreleaser with `-next` suffix

## Testing Releases Locally

To test the release process locally:

```bash
# Dry run
goreleaser release --snapshot --clean

# Check the dist/ directory for artifacts
ls -la dist/
```

## Troubleshooting

### Release Failed

If the GitHub Action fails:

1. Check the Actions tab for error logs
2. Common issues:
   - Tests failing
   - Invalid goreleaser configuration
   - Missing GITHUB_TOKEN

### Binary Not Working

If users report issues with binaries:

1. Check the build logs in GitHub Actions
2. Verify CGO_ENABLED=0 is set (for static binaries)
3. Test on the target platform

## Post-Release

After a successful release:

1. Announce the release (if major version)
2. Update any documentation referencing the version
3. Consider updating installation scripts/guides

## Semantic Versioning

We follow [Semantic Versioning](https://semver.org/):

- MAJOR: Incompatible API changes
- MINOR: New functionality in a backwards compatible manner
- PATCH: Backwards compatible bug fixes

## Commit Message Format

For better changelogs, use conventional commits:

- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `chore:` Maintenance tasks
- `refactor:` Code refactoring
- `test:` Test additions/changes
- `perf:` Performance improvements