# API Reference

## CLI Commands

### create

Create a new project:

```bash
go-web-build create <template> <project-name>
```

**Options:**

- `<template>`: The template to use (e.g., react, vue, basic)
- `<project-name>`: Name of your project

**Examples:**

```bash
# Create a React application
go-web-build create react my-app

# Create a basic Go web application
go-web-build create basic my-service
```

### dev

Start the development server with hot-reload:

```bash
go-web-build dev [options]
```

**Options:**

- `--port, -p`: Specify the port to run the development server (default: 3000)
- `--host`: Specify the host to bind the server to (default: localhost)
- `--no-hot`: Disable hot reloading

**Examples:**

```bash
# Start dev server on default port
go-web-build dev

# Start dev server on port 8080
go-web-build dev --port 8080
```

### build

Build the project for production:

```bash
go-web-build build [options]
```

**Options:**

- `--output, -o`: Specify the output directory (default: dist)
- `--no-minify`: Disable minification
- `--no-sourcemap`: Disable source maps

**Examples:**

```bash
# Build with default settings
go-web-build build

# Build with custom output directory
go-web-build build --output ./build
```

### deploy

Deploy the project to a hosting service:

```bash
go-web-build deploy [options]
```

**Options:**

- `--provider, -p`: Specify the deployment provider (e.g., aws, gcp, azure)
- `--env, -e`: Specify the environment to deploy to (default: production)

**Examples:**

```bash
# Deploy to default provider
go-web-build deploy

# Deploy to AWS in staging environment
go-web-build deploy --provider aws --env staging
```

## Configuration

The Go Web Build tool can be configured using a `gobuild.yaml` file in your project root. See the [Configuration](configuration.md) guide for more details.

```yaml
project:
  name: "my-app"
  version: "1.0.0"
  type: "react"

build:
  target: "web"
  output: "dist"
  optimization:
    minify: true
    splitChunks: true

dev:
  port: 3000
  hot: true
```