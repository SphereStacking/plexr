# Plexr Documentation

This directory contains the documentation for Plexr, built with VitePress.

## Development

```bash
# Install dependencies
npm install

# Start development server
npm run docs:dev

# Build for production
npm run docs:build

# Preview production build
npm run docs:preview
```

## Structure

```
docs/
├── .vitepress/         # VitePress configuration
├── guide/              # User guides
├── api/                # API reference
├── examples/           # Example configurations
├── ja/                 # Japanese translations
└── public/             # Static assets
```

## Languages

- English (default)
- Japanese (`/ja/`)

## Contributing

When adding new documentation:

1. Add English version first
2. Update navigation in `.vitepress/config.js`
3. Add Japanese translation if possible
4. Keep examples practical and tested