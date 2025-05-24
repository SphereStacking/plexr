# API Reference

Welcome to the Plexr API reference documentation. This section provides detailed technical information about all aspects of Plexr.

## Available References

### [CLI Commands](/api/cli-commands)

Complete reference for all command-line interface commands:
- `execute` - Run execution plans
- `validate` - Validate plan syntax
- `status` - Check execution status
- `reset` - Reset execution state
- And more...

### [Configuration Schema](/api/configuration-schema)

Detailed YAML configuration schema documentation:
- Plan structure and fields
- Step definitions
- Executor configurations
- Platform-specific settings
- Validation rules

### [Executors API](/api/executors)

Information about built-in and custom executors:
- Shell executor
- SQL executor (coming soon)
- Creating custom executors
- Executor interfaces

## Quick Links

### Command Line

- [execute command](/api/cli-commands#plexr-execute) - Main execution command
- [Global flags](/api/cli-commands#global-flags) - Flags available for all commands
- [Exit codes](/api/cli-commands#exit-codes) - Understanding return values

### Configuration

- [Root fields](/api/configuration-schema#root-fields) - Top-level configuration
- [Step schema](/api/configuration-schema#step) - Step configuration details
- [File config](/api/configuration-schema#fileconfig) - File execution options

### Development

- [Executor interface](/api/executors#executor-interface) - Implementing executors
- [State management](/api/executors#state-management) - Working with state
- [Error handling](/api/executors#error-handling) - Best practices

## Environment Variables

Plexr uses several environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `PLEXR_STATE_FILE` | State file location | `.plexr_state.json` |
| `PLEXR_LOG_LEVEL` | Logging level | `info` |
| `PLEXR_NO_COLOR` | Disable colors | `false` |
| `PLEXR_PLATFORM` | Override platform | auto-detect |

## File Formats

### State File Format

The state file (`.plexr_state.json`) tracks execution progress:

```json
{
  "version": "1.0",
  "plan_name": "Development Environment Setup",
  "plan_version": "1.0.0",
  "started_at": "2023-12-15T10:00:00Z",
  "updated_at": "2023-12-15T10:30:00Z",
  "current_step": "configure_app",
  "completed_steps": [
    {
      "id": "install_tools",
      "completed_at": "2023-12-15T10:10:00Z"
    }
  ],
  "failed_steps": [],
  "installed_tools": {
    "node": "20.10.0",
    "docker": "24.0.7"
  }
}
```

### Configuration File Format

See the [Configuration Schema](/api/configuration-schema) for complete YAML format documentation.

## Error Codes

Plexr uses consistent error codes across all operations:

| Code | Category | Description |
|------|----------|-------------|
| 0 | Success | Operation completed successfully |
| 1-99 | General | General errors |
| 100-199 | Validation | Configuration validation errors |
| 200-299 | Execution | Runtime execution errors |
| 300-399 | State | State management errors |
| 400-499 | Platform | Platform-specific errors |

## Versioning

Plexr follows semantic versioning:

- **Major:** Breaking changes to CLI or configuration format
- **Minor:** New features, backward compatible
- **Patch:** Bug fixes and minor improvements

## Support

For additional help:
- [GitHub Issues](https://github.com/SphereStacking/plexr/issues)
- [Examples](/examples/)
- [Troubleshooting Guide](/guide/troubleshooting)