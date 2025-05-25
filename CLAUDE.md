# CLAUDE.md - Plexr Development Context

## Project Overview

Plexr is a developer-friendly CLI tool for automating local development environment setup. It combines "Plan + Executor" to turn setup instructions into executable YAML configurations.

### Core Concepts

- **Plan**: YAML-based configuration files that define setup steps
- **Executor**: Components that run the steps (shell commands, scripts, etc.)
- **State Management**: Tracks execution progress and shares data between steps
- **Stateful Execution**: Can resume from where it left off if something fails

## Project Structure

```
plexr/
├── cmd/plexr/          # CLI entry point
├── internal/           # Core implementation
│   ├── cli/           # CLI commands
│   ├── config/        # Configuration loading
│   ├── core/          # Core logic (runner, state)
│   ├── executors/     # Step executors
│   └── utils/         # Utilities
├── examples/          # Example configurations
├── docs/              # VitePress documentation
└── go.mod            # Go module definition
```

## Key Commands

### Development
```bash
# Build the project
make build

# Run tests
make test

# Run linting
make lint

# Install locally
make install
```

### Documentation
```bash
# Change to docs directory
cd docs

# Install dependencies (first time)
npm install

# Run development server
npm run docs:dev

# Build documentation
npm run docs:build
```

## Code Style Guidelines

1. **Go Code**:
   - Follow standard Go conventions
   - Use meaningful variable and function names
   - Add comments for exported functions and types
   - Handle errors explicitly

2. **Documentation**:
   - Use clear, concise language
   - Include practical examples
   - Keep technical accuracy
   - Support both English and Japanese

3. **YAML Examples**:
   - Use descriptive step names
   - Include comments explaining complex parts
   - Show real-world use cases

## Current Implementation Status

### Implemented ✅
- Core execution engine
- Shell executor
- State management
- Basic CLI commands (execute, validate, status, reset)
- YAML configuration parser
- Dependency resolution
- Error handling

### Not Yet Implemented ❌
- SQL executor
- Full platform detection
- Transaction modes
- Advanced executors (HTTP, Docker, etc.)
- Parallel execution
- Conditional execution
- Remote state storage

## Common Tasks

### Adding a New Command

1. Create command file in `internal/cli/`
2. Implement command logic
3. Add command to root command in `internal/cli/root.go`
4. Add tests
5. Update documentation

### Adding a New Executor

1. Create executor file in `internal/executors/`
2. Implement Executor interface
3. Register executor in registry
4. Add tests
5. Document usage

### Updating Documentation

1. Edit Markdown files in `docs/`
2. For new pages, update `.vitepress/config.js`
3. Test locally with `npm run docs:dev`
4. Ensure all links work
5. Update both English and Japanese versions

## Testing Guidelines

- Write unit tests for all new functions
- Use table-driven tests where appropriate
- Mock external dependencies
- Test error cases
- Aim for >80% coverage

## Important Files

- `internal/core/runner.go` - Main execution logic
- `internal/config/plan.go` - Plan configuration structure
- `internal/core/state.go` - State management
- `internal/executors/shell.go` - Shell executor implementation
- `cmd/plexr/main.go` - CLI entry point

## Environment Variables

- `PLEXR_STATE_DIR` - Override state directory location
- `PLEXR_LOG_LEVEL` - Set logging level (debug, info, warn, error)
- `PLEXR_NO_COLOR` - Disable colored output

## Debugging Tips

1. Use `-v` flag for verbose output
2. Check `.plexr/state.json` for execution state
3. Use `plexr status` to see current state
4. Enable debug logging with `PLEXR_LOG_LEVEL=debug`

## Release Process

1. Update version in `version.go`
2. Update CHANGELOG.md
3. Run all tests
4. Build binaries for all platforms
5. Create GitHub release
6. Update documentation

## Contributing Guidelines

1. Fork the repository
2. Create feature branch
3. Make changes with tests
4. Ensure linting passes
5. Update documentation
6. Submit pull request

## Contact

- GitHub: https://github.com/SphereStacking/plexr
- Issues: https://github.com/SphereStacking/plexr/issues

## Commit Message Convention

### Basic Structure
<type>(<scope>): <subject>

<body>

<footer>

### Types
- feat: A new feature
- fix: A bug fix
- docs: Documentation only changes
- style: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- refactor: A code change that neither fixes a bug nor adds a feature
- test: Adding missing tests or correcting existing tests
- chore: Changes to the build process or auxiliary tools and libraries

### Scopes
- core: Changes related to core functionality
- config: Changes related to configuration
- test: Changes related to testing
- docs: Changes related to documentation
- ci: Changes related to CI

### Subject
- Use the imperative, present tense: "change" not "changed" nor "changes"
- Don't capitalize the first letter
- No dot (.) at the end
- 50 characters or less

### Body
- Explain the what and why of the change
- Use bullet points for details
- Be specific and clear
- Explain the motivation for the change

### Footer
- Reference related issue numbers
- Start with "BREAKING CHANGE:" for breaking changes
- For multiple issues, use format "Closes #123, #456"

### Examples
feat(core): add file path validation

- Add validation for empty file paths
- Add validation for absolute paths
- Add validation for paths containing '..'
- Add corresponding test cases

---

This file provides context for AI assistants helping with Plexr development. Keep it updated as the project evolves.
