# Quick Start Guide

This guide will help you get started with Go Web Build quickly. We'll cover creating a new project, running the development server, and building for production.

## Creating Your First Project

After [installing](installation.md) Go Web Build, you can create a new project using the `create` command:

```bash
go-web-build create <template> <project-name>
```

### Available Templates

- `react` - A React application with TypeScript support
- `vue` - A Vue.js application
- `basic` - A basic Go web application

### Example: Creating a React Application

```bash
# Create a new React application
go-web-build create react my-app

# Navigate to the project directory
cd my-app
```

## Project Structure

After creating a project, you'll have a structure similar to this (for a React template):

```
my-app/
├── gobuild.yaml        # Configuration file
├── go.mod              # Go module file
├── go.sum              # Go module checksum
├── src/                # Source code
│   ├── main.go         # Go entry point
│   ├── components/     # React components
│   ├── pages/          # React pages
│   └── assets/         # Static assets
└── public/             # Public files
    └── index.html      # HTML template
```

## Development

Start the development server with hot-reload:

```bash
cd my-app
go-web-build dev
```

This will start a development server (default: http://localhost:3000) with hot-reload enabled.

### Development Options

You can customize the development server:

```bash
# Start on a different port
go-web-build dev --port 8080

# Bind to a specific host
go-web-build dev --host 0.0.0.0

# Disable hot reloading
go-web-build dev --no-hot
```

## Building for Production

When you're ready to deploy your application, build it for production:

```bash
go-web-build build
```

This will create optimized files in the `dist` directory (or the output directory specified in your configuration).

### Build Options

Customize your build:

```bash
# Specify a custom output directory
go-web-build build --output ./build

# Skip minification
go-web-build build --no-minify

# Disable source maps
go-web-build build --no-sourcemap
```

## Deployment

Deploy your application to a hosting service:

```bash
go-web-build deploy
```

By default, this deploys to your configured provider in production mode. You can customize the deployment:

```bash
# Deploy to a specific provider
go-web-build deploy --provider aws

# Deploy to a specific environment
go-web-build deploy --env staging
```

## Next Steps

- Learn about [configuration options](configuration.md)
- Explore the [API reference](api.md)