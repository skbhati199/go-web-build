# Installation

## Prerequisites

Before installing Go Web Build, ensure you have:

- Go 1.19 or higher
- Node.js 16 or higher
- npm 7 or higher

## Installation Methods

### Using Go Install

```bash
go install github.com/skbhati199/go-web-build@latest
```

### From Source

```bash
git clone https://github.com/skbhati199/go-web-build.git
cd go-web-build
make install
```

# Quick Start Guide

## Creating a New Project

### React Application

```bash
go-web-build create react my-app
cd my-app
go-web-build dev
go-web-build build
```

# Configuration Reference

## Configuration File (gobuild.yaml)

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

# API Reference

## CLI Commands

### create

Create a new project:
```bash
go-web-build create <template> <project-name>
```

### dev

Start the development server:
```bash
go-web-build dev [options]
```

### build

Build the project for production:
```bash
go-web-build build [options]
```

### deploy

Deploy the project:
```bash
go-web-build deploy [options]
```